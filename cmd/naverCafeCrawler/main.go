package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/essemfly/internal-crawler/internal/crawling"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		fmt.Println("Error loading .env file:", err)
		return
	}

	cafeId := os.Getenv("NAVER_CAFE_ID") // 네이버 카페 ID 입력
	cookie := os.Getenv("NAVER_COOKIE")  // 환경 변수에서 쿠키 가져오기
	if cookie == "" {
		log.Fatal("NAVER_COOKIE 환경 변수가 설정되지 않았습니다.")
	}
	boardID := os.Getenv("NAVER_BOARD_ID") // 크롤링할 게시판 ID

	// 최대 페이지 수 설정 (0은 무제한)
	maxPages := 0
	// pageSize 설정 (기본값: 10)
	pageSize := 10

	if pageSizeStr := os.Getenv("NAVER_PAGE_SIZE"); pageSizeStr != "" {
		if size, err := strconv.Atoi(pageSizeStr); err == nil && size > 0 {
			pageSize = size
		} else {
			log.Printf("⚠️ 잘못된 NAVER_PAGE_SIZE 값입니다. 기본값(50)을 사용합니다.")
		}
	}

	fmt.Println("🚀 네이버 카페 크롤링 시작...")
	posts, err := crawling.CrawlBoard(cafeId, boardID, cookie, maxPages, pageSize)
	if err != nil {
		log.Fatal("❌ 크롤링 중 오류 발생:", err)
	}

	fmt.Printf("✅ 크롤링 완료! 총 %d개 게시글 수집\n", len(posts))

	// 콘솔에도 결과 출력
	for _, post := range posts {
		fmt.Printf("\n📌 [%d] %s\n", post["id"], post["title"])
		fmt.Printf("👤 작성자: %s (레벨: %s)\n", post["writer"], post["writer_level"])
		fmt.Printf("📅 작성일: %s\n", post["write_date"])
		fmt.Printf("📊 조회수: %d, 댓글: %d, 좋아요: %d\n", post["read_count"], post["comment_count"], post["like_count"])

		// 게시글 내용 출력
		if content, ok := post["content"].(string); ok {
			fmt.Printf("\n📝 내용:\n%s\n", content)
		}

		// 댓글 출력
		if comments, ok := post["comments"].([]map[string]interface{}); ok && len(comments) > 0 {
			fmt.Printf("\n💬 댓글 (%d개):\n", len(comments))
			for _, comment := range comments {
				fmt.Printf("  - [%s] %s (%s)\n",
					comment["writer"],
					comment["content"],
					comment["write_date"])
			}
		}
		fmt.Println("\n" + strings.Repeat("─", 80)) // 구분선
	}
}
