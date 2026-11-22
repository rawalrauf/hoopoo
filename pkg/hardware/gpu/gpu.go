// GPU detection and information

package gpu

import (
  "os/exec"
  "strings"
)

// IsHybrid checks if the system has hybrid graphics (2+ GPUs)
func IsHybrid() bool {
  if _, err := exec.LookPath("supergfxctl"); err == nil {
    cmd := exec.Command("supergfxctl", "-s")
    output, err := cmd.Output()
    if err == nil {
      return strings.Contains(string(output), "Hybrid")
    }
  }
  
  cmd := exec.Command("lspci")
  output, err := cmd.Output()
  if err != nil {
    return false
  }
  
  count := 0
  for _, line := range strings.Split(string(output), "\n") {
    lower := strings.ToLower(line)
    if strings.Contains(lower, "vga") || strings.Contains(lower, "3d") || strings.Contains(lower, "display") {
      count++
    }
  }
  return count >= 2
}

// IsNvidiaGSP checks if the system has an Nvidia GPU with GSP firmware support
func IsNvidiaGSP() bool {
  cmd := exec.Command("lspci")
  output, err := cmd.Output()
  if err != nil {
    return false
  }
  
  for _, line := range strings.Split(string(output), "\n") {
    if !strings.Contains(strings.ToLower(line), "nvidia") {
      continue
    }
    
    if strings.Contains(line, "GTX 16") ||
       strings.Contains(line, "RTX 2") ||
       strings.Contains(line, "RTX 3") ||
       strings.Contains(line, "RTX 4") ||
       strings.Contains(line, "RTX 5") ||
       strings.Contains(line, "RTX PRO") ||
       strings.Contains(line, "Quadro RTX") ||
       strings.Contains(line, "RTX A") ||
       strings.Contains(line, " A100") ||
       strings.Contains(line, " A30") ||
       strings.Contains(line, " A40") ||
       strings.Contains(line, " H100") ||
       strings.Contains(line, " T4") ||
       strings.Contains(line, " L4") ||
       strings.Contains(line, " L40") {
      return true
    }
  }
  return false
}

// IsNvidiaWithoutGSP checks if the system has an older Nvidia GPU without GSP
func IsNvidiaWithoutGSP() bool {
  cmd := exec.Command("lspci")
  output, err := cmd.Output()
  if err != nil {
    return false
  }
  
  for _, line := range strings.Split(string(output), "\n") {
    if !strings.Contains(strings.ToLower(line), "nvidia") {
      continue
    }
    
    if strings.Contains(line, "GTX 9") ||
       strings.Contains(line, "GTX 10") ||
       strings.Contains(line, "GT 10") ||
       strings.Contains(line, "Quadro P") ||
       strings.Contains(line, "Quadro M") ||
       strings.Contains(line, "Quadro GV100") ||
       strings.Contains(line, "MX") ||
       strings.Contains(line, "Titan X") ||
       strings.Contains(line, "Titan Xp") ||
       strings.Contains(line, "Titan V") ||
       strings.Contains(line, "Tesla V100") {
      return true
    }
  }
  return false
}

// HasVulkan checks if Vulkan ICD drivers are installed
func HasVulkan() bool {
  cmd := exec.Command("sh", "-c", "ls /usr/share/vulkan/icd.d/*.json 2>/dev/null")
  return cmd.Run() == nil
}

// IsIntelPantherLake checks if the system has Intel Panther Lake graphics
func IsIntelPantherLake() bool {
  cmd := exec.Command("lspci")
  output, err := cmd.Output()
  if err != nil {
    return false
  }
  
  lines := strings.Split(string(output), "\n")
  for _, line := range lines {
    lower := strings.ToLower(line)
    if (strings.Contains(lower, "vga") || strings.Contains(lower, "3d") || strings.Contains(lower, "display")) &&
       strings.Contains(lower, "panther lake") {
      return true
    }
  }
  return false
}
