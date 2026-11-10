// hoopoo:summary=Clear the internal-monitor-disable toggle if no external display is connected.

package main

import (
  "os"

  "hoopoo.net/internal/common/toggle"
  "hoopoo.net/internal/hw"
)

func main() {
  if toggle.NiriEnabled("internal-monitor-off") && !hw.HasExternalMonitors() {
    toggle.NiriClear("internal-monitor-off")
  }
  os.Exit(0)
}
