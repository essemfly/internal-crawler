package main

import (
	"github.com/essemfly/internal-crawler/internal"
)

func main() {
	projects := internal.CrawlWishket()
	for _, project := range projects {
		internal.SendToSlack(project)
	}
}
