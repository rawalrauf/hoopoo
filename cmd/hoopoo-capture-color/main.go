package main

import (
  "os"
  "os/exec"
  "strings"
)

// hoopoo:summary=Pick a color from the screen and copy hex to clipboard
// hoopoo:group=capture
// hoopoo:examples=hoopoo capture color

func main() {
  // Use niri's built-in color picker
  cmd := exec.Command("niri", "msg", "pick-color")
  output, err := cmd.Output()
  if err != nil {
    os.Exit(0) // User cancelled
  }
  
  // Parse hex color from output
  lines := strings.Split(string(output), "\n")
  var hex string
  
  for _, line := range lines {
    if strings.Contains(line, "Hex:") {
      parts := strings.Fields(line)
      if len(parts) >= 2 {
        hex = parts[1]
        break
      }
    }
  }
  
  if hex == "" {
    os.Exit(0)
  }
  
  // Copy to clipboard
  wlCopyCmd := exec.Command("wl-copy")
  wlCopyCmd.Stdin = strings.NewReader(hex)
  wlCopyCmd.Run()
  
  // Show notification
  exec.Command("hoopoo", "notification", "send",
    "--glyph=󰃉",
    "Color copied to clipboard",
    hex,
    "-u", "low",
  ).Run()
}
