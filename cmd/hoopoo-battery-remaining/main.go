package main

import (
  "fmt"
  "os"

  "github.com/rawalrauf/hoopoo/internal/sysfs"
)

func main() {
  capacity, err := sysfs.GetCombinedCapacity()
  if err != nil {
    os.Exit(1)
  }

  fmt.Println(capacity)
}
