// hoopoo:summary=Restart Walker and related user services

package main

import (
  "fmt"
  "os"
  "os/exec"
  "os/user"
  "path/filepath"

  "hoopoo/pkg/system/systemd"
)

func restartServices() {
  // Restart elephant if enabled
  if systemd.IsEnabled("elephant.service") {
    systemd.Restart("elephant.service")
  }

  // Restart walker if enabled
  if systemd.IsEnabled("app-walker@autostart.service") {
    systemd.Restart("app-walker@autostart.service")
  } else {
    fmt.Println("\033[31mUnable to restart Walker -- RESTART MANUALLY\033[0m")
  }
}

func main() {
  // Check if running as root
  if os.Geteuid() == 0 {
    // Get script owner
    scriptPath, _ := os.Executable()
    info, err := os.Stat(scriptPath)
    if err != nil {
      restartServices()
      return
    }

    // Get owner username and UID
    scriptOwner, _ := user.LookupId(fmt.Sprintf("%d", info.Sys()))
    if scriptOwner == nil {
      restartServices()
      return
    }

    runtimeDir := filepath.Join("/run/user", scriptOwner.Uid)

    // Run as user via systemd-run
    exec.Command("systemd-run",
      "--uid="+scriptOwner.Username,
      "--setenv=XDG_RUNTIME_DIR="+runtimeDir,
      "hoopoo", "restart", "walker",
    ).Run()
  } else {
    restartServices()
  }
}
