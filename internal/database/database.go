package database

import (
	"fmt"
	"gopher-social-backend-server/pkg/logger"
	"gopher-social-backend-server/pkg/utils"
	"sync"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	log        = logger.GetLogger()
	PostgresDB *gorm.DB
	once       sync.Once
)

func NewPostgresDB() (*gorm.DB, error) {
	var err error

	once.Do(func() {
		POSTGRES_HOST := utils.GetEnvAsString("POSTGRES_HOST", "postgres")
		POSTGRES_PORT := utils.GetEnvAsInt("POSTGRES_PORT", 5432)
		POSTGRES_USER := utils.GetEnvAsString("POSTGRES_USER", "postgres")
		POSTGRES_PASSWORD := utils.GetEnvAsString("POSTGRES_PASSWORD", "postgres")
		POSTGRES_DB := utils.GetEnvAsString("POSTGRES_DB", "gopher_social")

		MAX_OPEN_CONNS := utils.GetEnvAsInt("MAX_OPEN_CONNS", 10)
		MAX_IDLE_CONNS := utils.GetEnvAsInt("MAX_IDLE_CONNS", 5)
		MAX_IDLE_TIME := utils.GetEnvAsDuration("MAX_IDLE_TIME", "15m")

		connectionString := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=UTC",
			POSTGRES_HOST, POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DB, POSTGRES_PORT,
		)

		PostgresDB, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{})
		if err != nil {
			log.Error("failed to connect to the database", zap.Error(err))
			return
		}

		sqlDB, err := PostgresDB.DB()
		if err != nil {
			log.Error("failed to retrieve database object", zap.Error(err))
			return
		}

		sqlDB.SetMaxOpenConns(MAX_OPEN_CONNS)
		sqlDB.SetMaxIdleConns(MAX_IDLE_CONNS)
		sqlDB.SetConnMaxIdleTime(MAX_IDLE_TIME)

		log.Info("connected to the database successfully with connection pooling settings",
			zap.Int("max_open_conns", MAX_OPEN_CONNS),
			zap.Int("max_idle_conns", MAX_IDLE_CONNS),
			zap.Duration("max_idle_time", MAX_IDLE_TIME),
		)
	})

	return PostgresDB, err
}

func CloseDatabase() error {
	sqlDatabase, err := PostgresDB.DB()
	if err != nil {
		log.Error("failed to retrieve database object", zap.Error(err))
		return err
	}
	if err := sqlDatabase.Close(); err != nil {
		log.Error("failed to close database connection", zap.Error(err))
		return err
	}
	log.Info("database connection closed successfully")
	return nil
}

func MigrateModel(model interface{}) error {
	if err := PostgresDB.AutoMigrate(model); err != nil {
		log.Error("failed to migrate model", zap.Error(err))
		return err
	}

	log.Info("model migrated successfully", zap.String("model", fmt.Sprintf("%T", model)))
	return nil
}
