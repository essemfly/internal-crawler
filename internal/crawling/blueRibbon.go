package crawling

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/essemfly/internal-crawler/internal/domain"
)

type RibbonType string

const (
	RIBBON_ONE   RibbonType = "RIBBON_ONE"
	RIBBON_TWO   RibbonType = "RIBBON_TWO"
	RIBBON_THREE RibbonType = "RIBBON_THREE"
)

type RestaurantResponse struct {
	Embedded struct {
		Restaurants []Restaurant `json:"restaurants"`
	} `json:"_embedded"`
	Page struct {
		Size          int `json:"size"`
		TotalElements int `json:"totalElements"`
		TotalPages    int `json:"totalPages"`
		Number        int `json:"number"`
	} `json:"page"`
}

type Restaurant struct {
	ID         int64 `json:"id"`
	HeaderInfo struct {
		NameKR     string `json:"nameKR"`
		NameEN     string `json:"nameEN"`
		RibbonType string `json:"ribbonType"`
	} `json:"headerInfo"`
	DefaultInfo struct {
		Phone  string `json:"phone"`
		DayOff string `json:"dayOff"`
	} `json:"defaultInfo"`
	StatusInfo struct {
		Menu          string `json:"menu"`
		PriceRange    string `json:"priceRange"`
		BusinessHours string `json:"businessHours"`
	} `json:"statusInfo"`
	Juso struct {
		RoadAddrPart1 string `json:"roadAddrPart1"`
		RoadAddrPart2 string `json:"roadAddrPart2"`
	} `json:"juso"`
	GPS struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"gps"`
	Review struct {
		Review string `json:"review"`
	} `json:"review"`
	FoodTypes []string `json:"foodTypes"`
}

func CrawlBlueRibbon(ribbonType RibbonType) ([]*domain.PublicArticle, error) {
	var articles []*domain.PublicArticle
	client := &http.Client{Timeout: 30 * time.Second}

	page := 0
	size := 100

	for {
		url := buildURL(ribbonType, page, size)
		log.Printf("Fetching page %d: %s", page, url)

		resp, err := fetchAPI(client, url)
		if err != nil {
			log.Printf("Error fetching page %d: %v", page, err)
			break
		}

		if len(resp.Embedded.Restaurants) == 0 {
			log.Printf("No more restaurants found on page %d", page)
			break
		}

		for _, restaurant := range resp.Embedded.Restaurants {
			article := convertToArticle(restaurant, ribbonType)
			articles = append(articles, article)
		}

		if page >= resp.Page.TotalPages-1 {
			break
		}

		page++
		time.Sleep(3 * time.Second)
	}

	return articles, nil
}

func buildURL(ribbonType RibbonType, page, size int) string {
	baseURL := "https://www.bluer.co.kr/api/v1/restaurants"
	return fmt.Sprintf("%s?page=%d&size=%d&ribbonType=%s&tabMode=single&searchMode=ribbonType",
		baseURL, page, size, string(ribbonType))
}

func fetchAPI(client *http.Client, url string) (*RestaurantResponse, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "ko-KR,ko;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("Origin", "https://www.bluer.co.kr")
	req.Header.Set("Referer", "https://www.bluer.co.kr/")
	req.Header.Set("sec-ch-ua", `"Not A(Brand";v="99", "Google Chrome";v="121", "Chromium";v="121"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result RestaurantResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

func convertToArticle(restaurant Restaurant, ribbonType RibbonType) *domain.PublicArticle {
	now := time.Now()

	// 주소 조합
	address := restaurant.Juso.RoadAddrPart1
	if restaurant.Juso.RoadAddrPart2 != "" {
		address += " " + restaurant.Juso.RoadAddrPart2
	}

	// 메모 정보 조합
	var notes []string
	notes = append(notes, fmt.Sprintf("RibbonType: %s", ribbonType))
	notes = append(notes, fmt.Sprintf("Phone: %s", restaurant.DefaultInfo.Phone))
	notes = append(notes, fmt.Sprintf("DayOff: %s", restaurant.DefaultInfo.DayOff))
	notes = append(notes, fmt.Sprintf("PriceRange: %s", restaurant.StatusInfo.PriceRange))
	notes = append(notes, fmt.Sprintf("BusinessHours: %s", restaurant.StatusInfo.BusinessHours))
	notes = append(notes, fmt.Sprintf("FoodTypes: %v", restaurant.FoodTypes))
	notes = append(notes, fmt.Sprintf("Latitude: %f", restaurant.GPS.Latitude))
	notes = append(notes, fmt.Sprintf("Longitude: %f", restaurant.GPS.Longitude))

	return &domain.PublicArticle{
		Type:        domain.BlueRibbon,
		Title:       restaurant.HeaderInfo.NameKR,
		Description: restaurant.Review.Review,
		Address:     address,
		Link:        fmt.Sprintf("https://www.bluer.co.kr/restaurants/%d", restaurant.ID),
		Notes:       strings.Join(notes, " | "),
		Latitude:    &restaurant.GPS.Latitude,
		Longitude:   &restaurant.GPS.Longitude,
		PublishedAt: &now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
