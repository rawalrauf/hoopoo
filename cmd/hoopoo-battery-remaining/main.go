// hoopoo:summary=Returns the battery percentage remaining as an integer

package main

import (
  "fmt"
  "os"

  "hoopoo.net/internal/common/usage"
  "hoopoo.net/internal/sysfs/battery"
)

func main() {
  usage.Message("Returns the battery percentage remaining as an integer")

  capacity, err := battery.CombinedCapacity()
  if err != nil {
    os.Exit(1)
  }

  fmt.Println(capacity)
}
