package components

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudwego/eino/components/document"

	"github.com/cloudwego/eino/schema"

	"github.com/cloudwego/eino-ext/components/document/transformer/splitter/markdown"
	"github.com/joho/godotenv"
)

// 初始化mk文件格式的分割器
func NewTransMarkdown(ctx context.Context) (splitter document.Transformer, err error) {
	splitter, err = markdown.NewHeaderSplitter(ctx, &markdown.HeaderConfig{
		Headers: map[string]string{
			"#":   "h1",
			"##":  "h2",
			"###": "h3",
		},
		TrimHeaders: false,
	})
	return
}

func UseAloneTransformer() {
	// 初始化上下文
	ctx := context.Background()
	// 加载配置文件
	err := godotenv.Load("example.env")
	if err != nil {
		fmt.Println("godotenv.Load-err", err.Error())
		return
	}
	// 初始化分割器
	splitter, err := NewTransMarkdown(ctx)
	if err != nil {
		fmt.Println("NewTransMarkdown-err", err.Error())
		return
	}
	// 准备要分割的文档
	content, err := os.OpenFile("./components/document.md", os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		fmt.Println("os.OpenFile-err", err.Error())
		return
	}
	defer content.Close()
	bs, err := os.ReadFile("./components/document.md")
	if err != nil {
		fmt.Println("os.ReadFile-err", err.Error())
		return
	}
	docs := []*schema.Document{
		{
			ID:      "doc1",
			Content: string(bs),
		},
	}
	// 执行分割
	results, err := splitter.Transform(ctx, docs)
	if err != nil {
		fmt.Println("splitter.Transform-err", err.Error())
		return
	}

	// 处理分割结果
	for i, doc := range results {
		println("片段", i+1, ":", doc.Content)
		println("标题层级：")
		for k, v := range doc.MetaData {
			if k == "h1" || k == "h2" || k == "h3" {
				fmt.Println("  ", k, ":", v)
			}
		}
	}
}
