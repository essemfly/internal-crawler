package updating

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/essemfly/internal-crawler/internal/domain"
	"google.golang.org/api/sheets/v4"
)

func ReadLastIndex(service *sheets.Service, channel *domain.CrawlingSource) (int, error) {
	log.Println("SPreadSheetID: ", channel.SpreadSheetID)
	readRange := fmt.Sprintf("%s!B1", channel.SpreadSheetName)
	resp, err := service.Spreadsheets.Values.Get(channel.SpreadSheetID, readRange).Do()
	if err != nil {
		return 0, fmt.Errorf("unable to retrieve data from sheet: %v", err)
	}

	if len(resp.Values) == 0 || len(resp.Values[0]) == 0 {
		return 0, fmt.Errorf("no data found in B1")
	}

	lastIndex, err := strconv.Atoi(resp.Values[0][0].(string))
	if err != nil {
		return 0, fmt.Errorf("unable to convert last index to integer: %v", err)
	}

	return lastIndex, nil
}

func UpdateLastIndex(service *sheets.Service, channel *domain.CrawlingSource, lastIndex int) error {
	updateRange := fmt.Sprintf("%s!B1", channel.SpreadSheetName)
	valueRange := &sheets.ValueRange{
		Values: [][]interface{}{{lastIndex}},
	}

	_, err := service.Spreadsheets.Values.Update(channel.SpreadSheetID, updateRange, valueRange).ValueInputOption("RAW").Do()
	if err != nil {
		return fmt.Errorf("unable to update data in sheet: %v", err)
	}

	return nil
}

func ReadKeywords(service *sheets.Service, channel *domain.CrawlingSource) ([]string, error) {
	var strArr []string

	readRange := fmt.Sprintf("%s!D1", channel.SpreadSheetName)
	resp, err := service.Spreadsheets.Values.Get(channel.SpreadSheetID, readRange).Do()
	if err != nil {
		return strArr, fmt.Errorf("unable to retrieve data from sheet: %v", err)
	}

	if len(resp.Values) == 0 || len(resp.Values[0]) == 0 {
		return strArr, fmt.Errorf("no data found in D1")
	}

	keywords := resp.Values[0][0].(string)

	return strings.Split(keywords, ","), nil
}
