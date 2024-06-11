package updating

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/essemfly/internal-crawler/config"
	"github.com/essemfly/internal-crawler/internal/domain"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func SaveToSheetAtTop(service *sheets.Service, channel *domain.CrawlingSource, videos []*domain.YoutubeVideoStruct) error {
	var newData [][]interface{}
	for _, video := range videos {
		row := structToSlice(video)
		if channel.Constraint != nil {
			if !strings.Contains(video.Description, *channel.Constraint) && !strings.Contains(video.Title, *channel.Constraint) {
				continue
			}
		}
		newData = append(newData, row)
	}

	readRange := fmt.Sprintf("%s!A:J", *channel.SpreadSheetName)
	resp, err := service.Spreadsheets.Values.Get(*channel.SpreadSheetID, readRange).Do()
	if err != nil {
		return fmt.Errorf("unable to retrieve data from sheet: %v", err)
	}

	mergedData := append(newData, resp.Values...)
	updateRange := fmt.Sprintf("%s!A1", *channel.SpreadSheetName)
	vr := sheets.ValueRange{
		Values: mergedData,
	}
	_, err = service.Spreadsheets.Values.Update(*channel.SpreadSheetID, updateRange, &vr).ValueInputOption("RAW").Do()
	if err != nil {
		return fmt.Errorf("unable to update sheet with new data: %v", err)
	}

	return nil
}

func GetCurrentTopVideo(service *sheets.Service, spreadsheetID, sheetName string) (*domain.YoutubeVideoStruct, error) {
	readRange := fmt.Sprintf("%s!A1:J1", sheetName)
	resp, err := service.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve data from sheet: %v", err)
	}

	var existingLatestVideo *domain.YoutubeVideoStruct
	if len(resp.Values) <= 0 || len(resp.Values[0]) <= 0 {
		return nil, nil
	}

	existingLatestVideo = &domain.YoutubeVideoStruct{
		VideoID:      resp.Values[0][0].(string),
		IsProcessed:  resp.Values[0][1].(string) == "TRUE",
		NaverLink:    resp.Values[0][2].(string),
		Title:        resp.Values[0][3].(string),
		PublishedAt:  resp.Values[0][4].(string),
		Description:  resp.Values[0][5].(string),
		YouTubeLink:  resp.Values[0][6].(string),
		ThumbnailURL: resp.Values[0][7].(string),
	}

	return existingLatestVideo, nil
}

func ListUnProcessedVideos(service *sheets.Service, channel *domain.CrawlingSource) ([]*domain.YoutubeVideoStruct, error) {
	readRange := fmt.Sprintf("%s!A:J", *channel.SpreadSheetName)
	resp, err := service.Spreadsheets.Values.Get(*channel.SpreadSheetID, readRange).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve data from sheet: %v", err)
	}

	var videos []*domain.YoutubeVideoStruct
	for _, row := range resp.Values {
		if (row[1].(string)) == "TRUE" {
			continue
		}
		video := domain.YoutubeVideoStruct{
			VideoID:      row[0].(string),
			IsProcessed:  row[1].(string) == "TRUE",
			NaverLink:    row[2].(string),
			Title:        row[3].(string),
			PublishedAt:  row[4].(string),
			Description:  row[5].(string),
			YouTubeLink:  row[6].(string),
			ThumbnailURL: row[7].(string),
		}
		videos = append(videos, &video)
	}
	return videos, nil
}

func structToSlice(video *domain.YoutubeVideoStruct) []interface{} {
	return []interface{}{
		video.VideoID,
		video.IsProcessed,
		video.NaverLink,
		video.Title,
		video.PublishedAt,
		video.Description,
		video.YouTubeLink,
		video.ThumbnailURL,
	}
}

func GetLastProjectUrl(service *sheets.Service) string {

	// Define the spreadsheet ID and range
	spreadsheetID := os.Getenv("SPREADSHEET_ID")
	sheetName := os.Getenv("SPREADSHEET_NAME")
	readRange := sheetName + "!A:B"

	log.Println("Reading from sheet...", readRange, spreadsheetID, sheetName)

	// Make the API request to get the data
	resp, err := service.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	// Process the response and print the data
	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
		return ""
	}

	// Get the last record of columns A and B
	lastRecord := resp.Values[len(resp.Values)-1]
	// datetime := lastRecord[0].(string)
	projectURL := lastRecord[1].(string)

	return projectURL
}

func UpdateCheckpoint(projectURL string) {
	ctx := context.Background()
	creds := option.WithCredentialsFile(config.JsonKeyFilePath) // Replace with your credentials file path
	service, err := sheets.NewService(ctx, creds)
	if err != nil {
		log.Fatalf("Unable to create Sheets service: %v", err)
	}

	// Add a new record with the current datetime and project URL
	datetime := time.Now().Format("2006-01-02 15:04:05")
	newRow := []interface{}{datetime, projectURL}

	// Define the spreadsheet ID and range
	spreadsheetID := os.Getenv("SPREADSHEET_ID")
	sheetName := os.Getenv("SPREADSHEET_NAME")
	writeRange := sheetName + "!A:B"

	// Prepare the data for writing
	valueRange := &sheets.ValueRange{
		Values: [][]interface{}{newRow},
	}

	// Make the API request to update the data
	_, err = service.Spreadsheets.Values.Append(spreadsheetID, writeRange, valueRange).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		log.Fatalf("Unable to update data in sheet: %v", err)
	}

	fmt.Println("New record added successfully.")

}
