package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func getWebhookUrl(channelName string) string {
	webhookURL := "https://hooks.slack.com/services/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"
	return webhookURL
}

func SendToSlack(project *ProjectInfo) error {
	message := fmt.Sprintf("새 프로젝트: *%s*\n> URL: %s\n> 예상 금액: %s\n> 예상 기간: %s\n> 시작일: %s\n> 지원자 수: %s\n> 분야: %s\n> 위치: %s\n> 기술: %s",
		project.Title, project.URL, project.EstimatedAmount, project.EstimatedDuration,
		project.WorkStartDate, project.NumberOfApplicants, project.ProjectCategoryOrRole,
		project.Location, strings.Join(project.Skills, ", "))

	payload := map[string]string{
		"text": message,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	webhookURL := getWebhookUrl("wishket")
	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
