package repository

import (
	"time"

	"github.com/essemfly/internal-crawler/config"
	"github.com/essemfly/internal-crawler/internal/domain"
	"gorm.io/gorm"
)

type NaverBlogService struct {
	db *gorm.DB
}

func NewNaverBlogService() *NaverBlogService {
	db := config.DB()
	db.AutoMigrate(&domain.NaverBlogArticle{})
	return &NaverBlogService{
		db,
	}
}

func (ns *NaverBlogService) CreateNaverBlogArticle(article *domain.NaverBlogArticle) error {
	result := ns.db.Create(article)
	return result.Error
}

func (ns *NaverBlogService) CreateNaverBlogArticles(articles []*domain.NaverBlogArticle) error {
	result := ns.db.Create(articles)
	return result.Error
}

func (ns *NaverBlogService) GetLatestArticleDate(blogId string) (time.Time, error) {
	var article domain.NaverBlogArticle
	result := ns.db.Where("channel = ?", blogId).Order("created_at asc").First(&article)
	return article.CreatedAt, result.Error
}

func (ns *NaverBlogService) ListUnprocessedArticles() ([]*domain.NaverBlogArticle, error) {
	var articles []*domain.NaverBlogArticle
	result := ns.db.Where("is_processed = ? AND naver_places != ?", false, "").Find(&articles)
	return articles, result.Error
}

func (ns *NaverBlogService) UpdateArticleToProcessed(articleId uint) error {
	result := ns.db.Model(&domain.NaverBlogArticle{}).Where("id = ?", articleId).Update("is_processed", true)
	return result.Error
}

func (ns *NaverBlogService) ListByNaverPlaces(place string) ([]*domain.NaverBlogArticle, error) {
	var articles []*domain.NaverBlogArticle
	result := ns.db.Where("naver_places = ?", place).Order("created_at asc").Find(&articles) // 쿼리 실행 및 결과 저장
	return articles, result.Error                                                            // 결과와 에러 반환
}
