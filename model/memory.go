package model

import "time"

type MemorySource string

const (
	UserRole MemorySource = "user"
	AIRole   MemorySource = "assistant"
	ToolRole MemorySource = "tool"
)

// 原始记忆信息
type OriginalMemory struct {
	ID        int64 `gorm:"primaryKey"`
	Role      MemorySource
	Content   string
	CreatedAt time.Time // 内置默认时间
}

// 记忆上下文结构体
type ContextMemory struct {
	ID            int64
	Summary       string    // 用于管理记忆上下文，及智能体所处的环境,总结,对内容理解提供一个大致的方向性
	LastSummaryID int64     // 最后一次总结的id
	UpdatedAt     time.Time // 最近修改时间
	CreatedAt     time.Time // 内置默认时间
}

// 短期记忆结构体
type ShortMemory struct {
	Memorys []OriginalMemory
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

type MemoryEvent struct {
	ID    string            `json:"id"`
	Text  string            `json:"text"`
	Meta  map[string]string `json:"meta"`
	Event string            `json:"event"`
}
