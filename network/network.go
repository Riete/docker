package network

import (
	"context"

	"github.com/riete/convert/str"

	"github.com/riete/docker/common/filter"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type NetworkClient struct {
	c *client.Client
}

func (n NetworkClient) List(options ...ListOption) ([]types.NetworkResource, error) {
	o := types.NetworkListOptions{}
	for _, option := range options {
		option(&o)
	}
	return n.c.NetworkList(context.Background(), o)
}

// Inspect target can network name or id
func (n NetworkClient) Inspect(target string, options ...InspectOption) (types.NetworkResource, string, error) {
	o := types.NetworkInspectOptions{}
	for _, option := range options {
		option(&o)
	}
	r, b, err := n.c.NetworkInspectWithRaw(context.Background(), target, o)
	return r, str.FromBytes(b), err
}

func (n NetworkClient) Create(name string, options ...CreateOption) (types.NetworkCreateResponse, error) {
	o := types.NetworkCreate{CheckDuplicate: false}
	for _, option := range options {
		option(&o)
	}
	return n.c.NetworkCreate(context.Background(), name, o)
}

func (n NetworkClient) Remove(target string) error {
	return n.c.NetworkRemove(context.Background(), target)
}

// Prune remove unused network
func (n NetworkClient) Prune(options ...PruneOption) (types.NetworksPruneReport, error) {
	f := make(map[string]string)
	for _, option := range options {
		option(f)
	}
	return n.c.NetworksPrune(context.Background(), filter.NewFilterArgs(f))
}

func NewNetworkClient() (*NetworkClient, error) {
	var err error
	n := &NetworkClient{}
	n.c, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	return n, err
}
