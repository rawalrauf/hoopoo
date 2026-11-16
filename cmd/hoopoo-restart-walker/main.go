// hoopoo:summary=Restart Walker and related user services

package main

import (
  "fmt"
  "os"
  "os/exec"
  "os/user"
  "path/filepath"
)

func restartServices() {
  // Restart elephant if enabled
  if exec.Command("systemctl", "--user", "is-enabled", "elephant.service").Run() == nil {
    exec.Command("systemctl", "--user", "restart", "elephant.service").Run()
  }

  // Restart walker if enabled
  if exec.Command("systemctl", "--user", "is-enabled", "app-walker@autostart.service").Run() == nil {
    exec.Command("systemctl", "--user", "restart", "app-walker@autostart.service").Run()
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
