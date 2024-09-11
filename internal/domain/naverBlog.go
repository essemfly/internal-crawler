package domain

import "time"

type NaverBlogArticle struct {
	ID          uint `gorm:"primaryKey"`
	Channel     string
	ArticleID   string
	Title       string
	ArticleLink string
	Content     string
	PostDate    string
	NaverPlaces string
	CreatedAt   time.Time
}

func (NaverBlogArticle) TableName() string {
	return "naver_blogs"
}
