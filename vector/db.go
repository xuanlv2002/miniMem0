package vector

import (
	"context"
	"miniMem0/llm"

	"github.com/philippgille/chromem-go"
)

type Vector struct {
	DB         *chromem.DB
	Collection *chromem.Collection
	Embedding  *llm.Embedding
}

func NewVector(path string, collectionName string, embed *llm.Embedding) (*Vector, error) {
	db, err := chromem.NewPersistentDB(path, false)
	if err != nil {
		return nil, err
	}
	embeddingFunc := func(ctx context.Context, text string) ([]float32, error) {
		ebed, err := embed.Embedding(ctx, text)
		if err != nil {
			return nil, err
		}
		return ebed.Embedding, nil
	}

	collection, err := db.GetOrCreateCollection(collectionName, nil, embeddingFunc)

	return &Vector{
		DB:         db,
		Collection: collection,
		Embedding:  embed,
	}, nil
}

func (v *Vector) Add(ctx context.Context, documents []chromem.Document, concurrency int) error {
	return v.Collection.AddDocuments(ctx, documents, concurrency)
}

func (v *Vector) Search(ctx context.Context, search string, topK int) ([]chromem.Result, error) {
	ret, err := v.Collection.Query(ctx, search, topK, nil, nil)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// func main() {
// 	ctx := context.Background()
// 	// 新建数据库
// 	db := chromem.NewDB()
// 	// 创建集合
// 	c, err := db.CreateCollection("knowledge-base", nil, func(ctx context.Context, text string) ([]float32, error) {
// 		return nil, nil
// 	})
// 	if err != nil {
// 		panic(err)
// 	}

// 	// 添加文档
// 	err = c.AddDocuments(ctx, []chromem.Document{
// 		{
// 			ID:      "1",
// 			Content: "The sky is blue because of Rayleigh scattering.",
// 		},
// 		{
// 			ID:      "2",
// 			Content: "Leaves are green because chlorophyll absorbs red and blue light.",
// 		},
// 	}, runtime.NumCPU())
// 	if err != nil {
// 		panic(err)
// 	}
// 	// 查询
// 	res, err := c.Query(ctx, "Why is the sky blue?", 1, nil, nil)
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Printf("ID: %v\nSimilarity: %v\nContent: %v\n", res[0].ID, res[0].Similarity, res[0].Content)

// }
