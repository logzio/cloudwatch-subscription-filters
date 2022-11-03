package logs_processor

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type LogzioSender struct {
	Url        string
	HttpClient *http.Client
}

func initializeSender() (LogzioSender, error) {
	var logzioSender LogzioSender

	token, err := getToken()
	if err != nil {
		return logzioSender, err
	}

	listener, err := getListener()
	if err != nil {
		return logzioSender, err
	}

	client := &http.Client{
		Timeout: getTimeout(),
	}

	return LogzioSender{
		Url:        fmt.Sprintf("%s/token=%s", listener, token),
		HttpClient: client,
	}, nil
}

func (l *LogzioSender) SendToLogzio(bytesToSend []byte) error {
	var statusCode int
	var compressedBuf bytes.Buffer
	gzipWriter := gzip.NewWriter(&compressedBuf)
	_, err := gzipWriter.Write(bytesToSend)
	if err != nil {
		return err
	}

	err = gzipWriter.Close()
	if err != nil {
		return err
	}

	// retry logic
	backOff := time.Second * 2
	sendRetries := 4
	toBackOff := false
	for attempt := 0; attempt < sendRetries; attempt++ {
		if toBackOff {
			fmt.Printf("Failed to send logs, trying again in %v\n", backOff)
			time.Sleep(backOff)
			backOff *= 2
		}
		statusCode = l.makeHttpRequest(compressedBuf)
		if l.shouldRetry(statusCode) {
			toBackOff = true
		} else {
			break
		}
	}

	if statusCode != 200 {
		sugLog.Errorf("Error sending logs, status code is: %d", statusCode)
	}

	compressedBuf.Reset()
	return nil
}

func (l *LogzioSender) shouldRetry(statusCode int) bool {
	retry := true
	switch statusCode {
	case http.StatusBadRequest:
		retry = false
	case http.StatusNotFound:
		retry = false
	case http.StatusUnauthorized:
		retry = false
	case http.StatusForbidden:
		retry = false
	case http.StatusOK:
		retry = false
	}

	sugLog.Info("Got HTTP %d. Should retry? %t", statusCode, retry)

	return retry
}

func (l *LogzioSender) makeHttpRequest(data bytes.Buffer) int {
	req, err := http.NewRequest("POST", l.Url, &data)
	req.Header.Add("Content-Encoding", "gzip")
	resp, err := l.HttpClient.Do(req)
	if err != nil {
		fmt.Printf("Error sending logs to %s %s\n", l.Url, err)
		return 400
	}

	defer resp.Body.Close()
	statusCode := resp.StatusCode
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v", err)
	}

	sugLog.Debugf("Response body: %s", string(respBody))

	return statusCode
}
