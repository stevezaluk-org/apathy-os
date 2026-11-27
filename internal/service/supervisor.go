package service

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"

	"github.com/containerd/cgroups/v2/cgroup2"

	"golang.org/x/sys/unix"
)

type Supervisor struct {
	// service The service that this supervisor is controling
	service *Service

	// cgroup The services cgroup, branched from its parent
	cgroup *cgroup2.Manager

	// proc Process information about the running service
	proc *exec.Cmd

	// cancel Cancel function for the goroutine spawned by Supervisor.Start
	cancel context.CancelFunc
}

// PriorityGate Returns the priority of the service that this supervisor has been assigned
func (supervisor *Supervisor) PriorityGate() int {
	return supervisor.service.PriorityGate()
}

// Stop Stops a service provided that it was running
func (supervisor *Supervisor) Stop(ctx context.Context) {
	if supervisor.proc.Process == nil {
		return
	}

	// we have to cancel the context here before we signal since the goroutine started by Supervisor.Start
	// needs to know that it shouldn't try and restart the service. This alone does not actually stop the
	// service since supervisor.proc.Wait() blocks until it completes
	supervisor.cancel()

	// send a SIGTERM and try and wait for this process to exit gracefully
	if err := supervisor.proc.Process.Signal(unix.SIGTERM); err != nil {
		slog.Error("Failed to send SIGTERM to service", "svc", supervisor.service.Name, "err", err)
	}

	for {
		select {
		case <-ctx.Done():
			slog.Info("Service failed to respond to SIGTERM during the kill delay. Service will be SIGKILL'd", "svc", supervisor.service.Name)
			if err := supervisor.cgroup.Freeze(); err != nil {
				slog.Error("Failed to freeze cgroup for service during shutdown", "svc", supervisor.service.Name)
			}

			if err := supervisor.proc.Process.Kill(); err != nil {
				slog.Error("Failed to SIGKILL service", "svc", supervisor.service.Name, "err", err)
			}

			slog.Info("Service has been stopped (SIGKILL'd)", "svc", supervisor.service.Name, "exit_code", supervisor.proc.ProcessState.ExitCode())

			return
		default:
			// only ever populated when wait returns
			if supervisor.proc.ProcessState != nil {
				slog.Info("Service has been stopped (service responded to sigterm)", "svc", supervisor.service.Name, "exit_code", supervisor.proc.ProcessState.ExitCode())
				return
			}
		}
	}
}

// buildProc Builds a cmd.Execute for a given service and returns a pointer to it
func (supervisor *Supervisor) buildProc() *exec.Cmd {
	slog.Info("Attempting service start", "svc", supervisor.service.Name)

	proc := exec.Command(supervisor.service.Exec, supervisor.service.Args...)
	proc.Env = append(proc.Env, supervisor.service.Env...)

	// This is kind of a hacky way of doing this and we really should be using something like containerd/console for
	// creating a proper virtual console. To be replaced in the future
	// TODO: Replace with github.com/containerd/console
	if supervisor.service.AllocateTty {
		proc.Stdout = os.Stdout
		proc.Stdin = os.Stdin
		proc.Stderr = os.Stderr
	}

	return proc
}

// Start Starts the service and moves the service to a watch loop
func (supervisor *Supervisor) Start(ctx context.Context) {
	childCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	// needed by supervisor.Stop
	supervisor.cancel = cancel

	for i := range supervisor.service.MaxRestartCount {
		select {
		case <-childCtx.Done():
			slog.Info("Service stop requested (context cancelled)", "svc", supervisor.service.Name)
			return
		default:
			supervisor.proc = supervisor.buildProc()
			if err := supervisor.proc.Start(); err != nil {
				slog.Error("Failed to start service", "svc", supervisor.service.Name, "err", err)
				continue
			}

			slog.Info("Service has been started", "svc", supervisor.service.Name, "count", i)
			if err := supervisor.cgroup.AddProc(uint64(supervisor.proc.Process.Pid)); err != nil {
				slog.Error("Failed to migrate service to cgroup", "svc", supervisor.service.Name, "err", err)
				continue
			}

			// wait for the process to complete
			if err := supervisor.proc.Wait(); err != nil {
				slog.Info("Service exited with", "svc", supervisor.service.Name, "err", err)
			}
		}
	}
}

// NewSupervisor Initializes a new service supervisor and returns a pointer to it. Calling this function
// initializes the services cgroup under its parent
func NewSupervisor(service *Service, parent *cgroup2.Manager) (*Supervisor, error) {
	cgroup, err := parent.NewChild(service.Name+".service", new(cgroup2.Resources))
	if err != nil {
		return nil, fmt.Errorf("error during cgroup creation for service supervisor: %v", err)
	}

	return &Supervisor{
		service: service,
		cgroup:  cgroup,
	}, nil
}
