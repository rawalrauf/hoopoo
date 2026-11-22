package main

import (
  "os/exec"
  "path/filepath"
  
  "hoopoo/pkg/system/paths"
  "hoopoo/pkg/system/process"
)

// hoopoo:summary=Toggle swayidle idle locking

func main() {
  if process.IsRunning("swayidle") {
    // Stop idle locking
    process.KillByName("swayidle")
    exec.Command("hoopoo", "notification", "send",
      "--glyph=󱫖",
      "Stop locking computer when idle",
      "-u", "low",
    ).Run()
  } else {
    // Start idle locking
    configPath := filepath.Join(paths.UserConfig(), "swayidle", "swayidle.conf")
    cmd := exec.Command("swayidle", "-w", "-C", configPath)
    cmd.Start()
    
    exec.Command("hoopoo", "notification", "send",
      "--glyph=󱫖",
      "Now locking computer when idle",
      "-u", "low",
    ).Run()
  }
  
  // Update waybar
  process.SendSignal("RTMIN+9", "waybar")
}
