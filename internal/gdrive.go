package internal

import (
	"context"
	"log"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func CreateSheetsService(jsonKeyFilePath string) (*sheets.Service, error) {
	b, err := os.ReadFile(jsonKeyFilePath)
	if err != nil {
		log.Fatalf("Unable to read service account file: %v", err)
		return nil, err
	}

	config, err := google.JWTConfigFromJSON(b, sheets.SpreadsheetsScope)
	if err != nil {
		log.Printf("Unable to parse service account file to config: %v", err)
		return nil, err
	}

	ctx := context.Background()
	client := config.Client(ctx)
	sheetsService, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Printf("Unable to retrieve Sheets client: %v", err)
		return nil, err
	}

	return sheetsService, nil
}
