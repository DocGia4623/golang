package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testwire/config"
	"testwire/internal/models"

	"github.com/elastic/go-elasticsearch/v8"
	"gorm.io/gorm"
)

type ProductRepository interface {
	GetAll() ([]models.Product, error)
	Save(models.Product) error
	Find(string) (*models.Product, error)
	Delete(string) error
	FindByName(string) (*[]models.Product, error)
	MigrateToElastic() error
	GetSearchHistoryByUserId(userId int) ([]models.Product, error)
}

type ProductRepositoryImpl struct {
	Db       *gorm.DB
	ESClient *elasticsearch.Client
}

func NewProductRepositoryImpl(db *gorm.DB, esClient *elasticsearch.Client) ProductRepository {
	return &ProductRepositoryImpl{Db: db, ESClient: esClient}
}

// Struct ánh xạ `_source` từ Elasticsearch
type SourceData struct {
	Data []models.Product `json:"Data"` // Dữ liệu sản phẩm nằm trong "Data"
}

func (p *ProductRepositoryImpl) GetSearchHistoryByUserId(userId int) ([]models.Product, error) {
	ctx := context.Background()
	// Tạo truy vấn tìm kiếm
	query := fmt.Sprintf(`{
		"query": {
			"match": {
				"userId": %d
			}
		}
	}`, userId)
	// Thực hiện truy vấn đến Elasticsearch
	res, err := p.ESClient.Search(
		p.ESClient.Search.WithContext(ctx),
		p.ESClient.Search.WithIndex("logstash-docker-2025.03.20"), // Chỉ tìm trong index "products"
		p.ESClient.Search.WithBody(strings.NewReader(query)),
		p.ESClient.Search.WithTrackTotalHits(true),
		p.ESClient.Search.WithPretty(),
	)
	if err != nil {
		return nil, fmt.Errorf("error executing search query: %v", err)
	}
	defer res.Body.Close()
	// Đọc kết quả trả về
	if res.IsError() {
		return nil, fmt.Errorf("error response from Elasticsearch: %s", res.String())
	}
	// Giải mã kết quả JSON
	var result struct {
		Hits struct {
			Hits []struct {
				Source SourceData `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding Elasticsearch response: %v", err)
	}
	// Chuyển kết quả thành danh sách sản phẩm
	var products []models.Product
	for _, hit := range result.Hits.Hits {
		products = append(products, hit.Source.Data...)
	}
	// Nếu không có sản phẩm nào được tìm thấy, trả về nil
	if len(products) == 0 {
		return nil, fmt.Errorf("no products found for userId %d", userId)
	}
	return products, nil
	// Trả về danh sách sản phẩm tìm thấy
}
func (p *ProductRepositoryImpl) MigrateToElastic() error {
	// Migrate dữ liệu từ PostgreSQL sang Elasticsearch
	// Lấy tất cả dữ liệu từ PostgreSQL
	products, err := p.GetAll()
	if err != nil {
		return err
	}
	// Tạo kết nối tới Elasticsearch
	es, err := config.NewElasticClient()
	if err != nil {
		return err
	}
	// Lặp qua từng sản phẩm và lưu vào Elasticsearch
	for _, product := range products {
		// Chuyển product sang JSON
		productJSON, err := json.Marshal(product)
		if err != nil {
			return err
		}
		// Lưu vào Elasticsearch
		_, err = es.Index("products", bytes.NewReader(productJSON))
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *ProductRepositoryImpl) FindByName(name string) (*[]models.Product, error) {
	ctx := context.Background()

	// Tạo truy vấn tìm kiếm
	query := fmt.Sprintf(`{
		"query": {
			"match": {
				"name": "%s"
			}
		}
	}`, name)

	// Thực hiện truy vấn đến Elasticsearch
	res, err := p.ESClient.Search(
		p.ESClient.Search.WithContext(ctx),
		p.ESClient.Search.WithIndex("products"), // Chỉ tìm trong index "products"
		p.ESClient.Search.WithBody(strings.NewReader(query)),
		p.ESClient.Search.WithTrackTotalHits(true),
		p.ESClient.Search.WithPretty(),
	)
	if err != nil {
		return nil, fmt.Errorf("error executing search query: %v", err)
	}
	defer res.Body.Close()

	// Đọc kết quả trả về
	if res.IsError() {
		return nil, fmt.Errorf("error response from Elasticsearch: %s", res.String())
	}

	// Giải mã kết quả JSON
	var result struct {
		Hits struct {
			Hits []struct {
				Source models.Product `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding Elasticsearch response: %v", err)
	}

	// Chuyển kết quả thành danh sách sản phẩm
	var products []models.Product
	for _, hit := range result.Hits.Hits {
		products = append(products, hit.Source)
	}

	// Nếu không có sản phẩm nào được tìm thấy, trả về nil
	if len(products) == 0 {
		return nil, nil
	}

	return &products, nil
}

func (p *ProductRepositoryImpl) Save(product models.Product) error {
	result := p.Db.Save(&product)
	return result.Error
}
func (p *ProductRepositoryImpl) Find(name string) (*models.Product, error) {
	var product models.Product
	result := p.Db.Where("name = ?", name).First(&product)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil // Trả về nil thay vì lỗi nếu không tìm thấy
	}
	return &product, result.Error
}
func (p *ProductRepositoryImpl) Delete(name string) error {
	result := p.Db.Where("name = ?", name).Delete(&models.Product{})
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound // Nếu không xóa được thì trả về lỗi
	}
	return result.Error
}

func (p *ProductRepositoryImpl) GetAll() ([]models.Product, error) {
	var products []models.Product
	result := p.Db.Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}
	return products, nil
}
