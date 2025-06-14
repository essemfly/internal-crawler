package mylists

import (
	"github.com/essemfly/internal-crawler/config"
	"github.com/essemfly/internal-crawler/internal/domain/mylists"
	"gorm.io/gorm"
)

type ArticleService struct {
	db *gorm.DB
}

func NewArticleService() *ArticleService {
	db := config.DB()
	db.AutoMigrate(&mylists.Article{})
	return &ArticleService{
		db,
	}
}

func (s *ArticleService) CreateArticle(article *mylists.Article) error {
	return s.db.Create(article).Error
}

func (s *ArticleService) GetArticleByID(id uint) (*mylists.Article, error) {
	var article mylists.Article
	return &article, s.db.First(&article, id).Error
}

func (s *ArticleService) GetAllArticles() ([]*mylists.Article, error) {
	var articles []*mylists.Article
	return articles, s.db.Find(&articles).Error
}

func (s *ArticleService) UpdateArticle(article *mylists.Article) error {
	return s.db.Save(article).Error
}
