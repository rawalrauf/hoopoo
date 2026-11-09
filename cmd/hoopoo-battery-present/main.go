// hoopoo:summary=Returns true if a battery is present on the system

package main

import (
  "os"

  "github.com/rawalrauf/hoopoo/internal/usage"
  "github.com/rawalrauf/hoopoo/internal/sysfs/battery"
)

func main() {
  usage.Message("Returns true if a battery is present on the system")

  if battery.IsPresent() {
    os.Exit(0)
  }
  os.Exit(1)
}
