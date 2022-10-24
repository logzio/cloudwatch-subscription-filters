package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/logzio/logzio-go"
	"go.uber.org/zap"
	"io/ioutil"
	"main/aws_structures"
	lp "main/logger"
	"main/logs_processor"
)

var (
	logger       *zap.Logger
	logzioSender *logzio.LogzioSender
)

func HandleRequest(ctx context.Context, cwEventEncoded aws_structures.CWEventEncoded) {
	err := initialize(cwEventEncoded)
	if err != nil {
		return
	}
	defer logzioSender.Drain()
	logger.Debug(fmt.Sprintf("CW event encoded: %v", cwEventEncoded))
	cwEvent, err := decodeCwEvent(cwEventEncoded.Awslogs.Data)
	if err != nil {
		logger.Error("Aborting")
		return
	}

	logger.Debug(fmt.Sprintf("CW event: %v", cwEvent))
	logger.Debug(fmt.Sprintf("Detected %d logs in event", len(cwEvent.LogEvents)))
	logs_processor.ProcessLogs(cwEvent, logzioSender)
	logger.Info("Finished lambda run, draining Logzio Sender")
}

func main() {
	lambda.Start(HandleRequest)
}

func initialize(cwEventEncoded aws_structures.CWEventEncoded) error {
	var err error
	logger = lp.GetLogger()
	logger.Info("Starting handling event...")
	logger.Debug(fmt.Sprintf("Handling event: %+v", cwEventEncoded))
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

func decodeCwEvent(encodedData string) (aws_structures.CWEvent, error) {
	decoded, err := base64.StdEncoding.DecodeString(encodedData)
	var cwEvent aws_structures.CWEvent
	if err != nil {
		logger.Error(fmt.Sprintf("Error occurred while trying to decode data %s: %s", encodedData, err.Error()))
		return cwEvent, err
	}

	r, err := gzip.NewReader(bytes.NewBuffer(decoded))
	defer r.Close()
	if err != nil {
		logger.Error(fmt.Sprintf("Error occurred while trying to create new reader: %s", err.Error()))
		return cwEvent, err
	}

	data, err := ioutil.ReadAll(r)
	if err != nil {
		logger.Error(fmt.Sprintf("Error occurred while trying to read data: %s", err.Error()))
		return cwEvent, err
	}

	err = json.Unmarshal(data, &cwEvent)
	if err != nil {
		logger.Error(fmt.Sprintf("Error occurred while trying to unmarshal zipped data: %s", err.Error()))
		return cwEvent, err
	}

	return cwEvent, nil
}
