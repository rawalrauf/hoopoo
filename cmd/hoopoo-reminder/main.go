// hoopoo:summary=Set and show lightweight desktop notification reminders
// hoopoo:args=<minutes> [message] | show | clear
// hoopoo:examples=hoopoo reminder 5 | hoopoo reminder 30 "Check the oven" | hoopoo reminder show | hoopoo reminder clear

package main

import (
  "fmt"
  "os"
  "os/exec"
  "path/filepath"
  "strconv"
  "strings"
  "time"

  "hoopoo/pkg/system/paths"
)

func main() {
  if len(os.Args) < 2 {
    printUsage()
    os.Exit(1)
  }

  cmd := os.Args[1]

  switch cmd {
  case "show", "list":
    showReminders()
  case "clear":
    clearReminders()
  default:
    // Set a reminder
    minutes, err := strconv.Atoi(os.Args[1])
    if err != nil || minutes <= 0 {
      printUsage()
      os.Exit(1)
    }
    
    message := ""
    if len(os.Args) > 2 {
      message = strings.Join(os.Args[2:], " ")
    }
    
    setReminder(minutes, message)
  }
}

func printUsage() {
  fmt.Fprintln(os.Stderr, "Usage: hoopoo reminder <minutes> [message]")
  fmt.Fprintln(os.Stderr, "       hoopoo reminder show")
  fmt.Fprintln(os.Stderr, "       hoopoo reminder clear")
}

func setReminder(minutes int, customMessage string) {
  setAt := time.Now().Unix()
  remindAt := time.Now().Add(time.Duration(minutes) * time.Minute)
  unit := fmt.Sprintf("hoopoo-reminder-%dm-%d", minutes, setAt)
  
  reminderDir := filepath.Join(paths.Runtime(), "hoopoo-reminders")
  os.MkdirAll(reminderDir, 0755)
  
  message := fmt.Sprintf("Your %d minutes are up", minutes)
  confirmationTitle := fmt.Sprintf("Reminder set for %d minutes", minutes)
  confirmation := fmt.Sprintf("You'll be reminded at %s", remindAt.Format("15:04"))
  
  messageFile := filepath.Join(reminderDir, unit+".message")
  
  if customMessage != "" {
    message = customMessage
    os.WriteFile(messageFile, []byte(customMessage), 0644)
    confirmationTitle = fmt.Sprintf("%s in %d minutes", customMessage, minutes)
  }
  
  // Schedule systemd timer to send notification
  bashCmd := fmt.Sprintf(`hoopoo notification send --glyph=󰔛 "Reminder" "%s" --sound=/usr/share/sounds/freedesktop/stereo/complete.oga; rm -f "%s"`, message, messageFile)
  
  cmd := exec.Command("systemd-run",
    "--user",
    "--quiet",
    "--collect",
    fmt.Sprintf("--on-active=%dm", minutes),
    fmt.Sprintf("--unit=%s", unit),
    "bash", "-c", bashCmd,
  )
  
  if err := cmd.Run(); err != nil {
    fmt.Fprintf(os.Stderr, "Error scheduling reminder: %v\n", err)
    os.Exit(1)
  }
  
  // Send confirmation notification
  exec.Command("hoopoo", "notification", "send",
    "--glyph=󰔛",
    confirmationTitle,
    confirmation,
    "-u", "low",
  ).Run()
}

func showReminders() {
  cmd := exec.Command("systemctl", "--user", "list-timers", "--all", "--no-legend", "--no-pager", "hoopoo-reminder-*.timer")
  output, err := cmd.Output()
  if err != nil || strings.TrimSpace(string(output)) == "" {
    exec.Command("hoopoo", "notification", "send",
      "--glyph=󰔛",
      "Upcoming reminders",
      "No outstanding reminders",
      "-u", "low",
    ).Run()
    return
  }
  
  lines := strings.Split(strings.TrimSpace(string(output)), "\n")
  
  // TODO: Parse timer details and format remaining time
  body := fmt.Sprintf("Found %d active reminders", len(lines))
  
  exec.Command("hoopoo", "notification", "send",
    "--glyph=󰔛",
    "Upcoming reminders",
    body,
    "-u", "low",
  ).Run()
}

func clearReminders() {
  cmd := exec.Command("systemctl", "--user", "list-timers", "--all", "--no-legend", "--no-pager", "hoopoo-reminder-*.timer")
  output, _ := cmd.Output()
  
  lines := strings.Split(strings.TrimSpace(string(output)), "\n")
  for _, line := range lines {
    if line == "" {
      continue
    }
    fields := strings.Fields(line)
    if len(fields) >= 2 {
      timer := fields[len(fields)-2]
      service := fields[len(fields)-1]
      exec.Command("systemctl", "--user", "stop", timer, service).Run()
    }
  }
  
  // Clean up message files
  reminderDir := filepath.Join(paths.Runtime(), "hoopoo-reminders")
  files, _ := filepath.Glob(filepath.Join(reminderDir, "hoopoo-reminder-*.message"))
  for _, f := range files {
    os.Remove(f)
  }
  
  exec.Command("hoopoo", "notification", "send",
    "--glyph=󰔛",
    "All reminders have been cleared",
    "-u", "low",
  ).Run()
}
