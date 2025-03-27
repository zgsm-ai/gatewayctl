package postgres

import (
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"moul.io/zapgorm2"

	"github.com/zgsm-ai/gatewayctl/internal/pkg/config"
	"github.com/zgsm-ai/gatewayctl/internal/pkg/logger"
)

// New creates a new gorm db instance with the given options.
func NewDB() (*gorm.DB, error) {
	dbConf := config.App.Data.Database.Postgres
	zapLogger := zapgorm2.New(logger.NewZapLogger(logger.NewOptsFromConfig()))
	zapLogger.SetAsDefault()
	db, err := gorm.Open(
		postgres.Open(dbConf.Url),
		&gorm.Config{Logger: zapLogger},
	)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(dbConf.MaxOpenConnections)

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(dbConf.MaxIdleConnections)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}
