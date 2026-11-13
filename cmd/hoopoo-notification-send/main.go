// hoopoo:summary=Send a desktop notification with Hoopoo glyph and body spacing
// hoopoo:args=<headline> [description] [options]
// hoopoo:examples=hoopoo notification send --glyph="󰔛" "Reminder" "5 minutes are up" -u critical

package main

import (
  "fmt"
  "os"
  "strconv"

  "hoopoo/pkg/notification"
)

func main() {
  if len(os.Args) < 2 {
    fmt.Fprintln(os.Stderr, "Usage: hoopoo-notification-send [options] <headline> [description]")
    fmt.Fprintln(os.Stderr, "Options:")
    fmt.Fprintln(os.Stderr, "  --glyph=<icon>     Icon glyph to show before headline")
    fmt.Fprintln(os.Stderr, "  -u <urgency>       Urgency level (low, normal, critical)")
    fmt.Fprintln(os.Stderr, "  -i <icon>          Icon name or path")
    fmt.Fprintln(os.Stderr, "  -t <timeout>       Timeout in milliseconds")
    fmt.Fprintln(os.Stderr, "  --sound=<path>     Play sound file with notification")
    fmt.Fprintln(os.Stderr, "  --progress=<0-100> Show progress bar (replaces description)")
    fmt.Fprintln(os.Stderr, "  --replace-id=<id>  Replace existing notification")
    fmt.Fprintln(os.Stderr, "  --print-id         Print notification ID")
    os.Exit(1)
  }

  var positional []string
  glyph := ""
  urgency := ""
  icon := ""
  timeout := 0
  soundFile := ""
  progress := -1
  replaceID := ""
  printID := false
  var actions []string

  i := 1
  for i < len(os.Args) {
    arg := os.Args[i]
    
    if len(arg) > 8 && arg[:8] == "--glyph=" {
      glyph = arg[8:]
      i++
      continue
    }
    
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
    
    if len(arg) > 8 && arg[:8] == "--sound=" {
      soundFile = arg[8:]
      i++
      continue
    }
    
    if len(arg) > 11 && arg[:11] == "--progress=" {
      val, err := strconv.Atoi(arg[11:])
      if err == nil && val >= 0 && val <= 100 {
        progress = val
      }
      i++
      continue
    }
    
    if len(arg) > 13 && arg[:13] == "--replace-id=" {
      replaceID = arg[13:]
      i++
      continue
    }
    
    if arg == "--print-id" {
      printID = true
      i++
      continue
    }
    
    if arg == "-A" && i+1 < len(os.Args) {
      actions = append(actions, os.Args[i+1])
      i += 2
      continue
    }
    
    if len(arg) > 0 && arg[0] != '-' {
      positional = append(positional, arg)
      i++
      continue
    }
    
    i++
  }

  if len(positional) < 1 {
    fmt.Fprintln(os.Stderr, "Error: headline is required")
    os.Exit(1)
  }

  headline := positional[0]
  description := ""
  
  if len(positional) >= 2 {
    description = positional[1]
  }

  opts := notification.Options{
    Glyph:       glyph,
    Headline:    headline,
    Description: description,
    Urgency:     urgency,
    Icon:        icon,
    Timeout:     timeout,
    SoundFile:   soundFile,
    Progress:    progress,
    ReplaceID:   replaceID,
    Actions:     actions,
  }

  if printID {
    id, err := notification.SendWithID(opts)
    if err != nil {
      fmt.Fprintf(os.Stderr, "Error sending notification: %v\n", err)
      os.Exit(1)
    }
    fmt.Println(id)
  } else {
    err := notification.Send(opts)
    if err != nil {
      fmt.Fprintf(os.Stderr, "Error sending notification: %v\n", err)
      os.Exit(1)
    }
  }
}
