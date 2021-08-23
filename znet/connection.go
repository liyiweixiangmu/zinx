package znet

import (
	"fmt"
	"github.com/liyiweixiangmu/zinx/utils"
	"github.com/liyiweixiangmu/zinx/ziface"
	"net"
)

/**
 * @Description:链接模块
 */
type Connection struct {
	// 当前链接的socket TCP套接字
	Conn *net.TCPConn
	// 链接id
	ConnID uint32

	// 当前链接的状态
	isClosed bool

	// 当前链接所绑定的处理业务的方法API
	handleAPI ziface.HandleFunc

	// 告知当前链接已经退出的/停止的channel
	ExitChan chan bool

	// 该链接处理的方法Router
	Router ziface.IRouter
}

//初始化链接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
	c := &Connection{
		Conn:     conn,
		ConnID:   connID,
		Router:   router,
		isClosed: false,
		ExitChan: make(chan bool, 1),
	}
	return c
}

// 链接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running")
	defer fmt.Println("connID=", c.ConnID, " Reader is exit, remote addr is", c.RemoteAddr().String())
	defer c.Stop()

	for {
		// 读取客户端的数据到buf中， 最大512字节
		buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		_, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("read err", err)
			continue
		}

		//// 调用当前链接所绑定的HandleAPI
		//if err := c.handleAPI(c.Conn, buf,cnt);err != nil {
		//	fmt.Println("ConnID", c.ConnID, " handle is error", err)
		//	break
		//}

		// 得到当前conn数据的Request请求数据
		req := Request{
			conn: c,
			data: buf,
		}

		// 执行注册的路由方法
		go func(request ziface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)
		// 从路由中，找到注册绑定的Conn对应的router调用

	}

}

func (c *Connection) Start() {
	fmt.Println("Conn start()... ConnID", c.ConnID)
	//启动从当前链接的读取数据的业务
	go c.StartReader()
	// TODO 启动从当前链接写数据的业务

}
func (c *Connection) Stop() {
	fmt.Println(" Conn stop()..... ConnID", c.ConnID)

	if c.isClosed == true {
		return
	}
	c.isClosed = true
	// 关闭socket链接
	c.Conn.Close()

	close(c.ExitChan)
}
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}
func (c *Connection) Send(data []byte) error {
	return nil
}
