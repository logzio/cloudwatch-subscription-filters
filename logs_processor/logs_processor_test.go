package logs_processor

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"main/aws_structures"
	"os"
	"testing"
)

func TestHandleRequestText(t *testing.T) {
	logzioLog := make(map[string]interface{})
	message := "this is a text log"
	handleMessageField(logzioLog, message)
	assert.Contains(t, logzioLog, fieldMessage)
	assert.Equal(t, logzioLog[fieldMessage].(string), message)
}

func TestHandleRequestJson(t *testing.T) {
	logzioLog := make(map[string]interface{})
	message := "{\"key1\":\"val1\",\"key2\":\"val2\"}"
	handleMessageField(logzioLog, message)
	assert.Contains(t, logzioLog, "key1")
	assert.Contains(t, logzioLog, "key2")
	assert.Equal(t, logzioLog["key1"].(string), "val1")
	assert.Equal(t, logzioLog["key2"].(string), "val2")
}

func TestAddEventFields(t *testing.T) {
	logzioLog := make(map[string]interface{})
	cwEvent := aws_structures.CWEvent{
		MessageType:         "DATA_MESSAGE",
		Owner:               "12345678",
		LogGroup:            "/some/log/group",
		LogStream:           "/some/log/stream",
		SubscriptionFilters: []string{"some_sub"},
		LogEvents: []aws_structures.LogEvent{aws_structures.LogEvent{
			Id:        "987654321",
			Timestamp: 1666605290,
			Message:   "this is a message",
		}},
	}

	subscriptionFiltersField := fmt.Sprintf("%s_0", fieldSubscriptionFilters)
	addEventFields(logzioLog, cwEvent, 0)
	assert.Contains(t, logzioLog, fieldMessageType)
	assert.Equal(t, logzioLog[fieldMessageType].(string), cwEvent.MessageType)
	assert.Contains(t, logzioLog, fieldOwner)
	assert.Equal(t, logzioLog[fieldOwner].(string), cwEvent.Owner)
	assert.Contains(t, logzioLog, fieldLogGroup)
	assert.Equal(t, logzioLog[fieldLogGroup].(string), cwEvent.LogGroup)
	assert.Contains(t, logzioLog, fieldLogStream)
	assert.Equal(t, logzioLog[fieldLogStream].(string), cwEvent.LogStream)
	assert.Contains(t, logzioLog, subscriptionFiltersField)
	assert.Equal(t, logzioLog[subscriptionFiltersField].(string), cwEvent.SubscriptionFilters[0])
	assert.Contains(t, logzioLog, fieldLogEventId)
	assert.Equal(t, logzioLog[fieldLogEventId].(string), cwEvent.LogEvents[0].Id)
}

func TestAddLogzioFieldsTimestamp(t *testing.T) {
	ts := int64(1666605910)
	logzioLog := make(map[string]interface{})
	addLogzioFields(logzioLog, ts)
	assert.Contains(t, logzioLog, fieldLogEventTimestamp)
	assert.Equal(t, logzioLog[fieldLogEventTimestamp].(int64), ts)
}

func TestAddLogzioFieldsEmptyTimestamp(t *testing.T) {
	logzioLog := make(map[string]interface{})
	addLogzioFields(logzioLog, 0)
	assert.NotContains(t, logzioLog, fieldLogEventTimestamp)
}

func TestAddLogzioFieldsType(t *testing.T) {
	logType := "my_type"
	os.Setenv(envLogzioType, logType)
	logzioLog := make(map[string]interface{})
	addLogzioFields(logzioLog, 0)
	assert.Contains(t, logzioLog, fieldType)
	assert.Equal(t, logzioLog[fieldType].(string), logType)
}

func TestAddLogzioFieldsEmptyType(t *testing.T) {
	os.Setenv(envLogzioType, "")
	logzioLog := make(map[string]interface{})
	addLogzioFields(logzioLog, 0)
	assert.Contains(t, logzioLog, fieldType)
	assert.Equal(t, logzioLog[fieldType].(string), defaultType)
}

func TestAddAdditionalFields(t *testing.T) {
	afStr := "key1=val1;key2=val2"
	os.Setenv(envAdditionalFields, afStr)
	logzioLog := make(map[string]interface{})
	addAdditionalFields(logzioLog)
	assert.Contains(t, logzioLog, "key1")
	assert.Equal(t, logzioLog["key1"], "val1")
	assert.Contains(t, logzioLog, "key2")
	assert.Equal(t, logzioLog["key2"], "val2")
}

func TestAddAdditionalFieldsEmpty(t *testing.T) {
	os.Setenv(envAdditionalFields, "")
	logzioLog := make(map[string]interface{})
	addAdditionalFields(logzioLog)
	assert.Len(t, logzioLog, 0)
}
