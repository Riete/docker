package system

import (
	"encoding/json"
	"strings"

	"github.com/riete/convert/size"
	"github.com/riete/convert/str"

	"github.com/docker/docker/api/types"
)

type ImageUsage struct {
	ImageName  string `json:"image_name"`
	Size       string `json:"size"`
	SharedSize string `json:"shared_size"`
	Containers int64  `json:"containers"`
}

type ContainerUsage struct {
	ContainerName string `json:"container_name"`
	ImageName     string `json:"image_name"`
	Command       string `json:"command"`
	LocalVolumes  int64  `json:"local_volumes"`
	Size          string `json:"size"`
	Status        string `json:"status"`
}

type LocalVolumeUsage struct {
	VolumeName string `json:"volume_name"`
	Links      int64  `json:"links"`
	Size       string `json:"size"`
}

type DiskUsageSummary struct {
	Images       []ImageUsage       `json:"images"`
	Containers   []ContainerUsage   `json:"containers"`
	LocalVolumes []LocalVolumeUsage `json:"local_volumes"`
	diskUsage    *types.DiskUsage
}

func (d DiskUsageSummary) ToString() string {
	b, _ := json.Marshal(d)
	return str.FromBytes(b)
}

func (d *DiskUsageSummary) imageUsage() {
	for _, i := range d.diskUsage.Images {
		if len(i.RepoTags) > 0 {
			for _, imageName := range i.RepoTags {
				d.Images = append(
					d.Images,
					ImageUsage{
						ImageName:  imageName,
						Size:       size.ToHumanSize(i.Size),
						SharedSize: size.ToHumanSize(i.SharedSize),
						Containers: i.Containers,
					},
				)
			}
		} else {
			d.Images = append(
				d.Images,
				ImageUsage{
					ImageName:  "<none>:<none>",
					Size:       size.ToHumanSize(i.Size),
					SharedSize: size.ToHumanSize(i.SharedSize),
					Containers: i.Containers,
				},
			)
		}
	}
}

func (d *DiskUsageSummary) localVolumeUsage() {
	for _, i := range d.diskUsage.Volumes {
		d.LocalVolumes = append(
			d.LocalVolumes,
			LocalVolumeUsage{
				VolumeName: i.Name,
				Links:      i.UsageData.RefCount,
				Size:       size.ToHumanSize(i.UsageData.Size),
			},
		)
	}
}

func (d *DiskUsageSummary) containerUsage() {
	for _, i := range d.diskUsage.Containers {
		c := ContainerUsage{
			ContainerName: strings.TrimPrefix(i.Names[0], "/"),
			ImageName:     i.Image,
			Command:       i.Command,
			Size:          size.ToHumanSize(i.SizeRw),
			Status:        i.Status,
		}
		for _, mount := range i.Mounts {
			if mount.Type == "volume" {
				c.LocalVolumes += 1
			}
		}
		d.Containers = append(d.Containers, c)
	}
}

func (d *DiskUsageSummary) Usage() {
	d.imageUsage()
	d.containerUsage()
	d.localVolumeUsage()
}

func NewDiskUsageHuman(u *types.DiskUsage) DiskUsageSummary {
	return DiskUsageSummary{diskUsage: u}
}
