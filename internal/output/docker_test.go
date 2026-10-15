package output

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/RakhaYandra/shiro/internal/docker"
)

func dockerReport() docker.Report {
	return docker.Report{
		ImageItems:     []docker.ImageItem{{Name: "app:latest", Size: 500}},
		ContainerItems: []docker.ContainerItem{{Name: "old", Size: 50, Reclaimable: true}},
		VolumeItems:    []docker.VolumeItem{{Name: "dbdata", Size: 200, InUse: true, Risk: "high"}},
		CacheBytes:     1000,
		TotalBytes:     1750,
		Reclaimable:    1050,
		Stopped:        1,
	}
}

func TestRenderDockerText(t *testing.T) {
	out := RenderDockerText(dockerReport())
	for _, want := range []string{"app:latest", "dbdata", "high", "volumes may hold databases"} {
		if !strings.Contains(out, want) {
			t.Fatalf("missing %q:\n%s", want, out)
		}
	}
}

func TestRenderDockerJSON(t *testing.T) {
	var v map[string]any
	if err := json.Unmarshal(RenderDockerJSON(dockerReport()), &v); err != nil {
		t.Fatal(err)
	}
	if v["schema_version"] != float64(1) || v["reclaimable"] != float64(1050) {
		t.Fatalf("bad doc: %v", v)
	}
}
