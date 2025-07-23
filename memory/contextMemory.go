package memory

import (
	"miniMem0/config"
	"miniMem0/db/sqldb"
	"miniMem0/llm"
	"miniMem0/model"
	"time"
)

// 用于管理记忆上下文，及智能体所处的环境,总结,对内容理解提供一个大致的方向性
type ContextMemoryHandler struct {
	summary    string
	config     *config.MemoryContextConfig
	llmHandler *llm.LLM
	sqlHandler *sqldb.SqlHandler
}

func NewMemoryContext(config *config.MemoryContextConfig, sqlHander *sqldb.SqlHandler, llm *llm.LLM) *ContextMemoryHandler {
	return &ContextMemoryHandler{
		summary:    "",
		sqlHandler: sqlHander,
		llmHandler: llm,
		config:     config,
	}
}

// 添加一条上下文, 并记录时间
func (m *ContextMemoryHandler) AddMemoryContext(from model.MemorySource, content string) error {
	err := m.sqlHandler.AddOriginalMemory(
		&model.OriginalMemory{
			Role:      from,
			Content:   content,
			CreatedAt: time.Now(),
		},
	)
	if err != nil {
		return err
	}
	return nil
}

// 返回记忆上下文
func (m *ContextMemoryHandler) GetMemoryContext() string {
	return m.summary
}
