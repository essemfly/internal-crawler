package main

import (
	"fmt"
	"log"
	"os"

	"github.com/essemfly/internal-crawler/internal/crawling"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		fmt.Println("Error loading .env file:", err)
		return
	}

	cafeId := "27712248"                         // 네이버 카페 ID 입력
	cookie := "4336DE60662074685FE87A6F642E58D3" // 네이버 로그인 후 쿠키 가져와 입력
	boardID := "5"                               // 크롤링할 게시판 ID                        // 몇 페이지까지 크롤링할지 지정

	fmt.Println("🚀 네이버 카페 크롤링 시작...")
	posts, err := crawling.CrawlBoard(cafeId, boardID, cookie)
	if err != nil {
		log.Fatal("❌ 크롤링 중 오류 발생:", err)
	}

	fmt.Printf("✅ 크롤링 완료! 총 %d개 게시글 수집\n", len(posts))
	for _, post := range posts {
		fmt.Printf("📌 [%s] %s\n", post["id"], post["title"])
	}
}
