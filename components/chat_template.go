package components

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/joho/godotenv"

	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

// 单独使用
func UsedAloneTemplate() {
	err := godotenv.Load("example.env")
	if err != nil {
		panic(err)
	}
	// 初始化上下文
	ctx := context.Background()
	// 初始化ark模型
	model, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey: os.Getenv("ARK_API_KEY"),
		Model:  os.Getenv("MODEL"),
	})
	// 创建模板
	template := prompt.FromMessages(schema.FString,
		schema.SystemMessage("你是一个{role}"),
		&schema.Message{
			Role:    schema.User,
			Content: "请帮我解决{task}。",
		},
	)
	// 准备变量
	params := map[string]any{
		"role": "专业的助手",
		"task": "写一首诗",
	}
	// 格式化模板
	messages, err := template.Format(ctx, params)
	// 生成完整的模型响应--普通模式
	response, err := model.Generate(ctx, messages)
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
