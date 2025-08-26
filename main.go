package main

import "eino-sample/components"

// ChatModel 提供了大模型的对话能力
// Embedding 提供了基于语义的文本向量化能力
// Retriever 提供了关联内容召回的能力
// ToolsNode 提供了执行外部工具的能力
// 以上下游 类型对齐 为基本准则：前一个运行节点的输出值，可以作为下一个节点的输入值
func main() {
	// 生成完整的模型响应--普通模式
	//components.Generate()
	// 流式模式
	//components.Stream()

	// Template:Prompt 组件是一个用于处理和格式化提示模板的组件。它的主要作用是将用户提供的变量值填充到预定义的消息模板中，生成用于与语言模型交互的标准消息格式
	//components.UsedAloneTemplate()

	// Embedding:用于将文本转换为向量表示的组件.将文本内容映射到向量空间，使得语义相似的文本在向量空间中的距离较近
	/*
		文本相似度计算
		语义搜索
		文本聚类分析
	*/
	//components.UseAloneEmbedding()

	// Indexer:是一个用于存储和索引文档的组件，将文档及其向量表示存储到后端存储系统中，并提供高效的检索能力
	// 这个组件在以下场景中发挥重要作用：构建向量数据库，以用于语义关联搜索
	//components.InitClient()
	//components.UseAloneIndexer()

	// Retriever组件:把 Indexer 构建索引之后的内容进行召回，在 AI 应用中，一般使用 Embedding 进行语义相似性召回
	// Retriever 组件是一个用于从各种数据源检索文档的组件。它的主要作用是根据用户的查询（query）从文档库中检索出最相关的文档。这个组件在以下场景中特别有用：
	/*
		基于向量相似度的文档检索
		基于关键词的文档搜索
		知识库问答系统 (rag)
	*/
	//components.InitClient()
	//components.UseAloneRetriever()

	// Transformer组件
	//components.UseAloneTransformer()

	/*====================构建RAG--开始=========================*/
	//components.InitClient()
	//components.ConstructRAG()
	/*====================构建RAG--结束=========================*/

	// Tool工具：ToolsNode 组件是一个用于扩展模型能力的组件，它允许模型调用外部工具来完成特定的任务
	//components.BrowserUseTool()

	// 链式编排
	//components.CreateChain()
	//components.SimpleAgent()

	// 图式编排
	//components.CreateGraphOrc()      // 图式编排无model
	//components.GraphGatherARKModel() // 图式编排集合ark大模型
	//components.GraphWithState() // 图式编排集合ark大模型，切入state保留对话记忆功能+callback
	components.GraphWithNest() // 图式编排嵌套，把一整个graph作为一个lambda节点，嵌套入一个new graph

	// 应用开发工具链
}
