package memory

import (
	"miniMem0/config"
	"miniMem0/db/sqldb"
	"miniMem0/db/vector"
	"miniMem0/llm"
)

/*
	用于管理长期记忆
	长期记忆,及一些持久性的记忆,如用户的信息,一些历史记录等
	长期记忆通过向量数据库存储,并根据用户的输入进行检索,返回结果
*/

type LongMemoryHandler struct {
	config *config.LongMemoryConfig
}

func NewLongMemory(config *config.LongMemoryConfig, vector *vector.Vector, sqlHandler *sqldb.SqlHandler, llmModel *llm.LLM) *LongMemoryHandler {
	return &LongMemoryHandler{}
}
