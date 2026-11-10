// hoopoo:summary=Send a desktop notification with Hoopoo glyph and body spacing
// hoopoo:args=<glyph> <headline> [description] [notify-send options]
// hoopoo:examples=hoopoo notification send "󰔛" "Reminder" "5 minutes are up" -u critical

package main

import (
  "fmt"
  "os"

  "hoopoo.net/internal/notification"
)

func main() {
  if len(os.Args) < 3 {
    fmt.Fprintln(os.Stderr, "Usage: hoopoo-notification-send <glyph> <headline> [description] [notify-send options]")
    os.Exit(1)
  }

  glyph := os.Args[1]
  headline := os.Args[2]
  description := ""
  urgency := ""
  icon := ""
  timeout := 0

  i := 3
  for i < len(os.Args) {
    arg := os.Args[i]
    
    if arg == "-u" && i+1 < len(os.Args) {
      urgency = os.Args[i+1]
      i += 2
      continue
    }
    
    if arg == "-i" && i+1 < len(os.Args) {
      icon = os.Args[i+1]
      i += 2
      continue
    }
    
    if arg == "-t" && i+1 < len(os.Args) {
      fmt.Sscanf(os.Args[i+1], "%d", &timeout)
      i += 2
      continue
    }
    
    if arg[0] != '-' && description == "" {
      description = arg
      i++
      continue
    }
    
    i++
  }

  err := notification.Send(notification.Options{
    Glyph:       glyph,
    Headline:    headline,
    Description: description,
    Urgency:     urgency,
    Icon:        icon,
    Timeout:     timeout,
  })

  if err != nil {
    fmt.Fprintf(os.Stderr, "Error sending notification: %v\n", err)
    os.Exit(1)
  }
}
