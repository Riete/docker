package network

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
	"github.com/riete/docker/common/filter"
)

type NetworkDriver string
type NetworkScope string

const (
	ScopeLocal    NetworkScope  = "local"
	ScopeSwarm    NetworkScope  = "swarm"
	ScopeGlobal   NetworkScope  = "global"
	DriverBridge  NetworkDriver = "bridge"
	DriverOverlay NetworkDriver = "overlay"
	DriverIpvlan  NetworkDriver = "ipvlan"
	DriverMacvlan NetworkDriver = "macvlan"
)

type ListOption func(*types.NetworkListOptions)

func ListWithFilters(f map[string]string) ListOption {
	return func(o *types.NetworkListOptions) {
		o.Filters = filter.NewFilterArgs(f)

	}
}

type InspectOption func(*types.NetworkInspectOptions)

func InspectWithVerbose() InspectOption {
	return func(o *types.NetworkInspectOptions) {
		o.Verbose = true
	}
}

// InspectWithScope Filter the network by scope (swarm, global, or local)
func InspectWithScope(scope NetworkScope) InspectOption {
	return func(o *types.NetworkInspectOptions) {
		o.Scope = string(scope)
	}
}

type CreateOption func(*types.NetworkCreate)

func CreateWithDriver(driver NetworkDriver) CreateOption {
	return func(o *types.NetworkCreate) {
		o.Driver = string(driver)
	}
}

func CreateWithScope(scope NetworkScope) CreateOption {
	return func(o *types.NetworkCreate) {
		o.Scope = string(scope)
	}
}

// CreateWithIpNet
// gateway, ipv4 or ipv6 gateway for the master subnet, "" or 172.16.0.1
// subnet, in cidr format that represents a network segment, 172.16.0.0/24
// allocatableIpRange, allocate container ip from a sub-range, "" or 172.16.0.0/25
func CreateWithIpNet(gateway, subnet, allocatableIpRange string) CreateOption {
	return func(o *types.NetworkCreate) {
		o.IPAM = &network.IPAM{
			Config: []network.IPAMConfig{{Subnet: subnet, Gateway: gateway, IPRange: allocatableIpRange}},
		}
	}
}

func CreateWithEnableIPv6() CreateOption {
	return func(o *types.NetworkCreate) {
		o.EnableIPv6 = true
	}
}

func CreateWithInternal() CreateOption {
	return func(o *types.NetworkCreate) {
		o.Internal = true
	}
}

func CreateWithAttachable() CreateOption {
	return func(o *types.NetworkCreate) {
		o.Attachable = true
	}
}

func CreateWithIngress() CreateOption {
	return func(o *types.NetworkCreate) {
		o.Ingress = true
	}
}

func CreateWithOptions(options map[string]string) CreateOption {
	return func(o *types.NetworkCreate) {
		o.Options = options
	}
}

func CreateWithLabels(labels map[string]string) CreateOption {
	return func(o *types.NetworkCreate) {
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
