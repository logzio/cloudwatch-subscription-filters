package main

import (
	"fmt"
	"github.com/logzio/logzio-go"
	lp "main/logger"
	"os"
	"time"
)

const (
	envLogzioToken    = "LOGZIO_TOKEN"
	envLogzioListener = "LOGZIO_LISTENER"

	emptyString      = ""
	maxBulkSizeBytes = 10 * 1024 * 1024 // 10 MB
)

func getNewLogzioSender() (*logzio.LogzioSender, error) {
	token, err := getToken()
	if err != nil {
		return nil, err
	}
	listener, err := getListener()
	if err != nil {
		return nil, err
	}

	logLevel := lp.GetFuncLogLevel()
	var logzioLogger *logzio.LogzioSender
	if logLevel == lp.LogLevelDebug {
		logzioLogger, err = logzio.New(
			token,
			logzio.SetUrl(listener),
			logzio.SetInMemoryQueue(true),
			logzio.SetDebug(os.Stdout),
			logzio.SetinMemoryCapacity(maxBulkSizeBytes), //bytes
			logzio.SetDrainDuration(time.Second*5),
			logzio.SetDebug(os.Stdout),
			logzio.SetCompress(true),
		)
	} else {
		logzioLogger, err = logzio.New(
			token,
			logzio.SetUrl(listener),
			logzio.SetInMemoryQueue(true),
			logzio.SetDebug(os.Stdout),
			logzio.SetinMemoryCapacity(maxBulkSizeBytes), //bytes
			logzio.SetDrainDuration(time.Second*5),
			logzio.SetCompress(true),
		)
	}

	if err != nil {
		return nil, err
	}

	return logzioLogger, nil
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
