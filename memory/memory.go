package memory

import (
	"time"

	"github.com/sashabaranov/go-openai"
	"github.com/xuanlv2002/miniMem0/config"
	"github.com/xuanlv2002/miniMem0/db/sqldb"
	"github.com/xuanlv2002/miniMem0/db/vector"
	"github.com/xuanlv2002/miniMem0/llm"
	"github.com/xuanlv2002/miniMem0/model"
)

// 原始记忆结构体 包含角色 内容 时间

type MemorySystem struct {
	LongMemoryHandler    *LongMemoryHandler
	ShortMemoryHandler   *ShortMemroyHandler
	ContextMemoryHandler *ContextMemoryHandler
	sqlHandler           *sqldb.SqlHandler
	vectorHandler        *vector.Vector
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
		sqlHandler:           sqlHandler,
		vectorHandler:        vectorDB,
	}, nil
}

func (m *MemorySystem) FlushMemory() error {
	// 手动触发记忆更新
	m.ContextMemoryHandler.UpdateContextMemory()
	m.LongMemoryHandler.UpdateLongMemory()

	// 等待所有记忆处理完成
	m.ContextMemoryHandler.WaitDone()
	m.LongMemoryHandler.WaitDone()
	return nil
}

// 处理大模型输入内容
func (m *MemorySystem) ProcessInput(input string) (string, error) {
	// 传入激活内容
	activeMemory := &model.OriginalMemory{
		Role:      openai.ChatMessageRoleUser,
		Content:   input,
		CreatedAt: time.Now(),
	}

	// 获得完整短期记忆
	shortMemory, err := m.ShortMemoryHandler.GetShortMemory()
	if err != nil {
		return "", err
	}

	// 获得上下文记忆
	contextMemory, err := m.ContextMemoryHandler.GetContextMemory()
	if err != nil {
		return "", err
	}

	// 获得长期记忆
	longMemory, err := m.LongMemoryHandler.GetLongMemory(activeMemory.Content)
	if err != nil {
		return "", err
	}

	// 返回拼接后的prompt
	prompt := contextMemory.GetPrompt() + longMemory.GetPrompt() + shortMemory.GetPrompt() + "#用户输入: \n" + activeMemory.GetPrompt()

	// 将瞬时记忆存储 OriginalMemory
	err = m.sqlHandler.AddOriginalMemory(activeMemory)
	if err != nil {
		return "", err
	}
	return prompt, nil
}

// 处理大模型输出内容
func (m *MemorySystem) ProcessOutput(ouput string) error {
	// 将模型输出存储短期记忆
	outputMemory := &model.OriginalMemory{
		Role:      openai.ChatMessageRoleAssistant,
		Content:   ouput,
		CreatedAt: time.Now(),
	}
	err := m.sqlHandler.AddOriginalMemory(outputMemory)
	if err != nil {
		return err
	}

	// 更新上下文记忆
	m.ContextMemoryHandler.UpdateContextMemory()

	// 更新长期记忆
	m.LongMemoryHandler.UpdateLongMemory()

	// 等待所有记忆处理完成
	m.ContextMemoryHandler.WaitDone()
	m.LongMemoryHandler.WaitDone()
	return nil
}
