package memory

import (
	"miniMem0/model"
	"time"
)

// 用于管理记忆上下文，及智能体所处的环境,总结,对内容理解提供一个大致的方向性
type MemoryContext struct {
	Summary       string
	MemroyContext []model.OriginalMemory
	Length        int
}

func NewMemoryContext() *MemoryContext {
	return &MemoryContext{
		Summary:       "",
		MemroyContext: []model.OriginalMemory{},
	}
}

func (m *MemoryContext) AddMemoryContext(from model.MemorySource, content string) {
	m.MemroyContext = append(m.MemroyContext, model.OriginalMemory{
		Role:       from,
		Content:    content,
		CreateTime: time.Now(),
	})
}
