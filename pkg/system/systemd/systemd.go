// Systemd service management utilities

package systemd

import (
  "os/exec"
)

// IsEnabled checks if a systemd user service is enabled
func IsEnabled(service string) bool {
  cmd := exec.Command("systemctl", "--user", "is-enabled", service)
  return cmd.Run() == nil
}

// IsActive checks if a systemd user service is active (running)
func IsActive(service string) bool {
  cmd := exec.Command("systemctl", "--user", "is-active", service)
  return cmd.Run() == nil
}

// Restart restarts a systemd user service
func Restart(service string) error {
  cmd := exec.Command("systemctl", "--user", "restart", service)
  return cmd.Run()
}

// DaemonReload reloads systemd user daemon configuration
func DaemonReload() error {
  cmd := exec.Command("systemctl", "--user", "daemon-reload")
  return cmd.Run()
}

// ResetFailed resets the failed state of a systemd user service
func ResetFailed(service string) error {
  cmd := exec.Command("systemctl", "--user", "reset-failed", service)
  return cmd.Run()
}

// Stop stops a systemd user service
func Stop(service string) error {
  cmd := exec.Command("systemctl", "--user", "stop", service)
  return cmd.Run()
}

// Enable enables a systemd user service
func Enable(service string) error {
  cmd := exec.Command("systemctl", "--user", "enable", service)
  return cmd.Run()
}

// EnableNow enables and starts a systemd user service
func EnableNow(service string) error {
  cmd := exec.Command("systemctl", "--user", "enable", "--now", service)
  return cmd.Run()
}
