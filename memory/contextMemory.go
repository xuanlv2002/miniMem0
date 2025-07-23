package memory

import (
	"context"
	"fmt"
	"miniMem0/config"
	"miniMem0/db/sqldb"
	"miniMem0/llm"
	"miniMem0/model"
	"miniMem0/prompt"
	"sync"

	"github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
)

// 用于管理记忆上下文，及智能体所处的环境,总结,对内容理解提供一个大致的方向性
type ContextMemoryHandler struct {
	config     *config.ContextMemoryConfig
	llmHandler *llm.LLM
	sqlHandler *sqldb.SqlHandler
	mu         sync.Mutex     // 用来保证SummaryMemoryContext函数的串行
	wg         sync.WaitGroup // 用来等待所有任务完成
}

func NewContextMemoryHandler(config *config.ContextMemoryConfig, sqlHander *sqldb.SqlHandler, llm *llm.LLM) *ContextMemoryHandler {
	return &ContextMemoryHandler{
		sqlHandler: sqlHander,
		llmHandler: llm,
		config:     config,
	}
}

func (m *ContextMemoryHandler) WaitDone() {
	m.wg.Wait()
}

// 返回获得的上下文记忆
func (m *ContextMemoryHandler) GetMemoryContext() (*model.ContextMemory, error) {
	contextMemory, err := m.sqlHandler.GetLastContextMemory()
	if err != nil {
		return nil, err
	}
	return contextMemory, nil
}

// 异步总结记忆上下文 避免阻塞记忆主线程
func (m *ContextMemoryHandler) UpdateContextMemory() {
	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		err := m.SummaryMemoryContext()
		if err != nil {
			logrus.Errorf("SummaryMemoryContext error: %v", err)
		}
	}()
}

// 这个函数需要加锁串行 如果用户问的特别快 导致gap没有清0 导致问多次大模型,总结多次, 最新的summary 可能被老的覆盖掉
// 立即总结记忆上下文
func (m *ContextMemoryHandler) SummaryMemoryContext() error {
	// 加锁
	m.mu.Lock()
	defer m.mu.Unlock()
	// 获取上下文记忆
	contextMemory, err := m.sqlHandler.GetLastContextMemory()
	if err != nil {
		return err
	}
	// 如果小于gap值则不进行总结 并把gap++, 这里如果每锁 会导致gap值不正确
	if contextMemory.Gap < int64(m.config.SummaryGap) {
		contextMemory.Gap++
		err = m.sqlHandler.SaveContextMemory(contextMemory)
		if err != nil {
			return err
		}
		return nil
	}
	// 如果大于了gap值则进行总结
	originalMemories, _, err := m.sqlHandler.GetLastOriginalMemory(m.config.SummaryGap)
	content := "#已总结内容: \n" + contextMemory.Summary
	content += "\n#待总结对话: \n"
	for _, v := range originalMemories {
		content += fmt.Sprintf("%v:%v\n", v.Role, v.Content)
	}

	messages := []openai.ChatCompletionMessage{
		{
			Role:    "system",
			Content: prompt.CONTEXT_MEMORY_SUMMARY_PROMPT,
		},
		{
			Role:    "user",
			Content: content,
		},
	}

	fmt.Println("log:", content)

	// 使用大模型总结记忆
	summary, err := m.llmHandler.Chat(context.Background(), messages)
	if err != nil {
		return err
	}
	// 总结成功 这里把gap清空 后面的发现gap空之后 就++ 不进行总结了
	contextMemory.Summary = summary.Content
	contextMemory.Gap = 0
	// 更新数据库
	err = m.sqlHandler.SaveContextMemory(contextMemory)
	return nil
}
