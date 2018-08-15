package network

import (
	"github.com/gg001226/link/message"
	"net"
	"encoding/binary"
	"log"
	"github.com/gg001226/link/account"
	"fmt"
	"regexp"
)

type Server struct {
	id		uint8
	ip		string
	post	string			//服务器监听端口
	Peers	map[uint64]*Peer//服务器连接的账户
}

type Peer struct {
	id       uint64   //该账户的id
	name     string   //该账户的昵称
	password string   //密码
	conn     net.Conn //连接该账户的conn对象，通过这个对象才能实现read() send()等功能
}

var AllAccounts map[uint64]*Peer

func NewServer(IP string) *Server {
	s := Server{ip:IP}
	s.init(IP)
	return &s
}

func (s *Server) init(IP string){
	//这里要添加初始化参数
	s.id = 0

	s.ip = ""
	s.post = IP

	s.Peers = make(map[uint64]*Peer)
	log.Println("[SRV]服务器初始化完毕")
}

func (s *Server) Start() error {
	//开始监听
	go s.run()
	return nil
}

//该函数启动时会等待客户端的连接
func (s *Server) waitForConnect(listener net.Listener) (*Peer, error){
	//首先，等待连接
	conn, err := listener.Accept()
	if err != nil {
		log.Println("[ERR]未能接受连接请求: ", err)
		return nil, err
	}

	log.Println("[SRV]服务器与客户端已经建立连接，等待验证")
	p := &Peer{conn:conn}
	//连接成功后，需要验证，并且让客户端提供账号信息
	ok := false
	for !ok {
		msg, from, to, err := p.Read()
		if err != nil || to != account.SYSTEM_ACCOUNT {
			log.Println("[WARN]客户端连接失败，未发送验证信息/验证信息读取失败")
			continue
		}

		switch msg.MsgType {
		case message.TEXT_MESSAGE:
			log.Println("[WARN]账号未登录，无法发送消息")

		case message.REGISTER_MESSAGE:
			p.id = from

			temp, _ := regexp.Compile("\\n")
			p.name = temp.Split(msg.Content, 2)[0]
			p.password = temp.Split(msg.Content, 2)[1]

			err = s.RegisterAccount(p)

			if err == nil {ok = true}

		case message.LOGIN_MESSAGE:
			p.id = from
			p.password = msg.Content

			err = s.LoginAccount(p)

			if err == nil {ok = true}

		case message.CONTACT_MESSAGE:
			continue
		}

		if err != nil {
			s.contact(message.Message{Content:err.Error(), MsgType:message.ERROR_MESSAGE}, conn)
		} else {
			s.Peers[from] = p
			err = s.Transmit(message.Message{Content:"成功", MsgType:0}, from, account.SYSTEM_ACCOUNT)
		}
		if err != nil {
			log.Println(err)
		}
	}

	return p, nil
}

func (s *Server) run() error {
	//这边是监听循环，每次有新账户链接，都要把它加入peers[]列表中，并且新开一个线程loop()来处理
	listen, err := net.Listen("tcp", s.ip+s.post)
	defer listen.Close()
	if err != nil {
		log.Println("[ERR]服务器未能开启监听! ", err)
		return err
	}

	log.Println("[SRV]服务器开始监听")
	for {
		peer, err := s.waitForConnect(listen)
		if err != nil {
			return err
		}

		//_, in := s.Peers[peer.id]
		//if in {
		//	log.Println("[ERR]该账号已经在线! ")
		//	continue
		//} else {
		log.Println(fmt.Sprintf("[SRV]客户端%d连接至服务器! ", peer.id))
		s.Peers[peer.id] = peer
		go s.loop(peer.id)
		//}
	}

	return nil
}

func (s *Server) loop(id uint64) {
	//这边是针对特定账户的处理函数，监听从id发过来的消息，然后转发
	peer := s.Peers[id]
	for {
		msg, from, to, err := peer.Read()
		if err != nil {
			continue
		}
		if from != id {
			log.Println(fmt.Sprintf("[WARN]消息来源id不符! (%d, %d)", id, from))
			//continue
		}

		if to == account.SYSTEM_ACCOUNT {
			//服务器和客户端的底层通信
			log.Println("[SRV]服务器与", id, "通信")
			continue
		}

		log.Println("[SRV]服务器收到来自", id, "的信息, ", msg.Content)
		err = s.Transmit(msg, id, to)
		if err != nil {
			continue
		}
	}
}

func (s *Server) Save() error {
	//这是保存函数，Save()将程序的信息保存到服务端的电脑上
	//目前只用存储账户信息，即只需要调用s.savePeer(id)即可
	//后面需要加入群聊聊天记录等
	return nil
}

func (s *Server) savePeer(id uint64) error {
	//将指定账户的信息存储
	return nil
}

func (s *Server) RegisterAccount(p *Peer) error {
	_, in := AllAccounts[p.id]
	if !in {
		AllAccounts[p.id] = p
		log.Println("[ACC]账号注册成功! ")
		return nil
	} else {
		return fmt.Errorf("[ACC]该账号已经注册! ")
	}
}

func (s *Server) LoginAccount(p *Peer) error {
	a, registered := AllAccounts[p.id]
	if !registered {
		err := fmt.Errorf("[ACC]该账号未注册! ")
		log.Println(err)
		return err
	} else if a.password != p.password {
		err := fmt.Errorf("[ACC]密码错误! ")
		log.Println(err)
		return err
	} else {p = a}

	_, in := s.Peers[p.id]
	if in {
		err := fmt.Errorf("[ACC]该账号已在线! ")
		log.Println(err)
		return err
	}

	log.Println("[ACC]登录成功! ")
	return nil
}

func (s *Server) Transmit(msg message.Message, to uint64, from uint64) error {
	//这边是从服务器转发给目标账户的函数
	log.Println("[SRV]转发给", to, "的消息: ", msg.Content)
	p, ok := s.Peers[to]
	if !ok {
		err := fmt.Errorf("[ERR]消息目标不存在! ")
		log.Println(err)
		return err
	}

	b := make([]byte, 8)//b为要发送的[]byte数组
	binary.BigEndian.PutUint64(b, from)//先将消息来源id压入b中

	b = append(b, msg.Encode()...)//再将msg压入b中

	_, err := p.conn.Write(b)

	return err
}

func (p *Peer) Read() (message.Message, uint64, uint64, error){
	b := make([]byte, 256+8+1)
	n, err := p.conn.Read(b)
	if n<8 || err != nil {
		log.Println("[ERR]读取消息时出错!")
		return message.Message{}, 0, 0, err
	}

	to := binary.BigEndian.Uint64(b[:8])//取出byte中的头8个，是目标的id
	b = b[8:]

	from := binary.BigEndian.Uint64(b[:8])
	b = b[8:]//b的8-16位是来源的id

	return message.Decode(b), from, to, nil
}

//服务器与未登录client的底层通信
func (m *Server) contact(msg message.Message, conn net.Conn) error {
	log.Println("[SRV]与客户端的底层通信")

	b := make([]byte, 8)//b为要发送的[]byte数组
	binary.BigEndian.PutUint64(b, account.SYSTEM_ACCOUNT)//来自于服务器

	b = append(b, msg.Encode()...)//再将msg压入b中

	_, err := conn.Write(b)

	return err
}