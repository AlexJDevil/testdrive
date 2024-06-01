package alert

import (
	"log"

	"github.com/arun0009/testdrive/env"
	"github.com/slack-go/slack"
)

type (
	Notifier interface {
		Notify(message string)
	}

	SlackNotifer struct{}
)

func (s SlackNotifer) Notify(message string) {
	api := slack.New(env.SlackOAuthToken)
	_, _, err := api.PostMessage(env.SlackChannelID, slack.MsgOptionText(message, false), slack.MsgOptionUser("testdrive"))
	if err != nil {
		log.Fatalf("Error posting message to slack: %+v\n", err)
	}
}
