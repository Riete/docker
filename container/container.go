package container

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/riete/convert/str"

	"github.com/docker/docker/api/types/container"

	"github.com/riete/docker/common/filter"
	"github.com/riete/docker/common/reader"

	"github.com/riete/archive/tar"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type ContainerClient struct {
	c *client.Client
}

func (c ContainerClient) List(options ...ListOption) ([]types.Container, error) {
	o := types.ContainerListOptions{}
	for _, option := range options {
		option(&o)
	}
	return c.c.ContainerList(context.Background(), o)
}

// Inspect container can name or id
func (c ContainerClient) Inspect(container string) (types.ContainerJSON, string, error) {
	r, b, err := c.c.ContainerInspectWithRaw(context.Background(), container, false)
	return r, str.FromBytes(b), err
}

// Start container can name or id
func (c ContainerClient) Start(container string) error {
	return c.c.ContainerStart(context.Background(), container, types.ContainerStartOptions{})
}

// Stop target can name or id
func (c ContainerClient) Stop(target string, option TimeoutOption) error {
	o := container.StopOptions{}
	option(&o)
	return c.c.ContainerStop(context.Background(), target, o)
}

// Restart target can name or id
func (c ContainerClient) Restart(target string, option TimeoutOption) error {
	o := container.StopOptions{}
	option(&o)
	return c.c.ContainerRestart(context.Background(), target, o)
}

// Rename container can name or id
func (c ContainerClient) Rename(container, newName string) error {
	return c.c.ContainerRename(context.Background(), container, newName)
}

// Remove container can name or id
func (c ContainerClient) Remove(container string, options ...RemoveOption) error {
	o := types.ContainerRemoveOptions{}
	for _, option := range options {
		option(&o)
	}
	return c.c.ContainerRemove(context.Background(), container, o)
}

// Stats container can name or id
func (c ContainerClient) Stats(container string) (string, error) {
	r, err := c.c.ContainerStatsOneShot(context.Background(), container)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()
	b, err := io.ReadAll(r.Body)
	return str.FromBytes(b), err
}

// CopyFrom container can name or id, sourcePath is file path in container
// targetPath is path to save copied file, if unpack is false, save as targetPath/{sourcePath.PathStat.Name}.tar
// if unpack is true, will unpack items to targetPath
func (c ContainerClient) CopyFrom(container, sourcePath, targetPath string, unpack bool) error {
	r, s, err := c.c.CopyFromContainer(context.Background(), container, sourcePath)
	if err != nil {
		return err
	}
	defer r.Close()
	if unpack {
		t := tar.NewTarUnPackerFromReader(r, targetPath, false)
		if err = t.Unpack(); err != nil {
			return err
		}
		return nil
	}
	w, err := os.Create(filepath.Join(targetPath, s.Name+".tar"))
	if err != nil {
		return err
	}

	defer w.Close()
	_, err = io.Copy(w, r)
	return err
}

// CopyFromRaw return tar archived as io.ReadCloser
func (c ContainerClient) CopyFromRaw(container, path string) (io.ReadCloser, types.ContainerPathStat, error) {
	return c.c.CopyFromContainer(context.Background(), container, path)
}

// CopyTo sourcePath is file/folder path to be copied, container can name or id
// targetPath is a directory path in container(create if not exists)
// sourcePath first be archived as a tar file, then copy to container targetPath and extract it
func (c ContainerClient) CopyTo(sourcePath, container, targetPath string, options ...CopyToOption) error {
	o := types.CopyToContainerOptions{}
	for _, option := range options {
		option(&o)
	}
	s, err := os.Stat(sourcePath)
	if err != nil {
		return err
	}

	if tStat, found, err := c.PathStat(container, targetPath); !found {
		_, stderr, err := c.ExecOneShot(container, fmt.Sprintf("mkdir -p %s", targetPath))
		if err != nil {
			return err
		}
		if stderr != "" {
			return fmt.Errorf("create directory %s error: %s", targetPath, stderr)
		}
	} else if err != nil {
		return err
	} else if !tStat.Mode.IsDir() {
		return fmt.Errorf("%s is not a directory", targetPath)
	}

	tarFile := filepath.Join(os.TempDir(), fmt.Sprintf("%d.tar", time.Now().Unix()))
	cwd := filepath.Dir(sourcePath)
	source := filepath.Base(sourcePath)
	if s.IsDir() {
		cwd = sourcePath
		source = "."
	}
	t := tar.NewTarPacker(tarFile, cwd, false, source)
	if err := t.Pack(); err != nil {
		return err
	}
	defer os.Remove(tarFile)
	r, err := os.Open(tarFile)
	if err != nil {
		return err
	}
	defer r.Close()
	return c.c.CopyToContainer(context.Background(), container, targetPath, r, o)
}

func (c ContainerClient) execCreate(container string, options ...ExecConfigOptions) (string, error) {
	o := types.ExecConfig{AttachStdout: true, AttachStdin: true, AttachStderr: true, Cmd: []string{"bash"}}
	for _, option := range options {
		option(&o)
	}
	r, err := c.c.ContainerExecCreate(context.Background(), container, o)
	if err != nil {
		return "", err
	}
	return r.ID, nil
}

func (c ContainerClient) exec(execId string, options ...ExecStartOption) (types.HijackedResponse, error) {
	o := types.ExecStartCheck{}
	for _, option := range options {
		option(&o)
	}
	return c.c.ContainerExecAttach(context.Background(), execId, o)
}

