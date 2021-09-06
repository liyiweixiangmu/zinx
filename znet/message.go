package znet

type Message struct {
	DataLen uint32 //消息长度
	Id      uint32 //消息id
	Data    []byte //消息内容
}

// 创建一个Message消息包
func NewMsgPackage(id uint32, data []byte) *Message {
	return &Message{
		DataLen: uint32(len(data)),
		Id:      id,
		Data:    data,
	}
}

func (m *Message) GetMsgId() uint32 {
	return m.Id
}
func (m *Message) GetMsgLen() uint32 {
	return m.DataLen
}
func (m *Message) GetData() []byte {
	return m.Data
}
func (m *Message) SetMsgId(id uint32) {
	m.Id = id
}
func (m *Message) SetMsgLen(len uint32) {
	m.DataLen = len
}
func (m *Message) SetData(data []byte) {
	m.Data = data
}
