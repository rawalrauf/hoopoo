// hoopoo:summary=Returns the battery percentage remaining as an integer

package main

import (
  "fmt"
  "os"

  "github.com/rawalrauf/hoopoo/internal/usage"
  "github.com/rawalrauf/hoopoo/internal/sysfs/battery"
)

func main() {
  usage.Message("Returns the battery percentage remaining as an integer")

  capacity, err := battery.CombinedCapacity()
  if err != nil {
    os.Exit(1)
  }

  fmt.Println(capacity)
}
