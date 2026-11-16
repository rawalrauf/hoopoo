// hoopoo:summary=Reset the trackpad by unbinding and rebinding its driver.
// hoopoo:requires-sudo=true

package main

import (
  "fmt"
  "os/exec"
  "path/filepath"
  "strings"
  "time"
)

func writeToSysfs(path, value string) {
  cmd := exec.Command("sudo", "tee", path)
  cmd.Stdin = strings.NewReader(value)
  cmd.Stdout = nil
  cmd.Run()
}

func main() {
  // Reset i2c_hid_acpi devices
  pattern := "/sys/bus/i2c/drivers/i2c_hid_acpi/i2c-*"
  matches, _ := filepath.Glob(pattern)
  
  for _, dev := range matches {
    i2cDevice := filepath.Base(dev)
    fmt.Printf("Resetting %s via i2c_hid_acpi unbind/rebind...\n", i2cDevice)
    
    // Unbind
    unbindPath := "/sys/bus/i2c/drivers/i2c_hid_acpi/unbind"
    writeToSysfs(unbindPath, i2cDevice)
    
    time.Sleep(1 * time.Second)
    
    // Rebind
    bindPath := "/sys/bus/i2c/drivers/i2c_hid_acpi/bind"
    writeToSysfs(bindPath, i2cDevice)
    
    fmt.Println("Done")
  }

  // Check and reload intel_quicki2c if present
  lsmodCmd := exec.Command("lsmod")
  grepCmd := exec.Command("grep", "-q", "intel_quicki2c")
  
  lsmodOut, _ := lsmodCmd.StdoutPipe()
  grepCmd.Stdin = lsmodOut
  
  lsmodCmd.Start()
  if grepCmd.Run() == nil {
    fmt.Println("Reloading intel_quicki2c...")
    exec.Command("sudo", "modprobe", "-r", "intel_quicki2c").Run()
    exec.Command("sudo", "modprobe", "intel_quicki2c").Run()
  }
  lsmodCmd.Wait()
}
