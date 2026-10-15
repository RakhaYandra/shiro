package cli

import "testing"

func TestDockerCmd_Exists(t *testing.T) {
	root := NewRootCmd()
	for _, c := range root.Commands() {
		if c.Name() == "docker" {
			return
		}
	}
	t.Fatal("docker subcommand not registered")
}

func TestDocker_DaemonDown(t *testing.T) {
	cmd := NewDockerCmd()
	cmd.SetArgs([]string{"--host", "/does-not-exist-shiro/docker.sock"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("want error for missing daemon")
	}
}
