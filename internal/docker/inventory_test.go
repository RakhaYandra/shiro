package docker

import "testing"

func canned() *DiskUsage {
	return &DiskUsage{
		Images: []ImageSummary{
			{ID: "sha256:used", RepoTags: []string{"app:latest"}, Size: 500, SharedSize: 100, Containers: 2},
			{ID: "sha256:dangle", Size: 300, Containers: 0},
		},
		Containers: []ContainerSummary{
			{ID: "c1", Names: []string{"/run"}, Image: "app:latest", State: "running", SizeRw: 10},
			{ID: "c2", Names: []string{"/old"}, Image: "app:latest", State: "exited", SizeRw: 50},
		},
		Volumes: []VolumeSummary{
			{Name: "dbdata", UsageData: struct {
				Size     int64 `json:"Size"`
				RefCount int64 `json:"RefCount"`
			}{Size: 200, RefCount: 1}},
			{Name: "orphan", UsageData: struct {
				Size     int64 `json:"Size"`
				RefCount int64 `json:"RefCount"`
			}{Size: 100, RefCount: 0}},
		},
		BuildCache: []BuildCacheSummary{{Size: 1000, InUse: false}, {Size: 200, InUse: true}},
	}
}

func TestInventory_Reclaimable(t *testing.T) {
	rep := Summarize(canned())
	// dangling 300 + stopped container 50 + unused volume 100 + unused cache 1000
	if rep.Reclaimable != 1450 {
		t.Fatalf("reclaimable = %d, want 1450", rep.Reclaimable)
	}
	if rep.Running != 1 || rep.Stopped != 1 {
		t.Fatalf("containers = %+v", rep)
	}
}

func TestInventory_VolumeRisk(t *testing.T) {
	rep := Summarize(canned())
	for _, v := range rep.VolumeItems {
		if v.Name == "dbdata" && v.Risk != "high" {
			t.Fatalf("in-use volume must be high: %+v", v)
		}
		if v.Name == "orphan" && v.Risk != "review" {
			t.Fatalf("unused volume must be review: %+v", v)
		}
	}
}

func TestInventory_RunningExcluded(t *testing.T) {
	rep := Summarize(canned())
	for _, c := range rep.ContainerItems {
		if c.Name == "run" && c.Reclaimable {
			t.Fatal("running container must not be reclaimable")
		}
	}
}
