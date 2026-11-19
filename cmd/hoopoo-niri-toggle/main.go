// hoopoo:summary=Toggle niri flags by copying KDL snippets into the toggles/niri directory.
// hoopoo:args=[--enabled-notification <text>] [--disabled-notification <text>] [--invert] <flag-name>

package main

import (
  "fmt"
  "os"
  "os/exec"
  "path/filepath"
  "syscall"

  "hoopoo/pkg/system/paths"
  "hoopoo/pkg/system/toggle"
)

func main() {
  var enabledNotification, disabledNotification string
  var invert bool
  
  // Parse arguments
  args := os.Args[1:]
  for len(args) > 1 {
    if args[0] == "--enabled-notification" && len(args) > 1 {
      enabledNotification = args[1]
      args = args[2:]
    } else if args[0] == "--disabled-notification" && len(args) > 1 {
      disabledNotification = args[1]
      args = args[2:]
    } else if args[0] == "--invert" {
      invert = true
      args = args[1:]
    } else {
      break
    }
  }

  if len(args) < 1 {
    fmt.Fprintln(os.Stderr, "Usage: hoopoo-niri-toggle [--enabled-notification <text>] [--disabled-notification <text>] [--invert] <flag-name>")
    os.Exit(1)
  }

  flagName := args[0]
  flagFile := filepath.Join(paths.State(), "toggles", "niri", flagName+".kdl")
  flagSource := filepath.Join(paths.Hoopoo(), "default/niri/toggles", flagName+".kdl")

  // Check if flag is currently enabled
  if toggle.NiriEnabled(flagName) {
    // Disable it
    toggle.NiriClear(flagName)
    
    if invert {
      if enabledNotification != "" {
        exec.Command("hoopoo", "notification", "send", enabledNotification).Run()
      }
    } else {
      if disabledNotification != "" {
        exec.Command("hoopoo", "notification", "send", disabledNotification).Run()
      }
    }
  } else {
    // Enable it
    if _, err := os.Stat(flagSource); err != nil {
      fmt.Fprintf(os.Stderr, "Flag not found: %s\n", flagName)
      os.Exit(1)
    }
    
    os.MkdirAll(filepath.Dir(flagFile), 0755)
    
    content, _ := os.ReadFile(flagSource)
    os.WriteFile(flagFile, content, 0644)
    
    if invert {
      if disabledNotification != "" {
        exec.Command("hoopoo", "notification", "send", disabledNotification).Run()
      }
    } else {
      if enabledNotification != "" {
        exec.Command("hoopoo", "notification", "send", enabledNotification).Run()
      }
    }
  }

  // Restart niri
  cmd := exec.Command("hoopoo", "restart", "niri")
  cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
  cmd.Run()
}
