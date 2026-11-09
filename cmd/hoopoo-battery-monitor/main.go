// hoopoo:summary=Designed to be run by systemd timer every 30 seconds and alerts if battery is low
// hoopoo:hidden=true

package main

import (
  "fmt"
  "os"

  "github.com/rawalrauf/hoopoo/internal/usage"
  "github.com/rawalrauf/hoopoo/internal/paths"
  "github.com/rawalrauf/hoopoo/internal/notify"
  "github.com/rawalrauf/hoopoo/internal/hooks"
  "github.com/rawalrauf/hoopoo/internal/sysfs/battery"
)

const (
  batteryThreshold = 15
)

func main() {
  usage.Message("Designed to be run by systemd timer every 30 seconds and alerts if battery is low")

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
  glyph := "󱐋"
  headline := "Time to recharge!"
  description := fmt.Sprintf("Battery is down to %d%%", level)
  soundFile := "/usr/share/sounds/freedesktop/stereo/suspend-error.oga"
  
  notify.SendWithGlyphAndSound(glyph, headline, description, "critical", "battery-caution", 30000, soundFile)
  hooks.Run("battery-low", fmt.Sprintf("%d", level))
}
