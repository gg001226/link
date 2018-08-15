package main

import "github.com/gg001226/link/network"

func main() {

	network.AllAccounts = make(map[uint64]*network.Peer)

	server := network.NewServer(network.ServerIP)
	server.Start()

	select {}
}
