package ziface

/*
	irequest接口
	实际上是把客户端请求的链接信息，和请求的数据包装到了一个request中
*/

type IRequest interface {

	//得到当前链接
	GetConnection() IConnection

	GetData() []byte

	// 得到请求消息id
	GetMsgID() uint32
	// 得到请求消息长度
	GetMsgLen() uint32
}
