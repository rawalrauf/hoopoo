// hoopoo:summary=Close all open windows

package main

import (
  "encoding/json"
  "os/exec"
  "strconv"
)

type Window struct {
  ID int `json:"id"`
}

func main() {
  // Switch to workspace 1
  exec.Command("niri", "msg", "action", "focus-workspace", "1").Run()

  // Get window list
  out, err := exec.Command("niri", "msg", "-j", "windows").Output()
  if err != nil {
    return
  }

  var windows []Window
  if json.Unmarshal(out, &windows) != nil {
    return
  }

  // Close each window
  for _, w := range windows {
    exec.Command("niri", "msg", "action", "close-window", "--id", strconv.Itoa(w.ID)).Run()
  }
}
