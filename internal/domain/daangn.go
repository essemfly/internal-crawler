package domain

import (
	"time"
)

type DanggnStatus string

const (
	DANGGN_STATUS_ALL     DanggnStatus = "ALL"
	DANGGN_STATUS_SALE    DanggnStatus = "SALE"
	DANGGN_STATUS_SOLDOUT DanggnStatus = "SOLDOUT"
	DANGGN_STATUS_UNKNOWN DanggnStatus = "UNKNOWN"
)

type DaangnProduct struct {
	DanggnIndex       string       `json:"danggn_index"`
	Keyword           string       `json:"keyword"`
	KeywordGroup      string       `json:"keyword_group"`
	Name              string       `json:"name"`
	Description       string       `json:"description"`
	Price             int          `json:"price"`
	Images            []string     `json:"images"`
	Status            DanggnStatus `json:"status"`
	Url               string       `json:"url"`
	ViewCounts        int          `json:"view_counts"`
	LikeCounts        int          `json:"like_counts"`
	ChatCounts        int          `json:"chat_counts"`
	CrawlCategory     string       `json:"crawl_category"`
	SellerNickName    string       `json:"seller_nickname"`
	SellerRegionName  string       `json:"seller_region_name"`
	SellerTemperature string       `json:"seller_temperature"`
	WrittenAt         time.Time    `json:"written_at"`
	CreatedAt         time.Time    `json:"created_at"`
	UpdatedAt         time.Time    `json:"updated_at"`
}

type DaangnKeyword struct {
	Keyword   string    `json:"keyword"`
	IsLive    bool      `json:"is_live"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
