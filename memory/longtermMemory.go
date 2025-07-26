package memory

import (
	"context"
	"encoding/json"
	"fmt"
	"miniMem0/config"
	"miniMem0/db/sqldb"
	"miniMem0/db/vector"
	"miniMem0/llm"
	"miniMem0/model"
	"miniMem0/prompt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/philippgille/chromem-go"
	"github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
)

/*
	用于管理长期记忆
	长期记忆,及一些持久性的记忆,如用户的信息,一些历史记录等
	长期记忆通过向量数据库存储,并根据用户的输入进行检索,返回结果
*/

type LongMemoryHandler struct {
	config     *config.LongMemoryConfig
	vector     *vector.Vector
	llmHandler *llm.LLM
	sqlHandler *sqldb.SqlHandler
	mu         sync.Mutex     // 长期记忆锁
	wg         sync.WaitGroup // 用来等待所有任务完成
}

// 新建长期记忆系统
func NewLongMemory(config *config.LongMemoryConfig, vector *vector.Vector, sqlHandler *sqldb.SqlHandler, llmModel *llm.LLM) *LongMemoryHandler {
	return &LongMemoryHandler{
		vector:     vector,
		llmHandler: llmModel,
		sqlHandler: sqlHandler,
		config:     config,
	}
}

// 等待所有任务完成
func (l *LongMemoryHandler) WaitDone() {
	l.wg.Wait()
}

// 获得相关长期记忆
func (l *LongMemoryHandler) GetLongMemory(text string) (*model.LongMemory, error) {
	var LongMemory model.LongMemory
	// 搜索
	ret, err := l.vector.Search(context.Background(), text)
	if err != nil {
		return nil, err
	}

	var vectorMemory = make([]model.LongMemoryItem, 0)
	for _, v := range ret {
		vectorMemory = append(vectorMemory, model.LongMemoryItem{
			Text:     v.Content,
			Meta:     v.Metadata,
			Similary: v.Similarity,
		})
	}

	LongMemory.VectorMemorys = vectorMemory

	return &LongMemory, nil
}

// 更新长期记忆 异步更新 不对系统进行阻塞
func (l *LongMemoryHandler) UpdateLongMemory() {
	l.wg.Add(1)
	go func() {
		defer l.wg.Done()
		err := l.SaveLongMemory()
		if err != nil {
			logrus.Errorf("LongMemory error: %v", err)
		}
	}()
}

// 更新长期记忆 异步更新 不对系统进行阻塞
func (l *LongMemoryHandler) SaveLongMemory() error {
	// 加锁
	l.mu.Lock()
	defer l.mu.Unlock()
	// 获得长期记忆位置 获得长期记忆已经存储到的位置
	longMemory, err := l.sqlHandler.GetLastLongMemroy()
	if err != nil {
		logrus.Errorf("failed to get last long memory: %v", err)
		return err
	}
	// 判断是否需要更新记忆
	count, err := l.sqlHandler.GetUnExtractionMemoryCount(longMemory.LastExtractionID)
	if err != nil {
		logrus.Errorf("failed to get unextraction memory count: %v", err)
		return err
	}

	// 如果小于则不更新记忆
	if count < int64(l.config.LongGap) {
		logrus.Infof("No new memories to extract, current count: %d, required gap: %d", count, l.config.LongGap)
		return nil
	}

	// 获得上下文记忆
	contextMemory, err := l.sqlHandler.GetLastContextMemory()
	if err != nil {
		logrus.Errorf("failed to get last context memory: %v", err)
		return err
	}

	// 获得未抽取的记忆
	originalMemories, findCount, err := l.sqlHandler.GetLastOriginalMemory(int(count))
	if err != nil {
		logrus.Errorf("failed to get last original memory: %v", err)
		return err
	}

	if findCount <= 0 {
		// 如果没有未总结的记忆则不进行总结
		logrus.Info("No new memories to extract, find count is 0")
		return nil
	}

	// 组装信息
	var content string
	content += contextMemory.GetPrompt()
	content += "\n !！注解 !! 记忆上下文是基于大模型对整体对话的一个总结，可能与当前对话不完全相关，但可以作为参考。\n\n"

	content += "#待提取信息记忆: \n"
	for _, v := range originalMemories {
		content += string(v.Role) + ":" + v.Content
		content += "\n记忆元数据"
		content += "记忆时间:" + v.CreatedAt.Format("2006-01-02 15:04:05") + "\n\n"
	}

	// 抽取长期记忆
	facts, err := l.ExtractFacts(context.Background(), content)
	if err != nil {
		logrus.Errorf("failed to extract facts: %v", err)
		return err
	}
	if len(facts) == 0 {
		logrus.Info("no new facts found")
		longMemory.LastExtractionID = originalMemories[len(originalMemories)-1].ID
		longMemory.UpdatedAt = time.Now()
		err = l.sqlHandler.SaveLongMemoryLastExtractionID(longMemory)
		if err != nil {
			return err
		}
		return nil
	}

	// 抽取相关长期记忆
	var retrievedOldMemoriesMap = make(map[string]model.LongMemoryItem)
	for _, fact := range facts {
		memories, err := l.vector.Search(context.Background(), fact.Content)
		if err != nil {
			return fmt.Errorf("failed to search memories: %v", err)
		}
		// 抽取到的相关记忆
		for _, mem := range memories {
			retrievedOldMemoriesMap[mem.ID] = model.LongMemoryItem{
				ID:   mem.ID,       // 记忆ID
				Text: mem.Content,  // 记忆内容
				Meta: mem.Metadata, // 记忆元数据
			}
		}
	}
	// 去重
	var retrievedOldMemories = make([]model.LongMemoryItem, 0)
	for _, v := range retrievedOldMemoriesMap {
		retrievedOldMemories = append(retrievedOldMemories, v)
	}
	// 解决记忆冲突
	// 对记忆进行修改处理
	safeMemories, err := l.processMemory(context.Background(), facts, retrievedOldMemories)
	if err != nil {
		return fmt.Errorf("failed to process memories: %v", err)
	}
	fmt.Println("safeMemories:", safeMemories)
	// 更新长期记忆
	for _, mem := range safeMemories {
		event := mem.Event
		memoryID := mem.ID
		text := mem.Text
		meta := mem.Meta

		switch event {
		case "ADD":
			if _, err := l.addMemory(context.Background(), text, meta); err != nil {
				return fmt.Errorf("failed to add memory: %v", err)
			}
			logrus.Infof("Added memory: %s", text)
		case "UPDATE":
			if err := l.updateMemory(context.Background(), memoryID, text, meta); err != nil {
				return fmt.Errorf("failed to update memory: %v", err)
			}
			logrus.Infof("Updated memory: %s", text)
		case "DELETE":
			if err := l.deleteMemory(context.Background(), memoryID); err != nil {
				return fmt.Errorf("failed to delete memory: %v", err)
			}
			logrus.Infof("Deleted memory: %s", memoryID)
		case "NONE":
			logrus.Infof("Keeping memory unchanged: %s", text)
		}
	}

	// 更新长期记忆位置
	longMemory.LastExtractionID = originalMemories[len(originalMemories)-1].ID
	longMemory.UpdatedAt = time.Now()
	err = l.sqlHandler.SaveLongMemoryLastExtractionID(longMemory)
	if err != nil {
		return err
	}
	return nil
}

