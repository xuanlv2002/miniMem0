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
	contextMemoryHandler := NewMemoryContext(options.GetMemoryContextConfig(), sqlHandler, llmModel)
	// 初始化长期记忆系统。
	longMemoryHandler := NewLongMemory(options.GetLongMemoryConfig(), vectorDB, sqlHandler, llmModel)
	// 初始化短期记忆系统
	shortMemoryHandler := NewShortMemory(options.GetShortMemoryConfig(), sqlHandler, llmModel)

	return &MemorySystem{
		ContextMemoryHandler: contextMemoryHandler,
		LongMemoryHandler:    longMemoryHandler,
		ShortMemoryHandler:   shortMemoryHandler,
	}, nil
}

// 处理大模型输入内容
func ProcessInput(input string) string {
	return ""
}

// 处理大模型输出内容
func ProcessOutput(ouput string) error {
	return nil
}
