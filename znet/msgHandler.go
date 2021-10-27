package znet

import (
	"fmt"
	"github.com/liyiweixiangmu/zinx/utils"
	"github.com/liyiweixiangmu/zinx/ziface"
	"strconv"
)

type MsgHandler struct {
	Apis map[uint32]ziface.IRouter
	// 负责worker取任务的消息队列
	TaskQueue []chan ziface.IRequest
	// 业务工作worker池的worker的数量
	WorkerPoolSize uint32
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis: make(map[uint32]ziface.IRouter),
		// 从全局配置中获取
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

func (mh *MsgHandler) DoMsgHandler(request ziface.IRequest) {
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgID=", request.GetMsgID(), " IS NOT FOUND !")
		return
	}
	// 根据msgID实现
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

func (mh *MsgHandler) AddRouter(msgID uint32, router ziface.IRouter) {
	// 判断当前msgid方法是否存在
	if _, ok := mh.Apis[msgID]; ok {
		//id已经注册
		panic("receive api,msgID=" + strconv.Itoa(int(msgID)))
	}
	// 添加msg雨API关系
	mh.Apis[msgID] = router
	fmt.Println("添加路由成功 MsgID=", msgID)
}

// StartWorkerPool 启动一个worker工作池 (开启工作吃的动作只能发生一次，一个zinx框架只能有一个worker工作池
func (mh *MsgHandler) StartWorkerPool() {
	// 根据workerPollSize 分别开启worker， 每个worker用一个go承载
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		// 一个worker被启动
		// 1 当前的woeker对应的channel消息队列 开辟空间，第0哥worker就用第0个channel
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		// 2 启动当前的worker，阻塞等待消息从channel 传递进来
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

// StartOneWorker 启动一个worker工作流程
func (mh *MsgHandler) StartOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("Worker ID=", workerID, " is started")

	// 不断的阻塞等待对应消息队列的消息
	for {
		select {
		// 如果有消息过来，出列的第一个就是客户端的request，执行当前request所绑定的业务
		case requqst := <-taskQueue:
			mh.DoMsgHandler(requqst)
		}
	}

}

// SendMsgToTaskQueue 将消息交给TaskQueue， 由worker进行处理
func (mh *MsgHandler) SendMsgToTaskQueue(request ziface.IRequest) {
	// 将消息pingju你分配给不通过的worker
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add ConnID=", request.GetConnection().GetConnID(), " request MsgID = ", request.GetMsgID(), " to workerID=", workerID)

	// 将消息发送给对应的worker的TaskQueue即可
	mh.TaskQueue[workerID] <- request
}
