package container

import (
	"fmt"

	"github.com/riete/docker/common/filter"

	"github.com/docker/go-connections/nat"

	"github.com/docker/docker/api/types/container"

	"github.com/docker/docker/api/types"
)

type ListOption func(options *types.ContainerListOptions)

func ListWithAll() ListOption {
	return func(o *types.ContainerListOptions) {
		o.All = true
	}
}
func ListWithSize() ListOption {
	return func(o *types.ContainerListOptions) {
		o.Size = true
	}
}

func ListWithLatest() ListOption {
	return func(o *types.ContainerListOptions) {
		o.Latest = true
	}
}

func ListWithLimit(n int) ListOption {
	return func(o *types.ContainerListOptions) {
		o.Limit = n
	}
}

func ListWithFilters(f map[string]string) ListOption {
	return func(o *types.ContainerListOptions) {
		o.Filters = filter.NewFilterArgs(f)
	}
}

type TimeoutOption func(options *container.StopOptions)

func StopWithTimeout(t *int) TimeoutOption {
	return func(o *container.StopOptions) {
		o.Timeout = t
	}
}

func StopWithDefaultTimeout() TimeoutOption {
	return func(o *container.StopOptions) {
		o.Timeout = nil
	}
}

func StopWithNoTimeout() TimeoutOption {
	return func(o *container.StopOptions) {
		t := 0
		o.Timeout = &t
	}
}

type RemoveOption func(options *types.ContainerRemoveOptions)

func RemoveWithRemoveVolumes() RemoveOption {
	return func(o *types.ContainerRemoveOptions) {
		o.RemoveVolumes = true
	}
}

func RemoveWithRemoveLinks() RemoveOption {
	return func(o *types.ContainerRemoveOptions) {
		o.RemoveLinks = true
	}
}

func RemoveWithForce() RemoveOption {
	return func(o *types.ContainerRemoveOptions) {
		o.Force = true
	}
}

type CopyToOption func(options *types.CopyToContainerOptions)

func CopyToWithOverwriteDirWithFile() CopyToOption {
	return func(o *types.CopyToContainerOptions) {
		o.AllowOverwriteDirWithFile = true
	}
}

func CopyToWithCopyUidGid() CopyToOption {
	return func(o *types.CopyToContainerOptions) {
		o.CopyUIDGID = true
	}
}

type ExecConfigOptions func(*types.ExecConfig)

func ExecConfigWithUser(user string) ExecConfigOptions {
	return func(o *types.ExecConfig) {
		o.User = user
	}
}

func ExecConfigWithPrivileged() ExecConfigOptions {
	return func(o *types.ExecConfig) {
		o.Privileged = true
	}
}

func ExecConfigWithTty() ExecConfigOptions {
	return func(o *types.ExecConfig) {
		o.Tty = true
	}
}

func ExecConfigWithDetach() ExecConfigOptions {
	return func(o *types.ExecConfig) {
		o.Detach = true
	}
}

func ExecConfigWithEnv(env map[string]string) ExecConfigOptions {
	return func(o *types.ExecConfig) {
		for k, v := range env {
			o.Env = append(o.Env, fmt.Sprintf("%s=%s", k, v))
		}
	}
}

func ExecConfigWithWorkingDir(d string) ExecConfigOptions {
	return func(o *types.ExecConfig) {
		o.WorkingDir = d
	}
}

func ExecConfigWithCmd(cmd []string) ExecConfigOptions {
	return func(o *types.ExecConfig) {
		o.Cmd = cmd
	}
}

type ExecStartOption func(check *types.ExecStartCheck)

func ExecStartWithDetach() ExecStartOption {
	return func(o *types.ExecStartCheck) {
		o.Detach = true
	}
}

