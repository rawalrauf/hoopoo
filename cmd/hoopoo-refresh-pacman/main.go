// hoopoo:summary=Overwrite the package configuration for /etc/pacman with the Hoopoo default of using its dedicated mirrors and repositories, then update all packages.
// hoopoo:requires-sudo=true

package main

import (
  "fmt"
  "os"
  "os/exec"
  "path/filepath"

  "hoopoo/pkg/system/paths"
)

func main() {
  channel := "stable"
  if len(os.Args) > 1 {
    channel = os.Args[1]
  }

  if channel != "stable" && channel != "edge" {
    fmt.Fprintf(os.Stderr, "Error: Invalid channel '%s'. Must be one of: stable, edge\n", channel)
    os.Exit(1)
  }

  // Backup existing config
  exec.Command("sudo", "cp", "-f", "/etc/pacman.conf", "/etc/pacman.conf.bak").Run()

  fmt.Printf("Setting channel to %s\n\n", channel)

  // Copy channel-specific config
  configFile := fmt.Sprintf("pacman--%s.conf", channel)
  srcPath := filepath.Join(paths.Hoopoo(), "default/pacman", configFile)
  exec.Command("sudo", "cp", "-f", srcPath, "/etc/pacman.conf").Run()

  // Reset all package DBs and update
  exec.Command("sudo", "pacman", "-Syyuu", "--noconfirm").Run()
}
