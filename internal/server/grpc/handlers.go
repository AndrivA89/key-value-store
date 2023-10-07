package grpc

import (
	"context"
	"log"

	pb "github.com/AndrivA89/key-value-store/pkg/api"
)

func (s *server) Get(ctx context.Context, r *pb.GetRequest) (*pb.GetResponse, error) {
	log.Printf("Received GET key=%v", r.Key)

	value, err := s.store.Get(r.GetKey())

	return &pb.GetResponse{Value: value}, err
}

func (s *server) Put(ctx context.Context, r *pb.PutRequest) (*pb.PutResponse, error) {
	log.Printf("Received GET key=%v", r.Key)

	err := s.store.Put(r.GetKey(), r.GetValue())

	return &pb.PutResponse{}, err
}

func (s *server) Delete(ctx context.Context, r *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	log.Printf("Received GET key=%v", r.Key)

	err := s.store.Delete(r.GetKey())

	return &pb.DeleteResponse{}, err
}
