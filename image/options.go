package image

import (
	"encoding/base64"
	"fmt"

	"github.com/riete/convert/str"

	"github.com/riete/docker/common/filter"

	"github.com/docker/docker/api/types"
)

type BuildOption func(*types.ImageBuildOptions)

func BuildWithImageName(repo, tag string) BuildOption {
	return func(o *types.ImageBuildOptions) {
		o.Tags = []string{repo + ":" + tag}
	}
}

func BuildWithNoCache() BuildOption {
	return func(o *types.ImageBuildOptions) {
		o.NoCache = true
	}
}

func BuildWithNetworkMode(mode string) BuildOption {
	return func(o *types.ImageBuildOptions) {
		o.NetworkMode = mode
	}
}

func BuildWithRemoveIntermediateContainers(force bool) BuildOption {
	return func(o *types.ImageBuildOptions) {
		o.Remove = true
		o.ForceRemove = force
	}
}

func BuildWithAlwaysPullParent() BuildOption {
	return func(o *types.ImageBuildOptions) {
		o.PullParent = true
	}
}

// BuildWithDockerfile dockerfile is a relative path
func BuildWithDockerfile(dockerfile string) BuildOption {
	return func(o *types.ImageBuildOptions) {
		o.Dockerfile = dockerfile
	}
}

func BuildWithArgs(args map[string]*string) BuildOption {
	return func(o *types.ImageBuildOptions) {
		o.BuildArgs = args
	}
}

func BuildWithLabels(labels map[string]string) BuildOption {
	return func(o *types.ImageBuildOptions) {
		o.Labels = labels
	}
}

func BuildWithQuiet() BuildOption {
	return func(o *types.ImageBuildOptions) {
		o.SuppressOutput = true
	}
}

type ListOption func(options *types.ImageListOptions)

func ListWithAll() ListOption {
	return func(o *types.ImageListOptions) {
		o.All = true
	}
}

func ListWithFilters(f map[string]string) ListOption {
	return func(o *types.ImageListOptions) {
		o.Filters = filter.NewFilterArgs(f)
	}
}

type AuthOption func(*types.ImagePullOptions)

func PullPushWithAuth(username, password string) AuthOption {
	return func(o *types.ImagePullOptions) {
		auth := fmt.Sprintf(`{"username":"%s","password":"%s"}`, username, password)
		o.RegistryAuth = base64.URLEncoding.EncodeToString(str.ToBytes(auth))
	}
}

func PullPushWithAuthFunc(f types.RequestPrivilegeFunc) AuthOption {
	return func(o *types.ImagePullOptions) {
		o.PrivilegeFunc = f
	}
}

type RemoveOption func(options *types.ImageRemoveOptions)

func RemoveWithForce() RemoveOption {
	return func(o *types.ImageRemoveOptions) {
		o.Force = true
	}
}

func RemoveWithPruneChildren() RemoveOption {
	return func(o *types.ImageRemoveOptions) {
		o.PruneChildren = true
	}
}

type PruneOption func(map[string]string)

func PruneWithAllUnused() PruneOption {
	return func(m map[string]string) {
		m["dangling"] = "false"
	}
}

func PruneWithFilters(f map[string]string) PruneOption {
	return func(m map[string]string) {
		for k, v := range f {
			m[k] = v
		}
	}
}
