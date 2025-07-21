package llm

import (
	"context"
	"errors"
	"io"
	"miniMem0/config"
	"strings"
	"sync"

	"github.com/sashabaranov/go-openai"
)

/*
负责和底层大模型的交互，主要封装了 Chat 和 ChatAsync 两个方法
*/
func NewLLM(cfg *config.LLMConfig) *LLM {
	var openapiConfig openai.ClientConfig
	openapiConfig = openai.DefaultConfig(cfg.APIKey)
	openapiConfig.BaseURL = cfg.BaseURL

	return &LLM{
		lock:   sync.Mutex{},
		Config: cfg,
		Client: openai.NewClientWithConfig(openapiConfig),
	}
}

type LLM struct {
	lock   sync.Mutex
	Config *config.LLMConfig
	Client *openai.Client
}

func (l *LLM) Chat(ctx context.Context, messages []openai.ChatCompletionMessage) (*openai.ChatCompletionMessage, error) {
	l.lock.Lock()
	defer l.lock.Unlock()
	req := openai.ChatCompletionRequest{
		Model:       l.Config.Model,
		Temperature: float32(l.Config.Temperature),
		Stream:      false,
		Messages:    messages,
	}
	resp, err := l.Client.CreateChatCompletion(ctx, req)
	if err != nil {
		return nil, err
	}
	if len(resp.Choices) > 0 {
		return &resp.Choices[0].Message, nil
	}
	return nil, errors.New("no choices found")
}

func (l *LLM) ChatAsync(
	ctx context.Context,
	messages []openai.ChatCompletionMessage,
	caller ...func(body string)) (string, error) {
	l.lock.Lock()
	defer l.lock.Unlock()
	req := openai.ChatCompletionRequest{
		Model:       l.Config.Model,
		Temperature: l.Config.Temperature,
		Stream:      true,
		Messages:    messages,
	}
	resp, err := l.Client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return "", err
	}
	defer resp.Close()
	var builder strings.Builder
	for {
		data, err := resp.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return builder.String(), nil
			}
			return builder.String(), err
		}
		if len(data.Choices) > 0 {
			builder.WriteString(data.Choices[0].Delta.Content)
			if len(caller) > 0 {
				caller[0](data.Choices[0].Delta.Content)
			}
		}
	}
}

func (l *LLM) ChatWithTool(
	ctx context.Context,
	messages []openai.ChatCompletionMessage,
	tools []openai.Tool) (*openai.ChatCompletionMessage, error) {
	l.lock.Lock()
	defer l.lock.Unlock()
	req := openai.ChatCompletionRequest{
		Model:       l.Config.Model,
		Temperature: l.Config.Temperature,
		Stream:      true,
		Messages:    messages,
		Tools:       tools,
	}
	resp, err := l.Client.CreateChatCompletion(ctx, req)
	if err != nil {
		return nil, err
	}
	if len(resp.Choices) > 0 {
		return &resp.Choices[0].Message, nil
	}
	return nil, errors.New("no choices found")
}

func (l *LLM) ChatAsyncWithTool(
	ctx context.Context,
	messages []openai.ChatCompletionMessage,
	tools []openai.Tool,
	contentCaller func(body string),
	thinkCaller func(body string),
) (*openai.ChatCompletionMessage, error) {
	l.lock.Lock()
	defer l.lock.Unlock()
	req := openai.ChatCompletionRequest{
		Model:       l.Config.Model,
		Temperature: l.Config.Temperature,
		Stream:      true,
		Messages:    messages,
		Tools:       tools,
	}
	resp, err := l.Client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Close()
	var (
		ret            openai.ChatCompletionMessage
		contentBuilder strings.Builder
		currentToolID  string
		toolCalls      []openai.ToolCall
		toolCallMap    = make(map[string]*strings.Builder) // 用于跟踪正在构建的tool call
	)
	for {
		data, err := resp.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return &ret, err
		}
		if len(data.Choices) == 0 {
			continue
		}
		choice := data.Choices[0]
		// 处理思考流
		if choice.Delta.ReasoningContent != "" {
			if thinkCaller != nil {
				thinkCaller(choice.Delta.ReasoningContent)
			}
			// 思考内容不记录
		}
		// 处理内容流
		if choice.Delta.Content != "" {
			if contentCaller != nil {
				contentCaller(choice.Delta.Content)
			}
			contentBuilder.WriteString(choice.Delta.Content)
		}
		// 处理工具调用流
		if len(choice.Delta.ToolCalls) > 0 {
			deltaToolCall := choice.Delta.ToolCalls[0]
			// 新的工具调用开始
			if deltaToolCall.ID != "" && deltaToolCall.ID != currentToolID {
				currentToolID = deltaToolCall.ID
				toolCalls = append(toolCalls, deltaToolCall)
				toolCallMap[currentToolID] = &strings.Builder{}
			}
			// 累积参数
			if currentToolID != "" {
				if tc, ok := toolCallMap[currentToolID]; ok {
					tc.WriteString(deltaToolCall.Function.Arguments)
				}
			}
		}
	}
	ret.Content = contentBuilder.String()
	// 最后收集所有完成的工具调用
	for i, tc := range toolCalls {
		builder, ok := toolCallMap[tc.ID]
		if !ok {
			continue
		}
		toolCalls[i].Function.Arguments = strings.TrimSpace(builder.String())
	}
	if len(toolCalls) > 0 {
		ret.ToolCalls = toolCalls
	}
	return &ret, nil
}
