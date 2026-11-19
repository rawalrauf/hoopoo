// Process utilities

package process

import (
  "os/exec"
)

// IsRunning checks if a process with the given name is running
func IsRunning(name string) bool {
  return exec.Command("pgrep", "-x", name).Run() == nil
}

// Kill terminates a process by name
func Kill(name string) error {
  return exec.Command("pkill", "-x", name).Run()
}

// KillForce force terminates a process by name
func KillForce(name string) error {
  return exec.Command("pkill", "-9", "-x", name).Run()
}
