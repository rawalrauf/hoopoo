// hoopoo:summary=Designed to be run by systemd timer every 30 seconds and alerts if battery is low
// hoopoo:hidden=true

package main

import (
  "fmt"
  "os"

  "hoopoo.net/internal/common/paths"
  "hoopoo.net/internal/common/notify"
  "hoopoo.net/internal/common/hooks"
  "hoopoo.net/internal/sysfs/battery"
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
  notify.Send(notify.Options{
    Glyph:       "󱐋",
    Headline:    "Time to recharge!",
    Description: fmt.Sprintf("Battery is down to %d%%", level),
    Urgency:     "critical",
    Icon:        "battery-caution",
    Timeout:     30000,
    SoundFile:   "/usr/share/sounds/freedesktop/stereo/suspend-error.oga",
  })
  hooks.Run("battery-low", fmt.Sprintf("%d", level))
}
