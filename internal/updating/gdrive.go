package updating

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/essemfly/internal-crawler/config"
	"github.com/essemfly/internal-crawler/internal/domain"
	"github.com/essemfly/internal-crawler/pkg"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// For YOUTUBE CRAWLER
func SaveToSheetAtTop(service *sheets.Service, channel *domain.CrawlingSource, videos []*domain.YoutubeVideoStruct) error {
	var newData [][]interface{}
	for _, video := range videos {
		row := structToSlice(video)
		if len(channel.Constraint) > 0 {
			for _, constraint := range channel.Constraint {
				if !strings.Contains(video.Description, constraint) && !strings.Contains(video.Title, constraint) {
					continue
				}
			}
		}
		newData = append(newData, row)
	}

	readRange := fmt.Sprintf("%s!A:J", channel.SpreadSheetName)
	resp, err := service.Spreadsheets.Values.Get(channel.SpreadSheetID, readRange).Do()
	if err != nil {
		return fmt.Errorf("unable to retrieve data from sheet: %v", err)
	}

	mergedData := append(newData, resp.Values...)
	updateRange := fmt.Sprintf("%s!A1", channel.SpreadSheetName)
	vr := sheets.ValueRange{
		Values: mergedData,
	}
	_, err = service.Spreadsheets.Values.Update(channel.SpreadSheetID, updateRange, &vr).ValueInputOption("RAW").Do()
	if err != nil {
		return fmt.Errorf("unable to update sheet with new data: %v", err)
	}

	return nil
}

