package main

import (
  "fmt"
  "os"
  "os/exec"
  "strings"
  
  "hoopoo/pkg/system/paths"
)

// hoopoo:summary=Toggle dedicated vs integrated GPU mode via supergfxd (for hybrid gpu laptops)
// hoopoo:requires-sudo=true

const supergfxdConfig = `{
  "mode": "Hybrid",
  "vfio_enable": true,
  "vfio_save": false,
  "always_reboot": false,
  "no_logind": false,
  "logout_timeout_s": 180,
  "hotplug_type": "None"
}
`

func main() {
  // Check if supergfxctl is installed
  if exec.Command("command", "-v", "supergfxctl").Run() != nil {
    fmt.Println("Installing supergfxctl...")
    if err := exec.Command("hoopoo", "pkg", "add", "supergfxctl").Run(); err != nil {
      fmt.Fprintln(os.Stderr, "Error: Failed to install supergfxctl")
      os.Exit(1)
    }
    
    // Create initial config
    writeCmd := exec.Command("sudo", "tee", "/etc/supergfxd.conf")
    writeCmd.Stdin = strings.NewReader(supergfxdConfig)
    writeCmd.Stdout = nil
    if err := writeCmd.Run(); err != nil {
      fmt.Fprintln(os.Stderr, "Error: Failed to create supergfxd config")
      os.Exit(1)
    }
    
    // Enable and start service
    if err := exec.Command("sudo", "systemctl", "enable", "--now", "supergfxd").Run(); err != nil {
      fmt.Fprintln(os.Stderr, "Error: Failed to enable supergfxd service")
      os.Exit(1)
    }
  }
  
  // Get current GPU mode
  cmd := exec.Command("supergfxctl", "-g")
  output, err := cmd.Output()
  if err != nil {
    fmt.Fprintln(os.Stderr, "Error: Failed to get GPU mode")
    os.Exit(1)
  }
  
  gpuMode := strings.TrimSpace(string(output))
  
  switch gpuMode {
  case "Integrated":
    switchToHybrid()
  case "Hybrid":
    switchToIntegrated()
  default:
    fmt.Fprintln(os.Stderr, "Hybrid GPU not found or in unknown mode.")
    os.Exit(1)
  }
}

func switchToHybrid() {
  if !confirmAction("Enable dedicated GPU and reboot?") {
    os.Exit(0)
  }
  
  // Switch to hybrid mode in config
  if err := exec.Command("sudo", "sed", "-i", `s/"mode": ".*"/"mode": "Hybrid"/`, "/etc/supergfxd.conf").Run(); err != nil {
    fmt.Fprintln(os.Stderr, "Error: Failed to update config")
    os.Exit(1)
  }
  
  // Remove force-igpu sleep script
  exec.Command("sudo", "rm", "-rf", "/usr/lib/systemd/system-sleep/force-igpu").Run()
  
  // Remove startup delay override
  exec.Command("sudo", "rm", "-rf", "/etc/systemd/system/supergfxd.service.d/delay-start.conf").Run()
  
  // Reboot
  exec.Command("hoopoo", "system", "reboot").Run()
}

func switchToIntegrated() {
  if !confirmAction("Use only integrated GPU and reboot?") {
    os.Exit(0)
  }
  
  // Switch to integrated mode and enable vfio
  if err := exec.Command("sudo", "sed", "-i", `s/"mode": ".*"/"mode": "Integrated"/`, "/etc/supergfxd.conf").Run(); err != nil {
    fmt.Fprintln(os.Stderr, "Error: Failed to update config mode")
    os.Exit(1)
  }
  
  if err := exec.Command("sudo", "sed", "-i", `s/"vfio_enable": false/"vfio_enable": true/`, "/etc/supergfxd.conf").Run(); err != nil {
    fmt.Fprintln(os.Stderr, "Error: Failed to enable vfio")
    os.Exit(1)
  }
  
  // Install force-igpu sleep script
  exec.Command("sudo", "mkdir", "-p", "/usr/lib/systemd/system-sleep").Run()
  
  hoopooPath := paths.Hoopoo()
  forceIgpuSrc := hoopooPath + "/default/systemd/system-sleep/force-igpu"
  if err := exec.Command("sudo", "cp", "-p", forceIgpuSrc, "/usr/lib/systemd/system-sleep/").Run(); err != nil {
    fmt.Fprintln(os.Stderr, "Warning: Failed to install force-igpu sleep script")
  }
  
  // Install startup delay override
  exec.Command("sudo", "mkdir", "-p", "/etc/systemd/system/supergfxd.service.d").Run()
  
  delayConfSrc := hoopooPath + "/default/systemd/system/supergfxd.service.d/delay-start.conf"
  if err := exec.Command("sudo", "cp", "-p", delayConfSrc, "/etc/systemd/system/supergfxd.service.d/").Run(); err != nil {
    fmt.Fprintln(os.Stderr, "Warning: Failed to install delay-start config")
  }
  
  // Reboot
  exec.Command("hoopoo", "system", "reboot").Run()
}

func confirmAction(prompt string) bool {
  cmd := exec.Command("gum", "confirm", prompt)
  cmd.Stdin = os.Stdin
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr
  return cmd.Run() == nil
}
