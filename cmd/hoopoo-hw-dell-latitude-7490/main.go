// hoopoo:summary=Detect whether the computer is a Dell Latitude 7490.

package main

import (
  "os"

  "hoopoo.net/internal/sysfs/hw"
)

func main() {
  if hw.IsDellLatitude7490() {
    os.Exit(0)
  }
  os.Exit(1)
}