func GetCurrentTopVideo(service *sheets.Service, spreadsheetID, sheetName string) (*domain.YoutubeVideoStruct, error) {
	readRange := fmt.Sprintf("%s!A1:J1", sheetName)
	log.Println("Reading from sheet...", readRange, spreadsheetID, sheetName)
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

// Depreciated: IsProcessed = TRUE / FALSE로 구분하여 사용하였으나, 더이상 사용하지 않음
func ListUnProcessedVideos(service *sheets.Service, channel *domain.CrawlingSource) ([]*domain.YoutubeVideoStruct, error) {
	readRange := fmt.Sprintf("%s!A:J", channel.SpreadSheetName)
	resp, err := service.Spreadsheets.Values.Get(channel.SpreadSheetID, readRange).Do()
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

// FOR WISHKET
func GetLastProjectUrl(service *sheets.Service, channel *domain.CrawlingSource) string {

	// Define the spreadsheet ID and range
	spreadsheetID := channel.SpreadSheetID
	sheetName := channel.SpreadSheetName
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

func UpdateCheckpoint(projectURL string, channel *domain.CrawlingSource) {
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
	spreadsheetID := channel.SpreadSheetID
	sheetName := channel.SpreadSheetName
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

// For DAANGN CRAWLER
func SaveToSheetAppend(service *sheets.Service, channel *domain.CrawlingSource, products []*domain.DaangnProduct) error {
	var newData [][]interface{}
	for _, product := range products {
		row := []interface{}{
			product.DanggnIndex,
			product.Keyword,
			product.KeywordGroup,
			product.Name,
			product.Description,
			product.Price,
			strings.Join(product.Images, ","),
			product.Status,
			product.Url,
			product.ViewCounts,
			product.LikeCounts,
			product.ChatCounts,
			product.CrawlCategory,
			product.SellerNickName,
			product.SellerRegionName,
			product.SellerTemperature,
			product.WrittenAt.Format(time.RFC3339),
			product.CreatedAt.Format(time.RFC3339),
			product.UpdatedAt.Format(time.RFC3339),
		}
		if len(channel.Constraint) > 0 {
			for _, constraint := range channel.Constraint {
				if !strings.Contains(product.Description, constraint) && !strings.Contains(product.Name, constraint) {
					continue
				}
			}
		}
		newData = append(newData, row)
	}

	// Read existing data from the sheet to determine the last row
	readRange := fmt.Sprintf("%s!A3:S", channel.SpreadSheetName)
	resp, err := service.Spreadsheets.Values.Get(channel.SpreadSheetID, readRange).Do()
	if err != nil {
		return fmt.Errorf("unable to retrieve data from sheet: %v", err)
	}

	// Determine the next row to insert new data
	nextRow := len(resp.Values) + 3

	// Update the sheet starting from the next available row with new data
	updateRange := fmt.Sprintf("%s!A%d", channel.SpreadSheetName, nextRow)
	vr := sheets.ValueRange{
		Values: newData,
	}
	_, err = service.Spreadsheets.Values.Update(channel.SpreadSheetID, updateRange, &vr).ValueInputOption("RAW").Do()
	if err != nil {
		return fmt.Errorf("unable to update sheet with new data: %v", err)
	}

	return nil
}

// For DAUM CAFE CRAWLER
func GetLastGuestArticle(service *sheets.Service, channel *domain.CrawlingSource) (*domain.GuestArticle, error) {
	// Define the spreadsheet ID and range
	spreadsheetID := channel.SpreadSheetID
	sheetName := channel.SpreadSheetName
	readRange := sheetName + "!F:K"

	log.Println("Reading from sheet...", readRange, spreadsheetID, sheetName)

	// Make the API request to get the data
	resp, err := service.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve data from sheet: %v", err)
	}

	// Process the response and print the data
	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
		return nil, nil
	}

	// Get the last record of columns A and B
	lastRecord := resp.Values[len(resp.Values)-1]
	article := &domain.GuestArticle{
		TxtDetail:    lastRecord[0].(string),
		Username:     lastRecord[1].(string),
		WrittenAt:    lastRecord[2].(string),
		ViewCount:    lastRecord[3].(int),
		CommentCount: lastRecord[4].(int),
		URL:          lastRecord[5].(string),
	}

	return article, nil
}

func UpdateGuestArticleCheckpoint(service *sheets.Service, channel *domain.CrawlingSource, articles []*domain.GuestArticle) {
	spreadsheetID := channel.SpreadSheetID
	sheetName := channel.SpreadSheetName
	writeRange := sheetName + "!A:F"

	baseURL := "https://cafe.daum.net/dongarry"

	newData := [][]interface{}{}
	for _, article := range articles {
		writtenAt := pkg.ParseRelativeTime(article.WrittenAt)
		row := []interface{}{
			article.TxtDetail,
			article.Username,
			writtenAt,
			article.ViewCount,
			article.CommentCount,
			baseURL + article.URL,
		}
		newData = append(newData, row)
	}

	// Prepare the data for writing
	valueRange := &sheets.ValueRange{
		Values: newData,
	}

	// Make the API request to update the data
	_, err := service.Spreadsheets.Values.Append(spreadsheetID, writeRange, valueRange).ValueInputOption("RAW").Do()
	if err != nil {
		log.Fatalf("Unable to update data in sheet: %v", err)
	}
}

// For NAVER BLOG CRAWLER
func SaveToSheetAppendNaverBlog(service *sheets.Service, channel *domain.CrawlingSource, posts []*domain.NaverBlogArticle) error {
	spreadsheetID := channel.SpreadSheetID
	sheetName := channel.SpreadSheetName
	writeRange := sheetName + "!A:F"

	newData := [][]interface{}{}
	for _, article := range posts {
		row := []interface{}{
			article.ArticleID,
			article.ArticleLink,
			false,
			article.Title,
			strings.Join(article.NaverPlaces, ","),
			article.Content,
			article.PostDate,
		}
		newData = append(newData, row)
	}

	// Prepare the data for writing
	valueRange := &sheets.ValueRange{
		Values: newData,
	}

	// Make the API request to update the data
	_, err := service.Spreadsheets.Values.Append(spreadsheetID, writeRange, valueRange).ValueInputOption("RAW").Do()
	if err != nil {
		log.Fatalf("Unable to update data in sheet: %v", err)
	}

	return nil
}
