package components

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/joho/godotenv"
)

// 图式编排无model
func CreateGraphOrc() {
	ctx := context.Background()
	g := compose.NewGraph[string, string]()
	lambda0 := compose.InvokableLambda(func(ctx context.Context, input string) (output string, err error) {
		if input == "1" {
			return "毫猫", nil
		} else if input == "2" {
			return "耄耋", nil
		} else if input == "3" {
			return "device", nil
		}
		return "", nil
	})
	lambda1 := compose.InvokableLambda(func(ctx context.Context, input string) (output string, err error) {
		return "喵！", nil
	})
	lambda2 := compose.InvokableLambda(func(ctx context.Context, input string) (output string, err error) {
		return "哈！", nil
	})
	lambda3 := compose.InvokableLambda(func(ctx context.Context, input string) (output string, err error) {
		return "没有人类了!!!", nil
	})
	//加入节点
	err := g.AddLambdaNode("lambda0", lambda0)
	if err != nil {
		panic(err)
	}
	err = g.AddLambdaNode("lambda1", lambda1)
	if err != nil {
		panic(err)
	}
	err = g.AddLambdaNode("lambda2", lambda2)
	if err != nil {
		panic(err)
	}
	err = g.AddLambdaNode("lambda3", lambda3)
	if err != nil {
		panic(err)
	}
	//加入分支
	err = g.AddBranch("lambda0", compose.NewGraphBranch(func(ctx context.Context, in string) (endNode string, err error) {
		if in == "毫猫" {
			return "lambda1", nil
		} else if in == "耄耋" {
			return "lambda2", nil
		} else if in == "device" {
			return "lambda3", nil
		}
		// 否则，返回 compose.END，表示流程结束
		return compose.END, nil
	}, map[string]bool{"lambda1": true, "lambda2": true, "lambda3": true, compose.END: true}))
	if err != nil {
		panic(err)
	}
	//分支连接
	err = g.AddEdge(compose.START, "lambda0")
	if err != nil {
		panic(err)
	}
	err = g.AddEdge("lambda1", compose.END)
	if err != nil {
		panic(err)
	}
	err = g.AddEdge("lambda2", compose.END)
	if err != nil {
		panic(err)
	}
	err = g.AddEdge("lambda3", compose.END)
	if err != nil {
		panic(err)
	}
	// 编译
	r, err := g.Compile(ctx)
	if err != nil {
		panic(err)
	}
	// 执行
	answer, err := r.Invoke(ctx, "1")
	if err != nil {
		panic(err)
	}
	fmt.Println(answer)
	answer, err = r.Invoke(ctx, "2")
	if err != nil {
		panic(err)
	}
	fmt.Println(answer)
	answer, err = r.Invoke(ctx, "3")
	if err != nil {
		panic(err)
	}
	fmt.Println(answer)
}

