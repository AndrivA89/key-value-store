package entity

type LoggerType string

const (
	PostgresLogger LoggerType = "postgres"
	FileLogger     LoggerType = "file"
)

type ServerType string

const (
	ServerGRPC ServerType = "grpc"
	ServerHTTP ServerType = "http"
)
