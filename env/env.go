package env

import "os"

var (
	ApiURL          = os.Getenv("API_URL")
	Env             = os.Getenv("ENV")
	SlackOAuthToken = os.Getenv("SLACK_OAUTH_TOKEN")
	SlackChannelID  = os.Getenv("CHANNEL_ID")
)
