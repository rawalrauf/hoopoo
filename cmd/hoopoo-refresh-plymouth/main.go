// hoopoo:summary=Overwrite the user config for the Plymouth drive decryption and boot sequence with the Hoopoo default and rebuild it.
// hoopoo:requires-sudo=true

package main

import (
  "os/exec"
  "path/filepath"

  "hoopoo/pkg/system/paths"
)

func main() {
  // Copy plymouth theme files
  src := filepath.Join(paths.Hoopoo(), "default/plymouth")
  dst := "/usr/share/plymouth/themes/hoopoo/"
  
  exec.Command("sudo", "cp", "-r", src+"/.", dst).Run()

  // Set default theme
  exec.Command("sudo", "plymouth-set-default-theme", "hoopoo").Run()

  // Rebuild initramfs
  if exec.Command("command", "-v", "limine-mkinitcpio").Run() == nil {
    exec.Command("sudo", "limine-mkinitcpio").Run()
  } else {
    exec.Command("sudo", "mkinitcpio", "-P").Run()
  }
}
