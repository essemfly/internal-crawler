package crawling

import (
	"fmt"
	"log"
	"net/http"

	"github.com/essemfly/internal-crawler/internal/domain"
	"google.golang.org/api/youtube/v3"
)

func GetUploadsPlaylistID(service *youtube.Service, channelID string) (string, error) {
	call := service.Channels.List([]string{"contentDetails"}).Id(channelID).MaxResults(1)
	response, err := call.Do()
	if err != nil {
		log.Fatalf("Error fetching channel details: %v", err)
	}
	if len(response.Items) == 0 {
		log.Fatalf("No channel found")
	}
	uploadsPlaylistID := response.Items[0].ContentDetails.RelatedPlaylists.Uploads
	return uploadsPlaylistID, nil
}

func GetChannelVideos(service *youtube.Service, channelID string, latestVideo *domain.YoutubeVideoStruct) ([]*youtube.PlaylistItem, error) {
	uploadsPlaylistID, err := GetUploadsPlaylistID(service, channelID)
	if err != nil {
		log.Fatalf("Error fetching channel details: %v", err)
		return nil, err
	}

	var items []*youtube.PlaylistItem
	nextPageToken := ""
	for {
		call := service.PlaylistItems.List([]string{"snippet"}).PlaylistId(uploadsPlaylistID).MaxResults(50).PageToken(nextPageToken)
		response, err := call.Do()
		if err != nil {
			log.Fatalf("Error making API call: %v", err)
			return nil, err
		}

		meetLatestVideo := false
		var newItems []*youtube.PlaylistItem
		for _, item := range response.Items {
			if latestVideo != nil && item.Snippet.ResourceId.VideoId == latestVideo.VideoID {
				meetLatestVideo = true
				break
			}

			isShort, err := isYoutubeShort(item.Snippet.ResourceId.VideoId)
			if err != nil {
				log.Printf("Error checking if video is short: %v", err)
				// 에러가 나도 일단 비디오는 추가
				newItems = append(newItems, item)
			} else if !isShort { // shorts가 아닌 경우만 추가하거나, 필요에 따라 처리
				newItems = append(newItems, item)
			}
		}

		items = append(items, newItems...)

		if meetLatestVideo || response.NextPageToken == "" {
			break
		}

		nextPageToken = response.NextPageToken
	}

	return items, nil
}

// GetChannelStats fetches the statistics of a YouTube channel
func GetChannelStats(service *youtube.Service, channelID string) (*youtube.Channel, error) {
	call := service.Channels.List([]string{"statistics"}).
		Id(channelID)

	response, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("error fetching channel statistics: %v", err)
	}

	if len(response.Items) == 0 {
		return nil, fmt.Errorf("no channel found for ID: %s", channelID)
	}

	return response.Items[0], nil
}

func isYoutubeShort(videoID string) (bool, error) {
	url := fmt.Sprintf("https://www.youtube.com/shorts/%s", videoID)

	client := &http.Client{
		// 리다이렉트를 따르지 않도록 설정
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Head(url)
	if err != nil {
		return false, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// 200이면 shorts, 303이면 일반 동영상
	return resp.StatusCode == http.StatusOK, nil
}
