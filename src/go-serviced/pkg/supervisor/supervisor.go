package supervisor

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/containerd/cgroups/v2/cgroup2"
	v1 "github.com/immutable-os/immutable-os/src/go-serviced/proto/service/v1"
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

// Start Starts the respective service.
func (supervisor *Supervisor) Start(ctx context.Context) error {
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
