package logs_processor

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	fieldMessage             = "message"
	fieldMessageType         = "messageType"
	fieldOwner               = "owner"
	fieldLogGroup            = "logGroup"
	fieldLogStream           = "logStream"
	fieldSubscriptionFilters = "subscriptionFilters"
	fieldLogEventId          = "id"
	fieldLogEventTimestamp   = "@timestamp"
	fieldType                = "type"

	envLogzioType       = "LOGZIO_TYPE"
	envAdditionalFields = "ADDITIONAL_FIELDS"
	envSendAll          = "SEND_ALL"

	defaultType    = "logzio_cloudwatch_lambda"
	defaultSendAll = false

	prefixStart  = "START"
	prefixEnd    = "END"
	prefixReport = "REPORT"

	customFieldSeparator = ";"

	emptyString = ""
)

func getType() string {
	logzioType := os.Getenv(envLogzioType)
	if logzioType == emptyString {
		logzioType = defaultType
	}

	return logzioType
}

func getAdditionalFieldsStr() string {
	afStr := os.Getenv(envAdditionalFields)
	if afStr != emptyString {
		// remove possible whitespaces from string
		afStr = strings.ReplaceAll(afStr, " ", "")
	}

	return afStr
}

func getSendAll() bool {
	saStr := os.Getenv(envSendAll)
	if saStr == emptyString {
		return defaultSendAll
	}

	sendAll, err := strconv.ParseBool(saStr)
	if err != nil {
		logger.Info(fmt.Sprintf("Cannot handle user input for %s, error: %s", envSendAll, err.Error()))
		logger.Info(fmt.Sprintf("Reverting for default value %t", defaultSendAll))
		return defaultSendAll
	}

	return sendAll
}
