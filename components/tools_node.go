package components

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"

	"github.com/cloudwego/eino-ext/components/tool/browseruse"
)

/*
BrowserUse Tool:
Eino 的 BrowserUse Tool实现，实现Tool接口。这使得与 Eino 的 LLM 功能无缝集成，以增强自然语言处理和生成

BrowserUse 是一个开源的、AI 驱动的浏览器自动化工具。它的核心目标是让 AI 代理（AI Agent）能够像人类一样理解和操作网页浏览器，
你只需要用自然语言下达指令，它就能自动完成许多网页操作任务。
BrowserUse 的工作原理可以简单理解为：将你的自然语言指令“翻译”成浏览器能执行的一系列动作。

特点：BrowserUse 的优势在于其 “动口不动手”的自动化理念，大大降低了自动化任务的技术门槛，尤其适合非程序员或希望快速实现复杂流程自动化的用户
局限性：
稳定性与可靠性：AI 对复杂多变网页的理解和执行仍可能出错，技术仍在发展中。
设置门槛：本地部署需要一定的技术背景来配置 Python 环境（或golang环境）和 API 密钥。
成本：使用强大的商业 LLM 会产生 API 调用费用，复杂任务成本更高。
对特定网站无效：对于有人机验证（CAPTCHA）或反爬虫机制的网站，可能无法顺利工作
*/

func BrowserUseTool() {
	ctx := context.Background()
	// 初始化
	but, err := browseruse.NewBrowserUseTool(ctx, &browseruse.Config{})
	if err != nil {
		fmt.Println("browseruse.NewBrowserUseTool-err", err)
		return
	}

	// 执行
	url := "https://www.bilibili.com"
	result, err := but.Execute(&browseruse.Param{
		Action: browseruse.ActionGoToURL, // 浏览器处理
		URL:    &url,                     // 要跳转的URL
	})
	if err != nil {
		fmt.Println("but.Execute-err", err)
		return
	}
	fmt.Println(result)
	time.Sleep(10 * time.Second)
	but.Cleanup()
}

/*
创建一个tool：有4种：
1.直接实现接口
2.把本地函数转为 tool：使用 NewTool 方法；使用 InferTool 方法；使用 InferOptionableTool 方法
3.使用 eino-ext 中提供的 tool：上面函数（BrowserUseTool）就是属于这类
4.使用 MCP 协议
用的比较多的是第2种和第4种，以下是第2种的NewTool实现方式
*/
type Game struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type InputParams struct {
	Name string `json:"name" jsonschema:"description=the name of game"`
}

func GetGame(_ context.Context, params *InputParams) (string, error) {
	GameSet := []Game{
		{Name: "原神", Url: "https://ys.mihoyo.com/tool"},
		{Name: "鸣潮", Url: "https://mc.kurogames.com/tool"},
		{Name: "明日方舟", Url: "https://ak.hypergryph.com/tool"},
	}
	for _, game := range GameSet {
		if game.Name == params.Name {
			return game.Url, nil
		}
	}
	return "", nil
}

func CreateGameTool() tool.InvokableTool {
	// eino封装的基础工具包
	getGameTool := utils.NewTool(&schema.ToolInfo{
		Name: "get_game",
		Desc: "get a game url by name",
		ParamsOneOf: schema.NewParamsOneOfByParams(
			map[string]*schema.ParameterInfo{
				"name": &schema.ParameterInfo{
					Type:     schema.String,
					Desc:     "game's name",
					Required: true,
				},
			},
		),
	}, GetGame)
	return getGameTool
}
