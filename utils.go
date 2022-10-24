package main

import (
	"fmt"
	"github.com/logzio/logzio-go"
	lp "main/logger"
	"os"
	"strconv"
	"time"
)

const (
	envLogzioToken    = "LOGZIO_TOKEN"
	envLogzioListener = "LOGZIO_LISTENER"
	envCompress       = "COMPRESS"

	defaultCompress = true

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
	compress := getCompress()
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
			logzio.SetCompress(compress),
		)
	} else {
		logzioLogger, err = logzio.New(
			token,
			logzio.SetUrl(listener),
			logzio.SetInMemoryQueue(true),
			logzio.SetDebug(os.Stdout),
			logzio.SetinMemoryCapacity(maxBulkSizeBytes), //bytes
			logzio.SetDrainDuration(time.Second*5),
			logzio.SetCompress(compress),
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

func getCompress() bool {
	compressStr := os.Getenv(envCompress)
	if compressStr == emptyString {
		return defaultCompress
	}

	compress, err := strconv.ParseBool(compressStr)
	if err != nil {
		logger.Info(fmt.Sprintf("Cannot handle user input for %s, error: %s", envCompress, err.Error()))
		logger.Info(fmt.Sprintf("Reverting for default value %t", defaultCompress))
		return defaultCompress
	}

	return compress

}
