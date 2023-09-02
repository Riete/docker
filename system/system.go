package system

import (
	"context"
	"errors"
	"fmt"
	"os/exec"

	cmd "github.com/riete/exec"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/client"
)

type SystemClient struct {
	c *client.Client
}

// RegistryLogin validate credentials for a registry, "Login Succeeded" status for success
// note this function actually do not make docker daemon login to the registry
// User Login to login to the registry
func (s SystemClient) RegistryLogin(addr, username, password string) (registry.AuthenticateOKBody, error) {
	return s.c.RegistryLogin(
		context.Background(),
		registry.AuthConfig{
			Username:      username,
			Password:      password,
			ServerAddress: addr,
		},
	)
}

func (s SystemClient) Login(addr, username, password string) (string, error) {
	path, err := exec.LookPath("docker")
	if err != nil {
		return "", errors.New("docker command is not found")
	}
	r := cmd.NewCmdRunner(
		path,
		"login",
		fmt.Sprintf("--username=%s", username),
		fmt.Sprintf("--password=%s", password),
		addr,
	)
	return r.RunWithCombinedOutput()
}

func (s SystemClient) Ping() (types.Ping, error) {
	return s.PingContext(context.Background())
}

func (s SystemClient) PingContext(ctx context.Context) (types.Ping, error) {
	return s.c.Ping(ctx)
}

func (s SystemClient) Info() (types.Info, error) {
	return s.InfoContex(context.Background())
}

func (s SystemClient) InfoContex(ctx context.Context) (types.Info, error) {
	return s.c.Info(ctx)
}

// DiskUsage types.DiskUsage is original data, DiskUsageSummary show images, containers and local volumes usage
func (s SystemClient) DiskUsage() (types.DiskUsage, DiskUsageSummary, error) {
	r, err := s.c.DiskUsage(context.Background(), types.DiskUsageOptions{})
	if err != nil {
		return r, DiskUsageSummary{}, err
	}
	h := NewDiskUsageHuman(&r)
	h.Usage()
	return r, h, nil
}

func NewSystemClient() (*SystemClient, error) {
	var err error
	s := &SystemClient{}
	s.c, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	return s, err
}