func ExecStartWithTty() ExecStartOption {
	return func(o *types.ExecStartCheck) {
		o.Tty = true
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

type CommitOption func(options *types.ContainerCommitOptions)

func CommitWithAuthor(author string) CommitOption {
	return func(o *types.ContainerCommitOptions) {
		o.Author = author
	}
}

func CommitWithComment(comment string) CommitOption {
	return func(o *types.ContainerCommitOptions) {
		o.Comment = comment
	}
}

// CommitWithChanges Dockerfile instructions to apply while committing, i.e. "CMD echo"
func CommitWithChanges(instructions []string) CommitOption {
	return func(o *types.ContainerCommitOptions) {
		o.Changes = instructions
	}
}

func CommitWithUnpause() CommitOption {
	return func(o *types.ContainerCommitOptions) {
		o.Pause = false
	}
}

type LogsOption func(options *types.ContainerLogsOptions)

// LogsWithTail  number of lines to show from the end of the logs (default "all")
func LogsWithTail(n string) LogsOption {
	return func(o *types.ContainerLogsOptions) {
		o.Tail = n
	}
}

func LogsWithFollow() LogsOption {
	return func(o *types.ContainerLogsOptions) {
		o.Follow = true
	}
}

// LogsWithSince show logs since utc timestamp (e.g. 2013-01-02T13:23:37) or relative (e.g. 42m for 42 minutes)
func LogsWithSince(since string) LogsOption {
	return func(o *types.ContainerLogsOptions) {
		o.Since = since
	}
}

// LogsWithUntil show logs since utc timestamp (e.g. 2013-01-02T13:23:37) or relative (e.g. 42m for 42 minutes)
func LogsWithUntil(until string) LogsOption {
	return func(o *types.ContainerLogsOptions) {
		o.Until = until
	}
}

func LogsWithTimestamps() LogsOption {
	return func(o *types.ContainerLogsOptions) {
		o.Timestamps = true
	}
}

func LogsWithDetails() LogsOption {
	return func(o *types.ContainerLogsOptions) {
		o.Details = true
	}
}

type CreateOption func(config *ContainerCreateConfig)

func CreateWithHostname(hostname string) CreateOption {
	return func(o *ContainerCreateConfig) {
		o.Config.Hostname = hostname
	}
}

func CreateWithAttachStdin() CreateOption {
	return func(o *ContainerCreateConfig) {
		o.Config.AttachStdin = true
	}
}

func CreateWithAttachStdout() CreateOption {
	return func(o *ContainerCreateConfig) {
		o.Config.AttachStdout = true
	}
}

func CreateWithAttachStderr() CreateOption {
	return func(o *ContainerCreateConfig) {
		o.Config.AttachStderr = true
	}
}

func CreateWithTty() CreateOption {
	return func(o *ContainerCreateConfig) {
		o.Config.Tty = true
	}
}

func CreateWithUser(user string) CreateOption {
	return func(o *ContainerCreateConfig) {
		o.Config.User = user
	}
}

func CreateWithEnvMap(env map[string]string) CreateOption {
	return func(o *ContainerCreateConfig) {
		for k, v := range env {
			o.Config.Env = append(o.Config.Env, fmt.Sprintf("%s=%s", k, v))
		}
	}
}

// CreateWithEnvArray env string format is "key=value"
func CreateWithEnvArray(env []string) CreateOption {
	return func(o *ContainerCreateConfig) {
		o.Config.Env = env
	}
}

func CreateWithCmd(cmd []string) CreateOption {
	return func(o *ContainerCreateConfig) {
		o.Config.Cmd = cmd
	}
}

func CreateWithWorkingDir(d string) CreateOption {
	return func(o *ContainerCreateConfig) {
		o.Config.WorkingDir = d
	}
}

func CreateWithEntrypoint(entrypoint []string) CreateOption {
	return func(o *ContainerCreateConfig) {
		o.Config.Entrypoint = entrypoint
	}
}

func CreateWithLabels(labels map[string]string) CreateOption {
	return func(o *ContainerCreateConfig) {
		o.Config.Labels = labels
	}
}

func CreateWithStopTimeout(n int) CreateOption {
	return func(o *ContainerCreateConfig) {
		o.Config.StopTimeout = &n
	}
}

// CreateWithBindsMap binds key is host-src or volume, value is container-dest[:options]
// host-src, container-dest must be an absolute path
func CreateWithBindsMap(binds map[string]string) CreateOption {
	return func(o *ContainerCreateConfig) {
		for src, dest := range binds {
			o.HostConfig.Binds = append(o.HostConfig.Binds, fmt.Sprintf("%s:%s", src, dest))
		}
	}
}

// CreateWithBindsArray bind string value is host:container[:options]
func CreateWithBindsArray(binds []string) CreateOption {
	return func(o *ContainerCreateConfig) {
		o.HostConfig.Binds = binds
	}
}

func CreateWithNetworkMode(mode string) CreateOption {
	return func(o *ContainerCreateConfig) {
		o.HostConfig.NetworkMode = container.NetworkMode(mode)
	}
}

// CreateWithPortBindings port format is [ip:]hostPort:containerPort[/proto]
// user nat.ParsePortSpecs to get nat.PortMap
func CreateWithPortBindings(binds nat.PortMap) CreateOption {
	return func(o *ContainerCreateConfig) {
		o.HostConfig.PortBindings = binds
	}
}

// CreateWithRestartPolicy policy name can one of ["", "no", "always", "unless-stopped", "on-failure"]
// if policy is on-failure, MaximumRetryCount is the number of times to retry before giving up
// can use restart.AlwaysPolicy and etc., to get policy
func CreateWithRestartPolicy(policy container.RestartPolicy) CreateOption {
	return func(o *ContainerCreateConfig) {
		o.HostConfig.RestartPolicy = policy
	}
}

// CreateWithAutoRemove remove the container when it exits. has no effect if RestartPolicy is set.
func CreateWithAutoRemove() CreateOption {
	return func(o *ContainerCreateConfig) {
		o.HostConfig.AutoRemove = true
	}
}

func CreateWithPrivileged() CreateOption {
	return func(o *ContainerCreateConfig) {
		o.HostConfig.Privileged = true
	}
}

func CreateWithPublishAllPorts() CreateOption {
	return func(o *ContainerCreateConfig) {
		o.HostConfig.PublishAllPorts = true
	}
}

func CreateWithCpuNums(n float64) CreateOption {
	return func(o *ContainerCreateConfig) {
		o.HostConfig.NanoCPUs = int64(n * 1e9)
	}
}

func CreateWithMemoryLimit(n int64) CreateOption {
	return func(o *ContainerCreateConfig) {
		o.HostConfig.Memory = n
	}
}

func CreateWithPidMode(mode string) CreateOption {
	return func(o *ContainerCreateConfig) {
		o.HostConfig.PidMode = container.PidMode(mode)
	}
}
