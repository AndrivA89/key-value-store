package grpc

import (
	"net"

	"github.com/labstack/gommon/log"
	"google.golang.org/grpc"

	"github.com/AndrivA89/key-value-store/internal/store"
	pb "github.com/AndrivA89/key-value-store/pkg/api"
)

type server struct {
	pb.UnimplementedKeyValueServer

	store *store.Store
}

func NewServer() *server {
	return &server{}
}

func (s *server) Start(store *store.Store) error {
	newServer := grpc.NewServer()
	s.store = store
	pb.RegisterKeyValueServer(newServer, s)

	l, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	return newServer.Serve(l)
}
