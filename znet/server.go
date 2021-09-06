package znet

import (
	"fmt"
	"github.com/liyiweixiangmu/zinx/utils"
	"github.com/liyiweixiangmu/zinx/ziface"
	"net"
)

// IServer的接口实现，定义一个Server的服务器模块
type Server struct {
	// 服务器名称
	Name string
	// 服务器绑定的IP版本
	IPVersion string
	// 服务器监听的IP
	IP string
	// 服务器监听的端口
	Port int
	// 给当前的server添加一个router，server注册的链接对应的处理业务
	MsgHandler ziface.IMsgHandler
}

// 定义当前客户端链接的所绑定的handle api
//func CallBackToClient (conn *net.TCPConn,data []byte,cnt int) error {
//	//回显业务
//	fmt.Println("[Conn Handle] CallbackToclient...")
//	if _,err := conn.Write(data[:cnt]);err != nil {
//		fmt.Println("write back buf err",err)
//		return errors.New("CallbackToClient error")
//	}
//	return nil
//}

// 启动服务器
func (s *Server) Start() {
	fmt.Printf("[Zinx] Server name %s, listenner at IP : %s., Port : %d \n",
		utils.GlobalObject.Name,
		utils.GlobalObject.Host,
		utils.GlobalObject.TcpPort)
	fmt.Printf("[Zinx] Version: %s, Maxconn : %d , MaxPackageSize : %d \n",
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPackageSize)

	go func() {
		// 1 获取一个TCP的addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error", err)
			return
		}

		// 2 监听服务器的地址
		listenner, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen", s.IPVersion, " err ", err)
		}
		fmt.Println("start Zinx Server succ, ", s.Name, " succ, Listening...")
		// 3 阻塞的等待客户端链接，处理客户端链接业务（读写）

		var cid uint32
		cid = 0
		for {
			// 如果有客户端链接过来，阻塞会返回
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}

			// 已经与客户端建立链接（conn） 做一个业务，做一个最基本的最大512字节长度的回写业务

			// 将该处理链接的业务方法 和 conn 进行绑定，得到我们的链接模块
			dealConn := NewConnection(conn, cid, s.MsgHandler)

			cid++
			// 启动当前的链接业务处理
			go dealConn.Start()

		}
	}()

}

// 停止服务器
func (s *Server) Stop() {
	// TODO 江一些服务器的资源，状态或者一些已经开辟的链接信息，进行停止或者回收
}

// 运行服务器
func (s *Server) Serve() {
	// 启动server的服务器功能
	s.Start()

	// TODO 做一些启动服务器之后的额外业务

	// 阻塞状态
	select {}
}

func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("Add Router Success")
}

func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandler(),
	}
	return s
}
