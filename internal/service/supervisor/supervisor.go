package supervisor

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"math"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/containerd/cgroups/v2/cgroup2"
	v1 "github.com/stevezaluk-org/apathy-os/internal/service/proto/service/v1"
	"golang.org/x/sys/unix"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Supervisor An abstraction for controlling the lifecycle of a service. Defines the core functions needed for starting,
// stopping, restarting, and managing a service. The supervisor structure also directly manages the service's cgroup,
// usually a child cgroup of the overarching `serviced.slice` cgroup.
type Supervisor struct {
	// Service The service that this supervisor is managing
	Service *v1.Service

	// proc Stores process information about the service (if currently running)
	proc *exec.Cmd

	// cgroup The cgroup manager for this specific service. Usually a cgroup under the `serviced` slice
	cgroup *cgroup2.Manager

	// cancel
	cancel context.CancelCauseFunc
}

// PriorityGate Returns the given priority of a service
func (supervisor *Supervisor) PriorityGate() int {
	// Either the value is too large for a signed integer or this supervisor doesn't have a Service attached to it
	if supervisor.Service == nil || supervisor.Service.Priority > uint64(math.MaxInt) {
		return 0
	}

	return int(supervisor.Service.Priority)
}

// Stop Initiates the stop routine for the given service. When a stop is requested, the supervisor first sends a
// unix.SIGTERM to the process and waits for a maximum of 20 seconds for the process to gracefully end. If the process
// fails to respond to this signal, its cgroup is frozen and is unix.SIGKILL'd
//
// TODO: The context management is weird here. Kill exec.CommandContext so that all functions can use a child context directly (simplifies mass shutdown)
// TODO: Unclear when SIGKILL is being sent (since logic for it is buried under supervisor.Start)
func (supervisor *Supervisor) Stop(ctx context.Context) error {
	if supervisor.proc.Process == nil {
		return nil
	}

	// send a SIGTERM and try and wait for this process to exit gracefully
	if err := supervisor.proc.Process.Signal(unix.SIGTERM); err != nil {
		return err
	}
	defer supervisor.cancel(errors.New("service kill requested"))

	for {
		select {
		case <-ctx.Done():
			// hits when the context is canceled, at this point the process will be SIGKILL'd
			if err := supervisor.cgroup.Freeze(); err != nil {
				return err
			}

			return nil
		default:
			// proc.ProcessState is only ever populated when the service ends, so we can wait for this
			// to return non-nil in a loop
			if err := supervisor.proc.ProcessState; err != nil {
				return nil
			}
		}
	}
}

// Watch Initiates a Watch loop on the specified service and starts it. Provides an automatic restart loop for the service.
// A new child context is  derived from the parent context that is passed to this function in order to ensure that
// services do not attempt to restart after system shutdown has been initiated
// TODO: Service accepts input (when os.Stdin) is mapped even after Service.Stop is called
func (supervisor *Supervisor) Watch(ctx context.Context, cancel context.CancelCauseFunc) {
	for i := range 3 {
		select {
		case <-ctx.Done():
			// TODO: Need to stop here
			return
		default:
			slog.Info("Service started", "svc", supervisor.Service.Name, "pid", supervisor.proc.Process.Pid)
			if err := supervisor.proc.Wait(); err != nil {
				log.Printf("Service %s died with error: %v", supervisor.Service.Name, err)
			}

			slog.Info("Attempting to start service", "svc", supervisor.Service.Name, "count", i+1)
			if err := supervisor.Start(ctx, cancel); err != nil {
				log.Printf("Failed to start service %s: %v", supervisor.Service.Name, err)
			}
		}
	}

	// never hit on system shutdown, only ever hit if a service repeatedly refuses to start
	defer supervisor.cancel(errors.New("service reached its maximum restarts"))
	log.Printf("Service %s has tried to restart 3 times. Service will not be restarted", supervisor.Service.Name)
}

// Start Starts the respective service.
func (supervisor *Supervisor) Start(ctx context.Context, cancel context.CancelCauseFunc) error {
	supervisor.cancel = cancel

	supervisor.proc = exec.CommandContext(ctx, supervisor.Service.Exec, supervisor.Service.Args...)
	supervisor.proc.Cancel = func() error {
		cleanUpCtx, cleanUpCancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cleanUpCancel()

		if err := supervisor.Stop(cleanUpCtx); err != nil {
			return err
		}

		return nil
	}

	supervisor.proc.Env = append(supervisor.proc.Env, supervisor.Service.Env...)

	// This is kind of a hacky way of doing this and we really should be using something like containerd/console for
	// creating a proper virtual console. To be replaced in the future
	// TODO: Replace with github.com/containerd/console
	if supervisor.Service.AllocateTty {
		supervisor.proc.Stdout = os.Stdout
		supervisor.proc.Stdin = os.Stdin
		supervisor.proc.Stderr = os.Stderr
	}

	if err := supervisor.proc.Start(); err != nil {
		return err
	}

	// timestamps always taken after the service starts successfully
	supervisor.Service.StartedAt = timestamppb.Now()
	supervisor.Service.UpdatedAt = timestamppb.Now()

	// migrate the service to its cgroup
	if err := supervisor.cgroup.AddProc(uint64(supervisor.proc.Process.Pid)); err != nil {
		return err
	}

	return nil
}

// New Initializes a new Service (and its cgroup) and returns a pointer to it
func New(service *v1.Service, parent *cgroup2.Manager) (*Supervisor, error) {
	cgroup, err := parent.NewChild(service.Name+".service", new(cgroup2.Resources))
	if err != nil {
		return nil, err
	}

	return &Supervisor{
		Service: service,
		cgroup:  cgroup,
	}, nil
}

// FromFile Initializes a new Service directly from its file definition
func FromFile(path string, parent *cgroup2.Manager) (*Supervisor, error) {
	if !strings.HasSuffix(path, ".json") {
		return nil, fmt.Errorf("service file must be a .json file")
	}

	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	data := make([]byte, stat.Size())
	n, err := file.Read(data)
	if err != nil {
		return nil, err
	}

	service := new(v1.Service)
	err = protojson.Unmarshal(data[:n], service)
	if err != nil {
		return nil, err
	}

	return New(service, parent)
}
