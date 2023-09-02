package volume

import (
	volumetypes "github.com/docker/docker/api/types/volume"
	"github.com/riete/docker/common/filter"
)

type ListOption func(o *volumetypes.ListOptions)

func ListWithFilters(f map[string]string) ListOption {
	return func(o *volumetypes.ListOptions) {
		o.Filters = filter.NewFilterArgs(f)
	}
}

type CreateOption func(options *volumetypes.CreateOptions)

func CreateWithDriver(driver string) CreateOption {
	return func(o *volumetypes.CreateOptions) {
		o.Driver = driver
	}
}

func CreateWithDriverOpts(opts map[string]string) CreateOption {
	return func(o *volumetypes.CreateOptions) {
		o.DriverOpts = opts
	}
}

func CreateWithLabels(labels map[string]string) CreateOption {
	return func(o *volumetypes.CreateOptions) {
		o.Labels = labels
	}
}

type PruneOption func(map[string]string)

func PruneWithFilters(f map[string]string) PruneOption {
	return func(m map[string]string) {
		for k, v := range f {
			m[k] = v
		}
	}
}
