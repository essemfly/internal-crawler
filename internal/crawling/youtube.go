package crawling

import (
	"log"

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
			newItems = append(newItems, item)
		}

		items = append(items, newItems...)

		if meetLatestVideo || response.NextPageToken == "" {
			break
		}

		nextPageToken = response.NextPageToken
	}

	return items, nil
}
