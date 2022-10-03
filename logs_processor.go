package main

import (
	"encoding/json"
	"fmt"
)

func processLogs(cwEvent CWEvent) error {
	for index, logEvent := range cwEvent.LogEvents {
		logzioLog := make(map[string]interface{})

		handleMessageField(logzioLog, logEvent.Message)

		addEventFields(logzioLog, cwEvent, index)

		addLogzioFields(logzioLog, logEvent.Timestamp)

		// TODO: Check current abilities and add them

		// TODO: convert to byte array and add to the logzioSender

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
	if event.MessageType != "" {
		logzioLog[fieldMessageType] = event.MessageType
	}

	if event.Owner != "" {
		logzioLog[fieldOwner] = event.Owner
	}

	if event.LogGroup != "" {
		logzioLog[fieldLogGroup] = event.LogGroup
	}

	if event.LogStream != "" {
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
	if ts != "" {
		logzioLog[fieldLogEventTimestamp] = ts
	}

	logzioLog[fieldType] = getType()
}
