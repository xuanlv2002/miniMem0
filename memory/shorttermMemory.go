package memory

import (
	"miniMem0/config"
	"miniMem0/db/sqldb"
	"miniMem0/model"
)

/*
	用于管理短期记忆
	包括瞬时记忆 和 短时记忆
	其中瞬时记忆代表用户此刻的输入
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

	return &model.ShortMemory{
		Memorys: shortMemroy,
	}, nil
}
