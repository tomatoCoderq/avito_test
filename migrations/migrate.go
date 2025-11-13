package migrations

import (
	"embed"
	"log/slog"

	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
)

//go:embed *.sql
var embedMigrations embed.FS

func MigrateToLatest(db *gorm.DB, log *slog.Logger) error {
    sqlDB, err := db.DB()
    if err != nil {
        return err
    }

    goose.SetBaseFS(embedMigrations)

    if err := goose.SetDialect("postgres"); err != nil {
        return err
    }

    if err := goose.Up(sqlDB, "."); err != nil {
        return err
    }

    log.Info("Database migration completed successfully")
    return nil
}