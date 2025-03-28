package cmd

import (
	"database/sql"
	"fmt"
	"github.com/spf13/cobra"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"time"

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

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: MainServiceName,
	}
	cmd.AddCommand(serveCmd)
	cmd.AddCommand(migrateCmd)

	return cmd
}

func Execute() {
	pkg.InitLogger()

	rootCmd := NewRootCmd()

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
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
