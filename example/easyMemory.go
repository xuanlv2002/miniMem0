package example

/*
	记忆系统的简单实现原理,
	借鉴 https://github.com/TinyAgen/TinyMem0 的实现
	仅用于参考
*/

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"miniMem0/config"
	"miniMem0/llm"
	"miniMem0/prompt"
	"miniMem0/vector"
	"time"

	"github.com/google/uuid"
	"github.com/philippgille/chromem-go"
	"github.com/sashabaranov/go-openai"
)

// 记忆系统
type MemorySystem struct {
	LLM    *llm.LLM
	Vector *vector.Vector
	Option *config.Config
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

	return &MemorySystem{LLM: llmModel, Vector: vectorDB, Option: options}, nil
}

func (m *MemorySystem) String() string {
	return m.Option.String()
}

type Memory struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

// 事实提取 提取长期记忆内容
func (ms *MemorySystem) ExtractFacts(ctx context.Context, conversation string) ([]string, error) {
	// 调用LLM进行事实提取
	result, err := ms.LLM.Chat(ctx, []openai.ChatCompletionMessage{
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
		Facts []string `json:"facts"`
	}
	if err := json.Unmarshal([]byte(result.Content), &response); err != nil {
		return nil, fmt.Errorf("failed to parse LLM response: %v", err)
	}

	return response.Facts, nil
}

type MemoryItem struct {
	ID    string `json:"id"`
	Text  string `json:"text"`
	Event string `json:"event"`
}

// 记忆对比
func (ms *MemorySystem) ProcessMemory(ctx context.Context, newFacts []string, existingMemories []Memory) ([]MemoryItem, error) {
	input := map[string]interface{}{
		"new_facts":         newFacts,
		"existing_memories": existingMemories,
	}
	fmt.Println("new_facts:", newFacts)
	fmt.Println("existing_memories:", existingMemories)

	inputJSON, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal input: %v", err)
	}

	result, err := ms.LLM.Chat(ctx, []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: prompt.MEMORY_PROCESSING_PROMPT,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: string(inputJSON),
		},
	})
	if err != nil {
		return nil, err
	}

	var response struct {
		Memory []MemoryItem `json:"memory"`
	}

	fmt.Println("提取的JSON数据:", result.Content)
	if err := json.Unmarshal([]byte(result.Content), &response); err != nil {
		return nil, fmt.Errorf("failed to parse LLM response: %v", err)
	}

	return response.Memory, nil
}

// 添加记忆
func (ms *MemorySystem) AddMemory(ctx context.Context, text string, metadata map[string]string) (string, error) {
	// 持久记忆的ID
	memoryID := uuid.New().String()

	// 将文本转换为向量 并存入数据库
	err := ms.Vector.Add(ctx, []chromem.Document{
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
func (ms *MemorySystem) UpdateMemory(ctx context.Context, memoryID, newText string, metadata map[string]string) error {
	// 将文本转换为向量 并存入数据库
	err := ms.Vector.Add(ctx, []chromem.Document{
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
func (ms *MemorySystem) DeleteMemory(ctx context.Context, memoryID string) error {
	err := ms.Vector.Delete(context.Background(), []string{memoryID})
	if err != nil {
		return fmt.Errorf("failed to delete memory: %v", err)
	}

	return nil
}

// WriteMemory handles the complete memory writing process
func (ms *MemorySystem) WriteMemory(ctx context.Context, conversation string) error {
	// 提取事实
	newFacts, err := ms.ExtractFacts(ctx, conversation)
	if err != nil {
		return fmt.Errorf("failed to extract facts: %v", err)
	}
	if len(newFacts) == 0 {
		log.Println("No facts extracted")
		return nil
	}

	log.Printf("Extracted facts: %v", newFacts)

	// 查询事实相关的记忆
	var retrievedOldMemoriesMap = make(map[string]Memory)
	for _, fact := range newFacts {
		memories, err := ms.Vector.Search(ctx, fact, 10, 0)
		if err != nil {
			return fmt.Errorf("failed to search memories: %v", err)
		}

		for _, mem := range memories {
			retrievedOldMemoriesMap[mem.ID] = Memory{
				ID:   mem.ID,
				Text: mem.Content,
			}
		}
	}
	var retrievedOldMemories = make([]Memory, 0)
	for _, v := range retrievedOldMemoriesMap {
		retrievedOldMemories = append(retrievedOldMemories, v)
	}
	// 对记忆进行修改处理
	processedMemories, err := ms.ProcessMemory(ctx, newFacts, retrievedOldMemories)
	if err != nil {
		return fmt.Errorf("failed to process memories: %v", err)
	}

	for _, mem := range processedMemories {
		event := mem.Event
		memoryID := mem.ID
		text := mem.Text

		metadata := map[string]string{
			"created_at": time.Now().Format(time.RFC3339),
		}
		fmt.Println("event:", mem)
		switch event {
		case "ADD":
			if _, err := ms.AddMemory(ctx, text, metadata); err != nil {
				return fmt.Errorf("failed to add memory: %v", err)
			}
			log.Printf("Added memory: %s", text)
		case "UPDATE":
			if err := ms.UpdateMemory(ctx, memoryID, text, metadata); err != nil {
				return fmt.Errorf("failed to update memory: %v", err)
			}
			log.Printf("Updated memory: %s", text)
		case "DELETE":
			if err := ms.DeleteMemory(ctx, memoryID); err != nil {
				return fmt.Errorf("failed to delete memory: %v", err)
			}
			log.Printf("Deleted memory: %s", memoryID)
		case "NONE":
			log.Printf("Keeping memory unchanged: %s", text)
		}
	}

	return nil
}

type MemoryRet struct {
	ID       string
	Text     string
	Meta     map[string]string
	Similary float32
}

// SearchMemory searches for memories
func (ms *MemorySystem) SearchMemory(ctx context.Context, query string) ([]MemoryRet, error) {
	ret, err := ms.Vector.Search(ctx, query, 10, 0.5)
	if err != nil {
		return nil, err
	}
	memorys := make([]MemoryRet, 0, len(ret))
	for _, v := range ret {
		memorys = append(memorys, MemoryRet{
			ID:       v.ID,
			Text:     v.Content,
			Meta:     v.Metadata,
			Similary: v.Similarity,
		})
	}
	return memorys, nil
}
