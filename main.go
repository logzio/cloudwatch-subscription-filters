package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/logzio/logzio-go"
	"go.uber.org/zap"
)

var (
	logger       *zap.Logger
	logzioSender *logzio.LogzioSender
)

func HandleRequest(ctx context.Context, cwEvent CWEvent) {
	err := initialize(cwEvent)
	if err != nil {
		return
	}

	defer logzioSender.Drain()
	processLogs(cwEvent)
	logger.Info("Finished lambda run")
}

func main() {
	lambda.Start(HandleRequest)
}

func initialize(cwEvent CWEvent) error {
	var err error
	logger = getLogger()
	logger.Info("Starting handling event...")
	logger.Debug(fmt.Sprintf("Handling event: %+v", cwEvent))
	logger.Info("Setting up Logzio sender...")
	logzioSender, err = getNewLogzioSender()
	if err != nil {
		logger.Error(fmt.Sprintf("Error occurred while trying to setup Logzio sender: %s", err.Error()))
		logger.Error("Aborting")
		return err
	}

	logger.Info("Successfully initialized Logzio sender")
	return nil
}
