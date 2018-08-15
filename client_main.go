package main

import (
	"github.com/gg001226/link/manager"
	"github.com/gg001226/link/network"
	"github.com/gg001226/link/message"
	"time"
)

func main()  {
	/*client := network.NewClient()

	client.Init(network.ServerIP)
	client.Start()

	id := rand.Uint64()
	msg := message.Message{Content:"test data", MsgType:0}
	client.Send(message.Message{Content:"验证信息", MsgType:0}, 0, id)
	for {
		client.Send(msg, id, id)
		time.Sleep(5*time.Second)
	}*/
	network.AllAccounts = make(map[uint64]*network.Peer)

	server := network.NewServer(network.ServerIP)
	server.Start()

	mng, _ := manager.NewManager(network.ServerIP)
	mng.Start(network.ServerIP)

	mng.Register("test account", "111111")
	for {
		//mng.Login(mng.Account.GetID(), "123")
		//mng.Login(mng.Account.GetID(), "111111")
		mng.Send(message.Message{Content:"test data", MsgType:0}, mng.Account.GetID())
		time.Sleep(time.Second*5)
	}
}
