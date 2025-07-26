package example

import (
	"bufio"
	"context"
	"fmt"
	"miniMem0/config"
	"miniMem0/llm"
	"miniMem0/memory"
	"os"

	"github.com/sashabaranov/go-openai"
)

func Example1() {
	conf, err := config.LoadConfig("config/local.yaml")
	if err != nil {
		panic(err)
	}
	// Use the loaded configuration
	memSys, err := memory.NewMemorySystem(conf)
	if err != nil {
		panic(err)
	}
	// Use the memory system
	llmModel := llm.NewLLM(conf.GetChatConfig())

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("欢迎来到个性化智能服务,输入q退出:")

	for {
		fmt.Print("输入(按enter结束):")
		input, _ := reader.ReadString('\n')
		if input == "q\n" {
			fmt.Println("退出程序")
			return
		}
		prompt, err := memSys.ProcessInput(input)
		if err != nil {
			fmt.Println("处理输入时出错:", err)
			return
		}
		fmt.Println("处理后的提示词:", prompt)
		response, err := llmModel.Chat(context.Background(), []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "你的名字是个性化智能助手，你的任务是根据用户输入提供个性化的建议和回答。",
			}, {
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		})
		if err != nil {
			fmt.Println("获取响应时出错:", err)
			return
		}
		fmt.Println("大模型响应:", response.Content)
		memSys.ProcessOutput(response.Content)
	}
}
