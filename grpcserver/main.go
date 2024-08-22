package grpcserver

import "context"

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type GrpcServer interface {
	Serve(ctx context.Context) error
}

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

// Identfier of the  package found messages having the format "senzing-6204xxxx".
const ComponentID = 6204

// Log message prefix.
const Prefix = "serve-grpc.grpcserver."

// Default gRPC Observer port
const DefaultGrpcObserverPort = "8260"

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

// Message templates.
var IDMessages = map[int]string{
	2000: "Entry: %+v",
	2001: "SENZING_ENGINE_CONFIGURATION_JSON: %v",
	2002: "Enabling all services.",
	2003: "Server listening at %v",
	2004: "Serving avoided.",
	4001: "Call to net.Listen(tcp, %s) failed.",
	4002: "Call to Szdiagnostic.PurgeRepository() failed.",
	4003: "Call to Szengine.Destroy() failed.",
	5001: "Failed to serve.",
}

// Status strings for specific messages.
var IDStatuses = map[int]string{}
