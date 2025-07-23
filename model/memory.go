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
	ID int64
	// 用于管理记忆上下文，及智能体所处的环境,总结,对内容理解提供一个大致的方向性
	Summary   string
	UpdatedAt time.Time // 最近修改时间
	CreatedAt time.Time // 内置默认时间
}
