// hoopoo:summary=Toggle Waybar position between top and bottom
// hoopoo:examples=hoopoo toggle waybar position

package main

import (
  "os"
  "os/exec"
  "path/filepath"
  "strings"

  "hoopoo/pkg/system/paths"
  "hoopoo/pkg/system/process"
  "hoopoo/pkg/toggle"
)

const positionToggleFile = "waybar-position-bottom"

func main() {
  configPath := filepath.Join(paths.UserConfig(), "waybar/config.jsonc")
  
  // Read config file
  content, err := os.ReadFile(configPath)
  if err != nil {
    os.Exit(1)
  }
  
  config := string(content)
  
  // Check current position (default is top if no toggle file)
  isBottom := toggle.Enabled(positionToggleFile)
  
  var newConfig string
  var message string
  
  if isBottom {
    // Switch to top (remove toggle file)
    newConfig = strings.ReplaceAll(config, `"position": "bottom"`, `"position": "top"`)
    toggle.Clear(positionToggleFile)
    message = "󰍜  Waybar moved to top"
  } else {
    // Switch to bottom (create toggle file)
    newConfig = strings.ReplaceAll(config, `"position": "top"`, `"position": "bottom"`)
    toggle.Write(positionToggleFile, "")
    message = "󰍜  Waybar moved to bottom"
  }
  
  // Write updated config
  if err := os.WriteFile(configPath, []byte(newConfig), 0644); err != nil {
    os.Exit(1)
  }
  
  // Restart waybar (reload config with SIGUSR2)
  process.SendSignal("SIGUSR2", "waybar")
  
  // Notify
  exec.Command("hoopoo", "notification", "send", message).Run()
}
