package components

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/cloudwego/eino/schema"
	"github.com/joho/godotenv"
)

/*
构建RAG步骤：
1.初始化嵌入器
2.初始化分割器--准备要分割的文档--执行分割--处理分割结果
3.初始化存储组件：indexer-milvus--进行存储
4.创建一个检索组件 retriever--进行检索
*/
func ConstructRAG() {
	// 初始化上下文
	ctx := context.Background()
	// 加载配置文件
	err := godotenv.Load("example.env")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	InitClient()
	// 初始化嵌入器
	embedder, err := NewArkEmbedding(ctx)
	if err != nil {
		fmt.Println("NewArkEmbedding-err", err.Error())
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
		doc.ID = docs[0].ID + "_" + strconv.Itoa(i)
		fmt.Println(doc.ID)
	}

	// 初始化存储组件：indexer-milvus;进行存储
	indexer, err := NewIndexerMilvus(ctx, embedder)
	if err != nil {
		fmt.Println("NewIndexerMilvus-err", err.Error())
		return
	}
	// 进行存储
	ids, err := indexer.Store(ctx, docs)
	if err != nil {
		fmt.Println("indexer.Store", err.Error())
		return
	}
	fmt.Println("indexer.Store的id值", ids)

	// 创建一个检索组件 retriever
	retriever, err := NewRetrieverMilvus(ctx, embedder)
	if err != nil {
		fmt.Println("NewRetrieverMilvus-err", err.Error())
		return
	}
	// 检索组件返回值
	resultsRetrieve, err := retriever.Retrieve(ctx, "当解说称某警队总输时，弹幕调侃“我不喜欢这种警察一直输的比赛")
	if err != nil {
		fmt.Println("retriever.Retrieve", err.Error())
		return
	}
	for i, result := range resultsRetrieve {
		fmt.Println("检索组件返回值", i, result.ID, result.Content)
		fmt.Println("====================分割线=========================")
	}
}
