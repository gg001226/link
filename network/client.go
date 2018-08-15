package network

import (
	"net"
	"log"
	"fmt"
	"github.com/gg001226/link/message"
	"encoding/binary"
)

const ServerIP = ":25500"

type Client struct {
	add		string	//地址
	conn	net.Conn	//net对象，用于发送/接受消息

	ReceiveMsg    chan message.Message //接受服务器消息的channel
	NeedToSendMsg chan message.Message //Account发送消息的channel
}

func NewClient() *Client {
	c := Client{}
	return &c
}

func (c *Client) Init(targetIP string) error {
	log.Println("[CLN]客户端初始化中...")
	/*resp, err := http.Get("http://myexternalip.com/raw")
	defer resp.Body.Close()

	if err != nil {
		log.Println("[NET]获取公网IP失败: ", err)
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("[NET]获取公网IP失败：", err)
	}*/
	c.add = ""

	c.ReceiveMsg = make(chan message.Message, 128)   //从服务器接受最多128条等待中的消息
	c.NeedToSendMsg = make(chan message.Message, 64) //最多有64条消息等待发送

	//尝试链接server，链接5次，若5次都没有成功则无法链接
	for i:=0; i<5; i++ {
		conn, err := c.dial(targetIP)
		if err == nil {
			c.conn = conn
			return nil
		}
	}
	return fmt.Errorf("[ERR]无法连接至服务器! ")
}

func (c *Client) dial(targetIP string) (net.Conn, error) {
	conn, err := net.Dial("tcp", targetIP)
	if err != nil {
		log.Println("[ERR]拨号错误: ", err)
		return nil, err
	}
	log.Println("[CLN]成功连接服务器!")
	return conn, nil
}

func (c *Client) Start() error {

	go c.run()
	go c.loop()
	return nil
}

func (c *Client) run() error {
	for {
		select {
		case msg := <-c.NeedToSendMsg:
			c.Send(msg, 0, 0)
		}
	}
	return nil
}

//监听函数，每当有消息从服务器发过来，就接受并且发送给channel，让account来处理
func (c *Client) loop() error {
	for {
		msg, from, err := c.Read()
		if err != nil {
			continue
		}
		log.Println("[CLN]收到信息:", msg.Content, "来自:", from)
		c.ReceiveMsg <- msg
	}
	return nil
}

//Client.Send(消息, 目标)
func (c *Client) Send(msg message.Message, target uint64, from uint64) error {
	//log.Println("[CLN]客户端发送消息...")
	if len(msg.Encode()) > 256+1 {
		return fmt.Errorf("[CLN]消息过长，禁止发送")
	}

	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, target)//这里b是目标的id

	f := make([]byte, 8)
	binary.BigEndian.PutUint64(f, from)//这里b又加上了来源的id

	b = append(b, f...)
	b = append(b, msg.Encode()...)//这里b加上了消息内容

	_, err := c.conn.Write(b)
	if err != nil {
		return err
	}
	log.Println("[CLN]发送成功! ")
	return nil
}

func (c *Client) Read() (message.Message, uint64, error) {
	//log.Println("[CLN]读取来自服务器的消息...")

	b := make([]byte, 256+8+1)//最多读取256位消息，
	n, err := c.conn.Read(b)
	if n<8 || err != nil {
		log.Println("[ERR]读取消息时出错!")
		return message.Message{}, 0 , err
	}

	from := binary.BigEndian.Uint64(b[:8])//取出byte中的头8个，是来源的id
	b = b[8:]

	return message.Decode(b), from, nil
}