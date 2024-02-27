package internal

import (
	"io"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

type ProjectInfo struct {
	Title                 string
	URL                   string
	StatusMarks           []string
	EstimatedAmount       string
	EstimatedDuration     string
	WorkStartDate         string
	NumberOfApplicants    string
	ProjectCategoryOrRole string
	Location              string
	Skills                []string
}

func CrawlWishket() []*ProjectInfo {
	WISHKET_URL := "https://www.wishket.com/project/?d=M4JwLgvAdgpg7gMhgYwCYQCogK4yA%3D%3D%3D"
	resp, err := http.Get(WISHKET_URL)
	if err != nil {
		log.Fatal("Error making request: ", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response: ", err)
	}

	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		log.Fatal("Error parsing HTML: ", err)
	}

	var projects []*ProjectInfo
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "div" {
			for _, a := range n.Attr {
				if a.Key == "class" && a.Val == "project-info-box" {
					project := ProjectInfo{}
					extractProjectInfo(n, &project)
					projects = append(projects, &project)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)

	log.Println("projects", projects)

	return projects
}

func extractProjectInfo(n *html.Node, project *ProjectInfo) {
	if n.Type == html.ElementNode {
		for _, a := range n.Attr {
			if a.Key == "class" {
				switch a.Val {
				// 기존 케이스 유지
				case "project-classification-info mb32":
					for c := n.FirstChild; c != nil; c = c.NextSibling {
						if c.Type == html.ElementNode && c.Data == "p" {
							content := strings.TrimSpace(c.FirstChild.Data)
							if strings.Contains(c.Attr[0].Val, "project-category-or-role") {
								project.ProjectCategoryOrRole = content
							}
						}
					}
				case "project-minor-info":
					extractMinorInfo(n, project)
				}
			}
		}
	}

	if n.Type == html.ElementNode {
		switch n.Data {
		case "div":
			for _, a := range n.Attr {
				if a.Key == "class" {
					switch a.Val {
					case "project-status-label recruiting-status mb12":
						for c := n.FirstChild; c != nil; c = c.NextSibling {
							if c.Type == html.ElementNode && c.Data == "div" {
								for _, cAttr := range c.Attr {
									if cAttr.Key == "class" && strings.Contains(cAttr.Val, "status-mark") {
										project.StatusMarks = append(project.StatusMarks, c.FirstChild.Data)
									}
								}
							}
						}
					case "proposal-info":
						extractProposalInfo(n, project)
					case "project-core-info mb10":
						extractCoreInfo(n, project)
					case "project-skills-info":
						for c := n.FirstChild; c != nil; c = c.NextSibling {
							if c.Type == html.ElementNode && c.Data == "span" {
								project.Skills = append(project.Skills, c.FirstChild.Data)
							}
						}
					}
				}
			}
		case "a":
			for _, a := range n.Attr {
				if a.Key == "href" && strings.Contains(a.Val, "/project/") {
					project.URL = "https://www.wishket.com" + a.Val
				}
			}
			if n.FirstChild != nil && n.FirstChild.Type == html.ElementNode && n.FirstChild.Data == "p" {
				project.Title = n.FirstChild.FirstChild.Data
			}
		case "p":

			extractWorkStartDate(n, project)
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractProjectInfo(c, project)
	}
}

func extractCoreInfo(n *html.Node, project *ProjectInfo) {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "p" {
			for _, cAttr := range c.Attr {
				if cAttr.Key == "class" {
					switch cAttr.Val {
					case "budget body-1 text700":
						project.EstimatedAmount = strings.TrimSpace(c.FirstChild.NextSibling.FirstChild.Data)
					case "term body-1 text700":
						project.EstimatedDuration = strings.TrimSpace(c.FirstChild.NextSibling.FirstChild.Data)
					}
				}
			}
		}
	}
}

func extractProposalInfo(n *html.Node, project *ProjectInfo) {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "p" {
			for _, cAttr := range c.Attr {
				if cAttr.Key == "class" && strings.Contains(cAttr.Val, "info-detail") {
					var foundApplicantsText bool
					for cc := c.FirstChild; cc != nil; cc = cc.NextSibling {
						if cc.Type == html.TextNode && strings.Contains(cc.Data, "지원자") {
							foundApplicantsText = true
						}
						if foundApplicantsText && cc.Type == html.ElementNode && cc.Data == "span" {
							if cc.FirstChild != nil {
								project.NumberOfApplicants = strings.TrimSpace(cc.FirstChild.Data)
								break
							}
						}
					}
					if !foundApplicantsText {
						project.NumberOfApplicants = "정보 없음"
					}
				}
			}
		}
	}
}

func extractWorkStartDate(n *html.Node, project *ProjectInfo) {
	if n.Type == html.ElementNode && n.Data == "p" {
		for _, a := range n.Attr {
			if a.Key == "class" && strings.Contains(a.Val, "start-recruitment-data") {
				project.WorkStartDate = n.FirstChild.Data
				return
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractWorkStartDate(c, project)
	}
}

func extractMinorInfo(n *html.Node, project *ProjectInfo) {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "p" {
			for _, cAttr := range c.Attr {
				if cAttr.Key == "class" && strings.Contains(cAttr.Val, "location-data") {
					if c.FirstChild.NextSibling != nil {
						project.Location = strings.TrimSpace(c.FirstChild.NextSibling.Data)
					}
				} else if cAttr.Key == "class" && strings.Contains(cAttr.Val, "start-recruitment-data") {
					textContent := strings.TrimSpace(c.FirstChild.Data)
					splitContent := strings.Split(textContent, " ")
					if len(splitContent) > 1 {
						project.WorkStartDate = splitContent[len(splitContent)-1]
					}
				}
			}
		}
	}
}
