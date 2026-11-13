// hoopoo:summary=Detect whether Vulkan is available.

package main

import (
  "os"

  "hoopoo/pkg/hardware"
)

func main() {
  if hardware.HasVulkan() {
    os.Exit(0)
  }
  os.Exit(1)
}
