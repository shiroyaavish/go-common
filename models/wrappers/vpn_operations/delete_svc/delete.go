package delete_svc

import (
	"context"
	"github.com/shiroyaavish/go-common/config"
	"github.com/shiroyaavish/go-common/logger"
	"github.com/shiroyaavish/go-common/models/wrappers"
	"github.com/shiroyaavish/go-common/models/wrappers/vpn_operations/delete_svc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type IDeleteWrapper interface {
	DeleteClient(request *DeleteClientRequest) *DeleteClientResponse
}

type DeleteWrapper struct {
	defaultTimeout time.Duration
	grpcUrl        string
	serviceName    string
	conn           *grpc.ClientConn
	deleteClient   proto.DeleteClientServiceClient
}

func NewDeleteWrapper(config config.WrapperConfig) IDeleteWrapper {
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

	return &DeleteWrapper{
		defaultTimeout: timeout,
		grpcUrl:        config.GrpcUrl,
		serviceName:    "delete",
		conn:           conn,
		deleteClient:   proto.NewDeleteClientServiceClient(conn),
	}
}

func (d *DeleteWrapper) DeleteClient(request *DeleteClientRequest) *DeleteClientResponse {
	ctx, cancel := context.WithTimeout(context.Background(), d.defaultTimeout)
	defer cancel()

	resp, err := d.deleteClient.DeleteClient(ctx, request.DeleteClientRequest)
	if err != nil {
		return &DeleteClientResponse{
			GRPCError: &wrappers.GRPCError{
				Code: 500,
				Err:  err.Error(),
			},
			DeleteClientResponse: nil,
		}
	}

	return &DeleteClientResponse{
		GRPCError:            nil,
		DeleteClientResponse: resp,
	}
}
