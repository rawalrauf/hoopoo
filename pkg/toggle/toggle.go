// Toggle state management for generic toggle files

package toggle

import (
  "os"
  "path/filepath"

  "hoopoo/pkg/system/paths"
)

// Enabled checks if a toggle file exists in the state directory
// path is relative to ~/.local/state/hoopoo/toggles/
// Example: Enabled("niri/gaps-enabled.kdl") or Enabled("waybar-transparent")
func Enabled(path string) bool {
  fullPath := filepath.Join(paths.State(), "toggles", path)
  _, err := os.Stat(fullPath)
  return err == nil
}

// Disabled is the inverse of Enabled
func Disabled(path string) bool {
  return !Enabled(path)
}

// Clear removes a toggle file from the state directory
func Clear(path string) error {
  fullPath := filepath.Join(paths.State(), "toggles", path)
  err := os.Remove(fullPath)
  if os.IsNotExist(err) {
    return nil
  }
  return err
}

// Write creates or overwrites a toggle file with the given content
func Write(path, content string) error {
  fullPath := filepath.Join(paths.State(), "toggles", path)
  
  if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
    return err
  }
  
  return os.WriteFile(fullPath, []byte(content), 0644)
}
