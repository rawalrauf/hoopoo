// hoopoo:summary=Detect whether Vulkan is available.

package main

import (
  "os"

  "hoopoo.net/internal/hw"
)

func main() {
  if hw.HasVulkan() {
    os.Exit(0)
  }
  os.Exit(1)
}
