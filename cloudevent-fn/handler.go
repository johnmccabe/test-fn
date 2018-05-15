package function

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/nlopes/slack"
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

const MicrosoftStorageBlobCreatedType = "Microsoft.Storage.BlobCreated"
const OK = "OK"

// Handle a serverless request
func Handle(req []byte) string {
	// Handle Azure Subscription Validation event
	if resp := azureValidationEvent(req); resp != nil {
		return *resp
	}

	// Handle CloudEvent, logs error and exits if invalid
	c := cloudEvent(req)

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

func sendMessage(imgURL, eventType, cloudEvent string) {
	api := slack.New(getSlackToken())
	params := slack.PostMessageParameters{}
	attachment := slack.Attachment{
		ImageURL: imgURL,
		Color:    "#36a64f",
		Pretext:  fmt.Sprintf("Received Cloudevent Type: %s", eventType),
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "Raw Message",
				Value: fmt.Sprintf("```%s```", cloudEvent),
				Short: true,
			},
		},
	}
	params.Attachments = []slack.Attachment{attachment}
	channelID, timestamp, err := api.PostMessage(getSlackRoom(), "", params)
	if err != nil {
		log.Fatalf("%s\n", err)
	}
	fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
}

func azureValidationEvent(req []byte) *string {
	v := []SubscriptionValidationEvent{}
	if err := json.Unmarshal(req, &v); err == nil {
		r := SubscriptionValidationResp{}
		if len(v) > 0 && len(v[0].Data.ValidationCode) > 0 {
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

func cloudEvent(req []byte) *CloudEvent {
	c := CloudEvent{}
	if err := json.Unmarshal(req, &c); err != nil {
		log.Fatalf("Received an Unsupported Event: %s", string(req))
	}
	return &c
}

func getSlackToken() string {
	token, err := ioutil.ReadFile("/var/run/secrets/" + os.Getenv("slack_token"))
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimSpace(string(token))
}

func getSlackRoom() string {
	return strings.TrimSpace(os.Getenv("slack_room"))
}
