// hoopoo:summary=Control mirroring the internal laptop display to an external monitor via wl-mirror
// hoopoo:group=niri
// hoopoo:name=monitor internal mirror
// hoopoo:args=<on|off|toggle>

package main

import (
  "fmt"
  "os"
  "os/exec"
  "regexp"
  "syscall"
  "time"

  "hoopoo/pkg/niri"
  "hoopoo/pkg/system/paths"
  "hoopoo/pkg/system/process"
)

func hasExternalMonitors() bool {
  return exec.Command("hoopoo", "hw", "external", "monitors").Run() == nil
}

func enable(internal, external string) {
  if external == "" {
    exec.Command("hoopoo", "notification", "send", "-u", "low", "󰍹  No external monitor found").Run()
    return
  }
  if internal == "" {
    exec.Command("hoopoo", "notification", "send", "-u", "low", "󰍹  No internal display found").Run()
    return
  }
  if process.IsRunning("wl-mirror") {
    return
  }

  toggleContent := fmt.Sprintf(`window-rule {
    match app-id="at.yrlf.wl_mirror"
    open-on-output "%s"
}`, external)

  // Check if toggle file exists with different external monitor
  toggleFile := paths.State() + "/toggles/niri/wl-mirror-output.kdl"
  currentContent, _ := os.ReadFile(toggleFile)
  re := regexp.MustCompile(`open-on-output "([^"]+)"`)
  match := re.FindStringSubmatch(string(currentContent))
  
  // Write toggle file if it doesn't exist or external monitor changed
  if len(match) == 0 || match[1] != external {
    niri.WriteState("wl-mirror-output.kdl", toggleContent)
    exec.Command("hoopoo", "restart", "niri").Run()
    time.Sleep(300 * time.Millisecond)
  }

  cmd := exec.Command("wl-mirror", "--fullscreen-output", external, internal)
  cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
  cmd.Start()
  
  exec.Command("hoopoo", "notification", "send", "-u", "low", fmt.Sprintf("󰍹  Mirroring enabled on %s", external)).Run()
}

func disable() {
  // Get focused window
  focusedWindow, _ := niri.GetFocusedWindow()
  var focusedID int
  if focusedWindow != nil {
    focusedID = focusedWindow.ID
  }

  // Find and close wl-mirror window
  windows, _ := niri.GetWindows()
  mirrorClosed := false
  
  for _, w := range windows {
    if w.AppID == "at.yrlf.wl_mirror" {
      niri.CloseWindow(w.ID)
      mirrorClosed = true
      break
    }
  }

  if !mirrorClosed {
    process.Kill("wl-mirror")
  }

  // Restore focus
  if focusedID > 0 {
    niri.FocusWindow(focusedID)
  }

  // Clear toggle file
  exec.Command("hoopoo", "niri", "toggle", "clear", "wl-mirror-output.kdl").Run()
  
  exec.Command("hoopoo", "notification", "send", "-u", "low", "󰍹  Mirroring disabled").Run()
}

func main() {
  if len(os.Args) < 2 {
    fmt.Fprintln(os.Stderr, "Usage: hoopoo niri monitor internal mirror {on|off|toggle}")
    os.Exit(1)
  }

  internal, _ := niri.GetInternalDisplay()
  external, _ := niri.GetExternalDisplay()

  switch os.Args[1] {
  case "on":
    enable(internal, external)
  case "off":
    disable()
  case "toggle":
    if process.IsRunning("wl-mirror") {
      disable()
    } else {
      enable(internal, external)
    }
  default:
    fmt.Fprintln(os.Stderr, "Usage: hoopoo niri monitor internal mirror {on|off|toggle}")
    os.Exit(1)
  }
}
