package config

import (
	"log"

	"github.com/elastic/go-elasticsearch/v8"
)

// NewElasticClient tạo một Elasticsearch client
func NewElasticClient() (*elasticsearch.Client, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200", // Thay đổi nếu cần
		},
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Lỗi kết nối Elasticsearch: %v", err)
		return nil, err
	}

	// Kiểm tra xem Elasticsearch có đang chạy không
	res, err := client.Ping()
	if err != nil {
		log.Fatalf("Lỗi ping Elasticsearch: %v", err)
		return nil, err
	}
	defer res.Body.Close()

	log.Println("✅ Đã kết nối Elasticsearch thành công!")
	return client, nil
}
