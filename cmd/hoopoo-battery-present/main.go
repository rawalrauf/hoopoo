// hoopoo:summary=Returns true if a battery is present on the system.

package main

import (
  "os"

  "hoopoo.net/internal/battery"
)

func main() {
  if battery.IsPresent() {
    os.Exit(0)
  }
  os.Exit(1)
}
