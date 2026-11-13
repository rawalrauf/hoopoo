// hoopoo:summary=Clear the internal-monitor-disable toggle if no external display is connected.

package main

import (
  "os"

  "hoopoo/pkg/system/toggle"
  "hoopoo/pkg/hardware"
)

func main() {
  if toggle.NiriEnabled("internal-monitor-off") && !hardware.HasExternalMonitors() {
    toggle.NiriClear("internal-monitor-off")
  }
  os.Exit(0)
}