// 事实提取 提取长期记忆内容
func (l *LongMemoryHandler) ExtractFacts(ctx context.Context, conversation string) ([]model.Fact, error) {
	// 调用LLM进行事实提取
	result, err := l.llmHandler.Chat(ctx, []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: prompt.FACT_EXTRACTION_PROMPT,
		}, {
			Role:    openai.ChatMessageRoleUser,
			Content: conversation,
		},
	})
	if err != nil {
		return nil, err
	}

	var response struct {
		Facts []model.Fact `json:"facts"`
	}
	fmt.Println("大模型输出:", result.Content)
	if err := json.Unmarshal([]byte(result.Content), &response); err != nil {
		return nil, fmt.Errorf("failed to parse LLM response: %v", err)
	}

	return response.Facts, nil
}
func (l *LongMemoryHandler) processMemory(ctx context.Context, newFacts []model.Fact, oldMemory []model.LongMemoryItem) ([]model.MemoryEvent, error) {
	content := "#新获取的事实: \n"
	for _, fact := range newFacts {
		content += fmt.Sprintf("   -内容: %s, 出现时间: %s, 关于: %s\n", fact.Content, fact.AppearTime, fact.About)
	}

	content += "\n#可能相关的记忆: \n"
	for _, v := range oldMemory {
		content += fmt.Sprintf("   -ID: %s, 内容: %s, 元数据: %v\n", v.ID, v.Text, v.Meta)
	}

	result, err := l.llmHandler.Chat(ctx, []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: prompt.MEMORY_PROCESSING_PROMPT,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: content,
		},
	})
	if err != nil {
		return nil, err
	}

	var response struct {
		Memory []model.MemoryEvent `json:"memory"`
	}

	fmt.Println("大模型记忆更新输出:", result.Content)
	if err := json.Unmarshal([]byte(result.Content), &response); err != nil {
		return nil, fmt.Errorf("failed to parse LLM response: %v", err)
	}

	return response.Memory, nil
}

// 添加记忆
func (l *LongMemoryHandler) addMemory(ctx context.Context, text string, metadata map[string]string) (string, error) {
	// 持久记忆的ID
	memoryID := uuid.New().String()

	// 将文本转换为向量 并存入数据库
	err := l.vector.Add(ctx, []chromem.Document{
		{
			ID:       memoryID,
			Metadata: metadata,
			Content:  text,
		},
	}, 1)
	if err != nil {
		return "", err
	}

	return memoryID, nil
}

// 更新记忆
func (l *LongMemoryHandler) updateMemory(ctx context.Context, memoryID, newText string, metadata map[string]string) error {
	// 将文本转换为向量 并存入数据库
	err := l.vector.Add(ctx, []chromem.Document{
		{
			ID:       memoryID,
			Metadata: metadata,
			Content:  newText,
		},
	}, 1)
	if err != nil {
		return err
	}

	return nil

}

// 删除记忆
func (l *LongMemoryHandler) deleteMemory(ctx context.Context, memoryID string) error {
	err := l.vector.Delete(ctx, []string{memoryID})
	if err != nil {
		return fmt.Errorf("failed to delete memory: %v", err)
	}

	return nil
}
