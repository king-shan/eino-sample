package components

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino-ext/components/embedding/ark"
	"github.com/cloudwego/eino-ext/components/retriever/milvus"
	"github.com/joho/godotenv"
)

func NewRetrieverMilvus(ctx context.Context, embedder *ark.Embedder) (retriever *milvus.Retriever, err error) {
	retriever, err = milvus.NewRetriever(ctx, &milvus.RetrieverConfig{
		Client:      MilvusCli,
		Collection:  "test",
		Partition:   nil,
		VectorField: "vector",
		OutputFields: []string{
			"id",
			"content",
			"metadata",
		},
		TopK:      5,
		Embedding: embedder,
	})
	return
}

func UseAloneRetriever() {
	// 初始化上下文
	ctx := context.Background()
	// 加载配置文件
	err := godotenv.Load("example.env")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// 初始化嵌入器
	embedder, err := NewArkEmbedding(ctx)
	if err != nil {
		fmt.Println("NewArkEmbedding-err", err.Error())
		return
	}
	// 创建一个检索组件 retriever
	retriever, err := NewRetrieverMilvus(ctx, embedder)
	if err != nil {
		fmt.Println("NewRetrieverMilvus-err", err.Error())
		return
	}
	// 检索组件返回值
	results, err := retriever.Retrieve(ctx, "原神")
	if err != nil {
		fmt.Println("retriever.Retrieve", err.Error())
		return
	}
	for i, result := range results {
		fmt.Println("检索组件返回值", i, result.ID, result.Content)
	}
}
