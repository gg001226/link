package manager

import (
	"bufio"
	"fmt"
	"strconv"
	"github.com/gg001226/link/message"
)

type cmdFunc func([]string) error

func (m *Manager) loadCommands(r *bufio.Reader, w *bufio.Writer) {
	m.cmdList = make(map[string]cmdFunc)
	m.cmdList["help"] = m.help
	m.cmdList["login"] = m.login
	m.cmdList["register"] = m.register
	m.cmdList["send"] = m.send
	m.cmdList["info"] = m.info
}

func (m *Manager) help(args []string) error {
	var helps = []string{
		"help\t:指令系统，所有指令如下:\n",
		"register\t:注册账户\n",
		"login\t:登录\n",
		"send\t:发送消息\n",
		"info\t:查看账户信息\n",
	}

	for _, help := range helps {
		m.Writer.WriteString(help)
	}
	m.Writer.Flush()

	if len(args)>=1 {
		return fmt.Errorf("未知参数")
	}
	return nil
}

func (m *Manager) login(args []string) error {
	if m.Account != nil {
		return fmt.Errorf("该账号已经登录! ")
	}

	m.WriteStringAndFlush("ID\t:")
	id, err := m.ReadString('\n')
	if err != nil {
		return err
	}

	m.WriteStringAndFlush("密码\t:")
	psw, err := m.ReadString('\n')
	if err != nil {return err}

	i, err := strconv.ParseUint(id, 10, 64)
	err = m.Login(i, psw)
	return err
}

func (m *Manager) register(args []string) error {
	if m.Account != nil {
		return fmt.Errorf("该账号已经登录! ")
	}

	m.WriteStringAndFlush("昵称\t:")
	name, err := m.ReadString('\n')
	if err != nil {return err}

	m.WriteStringAndFlush("密码\t:")
	psw, err := m.ReadString('\n')
	if err != nil {return err}

	err = m.Register(name, psw)
	return err
}

func (m *Manager) send(args []string) error {
	if m.Account == nil {
		return fmt.Errorf("尚未登录，无法发送消息! ")
	}
	
	var t string
	if len(args)>=1 {
		t = args[0]
	} else {
		m.WriteStringAndFlush("目标\t:")
		temp, err := m.ReadString('\n')
		if err != nil {return err}
		t = temp
	}

	var content string
	if len(args)>= 2{
		content = args[1]
	} else {
		m.WriteStringAndFlush("内容\t:")
		c, err := m.ReadString('\n')
		if err != nil {return err}
		content = c
	}


	to, err := strconv.ParseUint(t, 10, 64)
	err = m.Send(message.Message{Content:content, MsgType:message.TEXT_MESSAGE}, to)
	return err
}

func (m *Manager) info(args []string) error {
	m.Writer.WriteString("ID\t:"+strconv.FormatUint(m.Account.GetID(), 10)+"\n")
	m.Writer.WriteString("昵称\t:"+m.Account.GetName()+"\n")
	err := m.Writer.Flush()
	return err
}