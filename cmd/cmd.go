package cmd

import (
	"database/sql"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

const (
	// SQL
	sqlHost                     = "localhost"
	port                        = "5432"
	user                        = "postgres"
	sqlPassword                 = "your_password"
	dbname                      = "your_db"
	sslmode                     = "disable"
	timezone                    = "UTC"
	maxOpenConns                = 100
	maxIdleConns                = 10
	connMaxLifetimeMilliseconds = 1000000

	// Redis
	redisAddress      = ""
	redisPassword     = ""
	redisPoolSize     = 1000
	redisMinIdleConns = 10
	redisDB           = 0
)

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
