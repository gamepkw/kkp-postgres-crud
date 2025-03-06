package database

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	Host      string
	Port      int
	User      string
	Password  string
	DBName    string
	Dialector gorm.Dialector
}

func New(ctx context.Context, cfg Config) (*gorm.DB, error) {
	var dialector gorm.Dialector
	if cfg.Dialector != nil {
		dialector = cfg.Dialector
	} else {
		dsn := fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Bangkok",
			cfg.Host,
			cfg.Port,
			cfg.User,
			cfg.Password,
			cfg.DBName,
		)
		dialector = postgres.Open(dsn)
	}
	db, err := gorm.Open(dialector, &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		return nil, errors.Wrap(err, "Can't initialize db session")
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, errors.Wrap(err, "Can't get PostgreSQL DB")
	}
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour * time.Duration(1))
	return db, err
}
