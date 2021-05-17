package main

import (
	"github.com/skeremidchiev/gRPC/app/client"
	"github.com/skeremidchiev/gRPC/app/server"
	"github.com/skeremidchiev/gRPC/app/storage/mapstorage"
)

func main() {
	setupLogger()
	setupRandSeed()

	s := mapstorage.NewStorage()

	go server.StartServer(s)

	client := client.NewClient(s)
	defer client.Close()
	client.Run()
}
