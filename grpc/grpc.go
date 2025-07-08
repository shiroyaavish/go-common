package grpc

import (
	"context"
	"fmt"
	"github.com/IntelXLabs-LLC/go-common/logger"
	"net"

	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"
)

// unaryServerInterceptor is a gRPC unary server interceptor that limits the number of inflight messages per user.
// Args:
//
//	r: Redis client for performing operations on Redis server.
//	maxMessages: Maximum number of inflight messages allowed per user.
//
// Returns:
//
//	A gRPC UnaryServerInterceptor function.
//
// Usage Example:
//
//	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(unaryServerInterceptor(redisClient, limit)))
func unaryServerInterceptor(r *redis.Client, maxMessages int64) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		userID := ctx.Value("userID").(string)

		userInflightCount, err := r.Incr(ctx, userID).Result()
		if err != nil {
			return nil, err
		}

		if userInflightCount > maxMessages {
			return nil, fmt.Errorf("Max inflight messages reached")
		}

		resp, err := handler(ctx, req)
		if err != nil {
			return nil, err
		}

		_, err = r.Decr(ctx, userID).Result()
		if err != nil {
			return nil, err
		}

		return resp, err
	}
}

// GrpcServer is a type that represents a gRPC server.
type GrpcServer struct {
	listener    net.Listener
	Srv         *grpc.Server
	port        int
	redisClient *redis.Client
}

// New is a function that creates a new instance of the GrpcServer struct.
// Args:
//
//	port: The port number on which the server will listen for incoming connections.
//	limit: The maximum number of inflight messages allowed per user.
//	redisClient: The Redis client used for performing operations on the Redis server.
//
// Returns:
//
//	A pointer to the newly created GrpcServer instance.
//
// Usage Example:
//
//	server := New(8080, 10, redisClient)
func New(port int, limit int64, redisClient *redis.Client) *GrpcServer {
	lis, err := net.Listen("tcp", fmt.Sprint(":", port))
	if err != nil {
		logger.Fatal(err)
	}
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(unaryServerInterceptor(redisClient, limit)))
	return &GrpcServer{
		listener:    lis,
		Srv:         grpcServer,
		port:        port,
		redisClient: redisClient,
	}
}

// StartAsync starts the gRPC server asynchronously.
//
// It launches a goroutine that calls Serve on the listener of the gRPC server.
//
// The function logs a fatal error if Serve returns an error, and logs an info message when the server successfully starts.
//
// This method should be called after setting all the necessary configurations and registering the server with gRPC.
//
// Example usage:
//
//	gs := &GrpcServer{
//	    listener:    listener,
//	    Srv:         srv,
//	    port:        port,
//	    redisClient: redisClient,
//	}
//	gs.StartAsync()
//
// Returns nothing.
func (g *GrpcServer) StartAsync() {
	go func() {
		err := g.Srv.Serve(g.listener)
		if err != nil {
			logger.Fatal(err, "failed to start grpc server")
			return
		}
	}()
	logger.Info("Started Server on %d", g.port)
}
