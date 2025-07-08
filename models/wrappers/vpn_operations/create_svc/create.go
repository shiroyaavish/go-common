package create_svc

import (
	"context"
	"github.com/shiroyaavish/go-common/config"
	"github.com/shiroyaavish/go-common/logger"
	"github.com/shiroyaavish/go-common/models/wrappers"
	"github.com/shiroyaavish/go-common/models/wrappers/vpn_operations/create_svc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type ICreateWrapper interface {
	CreateNewClient(request *CreateClientRequest) *CreateClientResponse
}

type CreateWrapper struct {
	defaultTimeout time.Duration
	grpcUrl        string
	serviceName    string
	conn           *grpc.ClientConn
	createClient   proto.CreateClientServiceClient
}

func NewCreateWrapper(config config.WrapperConfig) ICreateWrapper {
	timeout := 5 * time.Second

	if config.TimeoutSec > 0 {
		timeout = time.Duration(config.TimeoutSec) * time.Second
	}

	if config.GrpcUrl == "" {
		config.GrpcUrl = "vmua:50051"

		logger.Info("Api Url is missing for Search, defaulting to: %s\n", config.GrpcUrl)
	}

	conn, err := grpc.Dial(config.GrpcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal(err)
	}

	return &CreateWrapper{
		defaultTimeout: timeout,
		grpcUrl:        config.GrpcUrl,
		serviceName:    "create",
		conn:           conn,
		createClient:   proto.NewCreateClientServiceClient(conn),
	}
}

func (c *CreateWrapper) CreateNewClient(request *CreateClientRequest) *CreateClientResponse {
	ctx, cancel := context.WithTimeout(context.Background(), c.defaultTimeout)
	defer cancel()

	resp, err := c.createClient.CreateNewClient(ctx, &proto.CreateClientRequest{
		DeviceId:    request.DeviceId,
		OperationId: request.OperationId,
	})
	if err != nil {
		return &CreateClientResponse{
			GRPCError: &wrappers.GRPCError{
				Code: 500,
				Err:  err.Error(),
			},
			CreateClientResponse: nil,
		}
	}

	return &CreateClientResponse{
		GRPCError:            nil,
		CreateClientResponse: resp,
	}
}
