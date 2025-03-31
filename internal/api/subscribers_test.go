package api

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"webhook/internal/model"
	"webhook/internal/repository"
	"webhook/pkg"
)

const (
	MainServiceName = "webhook"

	// SQL
	sqlHost                     = "localhost"
	port                        = "5432"
	user                        = "postgres"
	sqlPassword                 = "postgres"
	dbname                      = "postgres"
	sslmode                     = "disable"
	timezone                    = "UTC"
	maxOpenConns                = 100
	maxIdleConns                = 10
	connMaxLifetimeMilliseconds = 1000000

	// Redis
	redisAddress      = "127.0.0.1:6379"
	redisPassword     = ""
	redisPoolSize     = 1000
	redisMinIdleConns = 10
	redisDB           = 0

	// Kafka
	kafkaBroker = "localhost:9092"
)

func setupTest() (SubscriberInterface, *gin.Context, *httptest.ResponseRecorder, error) {
	gin.SetMode(gin.TestMode)
	ctx := context.Background()
	pkg.InitLogger()

	// Init PostgreSQL
	var err error
	var db *gorm.DB
	db, err = NewDBConnection()
	if err != nil {
		fmt.Printf("failed to init db connection %s\n", err.Error())
		return nil, nil, nil, err
	}

	// Init redis
	cacheService, err := pkg.NewCacheCache(ctx, redisAddress, redisPassword, redisPoolSize, redisMinIdleConns, redisDB)
	if err != nil {
		fmt.Printf("failed to init redis client %s\n", err.Error())
		return nil, nil, nil, err
	}

	// Init kafka
	kafkaProducer, err := pkg.NewProducer(kafkaBroker)
	if err != nil {
		fmt.Printf("failed to init kafka client %s\n", err.Error())
		return nil, nil, nil, err
	}

	subscriberRepo := repository.NewSubscriber(db)

	// Initialize the controller.
	ctrl := NewSubscriberController(subscriberRepo, kafkaProducer, cacheService)

	// Create a new HTTP recorder and Gin context for testing.
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)

	return ctrl, c, recorder, nil
}

func TestSubscriberController_Create(t *testing.T) {
	ctrl, c, recorder, err := setupTest()
	assert.Nil(t, err)

	// Create whale first
	for i := 0; i < 10000; i++ {
		payload := model.Subscriber{
			TenantID:  "whale",
			Email:     gofakeit.Email(),
			FirstName: gofakeit.FirstName(),
			LastName:  gofakeit.LastName(),
		}
		jsonData, _ := json.Marshal(payload)
		req, err := http.NewRequest(http.MethodPost, "/subscribers", bytes.NewBuffer(jsonData))
		assert.Nil(t, err)
		req.Header.Set("Content-Type", "application/json")
		c.Request = req

		// Call the Create method of the controller.
		ctrl.Create(c)

		assert.Equal(t, http.StatusOK, recorder.Code, "Expected status code 200")
	}

	// Create small second
	for i := 0; i < 10; i++ {
		payload := model.Subscriber{
			TenantID:  "small",
			Email:     gofakeit.Email(),
			FirstName: gofakeit.FirstName(),
			LastName:  gofakeit.LastName(),
		}
		jsonData, _ := json.Marshal(payload)
		req, err := http.NewRequest(http.MethodPost, "/subscribers", bytes.NewBuffer(jsonData))
		assert.Nil(t, err)
		req.Header.Set("Content-Type", "application/json")
		c.Request = req

		// Call the Create method of the controller.
		ctrl.Create(c)

		assert.Equal(t, http.StatusOK, recorder.Code, "Expected status code 200")
	}
}

func BuildDsn() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		sqlHost, user, sqlPassword, dbname, port, sslmode, timezone)
}

func NewDBConnection() (*gorm.DB, error) {
	var err error

	// init database connection
	var orm *gorm.DB
	orm, err = gorm.Open(postgres.Open(BuildDsn()), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	var db *sql.DB
	db, err = orm.DB()
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(time.Duration(connMaxLifetimeMilliseconds) * time.Millisecond)

	return orm, nil
}
