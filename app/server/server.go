package server

import (
	context "context"
	"net"

	log "github.com/sirupsen/logrus"
	"github.com/skeremidchiev/gRPC/app/comm"
	"github.com/skeremidchiev/gRPC/app/storage"
	"google.golang.org/grpc"
)

type Server struct {
	storage storage.ServerStorage
	comm.UnimplementedCommServiceServer
}

func (s *Server) Create(ctx context.Context, message *comm.Data) (*comm.Reply, error) {
	log.Infof("[Server Create()] received message: %s\n", message.Body)

	// Random error is commented out because it makes testing inconsistent
	// if rand.Intn(9) == 0 {
	// 	return &comm.Reply{Error: "Random Error", Status: false}, nil
	// }

	err := s.storage.Store(message.Body)
	if err != nil {
		return &comm.Reply{Error: err.Error(), Status: false}, nil
	}

	return &comm.Reply{Error: "", Status: true}, nil
}

func (s *Server) Remove(ctx context.Context, message *comm.Data) (*comm.Reply, error) {
	log.Infof("[Server Remove()] received message: %s\n", message.Body)

	// Random error is commented out because it makes testing inconsistent
	// if rand.Intn(9) == 0 {
	// 	return &comm.Reply{Error: "Random Error", Status: false}, nil
	// }

	err := s.storage.Delete(message.Body)
	if err != nil {
		return &comm.Reply{Error: err.Error(), Status: false}, nil
	}

	return &comm.Reply{Error: "", Status: true}, nil
}

// Don't see a clear way to return error on random for List() method
func (s *Server) List(messages *comm.EmptyMessage, stream comm.CommService_ListServer) error {
	log.Infof("[Server List()]\n")

	// returning error will cause endless loop?
	data, _ := s.storage.GetAll()
	// if err != nil {
	// 	return err
	// }

	for _, d := range data {
		if err := stream.Send(&comm.Data{Body: d}); err != nil {
			return err
		}
	}

	return nil
}

func NewServer(ss storage.ServerStorage) *Server {
	s := &Server{}
	s.storage = ss

	return s
}

func StartServer(ss storage.ServerStorage) {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("[Server] Listener failed: %s", err.Error())
	}

	grpcServer := grpc.NewServer()
	comm.RegisterCommServiceServer(grpcServer, NewServer(ss))

	log.Infof("[Server] Listening for requests on port 8080 ...\n")
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("[Server] gRPC server failed: %s", err.Error())
	}
}
