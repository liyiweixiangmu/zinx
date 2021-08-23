package main

import "github.com/liyiweixiangmu/zinx/znet"

func main() {
	// 创建一个server句柄
	s := znet.NewServer("[Zinx v0.2]")
	s.Serve()
}
