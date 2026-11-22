// CPU detection and information

package cpu

import (
  "os"
  "strings"
)

// IsIntel checks if the system has an Intel CPU
func IsIntel() bool {
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

// IsAMD checks if the system has an AMD CPU
func IsAMD() bool {
  data, err := os.ReadFile("/proc/cpuinfo")
  if err != nil {
    return false
  }
  
  for _, line := range strings.Split(string(data), "\n") {
    if strings.HasPrefix(line, "vendor_id") {
      parts := strings.SplitN(line, ":", 2)
      if len(parts) == 2 {
        vendor := strings.TrimSpace(parts[1])
        return vendor == "AuthenticAMD"
      }
    }
  }
  return false
}
