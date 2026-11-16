// hoopoo:summary=Overwrite the user configs for the Walker application launcher (which also powers the Hoopoo Menu) and restart the services.

package main

import (
  "os"
  "os/exec"
  "path/filepath"
  "syscall"

  "hoopoo/pkg/system/paths"
  "hoopoo/pkg/system/systemd"
)

func main() {
  userConfig := paths.UserConfig()
  hoopooPath := paths.Hoopoo()

  // Create autostart directory and copy walker.desktop
  autostartDir := filepath.Join(userConfig, "autostart")
  os.MkdirAll(autostartDir, 0755)
  
  src := filepath.Join(hoopooPath, "default/walker/walker.desktop")
  dst := filepath.Join(autostartDir, "walker.desktop")
  srcContent, _ := os.ReadFile(src)
  os.WriteFile(dst, srcContent, 0644)

  // Create systemd override directory and copy restart.conf
  systemdOverrideDir := filepath.Join(userConfig, "systemd/user/app-walker@autostart.service.d")
  os.MkdirAll(systemdOverrideDir, 0755)
  
  src = filepath.Join(hoopooPath, "default/walker/restart.conf")
  dst = filepath.Join(systemdOverrideDir, "restart.conf")
  srcContent, _ = os.ReadFile(src)
  os.WriteFile(dst, srcContent, 0644)

  // Reload systemd daemon
  systemd.DaemonReload()

  // Refresh config files
  configs := []string{
    "walker/config.toml",
    "elephant/calc.toml",
    "elephant/desktopapplications.toml",
  }
  
  for _, config := range configs {
    cmd := exec.Command("hoopoo", "refresh", "config", config)
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
    cmd.Run()
  }

  // Link all elephant menus
  elephantMenusDir := filepath.Join(userConfig, "elephant/menus")
  os.MkdirAll(elephantMenusDir, 0755)
  
  menuPattern := filepath.Join(hoopooPath, "default/elephant/*.lua")
  menus, _ := filepath.Glob(menuPattern)
  
  for _, menu := range menus {
    menuName := filepath.Base(menu)
    linkPath := filepath.Join(elephantMenusDir, menuName)
    os.Remove(linkPath) // Remove old link if exists
    os.Symlink(menu, linkPath)
  }

  // Restart walker service
  cmd := exec.Command("hoopoo", "restart", "walker")
  cmd.Stdin = os.Stdin
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr
  cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
  cmd.Run()
}
