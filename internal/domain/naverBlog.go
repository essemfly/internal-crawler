package domain

type NaverBlogArticle struct {
	ID          uint `gorm:"primaryKey"`
	Channel     string
	ArticleID   string
	Title       string
	ArticleLink string
	Content     string
	PostDate    string
	NaverPlaces string
}

func (NaverBlogArticle) TableName() string {
	return "naver_blogs"
}
