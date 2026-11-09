// hoopoo:summary=Launch an app or focus an existing window matching a pattern
// hoopoo:args=<app-id> [launch-command]

package main

import (
  "encoding/json"
  "os"
  "os/exec"
  "fmt"
)

type Window struct {
  ID    int    `json:"id"`
  AppID string `json:"app_id"`
}

func main() {
  if len(os.Args) < 2 {
    fmt.Fprintln(os.Stderr, "Usage: hoopoo-launch-or-focus [app-id] [launch-command]")
    os.Exit(1)
  }

  appID := os.Args[1]
  launchCommand := fmt.Sprintf("uwsm-app -- %s", appID)
  if len(os.Args) > 2 {
    launchCommand = os.Args[2]
  }

  cmd := exec.Command("niri", "msg", "-j", "windows")
  output, err := cmd.Output()
  if err != nil {
    os.Exit(1)
  }

  var windows []Window
  if err := json.Unmarshal(output, &windows); err != nil {
    os.Exit(1)
  }

  for _, win := range windows {
    if win.AppID == appID {
      exec.Command("niri", "msg", "action", "focus-window", "--id", fmt.Sprintf("%d", win.ID)).Run()
      return
    }
  }

  cmd = exec.Command("sh", "-c", fmt.Sprintf("exec setsid %s", launchCommand))
  cmd.Start()
}
