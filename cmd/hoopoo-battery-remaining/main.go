// hoopoo:summary=Returns the battery percentage remaining as an integer.

package main

import (
  "fmt"
  "os"

  "hoopoo.net/internal/sysfs/battery"
)

func main() {
  capacity, err := battery.CombinedCapacity()
  if err != nil {
    os.Exit(1)
  }

  fmt.Println(capacity)
}
