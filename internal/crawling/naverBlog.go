package crawling

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/essemfly/internal-crawler/internal/domain"
	"github.com/tidwall/gjson"
	"golang.org/x/net/html"
)

type BlogPost struct {
	LogNo   string `json:"logNo"`
	Title   string `json:"title"`
	AddDate string `json:"addDate"`
}

type Response struct {
	ResultCode    string     `json:"resultCode"`
	ResultMessage string     `json:"resultMessage"`
	PostList      []BlogPost `json:"postList"`
	CountPerPage  string     `json:"countPerPage"`
	TotalCount    string     `json:"totalCount"`
}

type ModuleData struct {
	Data struct {
		LocationId string `json:"locationId"`
	} `json:"data"`
}

func FetchAllBlogPosts(source *domain.CrawlingSource, workers int) ([]*domain.NaverBlogArticle, error) {
	blogId := source.SourceID
	var totalPosts []*domain.NaverBlogArticle

	for _, categoryNo := range source.Constraint {
		log.Println("Source: "+source.SourceName+" Category: ", categoryNo)
		// 첫 번째 페이지를 가져와서 전체 페이지 수 계산
		firstPagePosts, totalPages, err := FetchBlogPosts(blogId, categoryNo, 1)
		log.Println("Total Pages #: " + strconv.Itoa(totalPages))
		if err != nil {
			log.Println("ERR HERE?", err)
			return nil, err
		}
		totalPosts = append(totalPosts, firstPagePosts...)

		for i := 2; i <= totalPages; i++ {
			newPosts, _, err := FetchBlogPosts(blogId, categoryNo, i)
			if err != nil {
				return nil, err
			}
			totalPosts = append(totalPosts, newPosts...)
		}

	}

	return totalPosts, nil
}

func FetchBlogPosts(blogId, categoryNo string, currentPage int) ([]*domain.NaverBlogArticle, int, error) {
	var posts []*domain.NaverBlogArticle

	baseURL := fmt.Sprintf("https://blog.naver.com/%s/", blogId)

	// URL 생성
	categoryUrl := fmt.Sprintf("https://blog.naver.com/PostTitleListAsync.naver?blogId=%s&currentPage=%d&categoryNo=%s&countPerPage=30",
		blogId, currentPage, categoryNo)

	log.Println("3-12", categoryUrl)
	// HTTP GET 요청
	resp, err := http.Get(categoryUrl)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch page %d: %v", currentPage, err)
	}
	defer resp.Body.Close()

	time.Sleep(1 * time.Second)
	// 응답 바디 읽기
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to read response body: %v", err)
	}

	cleanedBody := strings.ReplaceAll(string(body), "\\'", "'")
	body = []byte(cleanedBody)

	// JSON 파싱
	var respData Response
	err = json.Unmarshal(body, &respData)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	// URL 디코딩 및 결과 수집
	for _, post := range respData.PostList {
		decodedTitle, err := url.QueryUnescape(post.Title)
		if err != nil {
			fmt.Printf("Error decoding title: %v\n", err)
			continue
		}
		post.Title = decodedTitle
		places := ParseNaverMapUrl(blogId, post.LogNo)
		log.Println("URL: "+baseURL+post.LogNo+", PLACES", places)

		content := blogId + "- "
		if blogId == "mardukas" {
			content += getContentKeyValues(blogId, categoryNo)
		} else if blogId == "paperchan" {
			content += getContentKeyValues(blogId, categoryNo)
		} else {
			content += categoryNo
		}

		posts = append(posts, &domain.NaverBlogArticle{
			ArticleID:   post.LogNo,
			ArticleLink: baseURL + post.LogNo,
			Title:       post.Title,
			Content:     blogId + "- " + categoryNo,
			PostDate:    post.AddDate,
			NaverPlaces: ParseNaverMapUrl(blogId, post.LogNo),
		},
		)
	}

	// 첫 번째 페이지에서 전체 페이지 수 계산
	if currentPage == 1 {
		countPerPage, _ := strconv.Atoi(respData.CountPerPage)
		totalCount, _ := strconv.Atoi(respData.TotalCount)
		totalPages := int(math.Ceil(float64(totalCount) / float64(countPerPage)))
		return posts, totalPages, nil
	}

	return posts, 0, nil
}

