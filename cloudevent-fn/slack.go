package function

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/nlopes/slack"
)

// getSlackToken returns the Slack Bot User OAuth Access Token
// from the configured secret file
func getSlackToken() string {
	token, err := ioutil.ReadFile("/var/run/secrets/" + os.Getenv("slack_token"))
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimSpace(string(token))
}

// getSlackRoom returns the Slack room ID stored in the functions
// environment variables
func getSlackRoom() string {
	return strings.TrimSpace(os.Getenv("slack_room"))
}

// sendMessage to the Slack bot and room configured in stack.yml
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
