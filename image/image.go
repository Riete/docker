package image

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/riete/convert/str"

	"github.com/docker/docker/api/types/image"

	"github.com/riete/docker/common/filter"

	"github.com/riete/archive/tar"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type ImageClient struct {
	c *client.Client
}

func (i ImageClient) getImageByName(name string) (types.ImageSummary, error) {
	s, err := i.List(ListWithFilters(map[string]string{"reference": name}))
	if err != nil {
		return types.ImageSummary{}, err
	}
	if len(s) == 0 {
		return types.ImageSummary{}, errors.New("no such image: " + name)
	}
	return s[0], nil
}

func (i ImageClient) List(options ...ListOption) ([]types.ImageSummary, error) {
	o := types.ImageListOptions{}
	for _, option := range options {
		option(&o)
	}
	return i.c.ImageList(context.Background(), o)
}

// Inspect target can image name(repo:tag) or id
func (i ImageClient) Inspect(target string) (types.ImageInspect, string, error) {
	if strings.Contains(target, "/") {
		image, err := i.getImageByName(target)
		if err != nil {
			return types.ImageInspect{}, "", err
		}
		target = image.ID
	}
	inspect, b, err := i.c.ImageInspectWithRaw(context.Background(), target)
	return inspect, str.FromBytes(b), err
}

func (i ImageClient) Pull(ctx context.Context, image string, options ...AuthOption) (io.ReadCloser, error) {
	o := types.ImagePullOptions{}
	for _, option := range options {
		option(&o)
	}
	return i.c.ImagePull(ctx, image, o)
}

func (i ImageClient) Push(ctx context.Context, image string, options ...AuthOption) (io.ReadCloser, error) {
	o := types.ImagePullOptions{}
	for _, option := range options {
		option(&o)
	}
	return i.c.ImagePush(ctx, image, types.ImagePushOptions(o))
}

// Build path is directory containing a "Dockerfile" and other building related files
func (i ImageClient) Build(ctx context.Context, path string, options ...BuildOption) (io.ReadCloser, error) {
	o := types.ImageBuildOptions{}
	for _, option := range options {
		option(&o)
	}
	dockerfile := filepath.Join(path, o.Dockerfile)
	if o.Dockerfile == "" {
		dockerfile = filepath.Join(path, "Dockerfile")
	}
	if _, err := os.Stat(dockerfile); errors.Is(err, fs.ErrNotExist) {
		return nil, errors.New(fmt.Sprintf("%s is not exists", dockerfile))
	}

	// create build context
	tarFile := filepath.Join(os.TempDir(), fmt.Sprintf("%d.tar", time.Now().Unix()))
	packer := tar.NewTarPacker(tarFile, path, false, ".")
	if err := packer.Pack(); err != nil {
		return nil, err
	}
	defer os.Remove(tarFile)

	f, err := os.Open(tarFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	o.Context = f

	r, err := i.c.ImageBuild(ctx, f, o)
	if err != nil {
		return nil, err
	}
	return r.Body, nil
}

func (i ImageClient) Tag(src, tgt string) error {
	return i.c.ImageTag(context.Background(), src, tgt)
}

func (i ImageClient) Remove(target string, options ...RemoveOption) ([]types.ImageDeleteResponseItem, error) {
	if strings.Contains(target, "/") {
		image, err := i.getImageByName(target)
		if err != nil {
			return nil, err
		}
		target = image.ID
	}

	o := types.ImageRemoveOptions{}
	for _, option := range options {
		option(&o)
	}

	return i.c.ImageRemove(context.Background(), target, o)
}

// Prune remove unused image
func (i ImageClient) Prune(options ...PruneOption) (types.ImagesPruneReport, error) {
	f := make(map[string]string)
	for _, option := range options {
		option(f)
	}
	return i.c.ImagesPrune(context.Background(), filter.NewFilterArgs(f))
}

// Save save image as a tar file
func (i ImageClient) Save(image, saveTo string) error {
	r, err := i.c.ImageSave(context.Background(), []string{image})
	if err != nil {
		return err
	}
	defer r.Close()
	w, err := os.Create(saveTo)
	if err != nil {
		return err
	}
	_, err = io.Copy(w, r)
	return err
}

// Load load image form a tar file, ensure to close io.ReadCloser
func (i ImageClient) Load(loadFrom string) (io.ReadCloser, error) {
	f, err := os.Open(loadFrom)
	if err != nil {
		return nil, err
	}
	r, err := i.c.ImageLoad(context.Background(), f, false)
	if err != nil {
		return nil, err
	}
	return r.Body, err
}

func (i ImageClient) History(image string) ([]image.HistoryResponseItem, error) {
	return i.c.ImageHistory(context.Background(), image)
}

func NewImageClient() (*ImageClient, error) {
	var err error
	c := &ImageClient{}
	c.c, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	return c, err
}