func ParseNaverMapUrl(blogId, logNo string) []string {
	if blogId != "mardukas" {
		return []string{}
	}
	// URL 생성
	url := fmt.Sprintf("https://blog.naver.com/PostView.naver?blogId=%s&logNo=%s", blogId, logNo)

	log.Println("3-13", url)
	// HTTP GET 요청
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Failed to fetch the page: %v\n", err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Failed to fetch the page, status code: %d\n", resp.StatusCode)
		return nil
	}

	// HTML 파싱
	doc, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Printf("Failed to parse HTML: %v\n", err)
		return nil
	}

	// 스크립트 태그에서 v2_map을 찾고 placeId 추출
	var placeIds []string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "script" {
			for _, attr := range n.Attr {
				if attr.Key == "type" && attr.Val == "text/data" {
					// data-module 속성의 값을 찾음
					for _, a := range n.Attr {
						if a.Key == "data-module" && strings.Contains(a.Val, "\"type\":\"v2_map\"") {
							placeIds = extractPlaceIds(a.Val)
						} else if a.Key == "data-module" {
							locationId := extractLocationIdFromDataModule(a.Val)
							if locationId != "" {
								placeIds = append(placeIds, locationId)
							}
						}
					}
				}
			}
		}

		if n.Type == html.ElementNode && n.Data == "iframe" {
			log.Println("HOIT?!", n)
			var src string
			for _, attr := range n.Attr {
				if attr.Key == "src" {
					log.Println("SRC: ", attr.Val)
					src = attr.Val
					// Extract the number after @s in the src attribute
					locationId := extractLoacionIdFromSrc(src)
					if locationId != "" {
						placeIds = append(placeIds, locationId)
					}
				}
			}
		}

		// // "a" 태그에서 "data-linkdata" 속성을 찾아 locationId 추출
		// if n.Type == html.ElementNode && n.Data == "a" {
		// 	var dataLinkData string
		// 	for _, attr := range n.Attr {
		// 		if attr.Key == "data-linkdata" {
		// 			dataLinkData = attr.Val
		// 		}
		// 	}

		// 	// data-linkdata가 있다면 JSON 형태로 파싱하여 locationId 추출
		// 	if dataLinkData != "" {
		// 		locationId := extractLocationId(dataLinkData)
		// 		if locationId != "" {
		// 			placeIds = append(placeIds, locationId)
		// 		}
		// 	}
		// }

		// 재귀적으로 모든 노드를 탐색
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return placeIds
}

func extractPlaceIds(jsonStr string) []string {
	var placeIds []string

	// gjson을 사용하여 JSON 데이터에서 placeId 값을 추출
	result := gjson.Get(jsonStr, "data.places.#.placeId")
	result.ForEach(func(key, value gjson.Result) bool {
		placeIds = append(placeIds, value.String())
		return true // 계속 반복
	})

	var results []string

	for _, placeId := range placeIds {
		fullUrl := fmt.Sprintf("https://map.naver.com/p/entry/place/%s", placeId)
		results = append(results, fullUrl)
	}

	return results
}

// Deprecated: data-linkdata 속성을 사용하여 locationId 값을 추출하는 방법
/*
func extractLocationId(data string) string {
	// HTML 엔티티인 &quot;을 "로 변환
	data = strings.ReplaceAll(data, "&quot;", "\"")

	// locationId 값을 추출하기 위한 문자열 파싱
	start := strings.Index(data, "\"locationId\":\"")
	if start == -1 {
		return ""
	}
	start += len("\"locationId\":\"")
	end := strings.Index(data[start:], "\"")
	if end == -1 {
		return ""
	}

	// locationId 값 반환
	locationId := data[start : start+end]
	fullUrl := fmt.Sprintf("https://map.naver.com/p/entry/place/%s", locationId)
	return fullUrl
}
*/

func extractLocationIdFromDataModule(dataModule string) string {
	var moduleData ModuleData
	err := json.Unmarshal([]byte(dataModule), &moduleData)
	if err != nil {
		fmt.Printf("Failed to parse data-module JSON: %v\n", err)
		return ""
	}

	locationId := moduleData.Data.LocationId
	fullUrl := fmt.Sprintf("https://map.naver.com/p/entry/place/%s", locationId)
	return fullUrl
}

func extractLoacionIdFromSrc(src string) string {
	// Regex to find the number after @s
	re := regexp.MustCompile(`%40s(\d+)`)
	match := re.FindStringSubmatch(src)
	if len(match) > 1 {
		return fmt.Sprintf("https://map.naver.com/p/entry/place/%s", match[1])
	}
	return ""
}

func getContentKeyValues(blogId, categoryNo string) string {
	if blogId == "mardukas" {
		if categoryNo == "9" {
			return "서울추천맛집(한식)"
		} else if categoryNo == "61" {
			return "서울추천맛집(양식)"
		} else if categoryNo == "62" {
			return "서울추천맛집(기타)"
		} else if categoryNo == "111" {
			return "서울추천맛집(간식,차)"
		} else if categoryNo == "65" {
			return "지방추천맛집(인천경기)"
		} else if categoryNo == "66" {
			return "지방추천맛집(강원)"
		} else if categoryNo == "1" {
			return "지방추천맛집(경상)"
		} else if categoryNo == "68" {
			return "지방추천맛집(전라)"
		} else if categoryNo == "70" {
			return "지방추천맛집(제주)"
		} else if categoryNo == "67" {
			return "지방추천맛집(충청)"
		}
	} else if blogId == "paperchan" {
		if categoryNo == "5" {
			return "한식"
		} else if categoryNo == "2" {
			return "일식"
		} else if categoryNo == "3" {
			return "중식"
		} else if categoryNo == "4" {
			return "양식"
		} else if categoryNo == "15" {
			return "기타"
		}
	}

	return ""
}
