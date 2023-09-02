package volume

import (
	"context"
	"errors"
	"fmt"

	"github.com/riete/convert/str"

	"github.com/docker/docker/api/types"

	volumetypes "github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/riete/docker/common/filter"
)

type VolumeClient struct {
	c *client.Client
}

func (v VolumeClient) List(options ...ListOption) (volumetypes.ListResponse, error) {
	o := volumetypes.ListOptions{}
	for _, option := range options {
		option(&o)
	}
	return v.c.VolumeList(context.Background(), o)
}

func (v VolumeClient) Inspect(volume string) (volumetypes.Volume, string, error) {
	r, b, err := v.c.VolumeInspectWithRaw(context.Background(), volume)
	return r, str.FromBytes(b), err
}

func (v VolumeClient) Create(volumeName string, options ...CreateOption) (volumetypes.Volume, error) {
	if _, _, err := v.Inspect(volumeName); err == nil {
		return volumetypes.Volume{}, errors.New(fmt.Sprintf(`Conflict: the volume name "%s" is already exists`, volumeName))
	}
	o := volumetypes.CreateOptions{Name: volumeName}
	for _, option := range options {
		option(&o)
	}
	return v.c.VolumeCreate(context.Background(), o)
}

func (v VolumeClient) Remove(volume string, force bool) error {
	return v.c.VolumeRemove(context.Background(), volume, force)
}

// Prune remove used volumes
func (v VolumeClient) Prune(options ...PruneOption) (types.VolumesPruneReport, error) {
	f := make(map[string]string)
	for _, option := range options {
		option(f)
	}
	return v.c.VolumesPrune(context.Background(), filter.NewFilterArgs(f))
}

func NewVolumeClient() (*VolumeClient, error) {
	var err error
	v := &VolumeClient{}
	v.c, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	return v, err
}
