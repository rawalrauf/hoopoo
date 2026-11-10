// hoopoo:summary=Dismiss a mako notification on the basis of its summary. Used by the first-run notifications to dismiss them after clicking for action.
// hoopoo:args=<summary>

package main

import (
  "fmt"
  "os"

  "hoopoo.net/internal/notification"
)

func main() {
  if len(os.Args) < 2 {
    fmt.Fprintln(os.Stderr, "Usage: hoopoo-notification-dismiss <summary>")
    os.Exit(1)
  }

  err := notification.Dismiss(os.Args[1])
  if err != nil {
    fmt.Fprintf(os.Stderr, "Error dismissing notification: %v\n", err)
    os.Exit(1)
  }
}
