package logs_processor

import (
	"encoding/json"
	"fmt"
	"github.com/logzio/logzio-go"
	"main/aws_structures"
	lp "main/logger"
	"strings"
)

var logger = lp.GetLogger()

func ProcessLogs(cwEvent aws_structures.CWEvent, sender *logzio.LogzioSender) {
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

		logsWritten += sendLog(logzioLog, sender)
	}

	logger.Info(fmt.Sprintf("Wrote %d logs to the Logzio Sender", logsWritten))
}

// handleMessageField checks if the message field in the event is JSON
// If it's a JSON - it adds its fields to logzioLog.
// If not - add as a string
func handleMessageField(logzioLog map[string]interface{}, messageField string) {
	var tmpJson map[string]interface{}
	err := json.Unmarshal([]byte(messageField), &tmpJson)
	if err != nil {
		logger.Info(fmt.Sprintf("Message %s cannot be parsed to JSON. Will be sent as a string", messageField))
		logzioLog[fieldMessage] = messageField
	} else {
		logger.Debug("Successfully parsed message to JSON!")
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

	logzioLog[fieldType] = getType()
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
// returns the number of logs that successfully written to the logzio sender
func sendLog(logzioLog map[string]interface{}, sender *logzio.LogzioSender) int {
	logBytes, err := json.Marshal(logzioLog)
	if err != nil {
		logger.Error(fmt.Sprintf("Error occurred while processing %s: %s", logzioLog, err.Error()))
		logger.Error("Log will be dropped")
		return 0
	}

	if logBytes != nil && len(logBytes) > 0 {
		_, err = sender.Write(logBytes)
		if err != nil {
			logger.Debug(fmt.Sprintf("Error for log %s", string(logBytes)))
			logger.Error(fmt.Sprintf("Error occurred while writing log to logzio sender: %s", err.Error()))
			logger.Error("Log will be dropped")
			return 0
		}

		return 1
	}

	return 0
}

// processLog returns whether a log should be processed or not.
// Based on user input - we can filter out lambda platform logs (START, END, REPORT).
func shouldProcessLog(message string) bool {
	process := true
	prefixList := []string{prefixStart, prefixEnd, prefixReport}
	if getSendAll() {
		return process
	} else {
		for _, prefix := range prefixList {
			if strings.HasPrefix(message, prefix) {
				logger.Info("Found a Lambda platform log (START, END or REPORT). Ignoring.")
				return !process
			}
		}

		return process
	}
}
