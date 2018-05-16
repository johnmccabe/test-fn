package function

import (
	"encoding/json"
	"log"
)

const MicrosoftStorageBlobCreatedType = "Microsoft.Storage.BlobCreated"

// Handle a serverless request
func Handle(req []byte) string {
	if resp := azureValidationEvent(req); resp != nil {
		return *resp
	}

	c := getCloudEvent(req)

	switch c.EventType {
	case MicrosoftStorageBlobCreatedType:
		d := MicrosoftStorageBlobCreated{}
		if err := json.Unmarshal(c.Data, &d); err != nil {
			log.Fatalf("Unable to unmarshal object for: %s", MicrosoftStorageBlobCreatedType)
		}
		sendMessage(d.Url, MicrosoftStorageBlobCreatedType, string(req))
		return OK
	default:
		log.Fatalf("Unsupported eventType received: %s", c.EventType)
	}

	return OK
}

// azureValidationEvent handles a received Azure Subscription Validation
// event, returning the ValidationCode extracted from the request
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

// getCloudEvent returns a pointer to a CloudEvent extracted from the
// request submitted to the handler
func getCloudEvent(req []byte) *CloudEvent {
	c := CloudEvent{}
	if err := json.Unmarshal(req, &c); err != nil {
		log.Fatalf("Received an Unsupported Event: %s", string(req))
	}
	return &c
}
