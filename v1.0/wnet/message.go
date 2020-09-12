package wnet


type Message struct {
	Id uint32	//消息ID
	DataLen uint32 	//消息长度
	Data []byte	//消息内容
}

//创建Message
func NewMessage(id uint32,data []byte) *Message{
	return &Message{
		Id: id,
		DataLen: uint32(len(data)),
		Data: data,
	}
}

//获取消息ID
func (m *Message) GetMsgId() uint32	{
	return m.Id
}

//获取消息长度
func (m *Message) GetMsgLen() uint32 {
	return m.DataLen
}

//获取消息内容
func (m *Message) GetData() []byte	{
	return m.Data
}

//设置消息ID
func (m *Message) SetMsgId(id uint32) {
	m.Id=id
}

//设置消息长度
func (m *Message) SetMsgLen(msgLen uint32)  {
	m.DataLen = msgLen
}

//设置消息内容
func (m *Message) SetData(msg []byte) {
	m.Data = msg
}
