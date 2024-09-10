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
	ID                uint         `gorm:"primaryKey"` // Primary key
	DanggnIndex       string       `json:"danggn_index" gorm:"type:varchar(255)"`
	Keyword           string       `json:"keyword" gorm:"type:varchar(255)"`
	KeywordGroup      string       `json:"keyword_group" gorm:"type:varchar(255)"`
	Name              string       `json:"name" gorm:"type:varchar(255)"`
	Description       string       `json:"description" gorm:"type:text"`
	Price             int          `json:"price"`
	Images            []string     `json:"images" gorm:"-"` // Excluded from the DB
	Status            DanggnStatus `json:"status" gorm:"type:varchar(20)"`
	Url               string       `json:"url" gorm:"type:varchar(255)"`
	ViewCounts        int          `json:"view_counts"`
	LikeCounts        int          `json:"like_counts"`
	ChatCounts        int          `json:"chat_counts"`
	CrawlCategory     string       `json:"crawl_category" gorm:"type:varchar(255)"`
	SellerNickName    string       `json:"seller_nickname" gorm:"type:varchar(255)"`
	SellerRegionName  string       `json:"seller_region_name" gorm:"type:varchar(255)"`
	SellerTemperature string       `json:"seller_temperature" gorm:"type:varchar(255)"`
	WrittenAt         time.Time    `json:"written_at"`
	CreatedAt         time.Time    `json:"created_at"`
	UpdatedAt         time.Time    `json:"updated_at"`
}

// DaangnKeyword represents a keyword used for Daangn search
type DaangnKeyword struct {
	ID        uint      `gorm:"primaryKey"` // Primary key
	Keyword   string    `json:"keyword" gorm:"type:varchar(255)"`
	IsLive    bool      `json:"is_live"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DaangnConfig struct {
	ID         uint `gorm:"primaryKey"` // Primary key
	CurrentIdx int  `json:"current_idx"`
}

func (DaangnProduct) TableName() string {
	return "daangn_products"
}

func (DaangnKeyword) TableName() string {
	return "daangn_keywords"
}

func (DaangnConfig) TableName() string {
	return "daangn_configs"
}
