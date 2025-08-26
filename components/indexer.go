package components

import (
	"context"
	"fmt"

	"github.com/joho/godotenv"

	"github.com/cloudwego/eino/schema"

	"github.com/cloudwego/eino-ext/components/embedding/ark"

	"github.com/cloudwego/eino-ext/components/indexer/milvus"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

func NewIndexerMilvus(ctx context.Context, embedder *ark.Embedder) (indexer *milvus.Indexer, err error) {
	var collection = "test"

	var fields = []*entity.Field{
		{
			Name:     "id",
			DataType: entity.FieldTypeVarChar,
			TypeParams: map[string]string{
				"max_length": "256",
			},
			PrimaryKey: true,
		},
		{
			Name:     "vector", // 确保字段名匹配
			DataType: entity.FieldTypeBinaryVector,
			//DataType: entity.FieldTypeFloatVector,
			TypeParams: map[string]string{
				"dim": "81920",
				//"dim": "2560",
			},
		},
		{
			Name:     "content",
			DataType: entity.FieldTypeVarChar,
			TypeParams: map[string]string{
				"max_length": "8192",
			},
		},
		{
			Name:     "metadata",
			DataType: entity.FieldTypeJSON,
		},
	}
	indexer, err = milvus.NewIndexer(ctx, &milvus.IndexerConfig{
		Client:     MilvusCli,
		Collection: collection,
		Fields:     fields,
		Embedding:  embedder,
		//DocumentConverter: floatDocumentConverter,
	})
	return
}

func UseAloneIndexer() {
	// 初始化上下文
	ctx := context.Background()
	err := godotenv.Load("example.env")
	if err != nil {
		panic(err)
	}
	// 初始化嵌入器
	embedder, err := NewArkEmbedding(ctx)
	if err != nil {
		fmt.Println("NewArkEmbedding-err", err.Error())
		return
	}
	// 初始化存储组件：indexer-milvus
	indexer, err := NewIndexerMilvus(ctx, embedder)
	if err != nil {
		fmt.Println("NewIndexerMilvus-err", err.Error())
		return
	}
	docs := []*schema.Document{
		{
			ID:      "3",
			Content: "功能：根据查询检索相关文档",
			MetaData: map[string]any{
				"author": "木乔",
			},
		},
		{
			ID:      "4",
			Content: "ctx：上下文对象，用于传递请求级别的信息，同时也用于传递 Callback Manager",
			MetaData: map[string]any{
				"author": "王山",
			},
		},
	}
	ids, err := indexer.Store(ctx, docs)
	if err != nil {
		fmt.Println("indexer.Store", err.Error())
		return
	}

	fmt.Println("ids:", ids)
}

func floatDocumentConverter(ctx context.Context, docs []*schema.Document, vectors [][]float64) ([]interface{}, error) {
	rows := make([]interface{}, 0, len(docs))
	for i, doc := range docs {
		// float64 -> float32
		float32Vec := make([]float32, len(vectors[i]))
		for j, v := range vectors[i] {
			float32Vec[j] = float32(v)
		}
		row := map[string]interface{}{
			"id":       doc.ID,
			"content":  doc.Content,
			"vector":   float32Vec,
			"metadata": doc.MetaData,
		}
		rows = append(rows, row)
	}
	return rows, nil
}
