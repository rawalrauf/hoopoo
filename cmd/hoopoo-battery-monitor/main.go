package main

import (
  "fmt"
  "os"
  "os/exec"
  "path/filepath"

  "github.com/rawalrauf/hoopoo/internal/sysfs"
)

const (
  batteryThreshold = 15
)

func main() {
  // Get notification flag path
  uid := os.Getuid()
  flagPath := fmt.Sprintf("/run/user/%d/hoopoo_battery_notified", uid)

  // Get battery level
  batteryLevel, err := sysfs.GetCombinedCapacity()
  if err != nil {
    os.Exit(1)
  }

  // Get battery status
  batteries, err := sysfs.FindAllBatteries()
  if err != nil {
    os.Exit(1)
  }
  status, err := batteries[0].GetStatus()
  if err != nil {
    os.Exit(1)
  }

  // Check if we should send notification
  if status == "Discharging" && batteryLevel <= batteryThreshold {
    // Check if already notified
    if _, err := os.Stat(flagPath); os.IsNotExist(err) {
      sendNotification(batteryLevel)
      
      // Create flag file
      file, err := os.Create(flagPath)
      if err == nil {
        file.Close()
      }
    }
  } else {
    // Remove flag file if it exists
    os.Remove(flagPath)
  }
}

func sendNotification(level int) {
  // Play sound
  exec.Command("systemd-run", "--user", "--no-block", "paplay",
    "/usr/share/sounds/freedesktop/stereo/suspend-error.oga").Start()

  // Send notification
  message := fmt.Sprintf("Battery is down to %d%%", level)
  exec.Command("notify-send", "-u", "critical", "󱐋 Time to recharge!",
    message, "-i", "battery-caution", "-t", "30000").Run()

  // Run hook
  hoopoBin := os.Getenv("HOOPOO_PATH")
  if hoopoBin == "" {
    hoopoBin = filepath.Join(os.Getenv("HOME"), ".local/share/hoopoo")
  }
  hookPath := filepath.Join(hoopoBin, "bin", "hoopoo-hook")
  exec.Command(hookPath, "battery-low", fmt.Sprintf("%d", level)).Run()
}
