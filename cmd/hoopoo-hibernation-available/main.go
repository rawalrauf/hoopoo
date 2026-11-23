// hoopoo:summary=Check if hibernation is supported

package main

import (
  "bufio"
  "os"
  "strconv"
  "strings"
)

func main() {
  // Check if hibernation kernel interface exists
  if _, err := os.Stat("/sys/power/image_size"); err != nil {
    os.Exit(1)
  }
  
  // Sum all swap sizes (excluding zram)
  swapSizeKB := getSwapSize()
  swapSize := 1024 * swapSizeKB
  
  // Get hibernation image size
  imageSize := getHibernationImageSize()
  
  // Check if swap is large enough and mkinitcpio config exists
  if swapSize > imageSize && fileExists("/etc/mkinitcpio.conf.d/hoopoo_resume.conf") {
    os.Exit(0)
  }
  
  os.Exit(1)
}

// getSwapSize returns total swap size in KB (excluding zram)
func getSwapSize() int {
  file, err := os.Open("/proc/swaps")
  if err != nil {
    return 0
  }
  defer file.Close()
  
  var total int
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    line := scanner.Text()
    
    // Skip header and zram devices
    if strings.HasPrefix(line, "Filename") || strings.Contains(line, "zram") {
      continue
    }
    
    // Parse size from third column
    fields := strings.Fields(line)
    if len(fields) >= 3 {
      if size, err := strconv.Atoi(fields[2]); err == nil {
        total += size
      }
    }
  }
  
  return total
}

// getHibernationImageSize returns the hibernation image size in bytes
func getHibernationImageSize() int {
  data, err := os.ReadFile("/sys/power/image_size")
  if err != nil {
    return 0
  }
  
  size, err := strconv.Atoi(strings.TrimSpace(string(data)))
  if err != nil {
    return 0
  }
  
  return size
}

// fileExists checks if a file exists
func fileExists(path string) bool {
  _, err := os.Stat(path)
  return err == nil
}
