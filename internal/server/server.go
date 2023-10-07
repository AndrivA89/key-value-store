package server

import (
	"github.com/AndrivA89/key-value-store/internal/entity"
	"github.com/AndrivA89/key-value-store/internal/entity/errors"
	grpcServer "github.com/AndrivA89/key-value-store/internal/server/grpc"
	httpServer "github.com/AndrivA89/key-value-store/internal/server/http"
)

func NewServer(serverType entity.ServerType) (Server, error) {
	switch serverType {
	case entity.ServerGRPC:
		return grpcServer.NewServer(), nil
	case entity.ServerHTTP:
		return httpServer.NewServer(), nil
	default:
		return nil, errors.UnknownServerType
	}
}
