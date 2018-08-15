package message

type Message struct {
	Content	string
	MsgType	uint8
}

const (
	TEXT_MESSAGE = iota
	LOGIN_MESSAGE
	REGISTER_MESSAGE
	CONTACT_MESSAGE
	ERROR_MESSAGE
)

//下边这个interface是预定升级，暂且不用管他
type MessageInterface interface {
	String()
	MsgType()
	Encode()
}

//这是把Message转化为[]byte的函数
func (m Message) Encode() []byte {
	b := make([]byte, 1)
	b[0] = m.MsgType
	b = append(b, []byte(m.Content)...)
	return b
}

func Decode(b []byte) Message {
	//b的第一位是MsgType, 后边都是content转化成的[]byte
	return Message{Content:string(b[1:]), MsgType:b[0]}
}