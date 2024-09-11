package domain

import (
	"fmt"
	"regexp"
	"strings"

	"google.golang.org/api/youtube/v3"
)

type YoutubeVideoStruct struct {
	ID           uint `gorm:"primaryKey"` // Primary key field
	Channel      string
	VideoID      string
	IsProcessed  bool
	NaverLink    string
	Title        string
	PublishedAt  string
	Description  string
	YouTubeLink  string
	ThumbnailURL string
}

func (YoutubeVideoStruct) TableName() string {
	return "youtubes"
}

func ConvertToYoutubeVideoStruct(videos []*youtube.PlaylistItem, channel *CrawlingSource) []*YoutubeVideoStruct {
	var result []*YoutubeVideoStruct
	for _, item := range videos {
		videoID := item.Snippet.ResourceId.VideoId
		naverLink := extractNaverURL(item.Snippet.Description)
		row := YoutubeVideoStruct{
			Channel:      channel.SourceName,
			VideoID:      videoID,
			IsProcessed:  false,
			NaverLink:    naverLink,
			Title:        item.Snippet.Title,
			PublishedAt:  item.Snippet.PublishedAt,
			Description:  item.Snippet.Description,
			YouTubeLink:  fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID),
			ThumbnailURL: item.Snippet.Thumbnails.High.Url,
		}
		result = append(result, &row)
	}
	return result
}

func extractNaverURL(description string) string {
	re := regexp.MustCompile(`https?://naver\.me/[^\s\n]+`)
	matches := re.FindAllString(description, -1)
	return strings.Join(matches, ",")
}

func FilterWithChannelConstraints(videos []*YoutubeVideoStruct, channel *CrawlingSource) []*YoutubeVideoStruct {
	var result []*YoutubeVideoStruct
	for _, video := range videos {
		if len(channel.Constraint) < 1 || containsConstraint(video, channel.Constraint) {
			result = append(result, video)
		}
	}
	return result
}

func containsConstraint(video *YoutubeVideoStruct, constraints string) bool {
	constraintList := strings.Split(constraints, ",")

	for _, constraint := range constraintList {
		if strings.Contains(video.Title, constraint) {
			return true
		}
	}
	return false
}
