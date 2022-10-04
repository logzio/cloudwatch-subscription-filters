package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

func processLogs(cwEvent CWEvent) {
	for index, logEvent := range cwEvent.LogEvents {
		logzioLog := make(map[string]interface{})

		handleMessageField(logzioLog, logEvent.Message)

		addEventFields(logzioLog, cwEvent, index)

		addLogzioFields(logzioLog, logEvent.Timestamp)

		// TODO: Check current abilities and add them
		addAdditionalFields(logzioLog)

		sendLog(logzioLog)

	}
}

// handleMessageField checks if the message field in the event is JSON
// If it's a JSON - it adds its fields to logzioLog.
// If not - add as a string
func handleMessageField(logzioLog map[string]interface{}, messageField string) {
	var tmpJson map[string]interface{}
	err := json.Unmarshal([]byte(messageField), &tmpJson)
	if err != nil {
		logger.Info(fmt.Sprintf("Message %s cannot be parsed to JSON. Will be sent as a string"))
		logzioLog[fieldMessage] = messageField
	}

	logger.Debug("Successfully parsed message to JSON!")
	for key, value := range tmpJson {
		logzioLog[key] = value
	}
}

// addEventFields add to logzioLog the fields from the CW event, except for the timestamp field (which is handled by addLogzioFields)
func addEventFields(logzioLog map[string]interface{}, event CWEvent, logIndex int) {
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
func addLogzioFields(logzioLog map[string]interface{}, ts string) {
	if ts != emptyString {
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
func sendLog(logzioLog map[string]interface{}) {
	logBytes, err := json.Marshal(logzioLog)
	if err != nil {
		logger.Error(fmt.Sprintf("Error occurred while processing %s: %s", logzioLog, err.Error()))
		logger.Error("Log will be dropped")
	}

	if logBytes != nil && len(logBytes) > 0 {
		_, err = logzioSender.Write(logBytes)
		if err != nil {
			logger.Debug(fmt.Sprintf("Error for log %s", string(logBytes)))
			logger.Error(fmt.Sprintf("Error occurred while writing log to logzio sender: %s", err.Error()))
			logger.Error("Log will be dropped")
		}
	}
}
