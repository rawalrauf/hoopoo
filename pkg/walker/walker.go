// Walker launcher utilities

package walker

import (
  "fmt"
  "os"
  "os/exec"
  "strings"
  "syscall"

  "hoopoo/pkg/system/process"
)

// EnsureServices starts elephant and walker service if not running
func EnsureServices() {
  go func() {
    if !process.IsRunning("elephant") {
      cmd := exec.Command("elephant")
      cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
      cmd.Start()
    }
  }()

  go func() {
    if exec.Command("pgrep", "-f", "walker --gapplication-service").Run() != nil {
      cmd := exec.Command("walker", "--gapplication-service")
      cmd.Env = append(os.Environ(), "GSK_RENDERER=cairo")
      cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
      cmd.Start()
    }
  }()
}

// Launch starts walker with launcher defaults
func Launch(extraArgs []string) error {
  EnsureServices()

  // Check what dimension args user provided and convert --height to --maxheight/--minheight
  userHasWidth := false
  userHasHeight := false
  convertedArgs := []string{}
  
  for i := 0; i < len(extraArgs); i++ {
    arg := extraArgs[i]
    
    if arg == "--width" {
      userHasWidth = true
      convertedArgs = append(convertedArgs, arg)
      if i+1 < len(extraArgs) {
        convertedArgs = append(convertedArgs, extraArgs[i+1])
        i++
      }
    } else if arg == "--height" {
      userHasHeight = true
      if i+1 < len(extraArgs) {
        height := extraArgs[i+1]
        // Convert --height to --maxheight + --minheight
        convertedArgs = append(convertedArgs, "--maxheight", height, "--minheight", height)
        i++
      }
    } else if arg == "--maxheight" || arg == "--minheight" {
      userHasHeight = true
      convertedArgs = append(convertedArgs, arg)
      if i+1 < len(extraArgs) {
        convertedArgs = append(convertedArgs, extraArgs[i+1])
        i++
      }
    } else {
      convertedArgs = append(convertedArgs, arg)
    }
  }

  dims := GetLauncherDefaults() // Get font-aware launcher defaults
  args := []string{}
  
  // Only add defaults if user didn't provide them
  if !userHasWidth {
    args = append(args, "--width", dims.Width)
  }
  
  if !userHasHeight {
    args = append(args,
      "--maxheight", fmt.Sprintf("%d", dims.MaxHeight),
      "--minheight", fmt.Sprintf("%d", dims.MinHeight),
    )
  }
  
  // Append converted user args
  args = append(args, convertedArgs...)

  cmd := exec.Command("walker", args...)
  cmd.Stdin = os.Stdin
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr
  return cmd.Run()
}

// ShowDmenu displays a dmenu with dynamic height based on item count
func ShowDmenu(prompt, options string, itemCount int) (string, error) {
  EnsureServices()

  dims := Get()
  height := calculateHeight(itemCount, dims)

  cmd := exec.Command(
    "walker",
    "--dmenu",
    "--width", dims.Width,
    "--minheight", "1",
    "--maxheight", height,
    "-p", prompt+"…",
  )

  cmd.Stdin = strings.NewReader(options)
  output, err := cmd.Output()
  if err != nil {
    return "", err
  }

  return strings.TrimSpace(string(output)), nil
}

// calculateHeight returns dynamic menu height based on item count and dimensions
func calculateHeight(itemCount int, dims Dimensions) string {
  height := itemCount * dims.ItemHeight

  if height < dims.MinHeight {
    return fmt.Sprintf("%d", dims.MinHeight)
  }
  if height > dims.MaxHeight {
    return fmt.Sprintf("%d", dims.MaxHeight)
  }

  return fmt.Sprintf("%d", height)
}
