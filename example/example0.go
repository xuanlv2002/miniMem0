package example

import (
	"fmt"
	"miniMem0/config"
	"miniMem0/memory"
)

func Example0() {
	conf, err := config.LoadConfig("config/deafault.yaml")
	if err != nil {
		fmt.Println("Error loading configuration:", err)
		return
	}
	fmt.Println(conf)

	mem, err := memory.NewMemorySystem(conf)
	if err != nil {
		fmt.Println("Error initializing memory system:", err)
		return
	}
	prompt, err := mem.ProcessInput("我朋友是一个程序员，他的名字叫小明。")
	fmt.Println("Prompt:", prompt)
	if err != nil {
		fmt.Println("Error processing input:", err)
		return
	}
	err = mem.ProcessOutput("哦,小明是一个很棒的程序员！你能告诉我更多关于他的信息吗？")

	prompt, err = mem.ProcessInput("我也是一个程序员，我的名字叫小柴")
	fmt.Println("Prompt:", prompt)
	if err != nil {
		fmt.Println("Error processing input:", err)
		return
	}
	err = mem.ProcessOutput("哦,小柴也是一个很棒的程序员！你们两个都是程序员，真不错！")

	prompt, err = mem.ProcessInput("我感觉写代码很有趣，尤其是解决问题的时候。")
	fmt.Println("Prompt:", prompt)
	if err != nil {
		fmt.Println("Error processing input:", err)
		return
	}
	err = mem.ProcessOutput("是的，编程确实很有趣！解决问题的过程可以非常有成就感。")

	prompt, err = mem.ProcessInput("你喜欢编程吗？")
	fmt.Println("Prompt:", prompt)
	if err != nil {
		fmt.Println("Error processing input:", err)
		return
	}
	err = mem.ProcessOutput("我喜欢编程！它让我能够创造出有用的工具和应用程序。")

	prompt, err = mem.ProcessInput("我叫什么?")
	fmt.Println("Prompt:", prompt)
	if err != nil {
		fmt.Println("Error processing input:", err)
		return
	}
	err = mem.ProcessOutput("你叫小柴。你是一个程序员。")

	prompt, err = mem.ProcessInput("小明是谁?")
	fmt.Println("Prompt:", prompt)
	if err != nil {
		fmt.Println("Error processing input:", err)
		return
	}
	err = mem.ProcessOutput("小明是你的朋友，他也是一个程序员。")
}
