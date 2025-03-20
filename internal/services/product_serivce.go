package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"testwire/config"
	"testwire/internal/models"
	"testwire/internal/repository"
)

type ProductService interface {
	SaveProduct(models.Product) error
	FindProduct(string) (*models.Product, error)
	DeleteProduct(string) error
	GetAll() ([]models.Product, error)
	FindByName(string) (*[]models.Product, error)
	GetSearchHistory(string) ([]models.Product, error)
	GetSearchHistoryByUserId(int) ([]models.Product, error)
}

type ProductServiceImpl struct {
	ProductRepo repository.ProductRepository
}

func NewProductServiceImpl(productRepo repository.ProductRepository) ProductService {
	return &ProductServiceImpl{ProductRepo: productRepo}
}

func (p *ProductServiceImpl) GetSearchHistoryByUserId(userId int) ([]models.Product, error) {
	return p.ProductRepo.GetSearchHistoryByUserId(userId)
}

// GetSearchHistory lấy danh sách lịch sử tìm kiếm có "search successful" trong message
func (p *ProductServiceImpl) GetSearchHistory(userId string) ([]models.Product, error) {
	// Truy vấn Elasticsearch để chỉ lấy những bản ghi có `message` chứa "search successful"
	query := fmt.Sprintf(`{
		"_source": ["fields.level", "level", "message", "msg", "time", "timestamp", "userId"],
		"query": {
			"bool": {
				"must": [
					{ "match_phrase": { "message": "search successful" } },
					{ "term": { "userId": "%s" } }
				]
			}
		}
	}`, userId)
	ES, err := config.NewElasticClient()

	// Gửi request tìm kiếm
	res, err := ES.Search(
		ES.Search.WithContext(context.Background()),
		ES.Search.WithIndex("products"), // Thay đổi index nếu cần
		ES.Search.WithBody(strings.NewReader(query)),
		ES.Search.WithPretty(),
	)
	if err != nil {
		log.Printf("Lỗi tìm kiếm elasticsearch: %v", err)
		return nil, err
	}
	defer res.Body.Close()

	// Kiểm tra phản hồi
	if res.IsError() {
		return nil, fmt.Errorf("elasticsearch trả về lỗi: %s", res.String())
	}

	// Phân tích JSON phản hồi
	var searchResult struct {
		Hits struct {
			Hits []struct {
				Source models.Product `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&searchResult); err != nil {
		log.Printf("Lỗi khi giải mã JSON từ Elasticsearch: %v", err)
		return nil, err
	}

	// Lưu kết quả vào danh sách
	var results []models.Product
	for _, hit := range searchResult.Hits.Hits {
		results = append(results, hit.Source)
	}

	return results, nil
}

func (p *ProductServiceImpl) FindByName(name string) (*[]models.Product, error) {
	return p.ProductRepo.FindByName(name)
}

func (p *ProductServiceImpl) SaveProduct(product models.Product) error {
	return p.ProductRepo.Save(product)
}
func (p *ProductServiceImpl) FindProduct(name string) (*models.Product, error) {
	return p.ProductRepo.Find(name)
}
func (p *ProductServiceImpl) DeleteProduct(name string) error {
	return p.ProductRepo.Delete(name)
}
func (p *ProductServiceImpl) GetAll() ([]models.Product, error) {
	return p.ProductRepo.GetAll()
}
