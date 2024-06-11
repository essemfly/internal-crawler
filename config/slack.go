package config

import (
	"os"
)

func GetWebhookUrl(channelName string) string {
	if channelName == "YOUTUBE" {
		return os.Getenv("YOUTUBE_WEBHOOK_URL")
	} else if channelName == "WISHKET" {
		return os.Getenv("WISHKET_WEBHOOK_URL")
	}
	return ""
}
