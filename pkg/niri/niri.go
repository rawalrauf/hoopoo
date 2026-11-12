// Niri compositor utilities

package niri

import (
  "os/exec"
)

func Reload() error {
  return exec.Command("niri", "msg", "action", "load-config-file").Run()
}
