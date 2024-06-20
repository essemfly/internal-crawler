package seed

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/essemfly/internal-crawler/config"
	"github.com/essemfly/internal-crawler/internal/domain"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func getRegisteredSources() []*domain.CrawlingSource {
	spreadsheetID := os.Getenv("SPREADSHEET_ID")
	sheetName := "channels"
	readRange := sheetName + "!A:H"

	ctx := context.Background()
	creds := option.WithCredentialsFile(config.JsonKeyFilePath) // Replace with your credentials file path
	srv, err := sheets.NewService(ctx, creds)
	if err != nil {
		log.Fatalf("Unable to create Sheets service: %v", err)
	}

	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	var sources []*domain.CrawlingSource
	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
	} else {
		for _, row := range resp.Values {
			typeInString := row[1].(string)
			typeInStruct := domain.Youtube
			switch typeInString {
			case "wishket":
				typeInStruct = domain.Wishket
			case "daangn":
				typeInStruct = domain.Daangn
			case "daumcafe":
				typeInStruct = domain.DaumCafe
			}

			constraints := row[4].(string)

			source := &domain.CrawlingSource{
				SourceName:      row[0].(string),
				SourceID:        row[3].(string),
				SpreadSheetID:   spreadsheetID,
				SpreadSheetName: row[5].(string),
				Type:            typeInStruct,
				NaverListID:     "",
				NaverListName:   "",
				WebhookURL:      row[2].(string),
				Constraint:      strings.Split(constraints, ","),
			}

			// Check if row[6] and row[7] exist and set them appropriately
			if len(row) > 6 {
				source.NaverListID = row[6].(string)
			}
			if len(row) > 7 {
				source.NaverListName = row[7].(string)
			}
			sources = append(sources, source)
		}
	}
	return sources
}

func ListSources(sourceType domain.CrawlingSourceType) []*domain.CrawlingSource {
	sources := getRegisteredSources()
	var filteredSources []*domain.CrawlingSource
	for _, source := range sources {
		if source.Type == sourceType {
			filteredSources = append(filteredSources, source)
		}
	}
	return filteredSources
}
