package config

import (
	"fmt"
	"net/url"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDatabase() (*gorm.DB, error) {
	config := Config
	enCodedPassword := url.QueryEscape(config.Database.Password)
	uri := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		config.Database.Username,
		enCodedPassword,
		config.Database.Host,
		config.Database.Port,
		config.Database.Name,
	)

	db, err := gorm.Open(postgres.Open(uri), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(config.Database.MaxIdleConnections)
	sqlDB.SetMaxOpenConns(config.Database.maxOpenConnections)
	sqlDB.SetConnMaxLifetime(time.Duration(config.Database.MaxLifeTimeConnection) * time.Second)
	sqlDB.SetConnMaxIdleTime(time.Duration(config.Database.MaxIdleTime) + time.Second)
	return db, nil
}
