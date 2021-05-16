package server

import (
	context "context"
	"net"

	log "github.com/sirupsen/logrus"
	"github.com/skeremidchiev/gRPC/app/comm"
	"google.golang.org/grpc"
)

type Server struct {
	comm.UnimplementedCommServiceServer
}

func (s *Server) Create(ctx context.Context, message *comm.Data) (*comm.Reply, error) {
	log.Printf("Received message from client: %s\n", message.Body)

	return &comm.Reply{Error: "", Status: true}, nil
}

func (s *Server) Remove(ctx context.Context, message *comm.Data) (*comm.Reply, error) {
	log.Printf("Received message from client: %s\n", message.Body)

	return &comm.Reply{Error: "", Status: true}, nil
}

func (s *Server) List(messages *comm.EmptyMessage, stream comm.CommService_ListServer) error {
	log.Printf("Received message from client\n")

	return nil
}

func NewServer() *Server {
	return &Server{}
}

func StartServer() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("[Server] Listener failed: %s", err.Error())
	}

	grpcServer := grpc.NewServer()
	comm.RegisterCommServiceServer(grpcServer, NewServer())

	log.Infof("[Server] Listening for requests on port 8080 ...\n")
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("[Server] gRPC server failed: %s", err.Error())
	}
}
