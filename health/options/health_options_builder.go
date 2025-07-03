package options

import "frisboo-bank/pkg/health/contracts"

type HealthOption = func(options *HealthOptions)

// UseEndpointPath returns a HealthOption that sets the health endpoint path.
func UseEndpointPath(path string) HealthOption {
	return func(h *HealthOptions) {
		h.EndpointPath = path
	}
}

// UseStatusCodeUp returns a HealthOption that sets the status code for "up".
func UseStatusCodeUp(code uint) HealthOption {
	return func(h *HealthOptions) {
		h.StatusCodeUp = code
	}
}

// UseStatusUp returns a HealthOption that sets the status string for "up".
func UseStatusUp(status contracts.StatusType) HealthOption {
	return func(h *HealthOptions) {
		h.StatusUp = status
	}
}

// UseStatusCodeDown returns a HealthOption that sets the status code for "down".
func UseStatusCodeDown(code uint) HealthOption {
	return func(h *HealthOptions) {
		h.StatusCodeDown = code
	}
}

// UseStatusDown returns a HealthOption that sets the status string for "down".
func UseStatusDown(status contracts.StatusType) HealthOption {
	return func(h *HealthOptions) {
		h.StatusDown = status
	}
}

// UseServices returns a HealthOption that sets the health service checks.
func UseServices(services []contracts.HealthServiceCheck) HealthOption {
	return func(h *HealthOptions) {
		h.Services = services
	}
}
