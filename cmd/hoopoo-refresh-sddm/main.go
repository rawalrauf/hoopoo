// hoopoo:summary=Refresh the SDDM theme from default
// hoopoo:requires-sudo=true

package main

import (
  "os/exec"
  "path/filepath"

  "hoopoo/pkg/system/paths"
)

func main() {
  sddmThemeDir := "/usr/share/sddm/themes/hoopoo"
  defaultSddmTheme := filepath.Join(paths.Hoopoo(), "default/sddm/hoopoo")

  // Remove existing theme
  exec.Command("sudo", "rm", "-rf", sddmThemeDir).Run()

  // Copy default theme
  exec.Command("sudo", "cp", "-r", defaultSddmTheme, sddmThemeDir).Run()
}
