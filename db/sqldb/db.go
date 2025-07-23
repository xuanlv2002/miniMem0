package sqldb

import (
	"miniMem0/config"
	"miniMem0/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SqlDB struct {
	DB *gorm.DB
}

func NewSQL(cfg *config.SqlConfig) (*SqlDB, error) {
	db, err := gorm.Open(sqlite.Open(cfg.Path), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	// Migrate the schema
	db.AutoMigrate(&model.OriginalMemory{})
	return &SqlDB{DB: db}, nil
}

// 添加一条记忆
func (db *SqlDB) AddOriginalMemory(memory *model.OriginalMemory) error {
	return db.DB.Create(memory).Error
}

// 获得最近的n条记忆 最新的在前面
func (db *SqlDB) GetLastOriginalMemory(count int) ([]model.OriginalMemory, int64, error) {
	var ret []model.OriginalMemory
	var retCount int64
	err := db.DB.Order("id desc").Limit(count).Find(&ret).Count(&retCount).Error
	if err != nil {
		return nil, 0, err
	}
	return ret, retCount, nil
}

// 获得所有的记忆
func (db *SqlDB) GetTotalOriginalMemory() ([]model.OriginalMemory, int64, error) {
	var ret []model.OriginalMemory
	// 获得所有数据 没有数据返回空
	var count int64
	err := db.DB.Find(&ret).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	return ret, count, nil
}
