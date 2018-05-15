package function

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/nlopes/slack"
)

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

func sendMessage(imgURL, eventType, cloudEvent string) {
	api := slack.New(getSlackToken())
	params := slack.PostMessageParameters{}
	attachment := slack.Attachment{
		ImageURL: imgURL,
		Color:    "#36a64f",
		Pretext:  fmt.Sprintf("Received CloudEvent Type: %s", eventType),
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "Raw CloudEvent",
				Value: fmt.Sprintf("```%s```", cloudEvent),
				Short: false,
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
