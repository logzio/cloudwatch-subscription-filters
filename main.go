package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
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
	sugLog       *zap.SugaredLogger
	logzioSender *logzio.LogzioSender
)

func HandleRequest(ctx context.Context, cwEventEncoded aws_structures.CWEventEncoded) {
	err := initialize(cwEventEncoded)
	if err != nil {
		return
	}
	defer logger.Sync()
	defer logzioSender.Drain()
	sugLog.Debugf("CW event encoded: %v", cwEventEncoded)
	cwEvent, err := decodeCwEvent(cwEventEncoded.Awslogs.Data)
	if err != nil {
		sugLog.Error("Aborting")
		return
	}

	sugLog.Debugf("CW event: %v", cwEvent)
	sugLog.Debugf("Detected %d logs in event", len(cwEvent.LogEvents))
	logs_processor.ProcessLogs(cwEvent, logzioSender)
	sugLog.Info("Finished lambda run, draining Logzio Sender")
}

func main() {
	lambda.Start(HandleRequest)
}

func initialize(cwEventEncoded aws_structures.CWEventEncoded) error {
	var err error
	logger = lp.GetLogger()
	sugLog = logger.Sugar()
	sugLog.Info("Starting handling event...")
	sugLog.Debugf("Handling event: %+v", cwEventEncoded)
	sugLog.Info("Setting up Logzio sender...")
	logzioSender, err = getNewLogzioSender()
	if err != nil {
		sugLog.Errorf("Error occurred while trying to setup Logzio sender: %s", err.Error())
		sugLog.Error("Aborting")
		return err
	}

	sugLog.Info("Successfully initialized Logzio sender")
	return nil
}

func decodeCwEvent(encodedData string) (aws_structures.CWEvent, error) {
	decoded, err := base64.StdEncoding.DecodeString(encodedData)
	var cwEvent aws_structures.CWEvent
	if err != nil {
		sugLog.Errorf("Error occurred while trying to decode data %s: %s", encodedData, err.Error())
		return cwEvent, err
	}

	r, err := gzip.NewReader(bytes.NewBuffer(decoded))
	defer r.Close()
	if err != nil {
		sugLog.Errorf("Error occurred while trying to create new reader: %s", err.Error())
		return cwEvent, err
	}

	data, err := ioutil.ReadAll(r)
	if err != nil {
		sugLog.Error("Error occurred while trying to read data: %s", err.Error())
		return cwEvent, err
	}

	err = json.Unmarshal(data, &cwEvent)
	if err != nil {
		sugLog.Errorf("Error occurred while trying to unmarshal zipped data: %s", err.Error())
		return cwEvent, err
	}

	return cwEvent, nil
}
