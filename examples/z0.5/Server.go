package main

import (
	"fmt"
	"github.com/liyiweixiangmu/zinx/ziface"
	"github.com/liyiweixiangmu/zinx/znet"
)

type PingRouter struct {
	znet.BaseRouter
}

// Test Handle
func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle")
	// 先读取客户端数据，再回写 ping..ping..
	err := request.GetConnection().SendMsg(1, []byte("ping....ping..."))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	// 创建一个server句柄
	s := znet.NewServer("[Zinx v0.5]")

	// 给当前zinx框架添加一个自定义的路由
	//s.AddRouter(&PingRouter{})

	s.Serve()
}