// Exec like docker exec -it [container] bash, open an interactive connection, container can name or id
// default "shebang" is bash, use ExecConfigWithCmd() to overwrite it, i.e ExecConfigWithCmd([]string{"sh"})
// types.HijackedResponse.Conn.Write() send data
// types.HijackedResponse.Reader receive data
// call types.HijackedResponse.Close() to close connection
func (c ContainerClient) Exec(container string, options ...ExecConfigOptions) (string, types.HijackedResponse, error) {
	options = append(options, ExecConfigWithTty())
	execId, err := c.execCreate(container, options...)
	if err != nil {
		return execId, types.HijackedResponse{}, err
	}
	if r, err := c.exec(execId, ExecStartWithTty()); err != nil {
		return execId, types.HijackedResponse{}, err
	} else {
		return execId, r, nil
	}
}

func (c ContainerClient) ExecResizePty(execId string, height, width uint) error {
	o := types.ResizeOptions{Height: height, Width: width}
	return c.c.ContainerExecResize(context.Background(), execId, o)
}

func (c ContainerClient) execOneShot(container string, options ...ExecConfigOptions) (types.HijackedResponse, error) {
	execId, err := c.execCreate(container, options...)
	if err != nil {
		return types.HijackedResponse{}, err
	}
	return c.exec(execId)
}

// ExecOneShot run cmd once and return stdout and stderr, container can name or id
func (c ContainerClient) ExecOneShot(container string, cmd string, options ...ExecConfigOptions) (string, string, error) {
	if r, err := c.execOneShot(container, options...); err != nil {
		return "", "", err
	} else {
		defer r.Close()
		if !strings.HasSuffix(cmd, "\n") {
			cmd += "\n" + "exit\n"
		}
		if _, err := r.Conn.Write(str.ToBytes(cmd)); err != nil {
			return "", "", err
		}
		return reader.ParseToStdoutStderr(r.Reader)
	}
}

// ExecOneShotWithCombinedOutput run cmd once and return combined stdout and stderr
func (c ContainerClient) ExecOneShotWithCombinedOutput(container string, cmd string, options ...ExecConfigOptions) (string, error) {
	if r, err := c.execOneShot(container, options...); err != nil {
		return "", err
	} else {
		defer r.Close()
		if !strings.HasSuffix(cmd, "\n") {
			cmd += "\n" + "exit\n"
		}
		if _, err := r.Conn.Write(str.ToBytes(cmd)); err != nil {
			return "", err
		}
		return reader.ParseToCombinedOutput(r.Reader)
	}
}

// PathStat return path stat in target container
func (c ContainerClient) PathStat(container, path string) (types.ContainerPathStat, bool, error) {
	r, err := c.c.ContainerStatPath(context.Background(), container, path)
	return r, !client.IsErrNotFound(err), err
}

// Prune remove unused(not running) container
func (c ContainerClient) Prune(options ...PruneOption) (types.ContainersPruneReport, error) {
	f := make(map[string]string)
	for _, option := range options {
		option(f)
	}
	r, err := c.c.ContainersPrune(context.Background(), filter.NewFilterArgs(f))
	return r, err
}

// Commit create image from container，default is pause the container before committing
func (c ContainerClient) Commit(container, image string, options ...CommitOption) (string, error) {
	o := types.ContainerCommitOptions{Reference: image, Pause: true}
	for _, option := range options {
		option(&o)
	}
	r, err := c.c.ContainerCommit(context.Background(), container, o)
	return r.ID, err
}

// Export export container filesystem as a tar file
func (c ContainerClient) Export(container, path string) error {
	w, err := os.Create(path)
	if err != nil {
		return err
	}
	defer w.Close()
	r, err := c.c.ContainerExport(context.Background(), container)
	if err != nil {
		return err
	}
	defer r.Close()
	_, err = io.Copy(w, r)
	return err
}

// Kill send SIGKILL signal to container
func (c ContainerClient) Kill(container string) error {
	return c.c.ContainerKill(context.Background(), container, "SIGKILL")
}

// Terminate send SIGTERM signal to container
func (c ContainerClient) Terminate(container string) error {
	return c.c.ContainerKill(context.Background(), container, "SIGTERM")
}

// Logs return container logs as io.ReadCloser, can user reader.ParseToCombinedStreamOutput to get log message
func (c ContainerClient) Logs(container string, options ...LogsOption) (io.ReadCloser, error) {
	o := types.ContainerLogsOptions{ShowStderr: true, ShowStdout: true}
	for _, option := range options {
		option(&o)
	}
	return c.c.ContainerLogs(context.Background(), container, o)
}

func (c ContainerClient) Pause(container string) error {
	return c.c.ContainerPause(context.Background(), container)
}

func (c ContainerClient) Unpause(container string) error {
	return c.c.ContainerUnpause(context.Background(), container)
}

// Process processes info in container, ps -ef
func (c ContainerClient) Process(container string) (container.ContainerTopOKBody, error) {
	return c.c.ContainerTop(context.Background(), container, nil)
}

// Create container create, set replace to true to remove before create
func (c ContainerClient) Create(image, container string, replace bool, options ...CreateOption) (container.CreateResponse, error) {
	if replace {
		_ = c.Remove(container, RemoveWithForce())
	}

	o := NewContainerCreateConfig()
	for _, option := range options {
		option(o)
	}
	o.Config.Image = image
	return c.c.ContainerCreate(context.Background(), o.Config, o.HostConfig, o.NetworkConfig, o.Platform, container)
}

// Run create container and start it
func (c ContainerClient) Run(image, container string, replace bool, options ...CreateOption) (container.CreateResponse, error) {
	r, err := c.Create(image, container, replace, options...)
	if err != nil {
		return r, err
	}
	return r, c.Start(container)
}

func NewContainerClient() (*ContainerClient, error) {
	var err error
	c := &ContainerClient{}
	c.c, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	return c, err
}
