package domain

type GuestArticle struct {
	ID           uint   `gorm:"primaryKey"`        // ID is the primary key
	URL          string `gorm:"type:varchar(255)"` // URL of the article
	TxtDetail    string `gorm:"type:text"`         // Text detail for the article
	Username     string `gorm:"type:varchar(100)"` // Username of the guest user
	CreatedAt    string
	ViewCount    int // Number of views
	CommentCount int // Number of comments
}

func (GuestArticle) TableName() string {
	return "daum_cafe_articles"
}
