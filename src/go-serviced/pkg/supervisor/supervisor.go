package supervisor

import (
	"os/exec"

	"github.com/containerd/cgroups/v2/cgroup2"
	v1 "github.com/immutable-os/immutable-os/src/go-serviced/proto/service/v1"
)

type Supervisor struct {
	// Service The service that this supervisor is managing
	Service *v1.Service

	// proc Stores process information about the service (if currently running)
	proc *exec.Cmd

	// cgroup The cgroup manager for this specific service. Usually a cgroup under the `serviced` slice
	cgroup *cgroup2.Manager
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
