package docker

import "strings"

type ImageItem struct {
	Name        string
	Size        int64
	Shared      int64
	Dangling    bool
	Reclaimable bool
}

type ContainerItem struct {
	Name        string
	Image       string
	Size        int64
	Running     bool
	Reclaimable bool
}

type VolumeItem struct {
	Name        string
	Size        int64
	InUse       bool
	Risk        string
	Reclaimable bool
}

type Report struct {
	ImageItems     []ImageItem
	ContainerItems []ContainerItem
	VolumeItems    []VolumeItem
	CacheBytes     int64
	CacheReclaim   int64
	TotalBytes     int64
	Reclaimable    int64
	Running        int
	Stopped        int
}

func shortName(names []string) string {
	if len(names) == 0 {
		return "<none>"
	}
	return strings.TrimPrefix(names[0], "/")
}

func Summarize(df *DiskUsage) Report {
	var rep Report
	for _, im := range df.Images {
		name := "<dangling>"
		if len(im.RepoTags) > 0 && im.RepoTags[0] != "<none>:<none>" {
			name = im.RepoTags[0]
		}
		dangling := name == "<dangling>"
		rec := dangling && im.Containers == 0
		rep.ImageItems = append(rep.ImageItems, ImageItem{name, im.Size, im.SharedSize, dangling, rec})
		rep.TotalBytes += im.Size
		if rec {
			rep.Reclaimable += im.Size - im.SharedSize
		}
	}
	for _, c := range df.Containers {
		running := strings.EqualFold(c.State, "running")
		if running {
			rep.Running++
		} else {
			rep.Stopped++
		}
		rec := !running
		rep.ContainerItems = append(rep.ContainerItems, ContainerItem{shortName(c.Names), c.Image, c.SizeRw, running, rec})
		if rec {
			rep.Reclaimable += c.SizeRw
		}
	}
	for _, v := range df.Volumes {
		inUse := v.UsageData.RefCount > 0
		risk := "review"
		if inUse {
			risk = "high"
		}
		rep.VolumeItems = append(rep.VolumeItems, VolumeItem{v.Name, v.UsageData.Size, inUse, risk, !inUse})
		rep.TotalBytes += v.UsageData.Size
		if !inUse {
			rep.Reclaimable += v.UsageData.Size
		}
	}
	for _, b := range df.BuildCache {
		rep.CacheBytes += b.Size
		if !b.InUse {
			rep.CacheReclaim += b.Size
		}
	}
	rep.TotalBytes += rep.CacheBytes
	rep.Reclaimable += rep.CacheReclaim
	return rep
}
