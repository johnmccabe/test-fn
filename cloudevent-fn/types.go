package function

import (
	"encoding/json"
)

type SubscriptionValidationEvent struct {
	Id              string
	Topic           string
	Subject         string
	Data            SubscriptionValidationData
	EventType       string
	EventTime       string
	MetadataVersion string
	DataVersion     string
}

type SubscriptionValidationData struct {
	ValidationCode string
	ValidationUrl  string
}

type SubscriptionValidationResp struct {
	ValidationResponse string
}

type CloudEvent struct {
	EventType          string
	EventTypeVersion   string
	CloudEventsVersion string
	Source             string
	EventID            string
	EventTime          string
	Data               json.RawMessage
}

type MicrosoftStorageBlobCreated struct {
	Api                string
	ClientRequestId    string
	RequestId          string
	ETag               string
	ContentType        string
	ContentLength      int
	BlobType           string
	Url                string
	Sequencer          string
	StorageDiagnostics StorageDiagnostics
}

type StorageDiagnostics struct {
	BatchId string
}
