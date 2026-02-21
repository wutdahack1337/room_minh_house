package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"runtime/debug"

	log "github.com/InjectiveLabs/suplog"
	gateway "github.com/cosmos/gogogateway"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/rs/cors"
	"github.com/wutdahack1337/room_minh_house/backend/internal/database"
	types "github.com/wutdahack1337/room_minh_house/backend/pkg/pb"
	"go.elastic.co/apm/module/apmgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const GRPC_MAX_MSG_SIZE = 134217728 // 128MB

func NewServer(name string) *grpc.Server {
	logger := log.WithFields(log.Fields{name: "panic"})
	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			apmgrpc.NewUnaryServerInterceptor(apmgrpc.WithRecovery()),
			unaryLogStackTraceOnPanic(logger),
		),
		grpc.ChainStreamInterceptor(
			apmgrpc.NewStreamServerInterceptor(apmgrpc.WithRecovery()),
			streamLogStackTraceOnPanic(logger),
		),
		grpc.MaxRecvMsgSize(GRPC_MAX_MSG_SIZE),
		grpc.MaxSendMsgSize(GRPC_MAX_MSG_SIZE),
	}

	return grpc.NewServer(opts...)
}

func unaryLogStackTraceOnPanic(logger log.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				logger.
					WithField("method", info.FullMethod).
					WithField("stack", string(debug.Stack())).
					Error("panic on gRPC server")
				panic(r)
			}
		}()

		return handler(ctx, req)
	}
}

func streamLogStackTraceOnPanic(logger log.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		defer func() {
			if r := recover(); r != nil {
				logger.
					WithField("method", info.FullMethod).
					WithField("stack", string(debug.Stack())).
					Error("panic on gRPC server")
				panic(r)
			}
		}()
		return handler(srv, ss)
	}
}

func StartServer(grpcAddr string, httpAddr string, dbService database.DatabaseService) {
	ctx := context.Background()
	server := NewServer("apis")
	service := NewService(dbService)
	types.RegisterQueryServer(server, service)

	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &gateway.JSONPb{
			EmitDefaults: true,
			Indent:       "",
			OrigName:     true,
			AnyResolver:  nil,
		}),
	)

	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:3000",
			"http://localhost:5173",
		},
		AllowedMethods: []string{
			"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS",
		},
		AllowedHeaders: []string{
			"Content-Type", "Authorization", "Accept", "Origin", "X-Requested-With",
		},
		AllowCredentials: true,
		MaxAge:           86400,
	}).Handler(mux)

	if err := types.RegisterQueryHandlerFromEndpoint(ctx, mux, grpcAddr, []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(100_000000)),
	}); err != nil {
		panic(err)
	}

	listen, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		panic(err)
	}
	go server.Serve(listen)
	go func() {
		fmt.Println("serving http request:", httpAddr)
		err := http.ListenAndServe(httpAddr, corsHandler)
		if err != nil {
			panic(err)
		}
	}()
}
