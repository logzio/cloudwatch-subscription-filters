package logs_processor

import (
	"encoding/json"
	"fmt"
	"main/aws_structures"
	lp "main/logger"
	"strings"
)

var logger = lp.GetLogger()
var sugLog = logger.Sugar()

func ProcessLogs(cwEvent aws_structures.CWEvent) error {
	defer logger.Sync()
	logzioSender, err := initializeSender()
	if err != nil {
		return err
	}

	logsWritten := 0
	for index, logEvent := range cwEvent.LogEvents {
		if !shouldProcessLog(logEvent.Message) {
			continue
		}

		logzioLog := make(map[string]interface{})
		handleMessageField(logzioLog, logEvent.Message)

		addEventFields(logzioLog, cwEvent, index)

		addLogzioFields(logzioLog, logEvent.Timestamp)

		addAdditionalFields(logzioLog)

		err = sendLog(logzioLog, logzioSender)
		if err == nil {
			logsWritten += 1
		} else {
			sugLog.Error(err.Error())
		}
	}

	sugLog.Infof("Wrote %d logs to the Logzio Sender", logsWritten)
	return nil
}

// handleMessageField checks if the message field in the event is JSON
// If it's a JSON - it adds its fields to logzioLog.
// If not - add as a string
func handleMessageField(logzioLog map[string]interface{}, messageField string) {
	var tmpJson map[string]interface{}
	err := json.Unmarshal([]byte(messageField), &tmpJson)
	if err != nil {
		sugLog.Infof("Message %s cannot be parsed to JSON. Will be sent as a string", messageField)
		logzioLog[fieldMessage] = messageField
	} else {
		sugLog.Debug("Successfully parsed message to JSON!")
		for key, value := range tmpJson {
			logzioLog[key] = value
		}
	}
}

// addEventFields add to logzioLog the fields from the CW event, except for the timestamp field (which is handled by addLogzioFields)
func addEventFields(logzioLog map[string]interface{}, event aws_structures.CWEvent, logIndex int) {
	if event.MessageType != emptyString {
		logzioLog[fieldMessageType] = event.MessageType
	}

	if event.Owner != emptyString {
		logzioLog[fieldOwner] = event.Owner
	}

	if event.LogGroup != emptyString {
		logzioLog[fieldLogGroup] = event.LogGroup
	}

	if event.LogStream != emptyString {
		logzioLog[fieldLogStream] = event.LogStream
	}

	if len(event.SubscriptionFilters) > 0 {
		for filterIndex, filter := range event.SubscriptionFilters {
			key := fmt.Sprintf("%s_%d", fieldSubscriptionFilters, filterIndex)
			logzioLog[key] = filter
		}
	}

	logzioLog[fieldLogEventId] = event.LogEvents[logIndex].Id
}

// addLogzioFields adds to logzioLog the required Logzio fields
func addLogzioFields(logzioLog map[string]interface{}, ts int64) {
	if ts != 0 {
		logzioLog[fieldLogEventTimestamp] = ts
	}
}

// addAdditionalFields adds custom fields added by the user
func addAdditionalFields(logzioLog map[string]interface{}) {
	fieldsStr := getAdditionalFieldsStr()
	if fieldsStr != emptyString {
		keyValSeparator := "="
		indexKey := 0
		indexVal := 1

		fieldsArr := strings.Split(fieldsStr, customFieldSeparator)
		for _, field := range fieldsArr {
			keyVal := strings.Split(field, keyValSeparator)
			logzioLog[keyVal[indexKey]] = keyVal[indexVal]
		}
	}
}

// sendLog converts the log to a byte array ([]byte) and writes to the logzioSender
func sendLog(logzioLog map[string]interface{}, sender LogzioSender) error {
	logBytes, err := json.Marshal(logzioLog)
	if err != nil {
		return fmt.Errorf("Log will be dropped - error occurred while processing %s: %s", logzioLog, err.Error())
	}

	if len(logBytes) > maxLogBytesSize {
		sugLog.Debug("log dropped: %s", string(logBytes))
		return fmt.Errorf("Log will be dropped - log size is bigger than %d", maxLogBytesSize)
	}

	return sender.SendToLogzio(logBytes)
}

// processLog returns whether a log should be processed or not.
// Based on user input - we can filter out lambda platform logs (START, END, REPORT).
func shouldProcessLog(message string) bool {
	prefixList := []string{prefixStart, prefixEnd, prefixReport}
	if getSendAll() {
		return true
	} else {
		for _, prefix := range prefixList {
			if strings.HasPrefix(message, prefix) {
				sugLog.Debug("Found a Lambda platform log (START, END or REPORT). Ignoring.")
				return false
			}
		}

		return true
	}
}
