package main

import (
	"log"
	"net"
	"net/rpc"
)

func main() {
	rpc.Register(new(HelloRequest))
	listner, err := net.Listen("tcp", "127.0.0.1:10000")
	if err != nil {
		log.Fatal(err.Error())
	}
	// Close the listener whenever we stop
	defer listner.Close()
	// Wait for incoming connections
	rpc.Accept(listner)
}
