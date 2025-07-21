package main

import (
	"context"
	"fmt"
	"miniMem0/config"
	"miniMem0/llm"
	"miniMem0/vector"

	"github.com/philippgille/chromem-go"
)

func main() {
	fmt.Println("欢迎使用minimem0-go")
	conf, err := config.LoadConfig("./config/deafault.yaml")
	if err != nil {
		fmt.Println("加载配置文件失败：", err)
		return
	}
	fmt.Println("conf:", conf)

	// 初始化LLM
	// l := llm.NewLLM(conf.ChatConfig)
	// resp, err := l.Chat(context.Background(), []openai.ChatCompletionMessage{
	// 	{
	// 		Role:    openai.ChatMessageRoleUser,
	// 		Content: "你好",
	// 	},
	// })

	// if err != nil {
	// 	fmt.Println("调用LLM失败：", err)
	// 	return
	// }
	// fmt.Println("LLM响应：", resp.Content)

	// // 初始化Embedding
	l2 := llm.NewEmbedding(conf.EmbeddingConfig)
	// embed, err := l2.Embedding(context.Background(), []string{"你好"})
	// if err != nil {
	// 	fmt.Println("调用LLM失败：", err)
	// 	return
	// }
	// fmt.Println("LLM响应：", len(embed.Embedding))

	// 初始化向量数据库
	v, err := vector.NewVector("./data", "knowledge-base", l2)

	if err != nil {
		fmt.Println("初始化向量数据库失败：", err)
		return
	}
	err = v.Add(context.Background(), []chromem.Document{
		{
			ID:      "5",
			Content: "你啊啊啊",
		},
	}, 1)
	if err != nil {
		fmt.Println("添加向量失败：", err)
		return
	}
	c := v.Collection.Count()
	fmt.Println("count:", c)
	// 这里不可大于存储的文档数量

	ret, _ := v.Search(context.Background(), "你好111121312312", 3)
	fmt.Println("搜索结果：", len(ret))
	// 初始化记忆系统

	// 启动web服务器
}
