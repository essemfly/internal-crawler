package internal

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/essemfly/internal-crawler/config"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

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
