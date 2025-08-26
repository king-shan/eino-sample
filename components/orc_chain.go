package components

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"

	callbackHelpers "github.com/cloudwego/eino/utils/callbacks"
	"github.com/joho/godotenv"
)

// 链式编排
func CreateChain() {
	err := godotenv.Load("example.env")
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	// 初始化ark模型
	model, err := NewArkModel(ctx)
	if err != nil {
		fmt.Println("NewArkModel", err.Error())
		return
	}
	//编写lambda节点
	lambda := compose.InvokableLambda(func(ctx context.Context, input string) (output []*schema.Message, err error) {
		desuwa := input + "回答结尾加上 世界和平"
		output = []*schema.Message{
			{
				Role:    schema.User,
				Content: desuwa,
			},
		}
		return output, nil
	})
	// 注册链条
	chain := compose.NewChain[string, *schema.Message]()
	// 连接节点
	chain.AppendLambda(lambda).AppendChatModel(model)
	r, err := chain.Compile(ctx)
	if err != nil {
		panic(err)
	}
	// 编译运行
	answer, err := r.Invoke(ctx, "你好，请告诉我你的名字")
	if err != nil {
		panic(err)
	}
	fmt.Println(answer.Content)
}

// agent简单例子
func SimpleAgent() {

	getGameTool := CreateGameTool()
	ctx := context.Background()
	// 加载配置文件
	err := godotenv.Load("example.env")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//大模型回调函数
	modelHandler := &callbackHelpers.ModelCallbackHandler{
		OnEnd: func(ctx context.Context, info *callbacks.RunInfo, output *model.CallbackOutput) context.Context {
			// 1. output.Result 类型是 string
			fmt.Println("模型思考过程为：")
			fmt.Println(output.Message.Content)
			return ctx
		},
	}
	//工具回调函数
	toolHandler := &callbackHelpers.ToolCallbackHandler{
		OnStart: func(ctx context.Context, info *callbacks.RunInfo, input *tool.CallbackInput) context.Context {
			fmt.Printf("开始执行工具，参数: %s\n", input.ArgumentsInJSON)
			return ctx
		},
		OnEnd: func(ctx context.Context, info *callbacks.RunInfo, output *tool.CallbackOutput) context.Context {
			fmt.Printf("工具执行完成，结果: %s\n", output.Response)
			return ctx
		},
	}
	//构建实际回调函数Handler
	handler := callbackHelpers.NewHandlerHelper().
		ChatModel(modelHandler).
		Tool(toolHandler).
		Handler()
	// 初始化模型
	timeout := 30 * time.Second
	model, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey:  os.Getenv("ARK_API_KEY"),
		Model:   "doubao-1.5-pro-32k-250115",
		Timeout: &timeout,
	})
	if err != nil {
		panic(err)
	}
	//绑定工具
	info, err := getGameTool.Info(ctx)
	if err != nil {
		panic(err)
	}
	infos := []*schema.ToolInfo{
		info,
	}
	err = model.BindTools(infos)
	if err != nil {
		panic(err)
	}
	//创建tools节点
	ToolsNode, err := compose.NewToolNode(context.Background(), &compose.ToolsNodeConfig{
		Tools: []tool.BaseTool{
			getGameTool,
		},
	})
	if err != nil {
		panic(err)
	}
	//创建完整的处理链
	chain := compose.NewChain[[]*schema.Message, []*schema.Message]()
	chain.
		AppendChatModel(model, compose.WithNodeName("chat_model")).
		AppendToolsNode(ToolsNode, compose.WithNodeName("tools"))

	// 编译并运行 chain
	agent, err := chain.Compile(ctx)
	if err != nil {
		log.Fatal(err)
	}
	//运行Agent
	resp, err := agent.Invoke(ctx, []*schema.Message{
		{
			Role:    schema.User,
			Content: "请告诉我原神的URL是什么",
		},
	}, compose.WithCallbacks(handler))
	if err != nil {
		log.Fatal(err)
	}

	// 输出结果
	for _, msg := range resp {
		fmt.Println(msg.Content)
	}
}
