package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	// Daum 로그인 정보
	daumID := "lsm89712@hanmail.net"
	daumPassword := ""

	// Daum 로그인 URL
	loginURL := "https://accounts.kakao.com/login/?continue=https://cafe.daum.net/_c21_/bbs_list?grpid=IGaj&fldid=Dilr"

	// 로그인 폼 데이터
	loginData := url.Values{
		"id":  {daumID},
		"pw":  {daumPassword},
		"url": {"https://www.daum.net/"},
	}

	// HTTP 클라이언트 및 쿠키 저장소 생성
	jar, err := cookiejar.New(nil)
	if err != nil {
		fmt.Println("Error creating cookie jar:", err)
		return
	}
	client := &http.Client{Jar: jar}

	// 로그인 요청 보내기
	resp, err := client.PostForm(loginURL, loginData)
	if err != nil {
		fmt.Println("Error logging in:", err)
		return
	}
	defer resp.Body.Close()

	// 응답 본문 확인 (디버깅용)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading login response:", err)
		return
	}
	fmt.Println("Login response:", string(body))

	url := "https://cafe.daum.net/_c21_/bbs_list?grpid=IGaj&fldid=Dilr"

	// HTTP GET 요청 보내기
	resp, err = http.Get(url)
	if err != nil {
		fmt.Println("Error fetching the URL:", err)
		return
	}
	defer resp.Body.Close()

	// HTML 파싱
	doc, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Println("Error parsing the HTML:", err)
		return
	}

	// 게시글 제목과 링크 추출
	extractPosts(doc)
}

func extractPosts(n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "tr" {
		hasStateInfo := false
		for _, attr := range n.Attr {
			if attr.Key == "class" && strings.Contains(attr.Val, "state_info") {
				hasStateInfo = true
				break
			}
		}

		if !hasStateInfo {
			extractATagText(n)
		}
	}

	// 자식 노드 순회
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractPosts(c)
	}
}

func extractATagText(n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "a" {
		title := extractText(n)
		link := "https://cafe.daum.net" + getHref(n)
		fmt.Println("Title:", title)
		fmt.Println("Link:", link)
		fmt.Println()
	}

	// 자식 노드 순회
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractATagText(c)
	}
}

func extractText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	var text string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		text += extractText(c)
	}
	return strings.TrimSpace(text)
}

func getHref(n *html.Node) string {
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			return attr.Val
		}
	}
	return ""
}
