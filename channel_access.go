package cronjob

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func NewGrpcServer(
	healthCheckHandler grpc_health_v1.HealthServer,
) *GrpcServer {
	server := grpc.NewServer()
	grpc_health_v1.RegisterHealthServer(server, healthCheckHandler)
	return &GrpcServer{server}
}

type GrpcServer struct {
	*grpc.Server
}