package main

import (
  "errors"
  "fmt"
  "os"
  "strings"
  
  "hoopoo/pkg/portal"
)

// hoopoo:summary=Pick files with the desktop file chooser
// hoopoo:args=[--title <title>] [--multiple] [--directory] [--extensions "<ext ext...>"]
// hoopoo:examples=hoopoo file select --title "Send with LocalSend" --multiple | hoopoo file select --title "Pick image" --extensions "png svg" | hoopoo file select --title "Share folder" --directory

func main() {
  opts := portal.FileSelectOptions{}
  
  // Parse arguments
  args := os.Args[1:]
  for i := 0; i < len(args); i++ {
    switch args[i] {
    case "--title":
      if i+1 < len(args) {
        opts.Title = args[i+1]
        i++
      }
    case "--multiple":
      opts.Multiple = true
    case "--directory":
      opts.Directory = true
    case "--extensions":
      if i+1 < len(args) {
        opts.Extensions = splitExtensions(args[i+1])
        i++
      }
    case "--help", "-h":
      printUsage()
      os.Exit(0)
    }
  }
  
  // Call portal
  paths, err := portal.SelectFiles(opts)
  if err != nil {
    // Exit code 1 for user cancellation, 2 for errors
    if errors.Is(err, portal.ErrNoFilesSelected) || errors.Is(err, portal.ErrSelectionCancelled) {
      os.Exit(1)
    }
    fmt.Fprintf(os.Stderr, "hoopoo-file-select: %v\n", err)
    os.Exit(2)
  }
  
  // Print selected paths
  for _, path := range paths {
    fmt.Println(path)
  }
}

func splitExtensions(s string) []string {
  var exts []string
  for _, ext := range strings.Fields(s) {
    if ext != "" {
      exts = append(exts, ext)
    }
  }
  return exts
}

func printUsage() {
  fmt.Println(`hoopoo-file-select - Pick files with the desktop file chooser

Usage:
  hoopoo file select [OPTIONS]

Options:
  --title <title>         Dialog title
  --multiple              Allow multiple file selection
  --directory             Select directories instead of files
  --extensions "<exts>"   File extensions to filter (space-separated)
  -h, --help              Show this help

Examples:
  hoopoo file select --title "Send with LocalSend" --multiple
  hoopoo file select --title "Pick image" --extensions "png jpg svg"
  hoopoo file select --title "Share folder" --directory`)
}
