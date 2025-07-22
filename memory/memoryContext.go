package memory

import (
	"miniMem0/db/sqldb"
	"miniMem0/model"
)

// 用于管理记忆上下文，及智能体所处的环境,总结,对内容理解提供一个大致的方向性
type MemoryContext struct {
	Summary           string
	MemoryContextHand sqldb.SqlDB
}

func NewMemoryContext() *MemoryContext {
	return &MemoryContext{
		Summary: "",
	}
}

// 添加一条上下文, 并记录时间
func (m *MemoryContext) AddMemoryContext(from model.MemorySource, content string) {

}
