package memory

import (
	"miniMem0/config"
	"miniMem0/db/sqldb"
	"miniMem0/db/vector"
	"miniMem0/llm"
)

// 原始记忆结构体 包含角色 内容 时间

type MemorySystem struct {
	LongMemoryHandler    *LongMemoryHandler
	ShortMemoryHandler   *ShortMemroyHandler
	ContextMemoryHandler *ContextMemoryHandler
}

func NewMemorySystem(options *config.Config) (*MemorySystem, error) {
	// 初始化LLM
	llmModel := llm.NewLLM(options.GetChatConfig())
	// // 初始化Embedding
	embeddingModel := llm.NewEmbedding(options.GetEmbeddingConfig())
	// 初始化向量数据库
	vectorDB, err := vector.NewVector(options.GetVectorConfig(), embeddingModel.GetEmbeddingFunc())
	if err != nil {
		return nil, err
	}
	// 初始化SQL数据库
	sqlHandler, err := sqldb.NewSQL(options.GetSqlConfig())
	if err != nil {
		return nil, err
	}
	// 初始化记忆上下文系统
	contextMemoryHandler := NewContextMemoryHandler(options.GetMemoryContextConfig(), sqlHandler, llmModel)
	// 初始化长期记忆系统。
	longMemoryHandler := NewLongMemory(options.GetLongMemoryConfig(), vectorDB, sqlHandler, llmModel)
	// 初始化短期记忆系统
	shortMemoryHandler := NewShortMemoryHandler(options.GetShortMemoryConfig(), sqlHandler)

	return &MemorySystem{
		ContextMemoryHandler: contextMemoryHandler,
		LongMemoryHandler:    longMemoryHandler,
		ShortMemoryHandler:   shortMemoryHandler,
	}, nil
}

func (m *MemorySystem) InitMemory() error {
	return nil
}

// 处理大模型输入内容
func (m *MemorySystem) ProcessInput(input string) string {
	// 存储OriginalMemory
	// 获得上下文记忆
	// 获得长期记忆
	// 获得短期记忆
	// 拼接记忆内容
	// 返回拼接后的prompt
	return ""
}

// 处理大模型输出内容
func (m *MemorySystem) ProcessOutput(ouput string) error {
	// 存储OriginalMemory
	return nil
}
