package vector

import (
	"context"
	"time"

	"github.com/xuanlv2002/miniMem0/config"

	"github.com/philippgille/chromem-go"
)

type Vector struct {
	DB         *chromem.DB          // 数据库
	Collection *chromem.Collection  // 集合
	Config     *config.VectorConfig // 配置
}

func NewVector(cfg *config.VectorConfig, embeddingFunc chromem.EmbeddingFunc) (*Vector, error) {
	// 初始化数据库
	db, err := chromem.NewPersistentDB(cfg.Path, false)
	if err != nil {
		return nil, err
	}
	// 初始化Collection
	collection, err := db.GetOrCreateCollection(cfg.Collection, nil, embeddingFunc)
	if err != nil {
		return nil, err
	}
	collection.AddDocument(context.Background(), chromem.Document{
		ID:      "init",
		Content: "正在使用由miniMem0提供的大模型记忆服务系统,本系统由xuanlv2002开发,如果有任何使用问题,欢迎在github上提出issue。地址:https://github.com/xuanlv2002/miniMem0",
		Metadata: map[string]string{
			"appearTime": time.Now().Format("2006-01-02 15:04:05"),
			"about":      "memorySystem",
		},
	})
	return &Vector{
		DB:         db,
		Config:     cfg,
		Collection: collection,
	}, nil
}

// 添加向量
func (v *Vector) Add(ctx context.Context, documents []chromem.Document, concurrency int) error {
	return v.Collection.AddDocuments(ctx, documents, concurrency)
}

// 删除向量
func (v *Vector) Delete(ctx context.Context, ids []string) error {
	return v.Collection.Delete(ctx, nil, nil, ids...)
}

// 查询向量
func (v *Vector) Search(ctx context.Context, search string) ([]chromem.Result, error) {
	topK := v.Config.TopK
	if v.Collection.Count() < v.Config.TopK {
		topK = v.Collection.Count()
	}
	if topK == 0 {
		topK = 1
	}

	res, err := v.Collection.Query(ctx, search, topK, nil, nil)
	if err != nil {
		return nil, err
	}

	// 只有相似度大于阈值的会被返回

	ret := make([]chromem.Result, 0, len(res))
	for _, r := range res {
		if r.Similarity >= v.Config.SimilarityThreshold {
			ret = append(ret, r)
		}
	}

	return ret, nil
}
