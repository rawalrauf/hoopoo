// Hardware detection utilities

package hardware

import (
  "os"
  "path/filepath"
  "strings"
  
  "hoopoo/pkg/hardware/cpu"
  "hoopoo/pkg/hardware/dmi"
  "hoopoo/pkg/hardware/gpu"
  "hoopoo/pkg/hardware/input"
)

// System-specific hardware detection

// IsAsusROG checks if the system is an ASUS ROG laptop
func IsAsusROG() bool {
  vendor, err := dmi.Vendor()
  if err != nil || vendor != "ASUSTeK COMPUTER INC." {
    return false
  }
  return dmi.Match("product_family", "ROG")
}

// IsFramework16 checks if the system is a Framework 16 laptop
func IsFramework16() bool {
  vendor, err := dmi.Vendor()
  if err != nil || vendor != "Framework" {
    return false
  }
  return dmi.Match("product_name", "Laptop 16") || dmi.Match("product_family", "Laptop 16")
}

// IsSurface checks if the system is a Microsoft Surface device
func IsSurface() bool {
  vendor, err := dmi.Vendor()
  if err != nil || vendor != "Microsoft Corporation" {
    return false
  }
  return MatchHardware("Surface")
}

// MatchHardware checks if product name or family matches a pattern
func MatchHardware(pattern string) bool {
  return dmi.Match("product_name", pattern) || dmi.Match("product_family", pattern)
}

// Display detection

// HasExternalMonitors checks if external monitors are connected
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

// Dell XPS specific

// IsDellXPSOLED checks if Dell XPS has OLED display
func IsDellXPSOLED() bool {
  if !MatchHardware("XPS") || !gpu.IsIntelPantherLake() {
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

// Re-export sub-package functions for backward compatibility

func IsIntelCPU() bool {
  return cpu.IsIntel()
}

func IsIntelPantherLake() bool {
  return gpu.IsIntelPantherLake()
}

func IsHybridGPU() bool {
  return gpu.IsHybrid()
}

func IsNvidiaGSP() bool {
  return gpu.IsNvidiaGSP()
}

func IsNvidiaWithoutGSP() bool {
  return gpu.IsNvidiaWithoutGSP()
}

func HasVulkan() bool {
  return gpu.HasVulkan()
}

func HasTouchpad() bool {
  return input.HasTouchpad()
}

func HasTouchscreen() bool {
  return input.HasTouchscreen()
}
