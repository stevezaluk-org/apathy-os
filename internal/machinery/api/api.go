package api

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// Api An abstraction over a REST API served from a unix socket
type Api struct {
	// Socket The unix socket path that the API is available at
	Socket string

	// router The gin.Engine controlling the routing of this API
	router *gin.Engine

	// srv The underlying HTTP server for the API
	srv *http.Server
}

// RegisterHandler Register one or more handlers with the API
func (api *Api) RegisterHandler(handler ...Handler) {
	for _, h := range handler {
		h.Register(api.router)
	}
}

// EnsureSocket Ensures that the unix socket exists for this API
func (api *Api) EnsureSocket(path string) error {
	if _, err := os.Stat(path); errors.Is(err, fs.ErrNotExist) {
		if err := os.MkdirAll(filepath.Dir(path), 0600); err != nil {
			return fmt.Errorf("failed to create socket: %s (%v)", path, err)
		}
	}

	return nil
}

// Run Start the API and serve it at the provided socket path
func (api *Api) Run(ctx context.Context) error {
	listener, err := net.Listen("unix", api.Socket)
	if err != nil {
		return fmt.Errorf("failed to create unix listener: %v", err)
	}

	if err := api.srv.Serve(listener); err != nil {
		return fmt.Errorf("failed to start API server %s", err)
	}

	return nil
}

// New Initializes a new Api and returns a pointer to it
func New(path string) (*Api, error) {
	router := gin.Default()
	router.Use(
		gin.Recovery(),
	)

	gin.SetMode(gin.ReleaseMode)

	api := &Api{
		Socket: path,
		router: router,
		srv: &http.Server{
			Handler: router.Handler(),
			Addr:    fmt.Sprintf("unix://%s", path),
		},
	}

	if err := api.EnsureSocket(path); err != nil {
		return nil, err
	}

	return api, nil
}
