package cloudwatch_shipper

import "time"

type EventbridgeEvent struct {
	Version    string        `json:"version"`
	Id         string        `json:"id"`
	DetailType string        `json:"detail-type"`
	Source     string        `json:"source"`
	Account    string        `json:"account"`
	Time       time.Time     `json:"time"`
	Region     string        `json:"region"`
	Resources  []interface{} `json:"resources"`
	Detail     DetailObj     `json:"detail"`
}

type DetailObj struct {
	EventVersion       string               `json:"eventVersion"`
	UserIdentity       UserIdentityObj      `json:"userIdentity"`
	EventTime          time.Time            `json:"eventTime"`
	EventSource        string               `json:"eventSource"`
	EventName          string               `json:"eventName"`
	AwsRegion          string               `json:"awsRegion"`
	SourceIPAddress    string               `json:"sourceIPAddress"`
	UserAgent          string               `json:"userAgent"`
	RequestParameters  RequestParametersObj `json:"requestParameters"`
	ResponseElements   string               `json:"responseElements"`
	RequestID          string               `json:"requestID"`
	EventID            string               `json:"eventID"`
	ReadOnly           bool                 `json:"readOnly"`
	EventType          string               `json:"eventType"`
	ApiVersion         string               `json:"apiVersion"`
	ManagementEvent    bool                 `json:"managementEvent"`
	RecipientAccountId string               `json:"recipientAccountId"`
	EventCategory      string               `json:"eventCategory"`
	TlsDetails         TlsDetailsObj        `json:"tlsDetails"`
}

type UserIdentityObj struct {
	Type           string            `json:"type"`
	PrincipalId    string            `json:"principalId"`
	Arn            string            `json:"arn"`
	AccountId      string            `json:"accountId"`
	AccessKeyId    string            `json:"accessKeyId"`
	SessionContext SessionContextObj `json:"sessionContext"`
}

type RequestParametersObj struct {
	LogGroupName string `json:"logGroupName"`
}

type TlsDetailsObj struct {
	TlsVersion               string `json:"tlsVersion"`
	CipherSuite              string `json:"cipherSuite"`
	ClientProvidedHostHeader string `json:"clientProvidedHostHeader"`
}

type SessionContextObj struct {
	SessionIssuer       SessionIssuerObj       `json:"sessionIssuer"`
	WebIdFederationData map[string]interface{} `json:"webIdFederationData"`
	Attributes          AttributesObj          `json:"attributes"`
}

type SessionIssuerObj struct {
	Type        string `json:"type"`
	PrincipalId string `json:"principalId"`
	Arn         string `json:"arn"`
	AccountId   string `json:"accountId"`
	UserName    string `json:"userName"`
}

type AttributesObj struct {
	CreationDate     time.Time `json:"creationDate"`
	MfaAuthenticated string    `json:"mfaAuthenticated"`
}
