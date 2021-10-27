package ziface

// 消息管理抽象层
type IMsgHandler interface {
	// 调度/执行对应的router消息处理方法
	DoMsgHandler(request IRequest)

	// 为消息添加具体处理逻辑
	AddRouter(msgID uint32, router IRouter)

	// 启动一个工作池
	StartWorkerPool()

	SendMsgToTaskQueue(request IRequest)
}
