package client

import (
	"context"
	"io"
	"math/rand"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/skeremidchiev/gRPC/app/comm"
	"github.com/skeremidchiev/gRPC/app/storage"
	"google.golang.org/grpc"
)

const (
	letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func callCreate(cl comm.CommServiceClient) {
	message := comm.Data{
		Body: RandString(),
	}

	log.Infof("[Client Create()] call with message body: %s\n", message.Body)
	response, err := cl.Create(context.Background(), &message)
	if err != nil {
		log.Warningf("[Client Create()] server responded with error: %s\n", err.Error())
		return
	}

	log.Infof("[Client Create()] server responded:\n\tstatus: %t\n\terror: %s\n", response.Status, response.Error)
}

func callRemove(cl comm.CommServiceClient, cs storage.ClientStorage) {
	value, err := cs.GetRandom()
	if err != nil {
		log.Warningf("[Client Remove()] error: %s\n", err.Error())
		return
	}

	message := comm.Data{
		Body: value,
	}

	log.Infof("[Client Remove()] call with message body: %s\n", message.Body)
	response, err := cl.Create(context.Background(), &message)
	if err != nil {
		log.Warningf("[Client Remove()] server responded with error: %s\n", err.Error())
		return
	}

	log.Infof("[Client Remove()] server responded:\n\tstatus: %t\n\terror: %s\n", response.Status, response.Error)
}

func callList(cl comm.CommServiceClient) {
	log.Infof("[Client List()]\n")
	stream, err := cl.List(context.Background(), &comm.EmptyMessage{})
	if err != nil {
		log.Warningf("[Client List()] server responded with error: %s\n", err.Error())
		return
	}

	for {
		data, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Warningf("[Client List()] server responded with error: %s\n", err.Error())
		}
		log.Infof("[Client List()] data: %s\n", data)
	}
}

func StartClient(cs storage.ClientStorage) {
	conn, err := grpc.Dial(":8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("[Client] Failed to connect: %s\n", err)
	}

	defer conn.Close()

	cl := comm.NewCommServiceClient(conn)

	for {
		time.Sleep(time.Duration(1+rand.Intn(1)) * time.Second)
		switch rand.Intn(3) {
		case 0:
			callCreate(cl)
		case 1:
			callRemove(cl, cs)
		case 2:
			callList(cl)
		default:
		}
	}
}

func RandString() string {
	rs := make([]rune, 10)
	for i := 0; i < 10; i++ {
		rs[i] = []rune(letters)[rand.Intn(len(letters))]
	}
	return string(rs)
}
