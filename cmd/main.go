package main

import (
	"github.com/skeremidchiev/gRPC/app/client"
	"github.com/skeremidchiev/gRPC/app/server"
	"github.com/skeremidchiev/gRPC/app/storage/mapstorage"
)

func main() {
	setupLogger()

	s := mapstorage.NewStorage()

	go server.StartServer(s)

	client.StartClient(s)
}
