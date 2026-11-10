// hoopoo:summary=Reload mako configuration (used by the Hoopoo theme switching).

package main

import (
  "fmt"
  "os"

  "hoopoo.net/internal/notification"
)

func main() {
  err := notification.ReloadMako()
  if err != nil {
    fmt.Fprintf(os.Stderr, "Error reloading mako: %v\n", err)
    os.Exit(1)
  }
}
