package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gorm.io/gorm"

	"webhook/migration"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run the migration",
	Run:   runMigrate,
}

func runMigrate(cmd *cobra.Command, args []string) {
	// Init PostgreSQL
	var err error
	var orm *gorm.DB
	orm, err = NewDBConnection()
	if err != nil {
		panic(err)
	}
	err = migration.Migrate(orm)
	if err != nil {
		log.WithError(err).Error("migration failed")
	}
}
