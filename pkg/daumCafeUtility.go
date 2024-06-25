package pkg

import (
	"log"
	"regexp"
	"strconv"
	"time"
)

func ParseRelativeTime(relativeTime string) string {
	location, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		log.Fatalf("Failed to load location: %v", err)
	}
	now := time.Now().In(location)

	if relativeTime == "방금전" {
		return now.Format("2006-01-02 15:04:05")
	}

	regex := regexp.MustCompile(`(\d+)(시간|분)`)
	matches := regex.FindAllStringSubmatch(relativeTime, -1)

	if len(matches) == 0 {
		log.Fatalf("Invalid time format: %v", relativeTime)
	}

	for _, match := range matches {
		value, err := strconv.Atoi(match[1])
		if err != nil {
			log.Fatalf("Invalid number in time format: %v", match[1])
		}
		unit := match[2]

		switch unit {
		case "시간":
			now = now.Add(-time.Duration(value) * time.Hour)
		case "분":
			now = now.Add(-time.Duration(value) * time.Minute)
		}
	}

	return now.Format("2006-01-02 15:04:05")
}
