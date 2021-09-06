package main

import (
	"fmt"
	"github.com/liyiweixiangmu/zinx/znet"
	"io"
	"net"
	"time"
)

/**
 * @Description:模拟客户端
 */
func main() {
	fmt.Println("client start")
	time.Sleep(1 * time.Second)
	// 1 直接链接远程服务器，得到一个conn链接
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start err,exit!")
		return
	}

	for {
		dp := znet.NewDataPack()
		binaryMsg, err := dp.Pack(znet.NewMsgPackage(0, []byte("zinx0。5 框架消息 ")))
		if err != nil {
			fmt.Println("Pack msg err:", err)
			return
		}
		if _, err := conn.Write(binaryMsg); err != nil {
			fmt.Println("Write msg error:", err)
			return
		}

		// 接收服务端返回的ping
		binaryHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, binaryHead); err != nil {
			fmt.Println("read head error:", err)
			return
		}

		msgHead, err := dp.Unpack(binaryHead)
		if err != nil {
			fmt.Println("Unpack head err")
			break
		}
		if msgHead.GetMsgLen() > 0 {
			msg := msgHead.(*znet.Message)

			msg.Data = make([]byte, msgHead.GetMsgLen())
			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("Read body err:", err)
				return
			}
			fmt.Println("--> Recv Msg: ID=", msg.Id, " , len=", msg.DataLen, ", data=", string(msg.Data))
		}

		time.Sleep(1 * time.Second)

	}

}
