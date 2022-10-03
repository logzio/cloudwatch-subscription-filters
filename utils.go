package main

import (
	"fmt"
	"github.com/logzio/logzio-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strconv"
	"time"
)

const (
	LogLevelDebug = "debug"
	LogLevelInfo  = "info"
	LogLevelWarn  = "warn"
	LogLevelError = "error"
	LogLevelFatal = "fatal"
	LogLevelPanic = "panic"

	envLogLevel       = "LOG_LEVEL"
	envLogzioToken    = "LOGZIO_TOKEN"
	envLogzioListener = "LOGZIO_LISTENER"
	envLogzioType     = "LOGZIO_TYPE"
	envCompress       = "COMPRESS"

	maxBulkSizeBytes = 10 * 1024 * 1024 // 10 MB

	fieldMessage             = "message"
	fieldMessageType         = "messageType"
	fieldOwner               = "owner"
	fieldLogGroup            = "logGroup"
	fieldLogStream           = "logStream"
	fieldSubscriptionFilters = "subscriptionFilters"
	fieldLogEventId          = "id"
	fieldLogEventTimestamp   = "@timestamp"
	fieldType                = "type"

	defaultLogLevel = LogLevelInfo
	defaultType     = "logzio_cloudwatch_lambda"
	defaultCompress = true
)

func getLogger() *zap.Logger {
	logLevel := getLogLevel()
	cfg := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(logLevel),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:   "message",
			LevelKey:     "level",
			EncodeLevel:  zapcore.CapitalLevelEncoder,
			TimeKey:      "time",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}
	l, _ := cfg.Build()
	return l
}

func getLogLevel() zapcore.Level {
	logLevelStr := getHookLogLevel()
	levelsMap := map[string]zapcore.Level{
		LogLevelDebug: zapcore.DebugLevel,
		LogLevelInfo:  zapcore.InfoLevel,
		LogLevelWarn:  zapcore.WarnLevel,
		LogLevelError: zapcore.ErrorLevel,
		LogLevelPanic: zapcore.PanicLevel,
		LogLevelFatal: zapcore.FatalLevel,
	}

	return levelsMap[logLevelStr]
}

func getHookLogLevel() string {
	validLogLevels := []string{LogLevelDebug, LogLevelInfo, LogLevelWarn, LogLevelError, LogLevelFatal, LogLevelPanic}
	logLevel := os.Getenv(envLogLevel)
	for _, validLogLevel := range validLogLevels {
		if validLogLevel == logLevel {
			return validLogLevel
		}
	}

	return defaultLogLevel
}

func getNewLogzioSender() (*logzio.LogzioSender, error) {
	token, err := getToken()
	if err != nil {
		return nil, err
	}
	listener, err := getListener()
	if err != nil {
		return nil, err
	}

	logLevel := getHookLogLevel()
	compress := getCompress()
	var logzioLogger *logzio.LogzioSender
	if logLevel == LogLevelDebug {
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
		return "", fmt.Errorf("%s should be set", envLogzioToken)
	}

	return token, nil
}

func getListener() (string, error) {
	listener := os.Getenv(envLogzioListener)
	if len(listener) == 0 {
		return "", fmt.Errorf("%s must be set", envLogzioListener)
	}

	return listener, nil
}

func getCompress() bool {
	compressStr := os.Getenv(envCompress)
	if compressStr == "" {
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

func getType() string {
	logzioType := os.Getenv(envLogzioType)
	if logzioType == "" {
		logzioType = defaultType
	}

	return logzioType
}
