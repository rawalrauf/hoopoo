// Hardware detection utilities

package hardware

import (
  "os"
  "os/exec"
  "path/filepath"
  "strings"
)

func dmiRead(file string) (string, error) {
  data, err := os.ReadFile(filepath.Join("/sys/class/dmi/id", file))
  if err != nil {
    return "", err
  }
  return strings.TrimSpace(string(data)), nil
}

func dmiMatch(file, pattern string) bool {
  content, err := dmiRead(file)
  if err != nil {
    return false
  }
  return strings.Contains(strings.ToLower(content), strings.ToLower(pattern))
}

func IsAsusROG() bool {
  vendor, err := dmiRead("sys_vendor")
  if err != nil || vendor != "ASUSTeK COMPUTER INC." {
    return false
  }
  return dmiMatch("product_family", "ROG")
}

func IsDellLatitude7490() bool {
  vendor, err := dmiRead("sys_vendor")
  if err != nil || vendor != "Dell Inc." {
    return false
  }
  return dmiMatch("product_name", "Latitude 7490")
}

func IsFramework16() bool {
  vendor, err := dmiRead("sys_vendor")
  if err != nil || vendor != "Framework" {
    return false
  }
  return dmiMatch("product_name", "Laptop 16") || dmiMatch("product_family", "Laptop 16")
}

func IsIntelCPU() bool {
  data, err := os.ReadFile("/proc/cpuinfo")
  if err != nil {
    return false
  }
  
  for _, line := range strings.Split(string(data), "\n") {
    if strings.HasPrefix(line, "vendor_id") {
      parts := strings.SplitN(line, ":", 2)
      if len(parts) == 2 {
        vendor := strings.TrimSpace(parts[1])
        return vendor == "GenuineIntel"
      }
    }
  }
  return false
}

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

func IsLenovoThinkPad() bool {
  vendor, err := dmiRead("sys_vendor")
  if err != nil || vendor != "LENOVO" {
    return false
  }
  return dmiMatch("product_family", "ThinkPad")
}

func MatchHardware(pattern string) bool {
  return dmiMatch("product_name", pattern) || dmiMatch("product_family", pattern)
}

func IsSurface() bool {
  vendor, err := dmiRead("sys_vendor")
  if err != nil || vendor != "Microsoft Corporation" {
    return false
  }
  return MatchHardware("Surface")
}

func HasTouchpad() bool {
  data, err := os.ReadFile("/proc/bus/input/devices")
  if err != nil {
    return false
  }
  
  lower := strings.ToLower(string(data))
  return strings.Contains(lower, "touchpad") || 
         strings.Contains(lower, "trackpad") ||
         strings.Contains(lower, "synaptics")
}

func HasTouchscreen() bool {
  data, err := os.ReadFile("/proc/bus/input/devices")
  if err != nil {
    return false
  }
  
  lower := strings.ToLower(string(data))
  return strings.Contains(lower, "touchscreen") || 
         strings.Contains(lower, "multi-touch") ||
         strings.Contains(lower, "wacom")
}

func HasExternalMonitors() bool {
  pattern := "/sys/class/drm/card*-*/status"
  matches, err := filepath.Glob(pattern)
  if err != nil {
    return false
  }
  
  for _, path := range matches {
    if strings.Contains(path, "-eDP-") {
      continue
    }
    
    status, err := os.ReadFile(path)
    if err != nil {
      continue
    }
    
    if strings.TrimSpace(string(status)) == "connected" {
      return true
    }
  }
  return false
}

func HasDellXPSHapticTouchpad() bool {
  if !MatchHardware("XPS") {
    return false
  }
  
  _, err := os.Stat("/sys/bus/i2c/devices/i2c-VEN_06CB:00")
  return err == nil
}

func IsDellXPSOLED() bool {
  if !MatchHardware("XPS") || !IsIntelPantherLake() {
    return false
  }
  
  pattern := "/sys/class/drm/card*-eDP-*/edid"
  matches, err := filepath.Glob(pattern)
  if err != nil || len(matches) == 0 {
    return false
  }
  
  data, err := os.ReadFile(matches[0])
  if err != nil || len(data) < 10 {
    return false
  }
  
  return data[8] == 0x30 && data[9] == 0xe4
}

func IsHybridGPU() bool {
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

func HasVulkan() bool {
  dir := "/usr/share/vulkan/icd.d"
  entries, err := os.ReadDir(dir)
  if err != nil {
    return false
  }
  
  for _, entry := range entries {
    if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
      return true
    }
  }
  return false
}