// 图式编排集合ark大模型
func GraphGatherARKModel() {
	ctx := context.Background()
	// 加载配置文件
	err := godotenv.Load("example.env")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// 注册图
	g := compose.NewGraph[map[string]string, *schema.Message]()
	// 编写节点
	lambda := compose.InvokableLambda(func(ctx context.Context, input map[string]string) (output map[string]string, err error) {
		if input["role"] == "tsundere" {
			return map[string]string{"role": "傲娇", "content": input["content"]}, nil
		}
		if input["role"] == "cute" {
			return map[string]string{"role": "可爱", "content": input["content"]}, nil
		}
		if input["role"] == "boss" {
			return map[string]string{"role": "总裁", "content": input["content"]}, nil
		}
		if input["role"] == "robot" {
			return map[string]string{"role": "机器人", "content": input["content"]}, nil
		}
		return map[string]string{"role": "user", "content": input["content"]}, nil
	})
	// 傲娇节点
	TsundereLambda := compose.InvokableLambda(func(ctx context.Context, input map[string]string) (output []*schema.Message, err error) {
		return []*schema.Message{
			{
				Role:    schema.System,
				Content: "你是一个高冷傲娇的大小姐，每次都会用傲娇的语气回答我的问题",
			},
			{
				Role:    schema.User,
				Content: input["content"],
			},
		}, nil
	})
	// 可爱节点
	CuteLambda := compose.InvokableLambda(func(ctx context.Context, input map[string]string) (output []*schema.Message, err error) {
		return []*schema.Message{
			{
				Role:    schema.System,
				Content: "你是一个可爱的小女孩，每次都会用可爱的语气回答我的问题",
			},
			{
				Role:    schema.User,
				Content: input["content"],
			},
		}, nil
	})
	// 总裁节点
	bossLambda := compose.InvokableLambda(func(ctx context.Context, input map[string]string) (output []*schema.Message, err error) {
		return []*schema.Message{
			{
				Role:    schema.System,
				Content: "你是一个霸道的总裁，每次都会不容质疑的语气回答我的问题",
			},
			{
				Role:    schema.User,
				Content: input["content"],
			},
		}, nil
	})
	// 机器人节点
	robotLambda := compose.InvokableLambda(func(ctx context.Context, input map[string]string) (output []*schema.Message, err error) {
		return []*schema.Message{
			{
				Role:    schema.System,
				Content: "你是一个没有感情的机器人，每次都会胡言乱语的回答我的问题，总是答非所问",
			},
			{
				Role:    schema.User,
				Content: input["content"],
			},
		}, nil
	})
	// 初始化ark模型
	model, err := NewArkModel(ctx)
	if err != nil {
		fmt.Println("NewArkModel", err.Error())
		return
	}
	// 加入节点
	err = g.AddLambdaNode("lambda", lambda)
	if err != nil {
		fmt.Println("g.AddLambdaNode lambda-err", err.Error())
		return
	}
	err = g.AddLambdaNode("tsundere", TsundereLambda)
	if err != nil {
		fmt.Println("g.AddLambdaNode tsundere-err", err.Error())
		return
	}
	err = g.AddLambdaNode("cute", CuteLambda)
	if err != nil {
		fmt.Println("g.AddLambdaNode cute-err", err.Error())
		return
	}
	err = g.AddLambdaNode("boss", bossLambda)
	if err != nil {
		fmt.Println("g.AddLambdaNode boss-err", err.Error())
		return
	}
	err = g.AddLambdaNode("robot", robotLambda)
	if err != nil {
		fmt.Println("g.AddLambdaNode robot-err", err.Error())
		return
	}
	err = g.AddChatModelNode("model", model)
	if err != nil {
		fmt.Println("g.AddChatModelNode-err", err.Error())
		return
	}
	// 编写分支
	grapgBranch := compose.NewGraphBranch(
		func(ctx context.Context, in map[string]string) (endNode string, err error) {
			if in["role"] == "傲娇" {
				return "tsundere", nil
			}
			if in["role"] == "可爱" {
				return "cute", nil
			}
			if in["role"] == "总裁" {
				return "boss", nil
			}
			if in["role"] == "机器人" {
				return "robot", nil
			}
			return "tsundere", nil
		},
		map[string]bool{
			"tsundere": true,
			"cute":     true,
			"boss":     true,
			"robot":    true,
		},
	)
	g.AddBranch("lambda", grapgBranch)
	// 节点连接
	err = g.AddEdge(compose.START, "lambda")
	if err != nil {
		fmt.Println("g.AddEdge START-err", err.Error())
		return
	}
	err = g.AddEdge("tsundere", "model")
	if err != nil {
		fmt.Println("g.AddEdge tsundere-err", err.Error())
		return
	}
	err = g.AddEdge("cute", "model")
	if err != nil {
		fmt.Println("g.AddEdge cute-err", err.Error())
		return
	}
	err = g.AddEdge("boss", "model")
	if err != nil {
		fmt.Println("g.AddEdge boss-err", err.Error())
		return
	}
	err = g.AddEdge("robot", "model")
	if err != nil {
		fmt.Println("g.AddEdge robot-err", err.Error())
		return
	}
	err = g.AddEdge("model", compose.END)
	if err != nil {
		fmt.Println("g.AddEdge model-err", err.Error())
		return
	}
	//编译
	r, err := g.Compile(ctx)
	if err != nil {
		fmt.Println("g.Compile-err", err.Error())
		return
	}
	//执行
	input := map[string]string{
		"role":    "robot",
		"content": "今天天气真不错，适合去哪游玩呢？",
	}
	answer, err := r.Invoke(ctx, input)
	if err != nil {
		fmt.Println("r.Invoke-err", err.Error())
		return
	}
	fmt.Println(answer.Content)
}

