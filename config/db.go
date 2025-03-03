package config

import (
	"context"
	"fmt"
	"log"
	"testwire/internal/models"

	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB           *gorm.DB
	RedisClient  *redis.Client
	RabbitMQConn *amqp091.Connection
	AppConfig    Config
)

func ConnectDB(config *Config) *gorm.DB {
	// tạo chuỗi kết nối PostgreSQL
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.PostgresUser,
		config.PostgresPassword,
		config.DBHost,
		config.DBPort,
		config.PostgresDB,
	)

	// Mở kết nối tới postgreSql
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connec to database: %v", err)
	}
	// Migrate
	modelsToMigrate := []interface{}{
		&models.User{},
		&models.Role{},
		&models.Role{},
		&models.RefreshToken{},
	}

	for _, model := range modelsToMigrate {
		db.AutoMigrate(model)
	}
	DB = db
	return db
}

// ConnectRedis thiết lập kết nối Redis
func ConnectRedis(config *Config) *redis.Client {
	redisAddr := fmt.Sprintf("%s:%s", config.RedisHost, config.RedisPort)

	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
		DB:   config.RedisDB,
	})

	// Kiểm tra kết nối Redis
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	RedisClient = client
	return client
}

// ConnectRabbitMQ thiết lập kết nối RabbitMQ
func ConnectRabbitMQ(config *Config) *amqp091.Connection {
	// Chuỗi kết nối RabbitMQ
	RabbitMQConnStr := fmt.Sprintf("amqp://%s:%s@%s:%s/", config.RabbitMqUser, config.RabbitMQPassword, config.RabbitMQHost, config.RabbitMQPort)
	// Kết nối RabbitMQ
	conn, err := amqp091.Dial(RabbitMQConnStr)
	if err != nil {
		log.Fatalf("Failed to connect rabbitMQ: %v", err)
	}
	RabbitMQConn = conn
	return conn
}

func Connect(config *Config) {
	ConnectDB(config)
	ConnectRabbitMQ(config)
	ConnectRedis(config)
}
