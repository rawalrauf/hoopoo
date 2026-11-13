// hoopoo:summary=Designed to be run by systemd timer every 30 seconds and alerts if battery is low
// hoopoo:hidden=true

package main

import (
  "fmt"
  "os"
  "os/exec"

  "hoopoo/pkg/system/paths"
  "hoopoo/pkg/battery"
)

const (
  batteryThreshold = 15
)

func main() {
  flagPath := fmt.Sprintf("%s/hoopoo_battery_notified", paths.Runtime())

  batteryLevel, err := battery.CombinedCapacity()
  if err != nil {
    os.Exit(1)
  }

  batteries, err := battery.FindAll()
  if err != nil {
    os.Exit(1)
  }
  status, err := batteries[0].Status()
  if err != nil {
    os.Exit(1)
  }

  if status == "Discharging" && batteryLevel <= batteryThreshold {
    if _, err := os.Stat(flagPath); os.IsNotExist(err) {
      sendNotification(batteryLevel)
      
      file, err := os.Create(flagPath)
      if err == nil {
        file.Close()
      }
    }
  } else {
    os.Remove(flagPath)
  }
}

func sendNotification(level int) {
  cmd := exec.Command("hoopoo", "notification", "send",
    "--glyph=󱐋",
    "Time to recharge!",
    fmt.Sprintf("Battery is down to %d%%", level),
    "-u", "critical",
    "-i", "battery-caution",
    "-t", "30000",
    "--sound=/usr/share/sounds/freedesktop/stereo/suspend-error.oga",
  )
  cmd.Run()
}
