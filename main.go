package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	pb "grpcgatewaystarter/pb"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{
		Message: fmt.Sprintf("%#v\n", in)}, nil
}

// grpcHandlerFunc differentiates between grpc and http traffic.
func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("zap.NewDevelopment(): %v\n", err)
	}
	defer logger.Sync()

	// grpc
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			middleware.ChainUnaryServer(
				grpc_zap.UnaryServerInterceptor(logger),
				grpc_recovery.UnaryServerInterceptor())))
	pb.RegisterGreeterServer(grpcServer, &server{})
	reflection.Register(grpcServer) // grpcurl & grpc-client-cli

	// http
	ctx := context.Background()
	mux := runtime.NewServeMux()
	address := fmt.Sprintf(":8080")
	dopts := []grpc.DialOption{
		grpc.WithTransportCredentials(
			insecure.NewCredentials())}
	err = pb.RegisterGreeterHandlerFromEndpoint(ctx, mux, address, dopts)
	if err != nil {
		logger.Fatal("pb.RegisterGreeterHandlerFromEndpoint()",
			zap.String("address", address),
			zap.Error(err),
		)
	}

	// listen
	logger.Info("server listening", zap.String("address", address))
	err = http.ListenAndServe(address, grpcHandlerFunc(grpcServer, mux))
	if err != nil {
		logger.Fatal("http.ListenAndServe()",
			zap.String("address", address),
			zap.Error(err),
		)
	}
}
