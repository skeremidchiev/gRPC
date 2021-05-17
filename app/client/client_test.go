package client

import (
	"errors"
	"net"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/skeremidchiev/gRPC/app/comm"
	"github.com/skeremidchiev/gRPC/app/server"
	"google.golang.org/grpc"
)

type mockServerStorage struct{}

func (m *mockServerStorage) Store(string) error {
	return nil
}

func (m *mockServerStorage) Delete(string) error {
	return nil
}

func (m *mockServerStorage) GetAll() ([]string, error) {
	return []string{"1", "2"}, nil
}

type mockServerStorageErrors struct{}

func (m *mockServerStorageErrors) Store(string) error {
	return errors.New("Mocked Error")
}

func (m *mockServerStorageErrors) Delete(string) error {
	return errors.New("Mocked Error")
}

func (m *mockServerStorageErrors) GetAll() ([]string, error) {
	return nil, errors.New("Mocked Error")
}

type mockClientStorage struct{}

func (m *mockClientStorage) GetRandom() (string, error) {
	return "1", nil
}

type mockClientStorageErrors struct{}

func (m *mockClientStorageErrors) GetRandom() (string, error) {
	return "", errors.New("Mocked Error")
}

func Test_Success(t *testing.T) {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("[Server] Listener failed: %s", err.Error())
	}

	grpcServer := grpc.NewServer()
	comm.RegisterCommServiceServer(grpcServer, server.NewServer(&mockServerStorage{}))

	go func() {
		log.Infof("[Server] Listening for requests on port 8080 ...\n")
		err = grpcServer.Serve(listener)
		if err != nil {
			log.Fatalf("[Server] gRPC server failed: %s", err.Error())
		}
	}()

	clientCalls(t)
	grpcServer.GracefulStop()
	listener.Close()
}

func Test_Error_Cases(t *testing.T) {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("[Server] Listener failed: %s", err.Error())
	}

	grpcServer := grpc.NewServer()
	comm.RegisterCommServiceServer(grpcServer, server.NewServer(&mockServerStorageErrors{}))

	go func() {
		log.Infof("[Server] Listening for requests on port 8080 ...\n")
		err := grpcServer.Serve(listener)
		if err != nil {
			log.Fatalf("[Server] gRPC server failed: %s", err.Error())
		}
	}()

	clientCalls(t)
	grpcServer.GracefulStop()
	listener.Close()
}

func clientCalls(t *testing.T) {
	cl := NewClient(&mockClientStorage{})

	if err := cl.CallCreate(); err != nil {
		t.Error("Create: unexpected error")
	}

	if err := cl.CallRemove(); err != nil {
		t.Error("Remove: unexpected error")
	}

	if err := cl.CallList(); err != nil {
		t.Error("List: unexpected error")
	}

	cl = NewClient(&mockClientStorageErrors{})

	if err := cl.CallRemove(); err == nil {
		t.Error("Remove: method unexpectedly passed without errors")
	}
}