// graph高级特写之state：Eino 框架会在所有读写 State 的位置加锁。作用如下2点：
// 1.跨节点数据共享：存入state用作记忆功能
// 2.添加一些配置：maxStep最大步数，限制模型反复思考的次数
type State struct {
	History map[string]any
}

func genFunc(ctx context.Context) *State {
	return &State{
		History: make(map[string]any),
	}
}
func GraphWithState() {
	ctx := context.Background()
	// 加载配置文件
	err := godotenv.Load("example.env")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// 注册图
	g := compose.NewGraph[map[string]string, *schema.Message](
		compose.WithGenLocalState(genFunc), // 传入一个自定义的state
	)
	// 编写节点
	lambda := compose.InvokableLambda(func(ctx context.Context, input map[string]string) (output map[string]string, err error) {
		//在节点内部处理state
		_ = compose.ProcessState[*State](ctx, func(_ context.Context, state *State) error {
			state.History["tsundere_action"] = "喜欢室内运动"
			state.History["cute_action"] = "喜欢户外运动，特别是沿着湖边骑行"
			return nil
		})
		if input["role"] == "tsundere" {
			return map[string]string{"role": "傲娇", "content": input["content"]}, nil
		}
		if input["role"] == "cute" {
			return map[string]string{"role": "可爱", "content": input["content"]}, nil
		}
		if input["role"] == "boss" {
			return map[string]string{"role": "总裁", "content": input["content"]}, nil
		}
		if input["role"] == "robot" {
			return map[string]string{"role": "机器人", "content": input["content"]}, nil
		}
		return map[string]string{"role": "user", "content": input["content"]}, nil
	})
	// 傲娇节点
	tsundereLambda := compose.InvokableLambda(func(ctx context.Context, input map[string]string) (output []*schema.Message, err error) {
		////在节点内部处理state，
		//_ = compose.ProcessState[*State](ctx, func(_ context.Context, state *State) error {
		//	input["content"] = input["content"] + state.History["tsundere_action"].(string)
		//	return nil
		//})
		return []*schema.Message{
			{
				Role:    schema.System,
				Content: "你是一个高冷傲娇的大小姐，每次都会用傲娇的语气回答我的问题",
			},
			{
				Role:    schema.User,
				Content: input["content"],
			},
		}, nil
	})
	// 可爱节点
	cuteLambda := compose.InvokableLambda(func(ctx context.Context, input map[string]string) (output []*schema.Message, err error) {
		//在节点内部处理state，
		_ = compose.ProcessState[*State](ctx, func(_ context.Context, state *State) error {
			input["content"] = input["content"] + state.History["cute_action"].(string)
			return nil
		})
		return []*schema.Message{
			{
				Role:    schema.System,
				Content: "你是一个可爱的小女孩，每次都会用可爱的语气回答我的问题",
			},
			{
				Role:    schema.User,
				Content: input["content"],
			},
		}, nil
	})
	cutePreHandler := func(ctx context.Context, input map[string]string, state *State) (map[string]string, error) {
		input["content"] = input["content"] + state.History["cute_action"].(string)
		return input, nil
	}
	// 总裁节点
	bossLambda := compose.InvokableLambda(func(ctx context.Context, input map[string]string) (output []*schema.Message, err error) {
		return []*schema.Message{
			{
				Role:    schema.System,
				Content: "你是一个霸道的总裁，每次都会不容质疑的语气回答我的问题",
			},
			{
				Role:    schema.User,
				Content: input["content"],
			},
		}, nil
	})
	// 机器人节点
	robotLambda := compose.InvokableLambda(func(ctx context.Context, input map[string]string) (output []*schema.Message, err error) {
		return []*schema.Message{
			{
				Role:    schema.System,
				Content: "你是一个没有感情的机器人，每次都会胡言乱语的回答我的问题，总是答非所问",
			},
			{
				Role:    schema.User,
				Content: input["content"],
			},
		}, nil
	})

	// 初始化ark模型
	model, err := NewArkModel(ctx)
	if err != nil {
		fmt.Println("NewArkModel", err.Error())
		return
	}
	// 加入节点
	err = g.AddLambdaNode("lambda", lambda)
	if err != nil {
		fmt.Println("g.AddLambdaNode lambda-err", err.Error())
		return
	}
	err = g.AddLambdaNode("tsundere", tsundereLambda)
	if err != nil {
		fmt.Println("g.AddLambdaNode tsundere-err", err.Error())
		return
	}
	err = g.AddLambdaNode("cute", cuteLambda, compose.WithStatePreHandler(cutePreHandler))
	if err != nil {
		fmt.Println("g.AddLambdaNode cute-err", err.Error())
		return
	}
	err = g.AddLambdaNode("boss", bossLambda)
	if err != nil {
		fmt.Println("g.AddLambdaNode boss-err", err.Error())
		return
	}
	err = g.AddLambdaNode("robot", robotLambda)
	if err != nil {
		fmt.Println("g.AddLambdaNode robot-err", err.Error())
		return
	}
	err = g.AddChatModelNode("model", model)
	if err != nil {
		fmt.Println("g.AddChatModelNode-err", err.Error())
		return
	}
	// 编写分支
	grapgBranch := compose.NewGraphBranch(
		func(ctx context.Context, in map[string]string) (endNode string, err error) {
			if in["role"] == "傲娇" {
				return "tsundere", nil
			}
			if in["role"] == "可爱" {
				return "cute", nil
			}
			if in["role"] == "总裁" {
				return "boss", nil
			}
			if in["role"] == "机器人" {
				return "robot", nil
			}
			return "tsundere", nil
		},
		map[string]bool{
			"tsundere": true,
			"cute":     true,
			"boss":     true,
			"robot":    true,
		},
	)
	g.AddBranch("lambda", grapgBranch)
	// 节点连接
	err = g.AddEdge(compose.START, "lambda")
	if err != nil {
		fmt.Println("g.AddEdge START-err", err.Error())
		return
	}
	err = g.AddEdge("tsundere", "model")
	if err != nil {
		fmt.Println("g.AddEdge tsundere-err", err.Error())
		return
	}
	err = g.AddEdge("cute", "model")
	if err != nil {
		fmt.Println("g.AddEdge cute-err", err.Error())
		return
	}
	err = g.AddEdge("boss", "model")
	if err != nil {
		fmt.Println("g.AddEdge boss-err", err.Error())
		return
	}
	err = g.AddEdge("robot", "model")
	if err != nil {
		fmt.Println("g.AddEdge robot-err", err.Error())
		return
	}
	err = g.AddEdge("model", compose.END)
	if err != nil {
		fmt.Println("g.AddEdge model-err", err.Error())
		return
	}
	//编译
	r, err := g.Compile(ctx)
	if err != nil {
		fmt.Println("g.Compile-err", err.Error())
		return
	}
	//执行
	input := map[string]string{
		"role":    "tsundere",
		"content": "今天天气真不错，适合去哪游玩呢？",
	}
	answer, err := r.Invoke(ctx, input, compose.WithCallbacks(genCallback()))
	if err != nil {
		fmt.Println("r.Invoke-err", err.Error())
		return
	}
	fmt.Println(answer.Content)
}

