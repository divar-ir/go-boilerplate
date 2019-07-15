package main

import (
	"git.cafebazaar.ir/arcana261/golang-boilerplate/pkg/sql"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run necessary database migrations",
	Long:  `Migrate database to latest schema version`,
	Run:   migrateDatabase,
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}

func migrateDatabase(cmd *cobra.Command, args []string) {
	printVersion()
	config := loadConfigOrPanic(cmd)
	configureLoggerOrPanic(config.Logging)

	providerInstance := getProvider(config)
	migrater, ok := providerInstance.(sql.Migrater)
	if ok {
		err := migrater.Migrate()
		if err != nil {
			logrus.WithError(err).Panic("failed to migrate datbase")
		}
	}
}
