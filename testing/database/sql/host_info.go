package sql

import (
	"runtime"

	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
)

// getHostInfo tries to make both docker for mac and normal docker happy.  If it thinks the IP for the container is 0.0.0.0
// it'll return the "localhost" address.  Otherwise, it'll return the actual docker IP
func getHostInfo(resource *dockertest.Resource, service docker.Port) (res, port string) {
	// nolint: gocritic
	switch runtime.GOOS {
	// maybe we'll need this if someone wants to do windows development
	//nolint gocritic
	case "darwin":
		if runtime.GOOS == "darwin" {
			p := resource.Container.NetworkSettings.Ports[service]
			if p[0].HostIP == "0.0.0.0" {
				return "localhost", p[0].HostPort
			}
		}
		return resource.Container.NetworkSettings.IPAddress, service.Port()
	default:
		return resource.Container.NetworkSettings.IPAddress, service.Port()
	}
}
