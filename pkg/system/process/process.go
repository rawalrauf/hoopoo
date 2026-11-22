// Process utilities

package process

import (
  "os/exec"
)

// IsRunning checks if a process with the given name is running (exact match)
func IsRunning(name string) bool {
  return exec.Command("pgrep", "-x", name).Run() == nil
}

// IsRunningByPattern checks if a process matching a pattern is running
func IsRunningByPattern(pattern string) bool {
  return exec.Command("pgrep", "-f", pattern).Run() == nil
}

// IsRunningByName checks if a process with name in title/args is running
func IsRunningByName(name string) bool {
  return exec.Command("pgrep", "-f", name).Run() == nil
}

// Kill terminates a process by name (exact match)
func Kill(name string) error {
  return exec.Command("pkill", "-x", name).Run()
}

// KillByPattern terminates processes matching a pattern
func KillByPattern(pattern string) error {
  return exec.Command("pkill", "-SIGINT", "-f", pattern).Run()
}

// KillByName terminates processes with name in title/args
func KillByName(name string) error {
  return exec.Command("pkill", "-f", name).Run()
}

// KillForce force terminates a process by name (exact match)
func KillForce(name string) error {
  return exec.Command("pkill", "-9", "-x", name).Run()
}

// KillForceByPattern force terminates processes matching a pattern
func KillForceByPattern(pattern string) error {
  return exec.Command("pkill", "-9", "-f", pattern).Run()
}

// Start starts a process with given args in background
func Start(name string, args ...string) error {
  cmd := exec.Command(name, args...)
  return cmd.Start()
}

// SendSignal sends a signal to a process (exact name match)
func SendSignal(signal, name string) error {
  return exec.Command("pkill", "-"+signal, "-x", name).Run()
}

// SendSignalPattern sends a signal to processes matching a pattern
func SendSignalPattern(signal, pattern string) error {
  return exec.Command("pkill", "-"+signal, "-f", pattern).Run()
}
