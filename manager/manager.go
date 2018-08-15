package manager

import (
	"github.com/gg001226/link/account"
	"github.com/gg001226/link/network"
	"fmt"
	"github.com/gg001226/link/message"
	"log"
)

type Manager struct {
	Account *account.Account
	Client	*network.Client
}

func NewManager(serverIP string) (*Manager, error) {
	var m Manager
	m.Client = network.NewClient()
	return &m, nil
}

func (m *Manager) InputInfo() error {
	//输入账户信息
	return nil
}

func (m *Manager) Start(targetIP string) {
	m.Client.Init(targetIP)
	m.Client.Start()

	//go m.run()
}

func (m *Manager) run() {
	var input string
	fmt.Scanln(input)
}


//注册账号到服务端，包括创建+注册两部分
func (m *Manager) Register(name, password string) error {
	a := account.NewAccount(name, password)

	//log.Println("[ACC]注册账号中...")
	msg := message.Message{Content:name+"\n"+password, MsgType:message.REGISTER_MESSAGE}
	err := m.Client.Send(msg, account.SYSTEM_ACCOUNT, a.GetID())
	if err != nil {return err}

	reply := <-m.Client.ReceiveMsg
	if reply.MsgType == message.ERROR_MESSAGE {
		log.Println(reply.Content)
		return fmt.Errorf(reply.Content)
	}
	m.Account = a
	log.Println("[ACC]账号注册成功！")

	return nil
}

func (m *Manager) Login(id uint64, psw string) error {
	a := account.LoginAccount(id, psw)

	log.Println("[ACC]登录中...")
	msg := message.Message{Content:psw, MsgType:message.LOGIN_MESSAGE}
	err := m.Client.Send(msg, account.SYSTEM_ACCOUNT, id)
	if err != nil {return err}

	reply := <-m.Client.ReceiveMsg
	if reply.MsgType == message.ERROR_MESSAGE {
		log.Println(reply.Content)
		return fmt.Errorf(reply.Content)
	}
	m.Account = a
	log.Println("[ACC]登录成功! ")

	return nil
}

func (m *Manager) Send(msg message.Message, to uint64) error {
	return m.Client.Send(msg, to, m.Account.GetID())
}

func (m *Manager) Read() (message.Message, uint64, error) {
	return m.Client.Read()
}