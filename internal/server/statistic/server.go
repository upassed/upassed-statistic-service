package statistic

import (
	"github.com/upassed/upassed-statistic-service/internal/config"
	"github.com/upassed/upassed-statistic-service/pkg/client"
	"google.golang.org/grpc"
)

type statisticServerAPI struct {
	client.UnimplementedStatisticServer
	cfg     *config.Config
	service assignmentService
}

type assignmentService interface {
}

func Register(gRPC *grpc.Server, cfg *config.Config, service assignmentService) {
	client.RegisterStatisticServer(gRPC, &statisticServerAPI{
		cfg:     cfg,
		service: service,
	})
}
