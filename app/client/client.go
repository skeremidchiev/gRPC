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

type Client struct {
	csc comm.CommServiceClient
	cs  storage.ClientStorage
	cc  *grpc.ClientConn
}

func (c *Client) CallCreate() error {
	message := comm.Data{
		Body: RandString(),
	}

	log.Infof("[Client Create()] call with message body: %s\n", message.Body)
	response, err := c.csc.Create(context.Background(), &message)
	if err != nil {
		log.Warningf("[Client Create()] server responded with error: %s\n", err.Error())
		return err
	}

	log.Infof("[Client Create()] server responded:\n\tstatus: %t\n\terror: %s\n", response.Status, response.Error)
	return nil
}

func (c *Client) CallRemove() error {
	value, err := c.cs.GetRandom()
	if err != nil {
		log.Warningf("[Client Remove()] error: %s\n", err.Error())
		return err
	}

	message := comm.Data{
		Body: value,
	}

	log.Infof("[Client Remove()] call with message body: %s\n", message.Body)
	response, err := c.csc.Create(context.Background(), &message)
	if err != nil {
		log.Warningf("[Client Remove()] server responded with error: %s\n", err.Error())
		return err
	}

	log.Infof("[Client Remove()] server responded:\n\tstatus: %t\n\terror: %s\n", response.Status, response.Error)
	return nil
}

func (c *Client) CallList() error {
	log.Infof("[Client List()]\n")
	stream, err := c.csc.List(context.Background(), &comm.EmptyMessage{})
	if err != nil {
		log.Warningf("[Client List()] server responded with error: %s\n", err.Error())
		return err
	}

	for {
		data, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Warningf("[Client List()] server responded with error: %s\n", err.Error())
			continue
		}
		log.Infof("[Client List()] data: %s\n", data.Body)
	}

	return nil
}

func (c *Client) Run() {
	for {
		time.Sleep(time.Duration(1+rand.Intn(9)) * time.Second)
		switch rand.Intn(3) {
		case 0:
			c.CallCreate()
		case 1:
			c.CallRemove()
		case 2:
			c.CallList()
		default:
		}
	}
}

func NewClient(storage storage.ClientStorage) *Client {
	connection, err := grpc.Dial(":8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("[Client] Failed to connect: %s\n", err)
	}
	serviceClient := comm.NewCommServiceClient(connection)

	return &Client{
		csc: serviceClient,
		cs:  storage,
		cc:  connection,
	}
}

// Destructor
func (c *Client) Close() {
	c.cc.Close()
}

func RandString() string {
	rs := make([]rune, 10)
	for i := 0; i < 10; i++ {
		rs[i] = []rune(letters)[rand.Intn(len(letters))]
	}
	return string(rs)
}
