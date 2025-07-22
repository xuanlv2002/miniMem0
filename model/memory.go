package model

import "time"

type MemorySource string

const (
	UserRole MemorySource = "user"
	AIRole   MemorySource = "assistant"
	ToolRole MemorySource = "tool"
)

type OriginalMemory struct {
	ID        int64
	Role      MemorySource
	Content   string
	CreatedAt time.Time // 内置默认时间
}
