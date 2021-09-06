package znet

import (
	"fmt"
	"github.com/liyiweixiangmu/zinx/ziface"
	"strconv"
)

type MsgHandler struct {
	Apis map[uint32]ziface.IRouter
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis: make(map[uint32]ziface.IRouter),
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
