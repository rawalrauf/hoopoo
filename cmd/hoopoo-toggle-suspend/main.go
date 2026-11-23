// hoopoo:summary=Toggle suspend availability in the system menu

package main

import (
  "hoopoo/pkg/notification"
  "hoopoo/pkg/toggle"
)

func main() {
  if toggle.Enabled("suspend-off") {
    // Currently disabled, enable it
    toggle.Clear("suspend-off")
    notification.Send(notification.Options{
      Urgency:  "low",
      Glyph:    "󰒲",
      Headline: "Suspend now available in system menu",
    })
  } else {
    // Currently enabled, disable it
    toggle.Write("suspend-off", "")
    notification.Send(notification.Options{
      Urgency:  "low",
      Glyph:    "󰒲",
      Headline: "Suspend removed from system menu",
    })
  }
}
