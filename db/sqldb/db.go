package sqldb

import (
	"miniMem0/config"
	"miniMem0/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SqlDB struct {
	db *gorm.DB
}

func NewSQL(cfg *config.SqlConfig) (*SqlDB, error) {
	db, err := gorm.Open(sqlite.Open(cfg.Path), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	// Migrate the schema
	db.AutoMigrate(&model.OriginalMemory{})
	return &SqlDB{db: db}, nil
}
