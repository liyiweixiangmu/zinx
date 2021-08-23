package znet

import "github.com/liyiweixiangmu/zinx/ziface"

// 实现router时，先潜入这个baserouter基类，然后根据需要对这个基类的方法进行重写
type BaseRouter struct{}

func (br *BaseRouter) PreHandle(request ziface.IRequest) {

}
func (br *BaseRouter) Handle(request ziface.IRequest) {

}
func (br *BaseRouter) PostHandle(request ziface.IRequest) {

}