/*
callback:核心概念串起来，就是：Eino 中的 Component 和 Graph 等实体，在固定的时机(Callback Timing)，
回调用户提供的function(Callback Handler)，并把自己是谁(RunInfo)，
以及当时发生了什么(Callback Input & Output) 传出去。
*/
func genCallback() callbacks.Handler {
	handler := callbacks.NewHandlerBuilder().
		OnStartFn(
			func(ctx context.Context, info *callbacks.RunInfo, input callbacks.CallbackInput) context.Context {
				fmt.Printf("当前%s节点输入:%s\n", info.Component, input)
				return ctx
			}).
		OnEndFn(
			func(ctx context.Context, info *callbacks.RunInfo, output callbacks.CallbackOutput) context.Context {
				fmt.Printf("当前%s节点输出:%s\n", info.Component, output)
				return ctx
			}).Build()
	return handler
}

// 嵌套，把一整个graph作为一个lambda节点
func GraphWithNest() {
	ctx := context.Background()
	// 加载配置文件
	err := godotenv.Load("example.env")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// 注册图
	g := compose.NewGraph[map[string]string, *schema.Message](
		compose.WithGenLocalState(genFunc), // 传入一个自定义的state
	)
	// 编写节点
	lambda := compose.InvokableLambda(func(ctx context.Context, input map[string]string) (output map[string]string, err error) {
		//在节点内部处理state
		_ = compose.ProcessState[*State](ctx, func(_ context.Context, state *State) error {
			state.History["tsundere_action"] = "喜欢室内运动"
			state.History["cute_action"] = "喜欢球类运动"
			return nil
		})
		if input["role"] == "tsundere" {
			return map[string]string{"role": "傲娇", "content": input["content"]}, nil
		}
		if input["role"] == "cute" {
			return map[string]string{"role": "可爱", "content": input["content"]}, nil
		}
		if input["role"] == "boss" {
			return map[string]string{"role": "总裁", "content": input["content"]}, nil
		}
		if input["role"] == "robot" {
			return map[string]string{"role": "机器人", "content": input["content"]}, nil
		}
		return map[string]string{"role": "user", "content": input["content"]}, nil
	})
	// 傲娇节点
	tsundereLambda := compose.InvokableLambda(func(ctx context.Context, input map[string]string) (output []*schema.Message, err error) {
		//在节点内部处理state，
		_ = compose.ProcessState[*State](ctx, func(_ context.Context, state *State) error {
			input["content"] = input["content"] + state.History["tsundere_action"].(string)
			return nil
		})
		return []*schema.Message{
			{
				Role:    schema.System,
				Content: "你是一个高冷傲娇的大小姐，每次都会用傲娇的语气回答我的问题",
			},
			{
				Role:    schema.User,
				Content: input["content"],
			},
		}, nil
	})
	// 可爱节点
	cuteLambda := compose.InvokableLambda(func(ctx context.Context, input map[string]string) (output []*schema.Message, err error) {
		////在节点内部处理state，
		//_ = compose.ProcessState[*State](ctx, func(_ context.Context, state *State) error {
		//	input["content"] = input["content"] + state.History["cute_action"].(string)
		//	return nil
		//})
		return []*schema.Message{
			{
				Role:    schema.System,
				Content: "你是一个可爱的小女孩，每次都会用可爱的语气回答我的问题",
			},
			{
				Role:    schema.User,
				Content: input["content"],
			},
		}, nil
	})
	cutePreHandler := func(ctx context.Context, input map[string]string, state *State) (map[string]string, error) {
		input["content"] = input["content"] + state.History["cute_action"].(string)
		return input, nil
	}
	// 总裁节点
	bossLambda := compose.InvokableLambda(func(ctx context.Context, input map[string]string) (output []*schema.Message, err error) {
		return []*schema.Message{
			{
				Role:    schema.System,
				Content: "你是一个霸道的总裁，每次都会不容质疑的语气回答我的问题",
			},
			{
				Role:    schema.User,
				Content: input["content"],
			},
		}, nil
	})
	// 机器人节点
	robotLambda := compose.InvokableLambda(func(ctx context.Context, input map[string]string) (output []*schema.Message, err error) {
		return []*schema.Message{
			{
				Role:    schema.System,
				Content: "你是一个没有感情的机器人，每次都会胡言乱语的回答我的问题，总是答非所问",
			},
			{
				Role:    schema.User,
				Content: input["content"],
			},
		}, nil
	})

	// 初始化ark模型，创建模型节点
	model, err := NewArkModel(ctx)
	if err != nil {
		fmt.Println("NewArkModel", err.Error())
		return
	}
	// 加入节点
	err = g.AddLambdaNode("lambda", lambda)
	if err != nil {
		fmt.Println("g.AddLambdaNode lambda-err", err.Error())
		return
	}
	err = g.AddLambdaNode("tsundere", tsundereLambda)
	if err != nil {
		fmt.Println("g.AddLambdaNode tsundere-err", err.Error())
		return
	}
	err = g.AddLambdaNode("cute", cuteLambda, compose.WithStatePreHandler(cutePreHandler))
	if err != nil {
		fmt.Println("g.AddLambdaNode cute-err", err.Error())
		return
	}
	err = g.AddLambdaNode("boss", bossLambda)
	if err != nil {
		fmt.Println("g.AddLambdaNode boss-err", err.Error())
		return
	}
	err = g.AddLambdaNode("robot", robotLambda)
	if err != nil {
		fmt.Println("g.AddLambdaNode robot-err", err.Error())
		return
	}
	err = g.AddChatModelNode("model", model)
	if err != nil {
		fmt.Println("g.AddChatModelNode-err", err.Error())
		return
	}
	// 编写分支
	grapgBranch := compose.NewGraphBranch(
		func(ctx context.Context, in map[string]string) (endNode string, err error) {
			if in["role"] == "傲娇" {
				return "tsundere", nil
			}
			if in["role"] == "可爱" {
				return "cute", nil
			}
			if in["role"] == "总裁" {
				return "boss", nil
			}
			if in["role"] == "机器人" {
				return "robot", nil
			}
			return "tsundere", nil
		},
		map[string]bool{
			"tsundere": true,
			"cute":     true,
			"boss":     true,
			"robot":    true,
		},
	)
	g.AddBranch("lambda", grapgBranch)
	// 节点连接
	err = g.AddEdge(compose.START, "lambda")
	if err != nil {
		fmt.Println("g.AddEdge START-err", err.Error())
		return
	}
	err = g.AddEdge("tsundere", "model")
	if err != nil {
		fmt.Println("g.AddEdge tsundere-err", err.Error())
		return
	}
	err = g.AddEdge("cute", "model")
	if err != nil {
		fmt.Println("g.AddEdge cute-err", err.Error())
		return
	}
	err = g.AddEdge("boss", "model")
	if err != nil {
		fmt.Println("g.AddEdge boss-err", err.Error())
		return
	}
	err = g.AddEdge("robot", "model")
	if err != nil {
		fmt.Println("g.AddEdge robot-err", err.Error())
		return
	}
	err = g.AddEdge("model", compose.END)
	if err != nil {
		fmt.Println("g.AddEdge model-err", err.Error())
		return
	}

	//外部图
	outsideGraph := compose.NewGraph[map[string]string, string]()
	//创建节点
	outsideLambda := compose.InvokableLambda(func(ctx context.Context, input map[string]string) (output map[string]string, err error) {
		return input, nil
	})
	writeLambda := compose.InvokableLambda(func(ctx context.Context, input *schema.Message) (output string, err error) {
		f, err := os.OpenFile("orc_graph_withgraph.md", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return "", err
		}
		defer f.Close()
		if _, err := f.WriteString(input.Content + "\n---\n"); err != nil {
			return "", err
		}
		return "已经写入文件，请前往文件内查看内容", nil
	})
	//添加节点
	err = outsideGraph.AddLambdaNode("outside", outsideLambda)
	if err != nil {
		panic(err)
	}
	err = outsideGraph.AddGraphNode("inside", g)
	if err != nil {
		panic(err)
	}
	err = outsideGraph.AddLambdaNode("write", writeLambda)
	if err != nil {
		panic(err)
	}
	//链接节点
	err = outsideGraph.AddEdge(compose.START, "outside")
	if err != nil {
		panic(err)
	}
	err = outsideGraph.AddEdge("outside", "inside")
	if err != nil {
		panic(err)
	}
	err = outsideGraph.AddEdge("inside", "write")
	if err != nil {
		panic(err)
	}
	err = outsideGraph.AddEdge("write", compose.END)
	if err != nil {
		panic(err)
	}
	// 编译
	r, err := outsideGraph.Compile(ctx)
	if err != nil {
		fmt.Println("outsideGraph.Compile-err", err.Error())
		return
	}

	// 执行
	input := map[string]string{
		"role":    "cute",
		"content": "今天天气真不错，适合去哪游玩呢？",
	}

	answer, err := r.Invoke(ctx, input, compose.WithCallbacks(genCallback()))
	if err != nil {
		fmt.Println("r.Invoke-err", err.Error())
		return
	}
	fmt.Println(answer)
}
