package niri

import (
  "hoopoo/pkg/toggle"
)

// WriteState writes KDL content to a niri toggle file
// filename should include .kdl extension if needed
func WriteState(filename, content string) error {
  return toggle.Write("niri/"+filename, content)
}

// StateEnabled checks if a niri toggle file exists
func StateEnabled(filename string) bool {
  return toggle.Enabled("niri/" + filename)
}

// StateDisabled checks if a niri toggle file does not exist
func StateDisabled(filename string) bool {
  return toggle.Disabled("niri/" + filename)
}

// ClearState removes a niri toggle file
func ClearState(filename string) error {
  return toggle.Clear("niri/" + filename)
}
