// hoopoo:summary=Detect whether the computer is a Lenovo ThinkPad.

package main

import (
  "os"

  "hoopoo/pkg/hardware"
)

func main() {
  if hardware.IsLenovoThinkPad() {
    os.Exit(0)
  }
  os.Exit(1)
}
