// hoopoo:summary=Detect whether the computer is an Asus ROG machine.

package main

import (
  "os"

  "hoopoo/pkg/hardware"
)

func main() {
  if hardware.IsAsusROG() {
    os.Exit(0)
  }
  os.Exit(1)
}
