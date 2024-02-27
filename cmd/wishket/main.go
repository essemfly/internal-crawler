package main

import (
	"time"

	"github.com/essemfly/internal-crawler/internal"
)

func main() {

	var currentProjectURL string
	for {
		projects := internal.CrawlWishket()
		for _, project := range projects {
			if project.URL == currentProjectURL {
				break
			}
			internal.SendToSlack(project)
		}
		currentProjectURL = projects[0].URL
		time.Sleep(10 * time.Minute)
	}
}
