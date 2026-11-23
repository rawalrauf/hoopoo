// hoopoo:summary=Set power profile based on AC/battery state
// hoopoo:args=[autodetect|ac|battery]

package main

import (
  "os"
  "os/exec"
  "path/filepath"
  "strings"
)

func main() {
  action := "autodetect"
  if len(os.Args) > 1 {
    action = os.Args[1]
  }
  
  // Auto-detect power state if requested
  if action == "" || action == "autodetect" {
    action = detectPowerState()
  }
  
  // Get available profiles
  profiles := getAvailableProfiles()
  
  // Set profile based on power state
  var targetProfile string
  switch action {
  case "ac":
    // Prefer performance on AC, fall back to balanced
    if contains(profiles, "performance") {
      targetProfile = "performance"
    } else {
      targetProfile = "balanced"
    }
  case "battery":
    targetProfile = "balanced"
  default:
    return
  }
  
  // Set the profile
  exec.Command("powerprofilesctl", "set", targetProfile).Run()
}

// detectPowerState checks if system is on AC or battery
func detectPowerState() string {
  powerSupplyDir := "/sys/class/power_supply"
  
  entries, err := os.ReadDir(powerSupplyDir)
  if err != nil {
    return "battery"
  }
  
  for _, entry := range entries {
    psPath := filepath.Join(powerSupplyDir, entry.Name())
    
    // Check if this is a Mains or USB power supply
    typeData, err := os.ReadFile(filepath.Join(psPath, "type"))
    if err != nil {
      continue
    }
    psType := strings.TrimSpace(string(typeData))
    
    if psType != "Mains" && psType != "USB" {
      continue
    }
    
    // Check if it's online
    onlineData, err := os.ReadFile(filepath.Join(psPath, "online"))
    if err != nil {
      continue
    }
    
    if strings.TrimSpace(string(onlineData)) == "1" {
      return "ac"
    }
  }
  
  return "battery"
}

// getAvailableProfiles returns list of available power profiles
func getAvailableProfiles() []string {
  cmd := exec.Command("powerprofilesctl", "list")
  output, err := cmd.Output()
  if err != nil {
    return []string{"balanced"}
  }
  
  var profiles []string
  lines := strings.Split(string(output), "\n")
  
  for _, line := range lines {
    line = strings.TrimSpace(line)
    // Match lines like "* performance:" or "  balanced:"
    if strings.HasSuffix(line, ":") {
      // Remove leading * and spaces, and trailing :
      profile := strings.TrimSpace(strings.TrimPrefix(line, "*"))
      profile = strings.TrimSuffix(profile, ":")
      profile = strings.TrimSpace(profile)
      if profile != "" {
        profiles = append(profiles, profile)
      }
    }
  }
  
  return profiles
}

// contains checks if a string is in a slice
func contains(slice []string, item string) bool {
  for _, s := range slice {
    if s == item {
      return true
    }
  }
  return false
}
