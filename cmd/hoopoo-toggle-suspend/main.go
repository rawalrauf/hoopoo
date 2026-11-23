// hoopoo:summary=Toggle suspend availability in the system menu

package main

import (
  "os/exec"
  
  "hoopoo/pkg/toggle"
)

func main() {
  if toggle.Enabled("suspend-off") {
    // Currently disabled, enable it
    toggle.Clear("suspend-off")
    exec.Command("hoopoo", "notification", "send",
      "--glyph=󰒲",
      "Suspend now available in system menu",
      "-u", "low",
    ).Run()
  } else {
    // Currently enabled, disable it
    toggle.Write("suspend-off", "")
    exec.Command("hoopoo", "notification", "send",
      "--glyph=󰒲",
      "Suspend removed from system menu",
      "-u", "low",
    ).Run()
  }
}
