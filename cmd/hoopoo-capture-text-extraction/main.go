package main

import (
  "bytes"
  "os"
  "os/exec"
  "strings"
)

// hoopoo:summary=Extract text from a screenshot region with OCR
// hoopoo:group=capture
// hoopoo:examples=hoopoo capture text extraction

func main() {
  // Get selection region with slurp
  cmd := exec.Command("slurp")
  output, err := cmd.Output()
  if err != nil {
    os.Exit(0) // User cancelled
  }
  
  selection := strings.TrimSpace(string(output))
  if selection == "" {
    os.Exit(0)
  }
  
  // Capture screenshot with grim
  grimCmd := exec.Command("grim", "-g", selection, "-")
  screenshot, err := grimCmd.Output()
  if err != nil {
    os.Exit(1)
  }
  
  // Extract text with tesseract
  langs := os.Getenv("HOOPOO_OCR_LANGS")
  if langs == "" {
    langs = "eng"
  }
  
  tesseractCmd := exec.Command("tesseract", "stdin", "stdout",
    "--oem", "1",
    "--psm", "6",
    "-l", langs,
    "--dpi", "300",
    "-c", "preserve_interword_spaces=1",
  )
  tesseractCmd.Stdin = bytes.NewReader(screenshot)
  
  var stderr bytes.Buffer
  tesseractCmd.Stderr = &stderr
  
  text, err := tesseractCmd.Output()
  if err != nil || len(text) == 0 {
    os.Exit(1)
  }
  
  extractedText := strings.TrimSpace(string(text))
  if extractedText == "" {
    os.Exit(1)
  }
  
  // Copy to clipboard
  wlCopyCmd := exec.Command("wl-copy")
  wlCopyCmd.Stdin = strings.NewReader(extractedText)
  wlCopyCmd.Run()
  
  // Show notification
  exec.Command("hoopoo", "notification", "send",
    "--glyph=󰴑",
    "Copied text from selection to clipboard",
  ).Run()
}
