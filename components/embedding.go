package components

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"

	"github.com/cloudwego/eino-ext/components/embedding/ark"
)

func NewArkEmbedding(ctx context.Context) (embedder *ark.Embedder, err error) {
	// 初始化嵌入器
	timeout := 30 * time.Second
	embedder, err = ark.NewEmbedder(ctx, &ark.EmbeddingConfig{
		APIKey:  os.Getenv("ARK_API_KEY"),
		Model:   os.Getenv("EMBEDDER"),
		Timeout: &timeout,
	})
	return
}

func UseAloneEmbedding() {
	err := godotenv.Load("example.env")
	if err != nil {
		panic(err)
	}
	// 初始化上下文
	ctx := context.Background()
	// 初始化嵌入器
	embedder, err := NewArkEmbedding(ctx)
	if err != nil {
		fmt.Println("NewArkEmbedding-err", err.Error())
		return
	}
	// 生成文本向量
	texts := []string{
		"可爱的浙江有11个地级市",
		"杭州房价下跌到1万还需要2年",
		"你好",
	}
	embeddings, err := embedder.EmbedStrings(ctx, texts)
	if err != nil {
		fmt.Println("EmbedStrings", err.Error())
		return
	}
	// 使用生成的向量
	for i, embedding := range embeddings {
		println("文本", i+1, "的向量维度:", len(embedding))
	}
	fmt.Println(embeddings)
}
