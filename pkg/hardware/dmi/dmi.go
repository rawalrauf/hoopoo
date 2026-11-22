// DMI system information reading

package dmi

import (
  "os"
  "path/filepath"
  "strings"
)

// Read reads a DMI file from /sys/class/dmi/id/
func Read(file string) (string, error) {
  data, err := os.ReadFile(filepath.Join("/sys/class/dmi/id", file))
  if err != nil {
    return "", err
  }
  return strings.TrimSpace(string(data)), nil
}

// Match checks if a DMI file contains a pattern (case-insensitive)
func Match(file, pattern string) bool {
  content, err := Read(file)
  if err != nil {
    return false
  }
  return strings.Contains(strings.ToLower(content), strings.ToLower(pattern))
}

// Vendor returns the system vendor
func Vendor() (string, error) {
  return Read("sys_vendor")
}

// ProductName returns the product name
func ProductName() (string, error) {
  return Read("product_name")
}

// ProductFamily returns the product family
func ProductFamily() (string, error) {
  return Read("product_family")
}
