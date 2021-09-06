package znet

import (
	"errors"
	"fmt"
	"github.com/liyiweixiangmu/zinx/ziface"
	"io"
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
	MsgHandler ziface.IMsgHandler
}

//初始化链接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandler) *Connection {
	c := &Connection{
		Conn:       conn,
		ConnID:     connID,
		MsgHandler: msgHandler,
		isClosed:   false,
		ExitChan:   make(chan bool, 1),
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
		//buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		//_, err := c.Conn.Read(buf)
		//if err != nil {
		//	fmt.Println("read err", err)
		//	continue
		//}

		// 创建一个拆包解包对象
		dp := NewDataPack()
		// 读取客户端的Msg Head 二进制流 8个字节
		headData := make([]byte, dp.GetHeadLen())
		_, err := io.ReadFull(c.GetTCPConnection(), headData)
		if err != nil {
			fmt.Println("Read head error", err)
			break
		}

		// 拆包， 得到msgID 和 msgDatalen 放在msg消息中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("Unpack head error", err)
			break
		}

		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			_, err := io.ReadFull(c.GetTCPConnection(), data)
			if err != nil {
				fmt.Println("Read msg data error", err)
				break
			}
		}
		msg.SetData(data)

		//// 调用当前链接所绑定的HandleAPI
		//if err := c.handleAPI(c.Conn, buf,cnt);err != nil {
		//	fmt.Println("ConnID", c.ConnID, " handle is error", err)
		//	break
		//}

		// 得到当前conn数据的Request请求数据
		req := Request{
			conn: c,
			msg:  msg,
		}

		// 执行注册的路由方法
		// 从路由中，找到注册绑定的Conn对应的router调用
		go c.MsgHandler.DoMsgHandler(&req)

	}

}

// 提供一个SendMsg方法 将我们要发送给客户端的数据，先进行封包，再发送
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("connection closed when send msg")
	}
	// 将data 进行封包
	dp := NewDataPack()

	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgId)
		return errors.New("Pack error msg")
	}

	// 将数据发送给客户端
	if _, err := c.Conn.Write(binaryMsg); err != nil {
		fmt.Println("Write msg id", msgId, " error:", err)
		return errors.New("conn write error")
	}

	return nil
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
