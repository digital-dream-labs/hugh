package s3

import (
	"io"
	"os"

	docker "github.com/fsouza/go-dockerclient"
)

var (
	Logger           io.Writer
	DefaultEndpoint  = "unix:///var/run/docker.sock"
	DefaultClient, _ = NewClient(DefaultEndpoint)
)

func init() {
	if host := os.Getenv("DOCKER_HOST"); host != "" {
		DefaultClient, _ = NewClientFromEnv()
	}
}

// Client wraps the go-dockerclient Client type to add some convenience methods.
//
// It embeds that type so all properties and methods are supported.
type Client struct {
	*DockerClient
}

// NewClient creates a new Client object, wrapping the native
// go-dockerclient Client. See also NewClientFromEnv.
func NewClient(endpoint string) (*Client, error) {
	dc, err := DockerNewClient(endpoint)
	if err != nil {
		return nil, err
	}
	return &Client{dc}, nil
}

// NewClientFromEnv initializes a new Docker client object from
// environment variables. This is the best way to initialize a client
// on OS X, where communication with the Docker API is via a
// self-signed HTTPS endpoint, not a UNIX socket.
func NewClientFromEnv() (*Client, error) {
	dc, err := DockerNewClientFromEnv()
	if err != nil {
		return nil, err
	}
	return &Client{dc}, nil
}

// PullPublicIfRequired checks to see if the local docker installation has an image available
// and if not attempts to fetch it from the public docker registry.
//
// It sends any logging information to the Logger defined in this package (if any).
func (client *Client) PullPublicIfRequired(imageName string) error {
	img, _ := client.InspectImage(imageName)
	if img != nil {
		// already got it
		return nil
	}
	return client.PullImage(DockerPullImageOptions{
		Repository:   imageName,
		OutputStream: Logger,
	}, DockerAuthConfiguration{})
}

// CreateContainer is a convenience wrapper around Client.CreateContainer() that
// returns a wrapped Container type.
func (client *Client) CreateContainer(opts DockerCreateContainerOptions) (*Container, error) {
	container, err := client.DockerClient.CreateContainer(opts)
	if err != nil {
		return nil, err
	}
	return &Container{container, client}, nil
}

// InspectContainer is a convience wrapp around Client.InspectContainer() that
// returns a wrapped Container type.
func (client *Client) InspectContainer(id string) (*Container, error) {
	container, err := client.DockerClient.InspectContainerWithOptions(
		docker.InspectContainerOptions{
			ID: id,
		},
	)
	if err != nil {
		return nil, err
	}
	return &Container{container, client}, nil
}

// StartNewContainer creates a new container, starts it and retrieves detailed information about it.
func (client *Client) StartNewContainer(containerOpts DockerCreateContainerOptions) (c *Container, err error) {
	c = &Container{client: client}
	c.DockerContainer, err = client.DockerClient.CreateContainer(containerOpts)
	if err != nil {
		return nil, err
	}
	if err := client.DockerClient.StartContainer(c.DockerContainer.ID, containerOpts.HostConfig); err != nil {
		return nil, err
	}
	c.DockerContainer, _ = client.DockerClient.InspectContainerWithOptions(
		docker.InspectContainerOptions{
			ID: c.DockerContainer.ID,
		},
	)
	return c, nil
}
