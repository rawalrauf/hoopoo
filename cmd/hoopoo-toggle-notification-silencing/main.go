package main

import (
  "os/exec"
  "strings"
  
  "hoopoo/pkg/system/process"
)

// hoopoo:summary=Toggle notification do-not-disturb mode

func main() {
  // Toggle do-not-disturb mode
  exec.Command("makoctl", "mode", "-t", "do-not-disturb").Run()
  
  // Check current mode
  cmd := exec.Command("makoctl", "mode")
  output, err := cmd.Output()
  
  if err == nil && strings.Contains(string(output), "do-not-disturb") {
    // Do-not-disturb is ON (notifications silenced)
    exec.Command("hoopoo", "notification", "send",
      "--glyph=󰂛",
      "Silenced notifications",
      "-u", "low",
    ).Run()
  } else {
    // Do-not-disturb is OFF (notifications enabled)
    exec.Command("hoopoo", "notification", "send",
      "--glyph=󰂚",
      "Enabled notifications",
      "-u", "low",
    ).Run()
  }
  
  // Update waybar
  process.SendSignal("RTMIN+10", "waybar")
}
