package memory

import (
	"miniMem0/config"
	"miniMem0/db/sqldb"
	"miniMem0/model"
	"sort"
)

/*
	用于管理短期记忆
	包括瞬时记忆 和 短时记忆
	其中瞬时记忆代表用户此刻的输入 最近的一个 Q&A对
	短时记忆代表了 最近的m个 Q&A 对
	通过短时记忆可以让大模型的输出更加的符合用户的意图
	且短时记忆可以被抽取为长期记忆
*/

type ShortMemroyHandler struct {
	config     *config.ShortMemoryConfig
	sqlHandler *sqldb.SqlHandler
}

func NewShortMemoryHandler(config *config.ShortMemoryConfig, sqlHandler *sqldb.SqlHandler) *ShortMemroyHandler {
	return &ShortMemroyHandler{
		config:     config,
		sqlHandler: sqlHandler,
	}
}

// 拉取短期记忆
func (s *ShortMemroyHandler) GetShortMemory() (*model.ShortMemory, error) {
	shortMemroy, _, err := s.sqlHandler.GetLastOriginalMemory(s.config.ShortWindow)
	if err != nil {
		return nil, err
	}
	// 小到大排序
	sort.Slice(shortMemroy, func(i, j int) bool {
		return shortMemroy[i].ID < shortMemroy[j].ID
	})

	return &model.ShortMemory{
		Memorys: shortMemroy,
	}, nil
}
