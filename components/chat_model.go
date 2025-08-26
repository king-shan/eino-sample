package components

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/joho/godotenv"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/schema"
)

func NewArkModel(ctx context.Context) (model *ark.ChatModel, err error) {
	// 初始化ark模型
	model, err = ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey: os.Getenv("ARK_API_KEY"),
		Model:  os.Getenv("MODEL"),
	})
	return
}

// 普通模式
func Generate() {
	err := godotenv.Load("example.env")
	if err != nil {
		fmt.Println("godotenv.Load", err.Error())
		return
	}
	// 初始化上下文
	ctx := context.Background()
	// 初始化ark模型
	model, err := NewArkModel(ctx)
	if err != nil {
		fmt.Println("NewArkModel", err.Error())
		return
	}
	input := []*schema.Message{
		// 系统消息
		schema.SystemMessage("你是一个助手"),
		schema.UserMessage("介绍一下 Eino"),
	}
	// 生成完整的模型响应--普通模式
	response, err := model.Generate(ctx, input)
	if err != nil {
		panic(err)
	}
	// 响应处理
	fmt.Println(response.Content)
	// 获取 Token 使用情况
	if usage := response.ResponseMeta.Usage; usage != nil {
		fmt.Println("提示 Tokens:", usage.PromptTokens)
		fmt.Println("生成 Tokens:", usage.CompletionTokens)
		fmt.Println("总 Tokens:", usage.TotalTokens)
	}
}

// 流式:改善了交互延迟和用户等待时的不确定感
func Stream() {
	err := godotenv.Load("example.env")
	if err != nil {
		panic(err)
	}
	// 初始化上下文
	ctx := context.Background()
	// 初始化ark模型
	model, err := NewArkModel(ctx)
	if err != nil {
		fmt.Println("NewArkModel", err.Error())
		return
	}

	input := []*schema.Message{
		// 系统消息
		schema.SystemMessage("你是一个助手"),
		schema.UserMessage("介绍一下 Eino"),
	}
	// 获取流式回复
	reader, err := model.Stream(ctx, input)
	if err != nil {
		panic(err)
	}
	defer reader.Close()
	// 处理流式内容
	for {
		chunk, err := reader.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(chunk.Content)
	}
}
