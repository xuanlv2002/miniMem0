package main

import (
	"context"
	"fmt"
	"miniMem0/config"
	"miniMem0/memory"
)

func main() {
	fmt.Println("欢迎使用minimem0-go")
	conf, err := config.LoadConfig("./config/deafault.yaml")
	if err != nil {
		fmt.Println("加载配置文件失败：", err)
		return
	}
	fmt.Println(conf)
	// 初始化记忆系统
	memorySystem, err := memory.NewMemorySystem(conf)
	if err != nil {
		fmt.Println("初始化记忆系统失败：", err)
		return
	}
	fmt.Println(memorySystem)
	fmt.Println("初始化记忆系统成功")

	// err = memorySystem.WriteMemory(context.Background(), "你好，我是小明，我今年18岁，我来自中国，我是一名学生")
	// if err != nil {
	// 	fmt.Println("写入记忆失败：", err)
	// 	return
	// }
	// fmt.Println("写入记忆成功")
	mem, err := memorySystem.SearchMemory(context.Background(), "我的名称是什么")
	if err != nil {
		fmt.Println("搜索记忆失败：", err)
		return
	}
	fmt.Println(mem)
}
