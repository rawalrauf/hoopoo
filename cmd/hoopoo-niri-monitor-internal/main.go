// hoopoo:summary=Enable, disable, or toggle the internal laptop display
// hoopoo:group=niri
// hoopoo:name=monitor internal
// hoopoo:args=<on|off|toggle|recover>

package main

import (
  "fmt"
  "os"
  "os/exec"

  "hoopoo/pkg/niri"
  "hoopoo/pkg/niri/config"
)

const (
  configFile = "monitors.kdl"
  stateFile  = "internal-monitor-disabled"
)

func hasExternalMonitors() bool {
  return exec.Command("hoopoo", "hw", "external", "monitors").Run() == nil
}

func isDisabled(internal string) bool {
  // Check if 'off' property exists in monitors.kdl
  _, err := config.GetProperty(configFile, "output", internal, "off")
  return err == nil
}

func enable(internal string) {
  if err := config.RemoveProperty(configFile, "output", internal, "off"); err != nil {
    fmt.Fprintf(os.Stderr, "Failed to enable display: %v\n", err)
    os.Exit(1)
  }
  
  niri.ClearState(stateFile)
  exec.Command("hoopoo", "restart", "niri").Run()
  exec.Command("hoopoo", "notification", "send", "-u", "low", "󰍹  Laptop display enabled").Run()
}

func disable(internal string) {
  if !hasExternalMonitors() {
    exec.Command("hoopoo", "notification", "send", "-u", "low", "󰍹  Can't disable the only active display").Run()
    os.Exit(1)
  }
  
  // Add 'off' property to monitors.kdl
  if err := config.SetProperty(configFile, "output", internal, "off", "off"); err != nil {
    fmt.Fprintf(os.Stderr, "Failed to disable display: %v\n", err)
    os.Exit(1)
  }
  
  // Create state file with internal display name for boot-time recovery
  niri.WriteState(stateFile, internal)
  
  // Reload niri config to apply changes
  exec.Command("hoopoo", "restart", "niri").Run()
  
  exec.Command("hoopoo", "notification", "send", "-u", "low", "󰍹  Laptop display disabled").Run()
}

func recover(internal string) {
  if !hasExternalMonitors() && isDisabled(internal) {
    enable(internal)
  }
}

func main() {
  if len(os.Args) < 2 {
    fmt.Fprintln(os.Stderr, "Usage: hoopoo niri monitor internal {on|off|toggle|recover}")
    os.Exit(1)
  }

  internal, err := niri.GetInternalDisplay()
  if err != nil {
    exec.Command("hoopoo", "notification", "send", "-u", "low", "󰍹  No internal display found").Run()
    os.Exit(1)
  }

  switch os.Args[1] {
  case "on":
    enable(internal)
  case "off":
    disable(internal)
  case "toggle":
    if isDisabled(internal) {
      enable(internal)
    } else {
      disable(internal)
    }
  case "recover":
    recover(internal)
  default:
    fmt.Fprintln(os.Stderr, "Usage: hoopoo niri monitor internal {on|off|toggle|recover}")
    os.Exit(1)
  }
}