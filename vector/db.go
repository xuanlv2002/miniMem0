package vector

import (
	"context"
	"miniMem0/config"

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
func (v *Vector) Search(ctx context.Context, search string, topK int, threshold float32) ([]chromem.Result, error) {
	if v.Collection.Count() < topK {
		topK = v.Collection.Count()
	}
	if topK > v.Config.MaxTopK {
		topK = v.Config.MaxTopK
	}
	res, err := v.Collection.Query(ctx, search, topK, nil, nil)
	if err != nil {
		return nil, err
	}

	// 只有相似度大于阈值的会被返回
	if threshold < v.Config.SimilarityThreshold {
		threshold = v.Config.SimilarityThreshold
	}
	ret := make([]chromem.Result, 0, len(res))
	for _, r := range res {
		if r.Similarity >= threshold {
			ret = append(ret, r)
		}
	}

	return ret, nil
}
