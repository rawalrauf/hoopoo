// hoopoo:summary=Overwrite the user config for the Limine bootloader and rebuild it.
// hoopoo:requires-sudo=true

package main

import (
  "fmt"
  "os"
  "os/exec"
  "path/filepath"
  "strings"

  "hoopoo/pkg/system/paths"
)

func main() {
  // Cleanup extra UKI if both exist
  machineID, _ := os.ReadFile("/etc/machine-id")
  machineIDStr := strings.TrimSpace(string(machineID))
  
  hoopooUKI := "/boot/EFI/Linux/hoopoo_linux.efi"
  machineUKI := filepath.Join("/boot/EFI/Linux", machineIDStr+"_linux.efi")
  
  if _, err1 := os.Stat(hoopooUKI); err1 == nil {
    if _, err2 := os.Stat(machineUKI); err2 == nil {
      fmt.Println("Cleanup extra UKI")
      exec.Command("sudo", "rm", "-f", machineUKI).Run()
    }
  }

  fmt.Println("Resetting limine config")

  // Backup existing config
  exec.Command("sudo", "mv", "/boot/limine.conf", "/boot/limine.conf.bak").Run()

  // Copy new config
  defaultConfig := filepath.Join(paths.Hoopoo(), "default/limine/limine.conf")
  exec.Command("sudo", "cp", defaultConfig, "/boot/limine.conf").Run()

  // Update limine
  exec.Command("sudo", "limine-update").Run()
  exec.Command("sudo", "limine-snapper-sync").Run()
}
