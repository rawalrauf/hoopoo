// hoopoo:summary=Set the system monospace font
// hoopoo:group=font
// hoopoo:name=set
// hoopoo:args=<font-name>
// hoopoo:examples=hoopoo font list | hoopoo font set "CaskaydiaMono Nerd Font"

package main

import (
  "fmt"
  "os"
  "os/exec"

  "hoopoo/pkg/font"
  "hoopoo/pkg/system/hooks"
  "hoopoo/pkg/system/process"
)

func main() {
  if len(os.Args) < 2 {
    fmt.Fprintln(os.Stderr, "Usage: hoopoo font set <font-name>")
    os.Exit(1)
  }

  fontName := os.Args[1]
  if err := font.Set(fontName); err != nil {
    fmt.Fprintf(os.Stderr, "Error: %v\n", err)
    os.Exit(1)
  }

  // Restart services after font change
  exec.Command("hoopoo", "restart", "waybar").Run()
  exec.Command("hoopoo", "restart", "swayosd").Run()

  // Notify for manual restarts
  if process.IsRunning("ghostty") {
    exec.Command("hoopoo", "notification", "send",
      "Restart Ghostty to see font change").Run()
  }

  if process.IsRunning("foot") {
    exec.Command("hoopoo", "notification", "send",
      "Restart Foot to see font change").Run()
  }

  // Run hook
  hooks.Run("font-set", fontName)
}
