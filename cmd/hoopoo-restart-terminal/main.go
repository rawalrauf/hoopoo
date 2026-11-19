// hoopoo:summary=Reload supported terminal emulators after config changes

package main

import (
  "os"
  "os/exec"
  "path/filepath"
  "time"

  "hoopoo/pkg/system/paths"
  "hoopoo/pkg/system/process"
)

func main() {
  // Touch alacritty config to trigger reload
  alacrittyConfig := filepath.Join(paths.UserConfig(), "alacritty/alacritty.toml")
  if _, err := os.Stat(alacrittyConfig); err == nil {
    now := time.Now()
    os.Chtimes(alacrittyConfig, now, now)
  }

  // Reload kitty if running
  if process.IsRunning("kitty") {
    exec.Command("killall", "-SIGUSR1", "kitty").Run()
  }

  // Reload ghostty if running
  if process.IsRunning("ghostty") {
    exec.Command("killall", "-SIGUSR2", "ghostty").Run()
  }
}
