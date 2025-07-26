package memory

import (
	"context"
	"fmt"

	"sync"
	"time"

	"github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
	"github.com/xuanlv2002/miniMem0/config"
	"github.com/xuanlv2002/miniMem0/db/sqldb"
	"github.com/xuanlv2002/miniMem0/llm"
	"github.com/xuanlv2002/miniMem0/model"
	"github.com/xuanlv2002/miniMem0/prompt"
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
func (m *ContextMemoryHandler) GetContextMemory() (*model.ContextMemory, error) {
	contextMemory, err := m.sqlHandler.GetLastContextMemory()
	if err != nil {
		return nil, err
	}
	return contextMemory, nil
}

// 异步总结记忆上下文 避免阻塞记忆主线程
func (m *ContextMemoryHandler) UpdateContextMemory() {
	// 等待
	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		err := m.SummaryContextMemory()
		if err != nil {
			logrus.Errorf("SummaryContextMemory error: %v", err)
		}
	}()
}

// 这个函数需要加锁串行 如果用户问的特别快 导致gap没有清0 导致问多次大模型,总结多次, 最新的summary 可能被老的覆盖掉
// 立即总结记忆上下文
func (m *ContextMemoryHandler) SummaryContextMemory() error {
	// 加锁
	m.mu.Lock()
	defer m.mu.Unlock()

	// 获取上下文记忆
	contextMemory, err := m.sqlHandler.GetLastContextMemory()
	if err != nil {
		return err
	}

	// 获取未总结的记忆数量
	count, err := m.sqlHandler.GetUnSummarizedMemoryCount(contextMemory.LastSummaryID)
	if err != nil {
		return err
	}

	// 如果未总结数量小于等于gap值则不进行总结 当大于gap的第一个则总结 加入gap设置为5  当新增5次信息 则对5次信息统一进行总结
	if count < int64(m.config.SummaryGap) {
		return nil
	}

	// 如果大于了gap值则一次性进行总结 总结过程中 如果用户继续提问 会被阻塞 可以支持并发 如果程序挂断 重启后正常进行总结
	originalMemories, findCount, err := m.sqlHandler.GetLastOriginalMemory(int(count))
	if err != nil {
		return err
	}

	if findCount <= 0 {
		// 如果没有未总结的记忆则不进行总结
		return nil
	}

	content := "#已总结内容: \n" + contextMemory.Summary
	content += "\n#待总结对话: \n"

	lastSummaryId := originalMemories[len(originalMemories)-1].ID
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

	// 使用大模型总结记忆
	summary, err := m.llmHandler.Chat(context.Background(), messages)
	if err != nil {
		return err
	}
	// 总结成功 这里把gap清空 后面的发现gap空之后 就++ 不进行总结了
	contextMemory.Summary = summary.Content
	contextMemory.LastSummaryID = lastSummaryId
	contextMemory.UpdatedAt = time.Now()
	// 更新数据库
	err = m.sqlHandler.SaveContextMemory(contextMemory)
	if err != nil {
		return err
	}
	return nil
}
