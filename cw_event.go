package main

type CWEventEncoded struct {
	Awslogs AwsLogsObj `json:"awslogs"`
}

type AwsLogsObj struct {
	Data string `json:"data"`
}

type CWEvent struct {
	MessageType         string     `json:"messageType"`
	Owner               string     `json:"owner"`
	LogGroup            string     `json:"logGroup"`
	LogStream           string     `json:"logStream"`
	SubscriptionFilters []string   `json:"subscriptionFilters"`
	LogEvents           []LogEvent `json:"logEvents"`
}

type LogEvent struct {
	Id        string `json:"id"`
	Timestamp int64  `json:"timestamp"`
	Message   string `json:"message"`
}
