package llm

import (
	"context"
	"errors"
	"miniMem0/config"
	"sync"

	"github.com/philippgille/chromem-go"
	"github.com/sashabaranov/go-openai"
)

/*
负责和底层大模型的交互，主要封装了 Chat 和 ChatAsync 两个方法
*/
func NewEmbedding(cfg *config.EmbeddingConfig) *Embedding {
	var openapiConfig openai.ClientConfig
	openapiConfig = openai.DefaultConfig(cfg.APIKey)
	openapiConfig.BaseURL = cfg.BaseURL

	return &Embedding{
		lock:   sync.Mutex{},
		Config: cfg,
		Client: openai.NewClientWithConfig(openapiConfig),
	}
}

type Embedding struct {
	lock   sync.Mutex
	Config *config.EmbeddingConfig
	Client *openai.Client
}

func (l *Embedding) Embedding(ctx context.Context, messages string) (*openai.Embedding, error) {
	l.lock.Lock()
	defer l.lock.Unlock()
	req := openai.EmbeddingRequest{
		Dimensions: l.Config.Dimensions,
		Model:      l.Config.Model,
		Input:      []string{messages},
	}
	resp, err := l.Client.CreateEmbeddings(ctx, req)
	if err != nil {
		return nil, err
	}

	if len(resp.Data) > 0 {
		return &resp.Data[0], nil
	}
	return nil, errors.New("no choices found")
}

func (l *Embedding) GetEmbeddingFunc() chromem.EmbeddingFunc {
	return func(ctx context.Context, text string) ([]float32, error) {
		embedding, err := l.Embedding(ctx, text)
		if err != nil {
			return nil, err
		}
		return embedding.Embedding, nil
	}
}
