package main

import (
	"encoding/gob"
	"encoding/json"
	"log"
	"net"
	"net/rpc"
	"net/url"
	"os"
)

var config Config

func main() {

	gob.Register(url.Values{})
	// Getting Config File
	if configFileHandle, err := os.Open("./config.json"); err == nil {
		// Loading config from file to var
		configDecoder := json.NewDecoder(configFileHandle)
		if err := configDecoder.Decode(&config); err == nil {
			rpc.Register(new(ECOM))
			rpc.Register(new(Admin))
			listner, err := net.Listen("tcp", config.Address)
			if err != nil {
				log.Fatal(err.Error())
			}
			// Close the listener whenever we stop
			defer listner.Close()
			// Wait for incoming connections
			rpc.Accept(listner)
		} else {
			log.Fatalln("failed to decode config ")
		}
	} else {
		log.Fatalln("Failed to start Extension")
	}
}
