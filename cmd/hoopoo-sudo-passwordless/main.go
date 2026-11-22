package main

import (
  "fmt"
  "os"
  "os/exec"
  "strconv"
  "strings"
)

// hoopoo:summary=Toggle passwordless sudo for the current user
// hoopoo:args=[MINUTES]
// hoopoo:requires-sudo=true

func main() {
  user := os.Getenv("USER")
  nopasswdFile := "/etc/sudoers.d/99-hoopoo-nopasswd-" + user
  timerName := "hoopoo-nopasswd-expire-" + user
  flagFile := "/tmp/hoopoo-sudo-passwordless-" + user
  
  minutes := 15
  if len(os.Args) > 1 {
    var err error
    minutes, err = strconv.Atoi(os.Args[1])
    if err != nil {
      fmt.Fprintln(os.Stderr, "Usage: hoopoo sudo passwordless [MINUTES]")
      os.Exit(1)
    }
  }
  
  fmt.Println("Toggle passwordless sudo...")
  
  // Safety: if file exists but timer doesn't (e.g., after reboot), clean up
  fileExists := exec.Command("sudo", "test", "-f", nopasswdFile).Run() == nil
  timerActive := exec.Command("systemctl", "is-active", "--quiet", timerName+".timer").Run() == nil
  
  if fileExists && !timerActive {
    if err := exec.Command("sudo", "rm", nopasswdFile).Run(); err != nil {
      fmt.Fprintf(os.Stderr, "Warning: Failed to clean up orphaned sudoers file: %v\n", err)
    }
    os.Remove(flagFile)
    fileExists = false
  }
  
  if fileExists {
    // Passwordless sudo is enabled
    if len(os.Args) > 1 {
      // Update timer duration
      exec.Command("systemctl", "stop", timerName+".timer").Run()
      if !armExpiry(user, nopasswdFile, flagFile, timerName, minutes) {
        os.Exit(1)
      }
      fmt.Printf("Passwordless sudo timer updated. It will now automatically disable in %d minutes.\n", minutes)
    } else {
      // Disable passwordless sudo
      if err := exec.Command("sudo", "rm", nopasswdFile).Run(); err != nil {
        fmt.Fprintf(os.Stderr, "Error: Failed to remove sudoers file: %v\n", err)
        os.Exit(1)
      }
      exec.Command("systemctl", "stop", timerName+".timer").Run()
      os.Remove(flagFile)
      fmt.Println("Passwordless sudo has been DISABLED. Sudo will require a password again.")
    }
  } else {
    // Passwordless sudo is disabled, ask to enable
    fmt.Println()
    fmt.Println("⚠️ WARNING: This will allow ANY process running as your user to")
    fmt.Printf("execute ANY command as root WITHOUT a password for %d minutes.\n", minutes)
    fmt.Println()
    fmt.Println("This is useful for AI agents that need to run sudo commands,")
    fmt.Println("but it significantly weakens the security of your system.")
    fmt.Println("Anyone or anything with access to your user account gets full root.")
    fmt.Println()
    fmt.Printf("Passwordless sudo will automatically disable after %d minutes.\n", minutes)
    fmt.Println("Run this command again to disable it early.")
    fmt.Println()
    
    prompt := fmt.Sprintf("Enable passwordless sudo for %d minutes? This is a significant security risk!", minutes)
    cmd := exec.Command("gum", "confirm", prompt)
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    
    if cmd.Run() != nil {
      fmt.Println("Aborted. No changes made.")
      os.Exit(0)
    }
    
    // Create sudoers file
    content := fmt.Sprintf("%s ALL=(ALL) NOPASSWD: ALL\n", user)
    writeCmd := exec.Command("sudo", "tee", nopasswdFile)
    writeCmd.Stdin = strings.NewReader(content)
    writeCmd.Stdout = nil
    if err := writeCmd.Run(); err != nil {
      fmt.Fprintln(os.Stderr, "Error: Failed to create sudoers file")
      os.Exit(1)
    }
    
    if err := exec.Command("sudo", "chmod", "440", nopasswdFile).Run(); err != nil {
      fmt.Fprintf(os.Stderr, "Error: Failed to set sudoers file permissions: %v\n", err)
      exec.Command("sudo", "rm", nopasswdFile).Run()
      os.Exit(1)
    }
    
    if !armExpiry(user, nopasswdFile, flagFile, timerName, minutes) {
      os.Exit(1)
    }
    
    os.WriteFile(flagFile, []byte(""), 0644)
    
    fmt.Println()
    fmt.Printf("Passwordless sudo has been ENABLED. It will automatically disable in %d minutes.\n", minutes)
    fmt.Println("A restart removes the passwordless sudo rule as well.")
  }
}

func armExpiry(user, nopasswdFile, flagFile, timerName string, minutes int) bool {
  onActive := fmt.Sprintf("--on-active=%dm", minutes)
  cleanupCmd := fmt.Sprintf("rm -f -- '%s'; rm -f '%s'", nopasswdFile, flagFile)
  
  cmd := exec.Command("sudo", "systemd-run",
    onActive,
    "--timer-property=AccuracySec=1s",
    "--unit="+timerName,
    "bash", "-c", cleanupCmd,
  )
  
  if err := cmd.Run(); err != nil {
    fmt.Fprintln(os.Stderr, "Failed to schedule passwordless sudo expiry. Revoking access now.")
    if err := exec.Command("sudo", "rm", "-f", nopasswdFile).Run(); err != nil {
      fmt.Fprintf(os.Stderr, "CRITICAL: Could not remove %s. Remove it as root immediately.\n", nopasswdFile)
    }
    os.Remove(flagFile)
    return false
  }
  
  return true
}
