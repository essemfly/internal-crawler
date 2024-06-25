package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/texttheater/golang-levenshtein/levenshtein"
)

type NaverSearchResponse struct {
	Items []struct {
		Title       string `json:"title"`
		Link        string `json:"link"`
		RoadAddress string `json:"roadAddress"`
		Address     string `json:"address"`
		MapX        string `json:"mapx"`
		MapY        string `json:"mapy"`
	} `json:"items"`
}

type PlaceInfo struct {
	Name    string
	Address string
	URL     string
}

func getPlaceUrl(placeName, address string) (*PlaceInfo, error) {
	apiURL := "https://openapi.naver.com/v1/search/local.json"
	query := url.QueryEscape(placeName)
	requestURL := fmt.Sprintf("%s?query=%s&display=5&start=1&sort=random", apiURL, query)

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-Naver-Client-Id", os.Getenv("NAVER_CLIENT_ID"))
	req.Header.Add("X-Naver-Client-Secret", os.Getenv("NAVER_CLIENT_SECRET"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var searchResponse NaverSearchResponse
	err = json.Unmarshal(body, &searchResponse)
	if err != nil {
		return nil, err
	}

	var bestMatch *PlaceInfo
	var lowestDistance int

	for _, item := range searchResponse.Items {
		// HTML 태그 제거
		itemTitle := strings.ReplaceAll(item.Title, "<b>", "")
		itemTitle = strings.ReplaceAll(itemTitle, "</b>", "")

		dist := levenshtein.DistanceForStrings([]rune(address), []rune(item.RoadAddress), levenshtein.DefaultOptions)

		placeURL := fmt.Sprintf("https://map.naver.com/v5/search/%s,%s,%s", url.QueryEscape(itemTitle), item.MapX, item.MapY)

		log.Println("item", item.Title, item.Link)
		if bestMatch == nil || dist < lowestDistance {
			bestMatch = &PlaceInfo{
				Name:    itemTitle,
				Address: item.RoadAddress,
				URL:     placeURL,
			}
			lowestDistance = dist
		}
	}

	log.Println("best match", bestMatch)

	return bestMatch, nil
}

func main() {
	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		fmt.Println("Error loading .env file:", err)
		return
	}

	places := []struct {
		Name    string
		Address string
	}{
		{Name: "풍성감자탕", Address: "서울 광진구 자양로18길 5 (구의동 252-43)"},
		{Name: "효제루", Address: "서울 종로구 대학로 18 1층 (효제동 301-2)"},
	}

	for _, place := range places {
		placeInfo, err := getPlaceUrl(place.Name, place.Address)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
		log.Println("placeInfo", placeInfo)
		if placeInfo != nil {
			fmt.Printf("[%s] %s - %s\n", placeInfo.Name, placeInfo.Address, placeInfo.URL)
		} else {
			fmt.Printf("No matching place found for %s\n", place.Name)
		}
	}
}
