package network

import (
	"testing"
	"net"
	"log"
	"github.com/gg001226/link/message"
	"time"
)

func TestClient(t *testing.T) {
	listen, err := net.Listen("tcp", ServerIP)
	if err != nil {
		t.Error("[ERR]开启服务端失败")
	}

	client := NewClient()
	client.Start()
	go client.clientTest()

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Println("[WARN]连接失败")
			continue
		}

		go serverTest(conn)
	}
}

func (c *Client) clientTest() {
	msg := message.Message{Content:"test data from client", MsgType:0}
	for {
		c.Send(msg, 0, 0)
		time.Sleep(5*time.Second)
	}
}

func serverTest(conn net.Conn) {
	msg := message.Message{Content:"test data from server", MsgType:0}
	for {
		b := make([]byte, 256+8+1)
		conn.Read(b)
		log.Println("[NET]服务端收到消息: ", message.Decode(b))
		r := make([]byte, 8)
		conn.Write(append(r, msg.Encode()...))
	}
}