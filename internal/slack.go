package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func SendToSlack(project *ProjectInfo) error {
	message := fmt.Sprintf("프로젝트: *%s*\n> URL: %s\n> 형태: %s\n> 예상 금액: %s\n> 예상 기간: %s\n> %s\n> 지원자 수: %s\n> 분야: %s\n> 위치: %s\n> 기술: %s",
		project.Title, project.URL, project.StatusMarks, project.EstimatedAmount, project.EstimatedDuration,
		project.WorkStartDate, project.NumberOfApplicants, project.ProjectCategoryOrRole,
		project.Location, strings.Join(project.Skills, ", "))

	payload := map[string]string{
		"text": message,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	webhookURL := os.Getenv("WEBHOOK_URL")
	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
