package main

import (
	"fmt"
	"github.com/liyiweixiangmu/zinx/ziface"
	"github.com/liyiweixiangmu/zinx/znet"
)

type PingRouter struct {
	znet.BaseRouter
}

// Test PreRouter
func (this *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("Call Router PreHandle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping ... \n"))
	if err != nil {
		fmt.Println("call back before ping error")
	}
}

// Test Handle
func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte(" ping ...ping ...ping ... \n"))
	if err != nil {
		fmt.Println("call back ping ... erping ... error ")
	}
}

// Test PostRouter
func (this *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("Call Router PostHandle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte(" after ping ...\n"))
	if err != nil {
		fmt.Println("call back after ping ... error ")
	}
}

func main() {
	// 创建一个server句柄
	s := znet.NewServer("[Zinx v0.3]")

	// 给当前zinx框架添加一个自定义的路由
	s.AddRouter(&PingRouter{})

	s.Serve()
}
