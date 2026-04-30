package service

type Handler struct {
	// processServices Any basic process based services that should be started
	processServices []*ProcessService
}

// New Initializes a new service.Handler and returns a pointer to it
func New() *Handler {
	return &Handler{}
}
