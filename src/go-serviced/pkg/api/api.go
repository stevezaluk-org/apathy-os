package api

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/immutable-os/immutable-os/src/go-serviced/pkg/config"
	"github.com/immutable-os/immutable-os/src/go-serviced/pkg/supervisor"
	v1 "github.com/immutable-os/immutable-os/src/go-serviced/proto/service/v1"
	"google.golang.org/grpc"
)

type Api struct {
	v1.UnimplementedServicesServer

	// config A pointer to the config file that this API Server was started with
	config *config.Config

	// services A pointer to a map of service supervisors received from the daemon.Daemon
	services *map[string]*supervisor.Supervisor

	// server The underlying server for the gRCP API
	server *grpc.Server
}

// Get Fetch a single service by name
func (api *Api) Get(ctx context.Context, request *v1.GetServiceRequest) (*v1.GetServiceResponse, error) {
	superv, ok := (*api.services)[request.Name]
	if !ok {
		return &v1.GetServiceResponse{Service: nil}, nil
	}

	return &v1.GetServiceResponse{Service: superv.Service}, nil
}

// List All services the daemon is currently aware of
func (api *Api) List(ctx context.Context, request *v1.ListServiceRequest) (*v1.ListServiceResponse, error) {
	services := make([]*v1.Service, len(*api.services))
	for _, val := range *api.services {
		services = append(services, val.Service)
	}

	return &v1.ListServiceResponse{Service: services}, nil
}

func (api *Api) Start(ctx context.Context, request *v1.StartServiceRequest) (*v1.StartServiceResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (api *Api) Stop(ctx context.Context, request *v1.StopServiceRequest) (*v1.StopServiceResponse, error) {
	//TODO implement me
	panic("implement me")
}

// CleanUp Removes the dangling socket from /var/run/go-serviced and gracefully stops the gRPC server
func (api *Api) CleanUp() error {
	api.server.GracefulStop()

	if _, err := os.Stat(api.config.Socket); !errors.Is(err, os.ErrNotExist) {
		if err := os.RemoveAll(api.config.Socket); err != nil {
			return err
		}
	}

	return nil
}

// Serve Starts the API server. Blocks the current go-routine.
func (api *Api) Serve() error {
	if strings.HasPrefix(api.config.Socket, "~") {
		api.config.Socket = strings.Replace(api.config.Socket, "~", os.Getenv("HOME"), 1)
	}

	if _, err := os.Stat(api.config.Socket); errors.Is(err, fs.ErrNotExist) {
		if err := os.MkdirAll(filepath.Dir(api.config.Socket), 0644); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", api.config.Socket, err)
		}
	}

	listener, err := net.Listen("unix", api.config.Socket)
	if err != nil {
		return err
	}
	defer listener.Close()

	v1.RegisterServicesServer(api.server, api)
	slog.Info("Listening for requests at " + api.config.Socket)
	if err := api.server.Serve(listener); err != nil {
		return err
	}

	return nil
}

// New Initializes a new API and serves it through a Unix socket at the specified path
func New(config *config.Config, services *map[string]*supervisor.Supervisor) (*Api, error) {
	return &Api{
		config:   config,
		services: services,
		server:   grpc.NewServer(),
	}, nil
}
