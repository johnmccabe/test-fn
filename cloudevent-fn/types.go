package function

import (
	"encoding/json"
)

// CloudEvent v0.1
// https://github.com/cloudevents/spec/blob/v0.1/json-format.md
type CloudEvent struct {
	EventType          string
	EventTypeVersion   string
	CloudEventsVersion string
	Source             string
	EventID            string
	EventTime          string
	ContentType        string
	Extensions         map[string]string
	Data               json.RawMessage
}

// SubscriptionValidationEvent received from Azure EventGrid
// when it tests webhooks during event subscription creation
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

// SubscriptionValidationData contains a ValidationCode that
// must be returned in a response, and a ValidationUrl that
// can be used to manually validate the subscription
type SubscriptionValidationData struct {
	ValidationCode string
	ValidationUrl  string
}

// SubscriptionValidationResp returned to EventGrid in order
// to validate the event subscription
type SubscriptionValidationResp struct {
	ValidationResponse string
}

// MicrosoftStorageBlobCreated event used to demonstrate
// consumption of CloudEvents, this is included in the
// CloudEvent Data field when a blob is added to an Azure
// Storage container
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

// StorageDiagnostics occasionally included by the Azure
// Storage service. This property should be ignored by
// event consumers.
type StorageDiagnostics struct {
	BatchId string
}

// OK string to be returned
const OK = "OK"
