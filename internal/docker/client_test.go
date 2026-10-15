package docker

import (
	"net"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

const cannedDF = `{
  "Images": [{"Id": "sha256:abc", "RepoTags": ["app:latest"], "Size": 500, "SharedSize": 100, "Containers": 0}],
  "Containers": [{"Id": "c1", "Names": ["/web"], "Image": "app:latest", "State": "exited", "SizeRw": 50}],
  "Volumes": [{"Name": "dbdata", "UsageData": {"Size": 200, "RefCount": 0}}],
  "BuildCache": [{"Size": 1000, "InUse": false}]
}`

func fakeDaemon(t *testing.T) string {
	t.Helper()
	sock := filepath.Join(t.TempDir(), "docker.sock")
	ln, err := net.Listen("unix", sock)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { ln.Close() })
	go http.Serve(ln, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "read-only", http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path == "/v1.43/system/df" {
			w.Write([]byte(cannedDF))
			return
		}
		http.NotFound(w, r)
	}))
	return sock
}

func TestClient_GetDF(t *testing.T) {
	c := NewClient(fakeDaemon(t))
	df, err := c.SystemDF()
	if err != nil {
		t.Fatal(err)
	}
	if len(df.Images) != 1 || df.Images[0].Size != 500 {
		t.Fatalf("bad df: %+v", df)
	}
}

func TestClient_DaemonDown(t *testing.T) {
	c := NewClient(filepath.Join(t.TempDir(), "nope.sock"))
	if _, err := c.SystemDF(); err == nil {
		t.Fatal("want error for missing daemon")
	} else if !IsUnavailable(err) {
		t.Fatalf("want unavailable error, got %v", err)
	}
}

func TestFakeDaemon_RejectsPOST(t *testing.T) {
	sock := fakeDaemon(t)
	tr := &http.Transport{DialContext: dialUnix(sock)}
	c := &http.Client{Transport: tr}
	req, _ := http.NewRequest(http.MethodPost, "http://d/v1.43/containers/x/stop", nil)
	resp, err := c.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want 405", resp.StatusCode)
	}
	_ = os.Getenv
}
