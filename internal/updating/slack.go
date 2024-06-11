package updating

import (
	"fmt"
	"strings"

	"github.com/essemfly/internal-crawler/config"
	"github.com/essemfly/internal-crawler/internal/domain"
	"github.com/essemfly/internal-crawler/pkg"
)

func SendVideosToSlack(channel *domain.CrawlingSource, videos []*domain.YoutubeVideoStruct) error {
	for _, video := range videos {
		message := fmt.Sprintf("New Video: %s\n%s\n%s", channel.SourceName, video.Title, video.YouTubeLink)

		payload := map[string]string{
			"text": message,
		}
		webhookUrl := config.GetWebhookUrl("YOUTUBE")
		err := pkg.SendToSlack(webhookUrl, payload)
		if err != nil {
			return err
		}
	}
	return nil
}

func SendWishketProjectToSlack(project *domain.ProjectInfo) error {
	message := fmt.Sprintf("프로젝트: *%s*\n> URL: %s\n> 형태: %s\n> 예상 금액: %s\n> 예상 기간: %s\n> %s\n> 지원자 수: %s\n> 분야: %s\n> 위치: %s\n> 기술: %s",
		project.Title, project.URL, project.StatusMarks, project.EstimatedAmount, project.EstimatedDuration,
		project.WorkStartDate, project.NumberOfApplicants, project.ProjectCategoryOrRole,
		project.Location, strings.Join(project.Skills, ", "))

	payload := map[string]string{
		"text": message,
	}

	webhookUrl := config.GetWebhookUrl("WISHKET")
	err := pkg.SendToSlack(webhookUrl, payload)
	if err != nil {
		return err
	}
	return nil
}
