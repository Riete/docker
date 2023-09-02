package container

import (
	"github.com/docker/docker/api/types/network"

	"github.com/docker/docker/api/types/container"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
)

type ContainerCreateConfig struct {
	Config        *container.Config
	HostConfig    *container.HostConfig
	NetworkConfig *network.NetworkingConfig
	Platform      *specs.Platform
}

func NewContainerCreateConfig() *ContainerCreateConfig {
	return &ContainerCreateConfig{
		Config:        &container.Config{},
		HostConfig:    &container.HostConfig{},
		NetworkConfig: &network.NetworkingConfig{},
		Platform:      nil,
	}
}
