package model

import "time"

type MemorySource string

const (
	UserRole MemorySource = "user"
	AIRole   MemorySource = "assistant"
	ToolRole MemorySource = "tool"
)

type OriginalMemory struct {
	Role       MemorySource
	Content    string
	CreateTime time.Time // 内置默认时间
}
