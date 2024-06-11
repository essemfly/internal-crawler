package pkg

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func SendToSlack(webhookURL string, payload map[string]string) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
