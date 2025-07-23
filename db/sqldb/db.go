package sqldb

import (
	"miniMem0/config"
	"miniMem0/model"
	"sort"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SqlHandler struct {
	DB *gorm.DB
}

func NewSQL(cfg *config.SqlConfig) (*SqlHandler, error) {
	db, err := gorm.Open(sqlite.Open(cfg.Path), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	// Migrate the schema
	db.AutoMigrate(&model.OriginalMemory{}, &model.ContextMemory{})
	return &SqlHandler{DB: db}, nil
}

/* 原始记忆处理函数 */
// 添加一条记忆
func (db *SqlHandler) AddOriginalMemory(memory *model.OriginalMemory) error {
	return db.DB.Create(memory).Error
}

// 获得最近的n条记忆 最新的在前面
func (db *SqlHandler) GetLastOriginalMemory(count int) ([]model.OriginalMemory, int64, error) {
	var ret []model.OriginalMemory
	var retCount int64
	err := db.DB.Order("id desc").Limit(count).Find(&ret).Count(&retCount).Error
	if err != nil {
		return nil, 0, err
	}
	// 小到大排序
	sort.Slice(ret, func(i, j int) bool {
		return ret[i].ID < ret[j].ID
	})

	return ret, retCount, nil
}

// 获得所有的记忆
func (db *SqlHandler) GetTotalOriginalMemory() ([]model.OriginalMemory, int64, error) {
	var ret []model.OriginalMemory
	// 获得所有数据 没有数据返回空
	var count int64
	err := db.DB.Find(&ret).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	return ret, count, nil
}

/* 上下文记忆处理函数 */
// 获得上下文记忆
func (db *SqlHandler) GetLastContextMemory() (*model.ContextMemory, error) {
	var ret model.ContextMemory
	err := db.DB.Last(&ret).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &ret, nil
}

// 更新上下文记忆
func (db *SqlHandler) SaveContextMemory(memory *model.ContextMemory) error {
	return db.DB.Save(memory).Error
}

// 获得未总结的记忆的个数
func (db *SqlHandler) GetUnSummarizedMemoryCount(lastSummaryID int64) (int64, error) {
	var count int64
	err := db.DB.Model(&model.OriginalMemory{}).Where("id > ?", lastSummaryID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
