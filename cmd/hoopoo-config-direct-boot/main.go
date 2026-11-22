package main

import (
  "fmt"
  "os"
  "os/exec"
  "path/filepath"
  "strings"
  
  "hoopoo/pkg/hardware/dmi"
  "hoopoo/pkg/toggle"
)

// hoopoo:summary=Add or remove an EFI boot entry for the Hoopoo UKI, allowing the system to boot directly
// hoopoo:requires-sudo=true

func main() {
  // Check if system is booted in UEFI mode
  if _, err := os.Stat("/sys/firmware/efi"); os.IsNotExist(err) {
    fmt.Fprintln(os.Stderr, "Error: System is not booted in UEFI mode")
    os.Exit(1)
  }
  
  // Check if efibootmgr is available
  cmd := exec.Command("efibootmgr")
  cmd.Stdout = nil
  cmd.Stderr = nil
  if err := cmd.Run(); err != nil {
    fmt.Fprintln(os.Stderr, "Error: efibootmgr is not available or not functional")
    os.Exit(1)
  }
  
  // Check for American Megatrends firmware
  if dmi.Match("bios_vendor", "american megatrends") {
    fmt.Fprintln(os.Stderr, "Error: American Megatrends firmware may not safely support custom EFI entries")
    os.Exit(1)
  }
  
  // Check for Apple firmware
  if dmi.Match("bios_vendor", "apple") {
    fmt.Fprintln(os.Stderr, "Error: Apple firmware uses its own boot manager")
    os.Exit(1)
  }
  
  // Check for existing Hoopoo EFI entry
  cmd = exec.Command("efibootmgr")
  output, err := cmd.Output()
  if err != nil {
    fmt.Fprintln(os.Stderr, "Error: Failed to read EFI boot entries")
    os.Exit(1)
  }
  
  var existingBootNum string
  for _, line := range strings.Split(string(output), "\n") {
    if strings.Contains(line, "Hoopoo") && strings.HasPrefix(line, "Boot") {
      // Extract boot number (e.g., "Boot0001" -> "0001")
      parts := strings.Fields(line)
      if len(parts) > 0 {
        bootEntry := strings.TrimPrefix(parts[0], "Boot")
        bootEntry = strings.TrimSuffix(bootEntry, "*")
        existingBootNum = bootEntry
        break
      }
    }
  }
  
  if existingBootNum != "" {
    // Direct boot is enabled, ask to disable
    if !confirmAction("Disable direct boot (remove Hoopoo EFI entry)?") {
      os.Exit(0)
    }
    
    fmt.Println("Removing EFI boot entry", existingBootNum)
    cmd := exec.Command("sudo", "efibootmgr", "--bootnum", existingBootNum, "--delete-bootnum")
    cmd.Stdout = nil
    if err := cmd.Run(); err != nil {
      fmt.Fprintln(os.Stderr, "Error: Failed to remove EFI entry")
      os.Exit(1)
    }
    
    toggle.Clear("direct-boot-enabled")
    os.Exit(0)
  }
  
  // Direct boot is disabled, ask to enable
  ukiFiles, err := filepath.Glob("/boot/EFI/Linux/hoopoo*.efi")
  if err != nil || len(ukiFiles) == 0 {
    fmt.Fprintln(os.Stderr, "Error: No Hoopoo UKI found in /boot/EFI/Linux/")
    os.Exit(1)
  }
  
  ukiFile := filepath.Base(ukiFiles[0])
  
  // Find boot partition
  cmd = exec.Command("findmnt", "-n", "-o", "SOURCE", "/boot")
  bootSource, err := cmd.Output()
  if err != nil {
    fmt.Fprintln(os.Stderr, "Error: Failed to find boot partition")
    os.Exit(1)
  }
  
  source := strings.TrimSpace(string(bootSource))
  
  // Extract disk and partition number
  disk := strings.TrimRight(source, "0123456789")
  disk = strings.TrimSuffix(disk, "p")
  
  partNum := strings.TrimPrefix(source, disk)
  partNum = strings.TrimPrefix(partNum, "p")
  
  if !confirmAction("Setup direct boot (so snapshot booting must be done via bios)?") {
    os.Exit(0)
  }
  
  fmt.Println("Creating EFI boot entry for", ukiFile)
  
  // Create EFI boot entry
  loaderPath := "\\EFI\\Linux\\" + ukiFile
  cmd = exec.Command("sudo", "efibootmgr",
    "--create",
    "--disk", disk,
    "--part", partNum,
    "--label", "Hoopoo",
    "--loader", loaderPath,
  )
  if err := cmd.Run(); err != nil {
    fmt.Fprintln(os.Stderr, "Error: Failed to create EFI entry")
    os.Exit(1)
  }
  
  // Create toggle state (toggle.Write creates directories automatically)
  toggle.Write("direct-boot-enabled", "")
}

func confirmAction(prompt string) bool {
  cmd := exec.Command("gum", "confirm", prompt)
  cmd.Stdin = os.Stdin
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr
  return cmd.Run() == nil
}
