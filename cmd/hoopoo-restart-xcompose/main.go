// hoopoo:summary=Restart the XCompose input method service (fcitx5) to apply new compose key settings.

package main

import (
  "os"
  "os/exec"
  "syscall"
)

func main() {
  // Call hoopoo restart app with fcitx5
  cmd := exec.Command(
    "hoopoo",
    "restart",
    "app",
    "fcitx5",
    "--disable",
    "notificationitem",
  )
  
  cmd.Stdin = os.Stdin
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr
  cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
  
  cmd.Run()
}
