package main

import (
  "os"

  "github.com/rawalrauf/hoopoo/internal/sysfs"
)

func main() {
  if sysfs.IsPresent() {
    os.Exit(0)
  }
  os.Exit(1)
}
