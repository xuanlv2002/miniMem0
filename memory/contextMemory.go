package memory

import (
	"miniMem0/config"
	"miniMem0/db/sqldb"
	"miniMem0/llm"
	"miniMem0/model"
)

// 用于管理记忆上下文，及智能体所处的环境,总结,对内容理解提供一个大致的方向性
type ContextMemoryHandler struct {
	config     *config.ContextMemoryConfig
	llmHandler *llm.LLM
	sqlHandler *sqldb.SqlHandler
}

func NewContextMemoryHandler(config *config.ContextMemoryConfig, sqlHander *sqldb.SqlHandler, llm *llm.LLM) *ContextMemoryHandler {

	return &ContextMemoryHandler{
		sqlHandler: sqlHander,
		llmHandler: llm,
		config:     config,
	}
}

// 返回记忆上下文
func (m *ContextMemoryHandler) GetMemoryContext() (*model.ContextMemory, error) {
	contextMemory, err := m.sqlHandler.GetLastContextMemory()
	if err != nil {
		return nil, err
	}
	return contextMemory, nil
}

// 总结记忆上下文
func (m *ContextMemoryHandler) SummaryMemoryContext() error {
	go func() {
		// 启动一个线程去异步总结记忆上下文
	}()
	return nil
}
