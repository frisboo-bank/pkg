package contracts

type HealthEndpoint interface {
	WithEndpointPath(endpointPath string) HealthEndpoint
	WithStatusCodeDown(statusCodeDown int) HealthEndpoint
	WithStatusCodeUp(statusCodeUp int) HealthEndpoint

	RegisterEndpoints()
}
