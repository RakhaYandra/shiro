package docker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

const apiVersion = "v1.43"

type UnavailableError struct{ Err error }

func (e *UnavailableError) Error() string { return fmt.Sprintf("docker unavailable: %v", e.Err) }
func (e *UnavailableError) Unwrap() error { return e.Err }

func IsUnavailable(err error) bool {
	var u *UnavailableError
	return errors.As(err, &u)
}

type Client struct {
	sock string
	http *http.Client
}

func DefaultSocket() string {
	if h := os.Getenv("DOCKER_HOST"); strings.HasPrefix(h, "unix://") {
		return strings.TrimPrefix(h, "unix://")
	}
	return "/var/run/docker.sock"
}

func NewClient(sock string) *Client {
	return &Client{
		sock: sock,
		http: &http.Client{
			Transport: &http.Transport{DialContext: dialUnix(sock)},
			Timeout:   15 * time.Second,
		},
	}
}

func dialUnix(sock string) func(context.Context, string, string) (net.Conn, error) {
	return func(ctx context.Context, _, _ string) (net.Conn, error) {
		return (&net.Dialer{}).DialContext(ctx, "unix", sock)
	}
}

func (c *Client) get(path string, out interface{}) error {
	req, err := http.NewRequest(http.MethodGet, "http://docker/"+apiVersion+path, nil)
	if err != nil {
		return err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return &UnavailableError{err}
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return &UnavailableError{fmt.Errorf("endpoint %s not found (daemon too old?)", path)}
	}
	if resp.StatusCode != http.StatusOK {
		return &UnavailableError{fmt.Errorf("daemon status %s", resp.Status)}
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

type DiskUsage struct {
	Images     []ImageSummary      `json:"Images"`
	Containers []ContainerSummary  `json:"Containers"`
	Volumes    []VolumeSummary     `json:"Volumes"`
	BuildCache []BuildCacheSummary `json:"BuildCache"`
}

type ImageSummary struct {
	ID         string   `json:"Id"`
	RepoTags   []string `json:"RepoTags"`
	Size       int64    `json:"Size"`
	SharedSize int64    `json:"SharedSize"`
	Containers int64    `json:"Containers"`
}

type ContainerSummary struct {
	ID     string   `json:"Id"`
	Names  []string `json:"Names"`
	Image  string   `json:"Image"`
	State  string   `json:"State"`
	SizeRw int64    `json:"SizeRw"`
}

type VolumeSummary struct {
	Name      string `json:"Name"`
	UsageData struct {
		Size     int64 `json:"Size"`
		RefCount int64 `json:"RefCount"`
	} `json:"UsageData"`
}

type BuildCacheSummary struct {
	Size  int64 `json:"Size"`
	InUse bool  `json:"InUse"`
}

func (c *Client) SystemDF() (*DiskUsage, error) {
	var df DiskUsage
	if err := c.get("/system/df", &df); err != nil {
		return nil, err
	}
	return &df, nil
}
