package output

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/RakhaYandra/shiro/internal/docker"
)

func RenderDockerText(rep docker.Report) string {
	var b strings.Builder
	b.WriteString("Shiro · Docker (read-only, nothing pruned)\n\n")
	fmt.Fprintf(&b, "Total %s   reclaimable %s (no prune executed)\n",
		Human(rep.TotalBytes), Human(rep.Reclaimable))
	b.WriteString("\nImages\n")
	for _, im := range rep.ImageItems {
		flag := ""
		if im.Reclaimable {
			flag = "  [reclaimable]"
		}
		fmt.Fprintf(&b, "  %8s  %s%s\n", Human(im.Size), im.Name, flag)
	}
	b.WriteString("\nContainers\n")
	for _, c := range rep.ContainerItems {
		flag := ""
		if c.Reclaimable {
			flag = "  [reclaimable]"
		}
		fmt.Fprintf(&b, "  %8s  %-20s %s%s\n", Human(c.Size), c.Name, c.Image, flag)
	}
	b.WriteString("\nVolumes (volumes may hold databases — review before any action)\n")
	for _, v := range rep.VolumeItems {
		flag := ""
		if v.Reclaimable {
			flag = "  [reclaimable]"
		}
		fmt.Fprintf(&b, "  %8s  %-20s [%s]%s\n", Human(v.Size), v.Name, v.Risk, flag)
	}
	fmt.Fprintf(&b, "\nBuild cache  %s (%s reclaimable)\n",
		Human(rep.CacheBytes), Human(rep.CacheReclaim))
	return b.String()
}

func RenderDockerJSON(rep docker.Report) []byte {
	raw, _ := json.MarshalIndent(map[string]any{
		"schema_version": 1,
		"images":         rep.ImageItems,
		"containers":     rep.ContainerItems,
		"volumes":        rep.VolumeItems,
		"cache_bytes":    rep.CacheBytes,
		"total_bytes":    rep.TotalBytes,
		"reclaimable":    rep.Reclaimable,
	}, "", "  ")
	return append(raw, '\n')
}
