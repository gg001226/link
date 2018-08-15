package main

import "github.com/gg001226/link/network"

func main() {
	server := network.NewServer(network.ServerIP)
	server.Start()
}
