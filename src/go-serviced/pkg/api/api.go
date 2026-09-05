package api

import (
	"context"

	v1 "github.com/immutable-os/immutable-os/src/go-serviced/proto/service/v1"
	"google.golang.org/grpc"
)

type Api struct {
	v1.UnimplementedServicesServer
}

// Get Retrieve a services metadata
func (a *Api) Get(ctx context.Context, in *v1.GetServiceRequest, opts ...grpc.CallOption) (*v1.GetServiceResponse, error) {
	//TODO implement me
	panic("implement me")
}

// List Lists all currently loaded services
func (a *Api) List(ctx context.Context, in *v1.ListServiceRequest, opts ...grpc.CallOption) (*v1.ListServiceResponse, error) {
	//TODO implement me
	panic("implement me")
}

// Start Starts a service by name
func (a *Api) Start(ctx context.Context, in *v1.StartServiceRequest, opts ...grpc.CallOption) (*v1.StartServiceResponse, error) {
	//TODO implement me
	panic("implement me")
}

// Stop Stops a service by name
func (a *Api) Stop(ctx context.Context, in *v1.StopServiceRequest, opts ...grpc.CallOption) (*v1.StopServiceResponse, error) {
	//TODO implement me
	panic("implement me")
}

// New Initializes a new API and serves it through a Unix socket at the specified path
func New(path string) (*Api, error) {
	return &Api{}, nil
}
