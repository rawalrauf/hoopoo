// hoopoo:summary=Restart tmux if running with the latest configuration

package main

import (
  "os"
  "os/exec"
  "path/filepath"
)

func main() {
  // Check if tmux is running
  if exec.Command("pgrep", "-x", "tmux").Run() == nil {
    configPath := filepath.Join(os.Getenv("HOME"), ".config/tmux/tmux.conf")
    exec.Command("tmux", "source-file", configPath).Run()
  }
}
