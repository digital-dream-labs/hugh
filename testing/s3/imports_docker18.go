package s3

import (
	docker "github.com/fsouza/go-dockerclient"
)

type DockerContainer = docker.Container
type DockerClient = docker.Client
type DockerAuthConfiguration = docker.AuthConfiguration
type DockerPullImageOptions = docker.PullImageOptions
type DockerCreateContainerOptions = docker.CreateContainerOptions
type DockerRemoveContainerOptions = docker.RemoveContainerOptions
type DockerPort = docker.Port
type DockerConfig = docker.Config
type DockerHostConfig = docker.HostConfig
type DockerBuildImageOptions = docker.BuildImageOptions
type DockerHostPortBinding = docker.PortBinding

func DockerNewClient(endpoint string) (*DockerClient, error) {
	return docker.NewClient(endpoint)
}

func DockerNewClientFromEnv() (*DockerClient, error) {
	return docker.NewClientFromEnv()
}
