package grpcserver

import (
	"context"
	"fmt"
	"net"

	"github.com/senzing/g2-sdk-proto/go/g2config"
	"github.com/senzing/g2-sdk-proto/go/g2configmgr"
	"github.com/senzing/g2-sdk-proto/go/g2diagnostic"
	"github.com/senzing/g2-sdk-proto/go/g2engine"
	"github.com/senzing/g2-sdk-proto/go/g2product"
	"github.com/senzing/go-logging/logger"
	"github.com/senzing/go-logging/messagelogger"
	"github.com/senzing/go-observing/observer"
	"github.com/senzing/serve-grpc/g2configmgrserver"
	"github.com/senzing/serve-grpc/g2configserver"
	"github.com/senzing/serve-grpc/g2diagnosticserver"
	"github.com/senzing/serve-grpc/g2engineserver"
	"github.com/senzing/serve-grpc/g2productserver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// GrpcServerImpl is the default implementation of the GrpcServer interface.
type GrpcServerImpl struct {
	EnableG2config                 bool
	EnableG2configmgr              bool
	EnableG2diagnostic             bool
	EnableG2engine                 bool
	EnableG2product                bool
	LogLevel                       logger.Level
	Observers                      []observer.Observer
	Port                           int
	SenzingEngineConfigurationJson string
	SenzingModuleName              string
	SenzingVerboseLogging          int
	// FIXME:  All variables go in here. Cobra will be used outside to set the variables.
}

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

// Add G2Config service to gRPC server.
func (grpcServer *GrpcServerImpl) enableG2config(ctx context.Context, serviceRegistrar grpc.ServiceRegistrar) {
	server := &g2configserver.G2ConfigServer{}
	server.SetLogLevel(ctx, grpcServer.LogLevel)
	g2configserver.GetSdkG2config().Init(ctx, grpcServer.SenzingModuleName, grpcServer.SenzingEngineConfigurationJson, grpcServer.SenzingVerboseLogging)
	if grpcServer.Observers != nil {
		for _, observer := range grpcServer.Observers {
			server.RegisterObserver(ctx, observer)
		}
	}
	g2config.RegisterG2ConfigServer(serviceRegistrar, server)
}

// Add G2Configmgr service to gRPC server.
func (grpcServer *GrpcServerImpl) enableG2configmgr(ctx context.Context, serviceRegistrar grpc.ServiceRegistrar) {
	server := &g2configmgrserver.G2ConfigmgrServer{}
	server.SetLogLevel(ctx, grpcServer.LogLevel)
	g2configmgrserver.GetSdkG2configmgr().Init(ctx, grpcServer.SenzingModuleName, grpcServer.SenzingEngineConfigurationJson, grpcServer.SenzingVerboseLogging)
	if grpcServer.Observers != nil {
		for _, observer := range grpcServer.Observers {
			server.RegisterObserver(ctx, observer)
		}
	}
	g2configmgr.RegisterG2ConfigMgrServer(serviceRegistrar, server)
}

// Add G2Diagnostic service to gRPC server.
func (grpcServer *GrpcServerImpl) enableG2diagnostic(ctx context.Context, serviceRegistrar grpc.ServiceRegistrar) {
	server := &g2diagnosticserver.G2DiagnosticServer{}
	server.SetLogLevel(ctx, grpcServer.LogLevel)
	g2diagnosticserver.GetSdkG2diagnostic().Init(ctx, grpcServer.SenzingModuleName, grpcServer.SenzingEngineConfigurationJson, grpcServer.SenzingVerboseLogging)
	if grpcServer.Observers != nil {
		for _, observer := range grpcServer.Observers {
			server.RegisterObserver(ctx, observer)
		}
	}
	g2diagnostic.RegisterG2DiagnosticServer(serviceRegistrar, server)
}

// Add G2Engine service to gRPC server.
func (grpcServer *GrpcServerImpl) enableG2engine(ctx context.Context, serviceRegistrar grpc.ServiceRegistrar) {
	server := &g2engineserver.G2EngineServer{}
	server.SetLogLevel(ctx, grpcServer.LogLevel)
	g2engineserver.GetSdkG2engine().Init(ctx, grpcServer.SenzingModuleName, grpcServer.SenzingEngineConfigurationJson, grpcServer.SenzingVerboseLogging)
	if grpcServer.Observers != nil {
		for _, observer := range grpcServer.Observers {
			server.RegisterObserver(ctx, observer)
		}
	}
	g2engine.RegisterG2EngineServer(serviceRegistrar, server)
}

// Add G2Product service to gRPC server.
func (grpcServer *GrpcServerImpl) enableG2product(ctx context.Context, serviceRegistrar grpc.ServiceRegistrar) {
	server := &g2productserver.G2ProductServer{}
	server.SetLogLevel(ctx, grpcServer.LogLevel)
	g2productserver.GetSdkG2product().Init(ctx, grpcServer.SenzingModuleName, grpcServer.SenzingEngineConfigurationJson, grpcServer.SenzingVerboseLogging)
	if grpcServer.Observers != nil {
		for _, observer := range grpcServer.Observers {
			server.RegisterObserver(ctx, observer)
		}
	}
	g2product.RegisterG2ProductServer(serviceRegistrar, server)
}

// ----------------------------------------------------------------------------
// Main
// ----------------------------------------------------------------------------

func (grpcServer *GrpcServerImpl) Serve(ctx context.Context) error {
	logger, _ := messagelogger.NewSenzingApiLogger(ProductId, IdMessages, IdStatuses, grpcServer.LogLevel)

	// Log entry parameters.

	logger.Log(2000, grpcServer)

	// Determine which services to start. If no services are explicitly set, then all services are started.

	if !grpcServer.EnableG2config && !grpcServer.EnableG2configmgr && !grpcServer.EnableG2diagnostic && !grpcServer.EnableG2engine && !grpcServer.EnableG2product {
		logger.Log(2002)
		grpcServer.EnableG2config = true
		grpcServer.EnableG2configmgr = true
		grpcServer.EnableG2diagnostic = true
		grpcServer.EnableG2engine = true
		grpcServer.EnableG2product = true
	}

	// Set up socket listener.

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcServer.Port))
	if err != nil {
		logger.Log(4001, grpcServer.Port, err)
	}
	logger.Log(2003, listener.Addr())

	// Create server.

	aGrpcServer := grpc.NewServer()

	// Register services with gRPC server.

	if grpcServer.EnableG2config {
		grpcServer.enableG2config(ctx, aGrpcServer)
	}
	if grpcServer.EnableG2configmgr {
		grpcServer.enableG2configmgr(ctx, aGrpcServer)
	}
	if grpcServer.EnableG2diagnostic {
		grpcServer.enableG2diagnostic(ctx, aGrpcServer)
	}
	if grpcServer.EnableG2engine {
		grpcServer.enableG2engine(ctx, aGrpcServer)
	}
	if grpcServer.EnableG2product {
		grpcServer.enableG2product(ctx, aGrpcServer)
	}

	// Enable reflection.

	reflection.Register(aGrpcServer)

	// Run server.

	err = aGrpcServer.Serve(listener)
	if err != nil {
		logger.Log(5001, err)
	}

	return err
}
