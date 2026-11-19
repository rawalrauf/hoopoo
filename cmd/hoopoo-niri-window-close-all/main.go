// hoopoo:summary=Close all open windows

package main

import (
  "hoopoo/pkg/niri"
)

func main() {
  // Switch to workspace 1
  niri.FocusWorkspace(1)

  // Get window list
  windows, err := niri.GetWindows()
  if err != nil {
    return
  }

  // Close each window
  for _, w := range windows {
    niri.CloseWindow(w.ID)
  }
}
