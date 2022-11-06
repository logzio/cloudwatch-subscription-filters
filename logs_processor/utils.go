package logs_processor

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
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

	envLogzioToken      = "LOGZIO_TOKEN"
	envLogzioListener   = "LOGZIO_LISTENER"
	envLogzioType       = "LOGZIO_TYPE"
	envAdditionalFields = "ADDITIONAL_FIELDS"
	envSendAll          = "SEND_ALL"
	envTimeout          = "TIMEOUT"

	defaultType    = "logzio_cloudwatch_lambda"
	defaultSendAll = false
	defaultTimeout = 10

	prefixStart  = "START"
	prefixEnd    = "END"
	prefixReport = "REPORT"

	maxLogBytesSize = 500000

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

func getToken() (string, error) {
	token := os.Getenv(envLogzioToken)
	if len(token) == 0 {
		return emptyString, fmt.Errorf("%s should be set", envLogzioToken)
	}

	return token, nil
}

func getListener() (string, error) {
	listener := os.Getenv(envLogzioListener)
	if len(listener) == 0 {
		return emptyString, fmt.Errorf("%s must be set", envLogzioListener)
	}

	return listener, nil
}

func getTimeout() time.Duration {
	timeoutStr := os.Getenv(envTimeout)
	if timeoutStr == emptyString {
		return defaultTimeout
	}

	timeoutNum, err := strconv.Atoi(timeoutStr)
	if err != nil {
		sugLog.Warnf("Could not convert properly timeout entered by user: %s", err.Error())
		sugLog.Infof("Reverting to default timeout: %d", defaultTimeout)
		timeoutNum = defaultTimeout
	}

	return time.Second * time.Duration(timeoutNum)
}
