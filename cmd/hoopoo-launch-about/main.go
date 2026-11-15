// hoopoo:summary=Launch the fastfetch TUI that gives information about the current system.

package main

import (
  "os"
  "os/exec"
  "syscall"
)

func main() {
  // Call hoopoo launch tui with fastfetch wrapper
  // tput civis = hide cursor, tput cnorm = show cursor
  // read -n 1 -s = wait for any key press
  cmd := exec.Command(
    "hoopoo",
    "launch",
    "tui",
    "--app-id=org.hoopoo.about",
    "--title=Hoopoo",
    "bash",
    "-c",
    "tput civis; fastfetch; read -n 1 -s; tput cnorm",
  )
  
  cmd.Stdin = os.Stdin
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr
  cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
  
  cmd.Run()
}
