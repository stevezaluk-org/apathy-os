package api

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/immutable-os/immutable-os/src/go-serviced/pkg/supervisor"
	v1 "github.com/immutable-os/immutable-os/src/go-serviced/proto/service/v1"
	"google.golang.org/grpc"
)

type Api struct {
	v1.UnimplementedServicesServer

	// services A pointer to a map of service supervisors received from the daemon.Daemon
	services *map[string]*supervisor.Supervisor
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
func New(path string, services *map[string]*supervisor.Supervisor) (*Api, error) {
	if strings.HasPrefix(path, "~") {
		path = strings.Replace(path, "~", os.Getenv("HOME"), 1)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Base(path), 0644); err != nil {
			return nil, fmt.Errorf("failed to create directory %s: %v", path, err)
		}
	}

	api := &Api{
		services: services,
	}

	return api, nil
}
