package memory

import (
	"context"
	"miniMem0/config"
	"miniMem0/db/sqldb"
	"miniMem0/db/vector"
	"miniMem0/llm"
	"miniMem0/model"
	"sync"

	"github.com/sirupsen/logrus"
)

/*
	用于管理长期记忆
	长期记忆,及一些持久性的记忆,如用户的信息,一些历史记录等
	长期记忆通过向量数据库存储,并根据用户的输入进行检索,返回结果
*/

type LongMemoryHandler struct {
	config     *config.LongMemoryConfig
	vector     *vector.Vector
	llmHandler *llm.LLM
	sqlHandler *sqldb.SqlHandler
	mu         sync.Mutex     // 长期记忆锁
	wg         sync.WaitGroup // 用来等待所有任务完成
}

func NewLongMemory(config *config.LongMemoryConfig, vector *vector.Vector, sqlHandler *sqldb.SqlHandler, llmModel *llm.LLM) *LongMemoryHandler {
	return &LongMemoryHandler{
		vector:     vector,
		llmHandler: llmModel,
		sqlHandler: sqlHandler,
		config:     config,
	}
}

func (l *LongMemoryHandler) WaitDone() {
	l.wg.Done()
}

// 获得长期记忆
func (l *LongMemoryHandler) GetLongMemory(text string) (*model.LongMemory, error) {
	var LongMemory model.LongMemory

	ret, err := l.vector.Search(context.Background(), text)
	if err != nil {
		return nil, err
	}

	var vectorMemory = make([]model.LongMemoryItem, 0)
	for _, v := range ret {
		vectorMemory = append(vectorMemory, model.LongMemoryItem{
			Text:     v.Content,
			Meta:     v.Metadata,
			Similary: v.Similarity,
		})
	}

	LongMemory.VectorMemorys = vectorMemory

	return &LongMemory, nil
}

// 更新长期记忆 异步更新 不对系统进行阻塞
func (l *LongMemoryHandler) UpdateLongMemory() {
	l.wg.Add(1)
	go func() {
		defer l.wg.Done()
		err := l.SaveLongMemory()
		if err != nil {
			logrus.Errorf("SummaryMemoryContext error: %v", err)
		}
	}()
}

// 更新长期记忆 异步更新 不对系统进行阻塞
func (l *LongMemoryHandler) SaveLongMemory() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	return nil
}
