package config

import (
	"os"
)

func GetWebhookUrl(channelName string) string {
	webhookURL := os.Getenv("WEBHOOK_URL")

	return webhookURL
}
