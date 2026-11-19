// hoopoo:summary=Launch an app or focus an existing window matching a pattern
// hoopoo:args=<app-id> [launch-command]

package main

import (
  "os"
  "os/exec"
  "fmt"
  "syscall"

  "hoopoo/pkg/niri"
)

func main() {
  if len(os.Args) < 2 {
    fmt.Fprintln(os.Stderr, "Usage: hoopoo launch or focus <app-id> [launch-command]")
    os.Exit(1)
  }

  appID := os.Args[1]
  launchCommand := appID
  if len(os.Args) > 2 {
    launchCommand = os.Args[2]
  }

  windows, err := niri.GetWindows()
  if err == nil {
    for _, win := range windows {
      if win.AppID == appID {
        niri.FocusWindow(win.ID)
        return
      }
    }
  }

  // Use shell for flexible command parsing
  cmd := exec.Command("sh", "-c", launchCommand)
  cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
  cmd.Start()
}
