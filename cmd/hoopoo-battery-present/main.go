// hoopoo:summary=Returns true if a battery is present on the system

package main

import (
  "os"

  "hoopoo.net/internal/common/usage"
  "hoopoo.net/internal/sysfs/battery"
)

func main() {
  usage.Message("Returns true if a battery is present on the system")

  if battery.IsPresent() {
    os.Exit(0)
  }
  os.Exit(1)
}
