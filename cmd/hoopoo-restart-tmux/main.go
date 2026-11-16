// hoopoo:summary=Restart tmux if running with the latest configuration

package main

import (
  "os/exec"
  "path/filepath"

  "hoopoo/pkg/system/paths"
)

func main() {
  // Check if tmux is running
  if exec.Command("pgrep", "-x", "tmux").Run() == nil {
    configPath := filepath.Join(paths.UserConfig(), "tmux/tmux.conf")
    exec.Command("tmux", "source-file", configPath).Run()
  }
}
