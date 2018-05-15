package function

import (
	"encoding/json"
	"log"
)

const MicrosoftStorageBlobCreatedType = "Microsoft.Storage.BlobCreated"
const OK = "OK"

// Handle a serverless request
func Handle(req []byte) string {
	// Handle Azure Subscription Validation event
	if resp := azureValidationEvent(req); resp != nil {
		return *resp
	}

	// Handle CloudEvent, logs error and exits if invalid
	c := getCloudEvent(req)

	switch c.EventType {
	case MicrosoftStorageBlobCreatedType:
		d := MicrosoftStorageBlobCreated{}
		if err := json.Unmarshal(c.Data, &d); err != nil {
			log.Fatal("Unable to unmarshal MicrosoftStorageBlobCreated object")
		}
		sendMessage(d.Url, MicrosoftStorageBlobCreatedType, string(req))
		return OK
	default:
		log.Fatalf("Unsupported eventType received: %s", c.EventType)
	}

	return OK
}

func azureValidationEvent(req []byte) *string {
	v := []SubscriptionValidationEvent{}
	if err := json.Unmarshal(req, &v); err == nil {
		if len(v) > 0 && len(v[0].Data.ValidationCode) > 0 {
			r := SubscriptionValidationResp{}
			r.ValidationResponse = v[0].Data.ValidationCode
			if b, err := json.Marshal(r); err == nil {
				resp := string(b)
				return &resp
			}
			log.Fatalf("Unable to marshal Azure validation response")
		}
	}
	return nil
}

func getCloudEvent(req []byte) *CloudEvent {
	c := CloudEvent{}
	if err := json.Unmarshal(req, &c); err != nil {
		log.Fatalf("Received an Unsupported Event: %s", string(req))
	}
	return &c
}
