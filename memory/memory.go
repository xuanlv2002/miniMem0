package memory

import "miniMem0/config"

// 原始记忆结构体 包含角色 内容 时间

type MemorySystem struct {
	LongMemorySystem  *LongMemory
	ShortMemorySystem *ShortMemroy
	MemoryContext     *MemoryContext
}

func NewMemorySystem(option *config.Config) (*MemorySystem, error) {
	return nil, nil
}

// 处理大模型输入内容
func ProcessInput(input string) string {
	return ""
}

// 处理大模型输出内容
func ProcessOutput(ouput string) error {
	return nil
}
