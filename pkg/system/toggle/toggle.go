// Toggle state management

package toggle

import (
  "os"
  "path/filepath"

  "hoopoo/pkg/system/paths"
)

func HoopooEnabled(name string) bool {
  path := filepath.Join(paths.State(), "toggles", name)
  _, err := os.Stat(path)
  return err == nil
}

func HoopooDisabled(name string) bool {
  return !HoopooEnabled(name)
}

func HoopooClear(name string) error {
  path := filepath.Join(paths.State(), "toggles", name)
  err := os.Remove(path)
  if os.IsNotExist(err) {
    return nil
  }
  return err
}

func NiriEnabled(name string) bool {
  path := filepath.Join(paths.State(), "toggles", "niri", name+".kdl")
  _, err := os.Stat(path)
  return err == nil
}

func NiriDisabled(name string) bool {
  return !NiriEnabled(name)
}

func NiriClear(name string) error {
  path := filepath.Join(paths.State(), "toggles", "niri", name+".kdl")
  err := os.Remove(path)
  if os.IsNotExist(err) {
    return nil
  }
  return err
}
