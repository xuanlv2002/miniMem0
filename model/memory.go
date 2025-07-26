package model

import (
	"fmt"
	"time"
)

type MemorySource string

// 原始记忆信息
type OriginalMemory struct {
	ID        int64 `gorm:"primaryKey"`
	Role      MemorySource
	Content   string
	CreatedAt time.Time // 内置默认时间
}

func (o *OriginalMemory) GetPrompt() string {
	return fmt.Sprintf("%s: %s", o.Role, o.Content)
}

// 记忆上下文结构体
type ContextMemory struct {
	ID            int64
	Summary       string    // 用于管理记忆上下文，及智能体所处的环境,总结,对内容理解提供一个大致的方向性
	LastSummaryID int64     // 最后一次总结的id
	UpdatedAt     time.Time // 最近修改时间
	CreatedAt     time.Time // 内置默认时间
}

func (c *ContextMemory) GetPrompt() string {
	if c.Summary == "" {
		return "#上下文摘要记忆: \n 暂无摘要信息\n"
	}
	return "#上下文摘要记忆: \n" + c.Summary + "\n"
}

// 短期记忆结构体
type ShortMemory struct {
	Memorys []OriginalMemory
}

func (s *ShortMemory) GetPrompt() string {
	content := "#短期记忆: \n"
	for _, memory := range s.Memorys {
		content += fmt.Sprintf("%v:%v\n", memory.Role, memory.Content)
	}
	return content
}

type LongMemoryItem struct {
	ID       string            `json:"id"`
	Text     string            `json:"text"`
	Meta     map[string]string `json:"meta"`
	Similary float32           `json:"similary,omitempty"` // 基于文本内容符合度搜索   基于元数据搜索
}

// 长期记忆结构体
type LongMemory struct {
	ID               int64
	LastExtractionID int64            // 最近一次抽取长期记忆ID
	VectorMemorys    []LongMemoryItem `gorm:"-"` // 基于语义相似搜索
	// 基于模型来把自然语言转为结构化查询 来获得更全面的关系数据 暂未实现
	// 通过混合长期记忆搜索的方式 获得更全面的消息信息(function call?)
	UpdatedAt time.Time // 最近修改时间
}

func (l *LongMemory) GetPrompt() string {
	if len(l.VectorMemorys) == 0 {
		return "#长期记忆: \n 暂无长期记忆信息"
	}
	content := "#长期记忆: \n"
	for _, item := range l.VectorMemorys {
		content += fmt.Sprintf("记忆内容:%v \n记忆元信息:%v \n记忆相关度:%v \n\n", item.Text, item.Meta, item.Similary)
	}
	return content
}

type MemoryEvent struct {
	ID    string            `json:"id"`
	Text  string            `json:"text"`
	Meta  map[string]string `json:"meta"`
	Event string            `json:"event"`
}

type Fact struct {
	Content    string `json:"content"`
	AppearTime string `json:"appearTime"`
	About      string `json:"about"`
}
