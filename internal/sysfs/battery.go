package sysfs

import (
  "fmt"
  "os"
  "path/filepath"
  "strconv"
  "strings"
)

// PowerSupplyPath can be overridden via HOOPOO_POWER_SUPPLY_PATH for testing
var PowerSupplyPath = getPowerSupplyPath()

func getPowerSupplyPath() string {
  if path := os.Getenv("HOOPOO_POWER_SUPPLY_PATH"); path != "" {
    return path
  }
  return "/sys/class/power_supply"
}

// Battery represents a system battery
type Battery struct {
  Path string
  Name string
}

// cachedBatteries holds the battery scan result to avoid repeated directory scans
var cachedBatteries []*Battery
var batteriesScanned bool

// FindAllBatteries returns all system batteries (cached after first call)
func FindAllBatteries() ([]*Battery, error) {
  // Return cached result if already scanned
  if batteriesScanned {
    if len(cachedBatteries) == 0 {
      return nil, fmt.Errorf("no batteries found")
    }
    return cachedBatteries, nil
  }

  entries, err := os.ReadDir(PowerSupplyPath)
  if err != nil {
    return nil, fmt.Errorf("failed to read power supply: %w", err)
  }

  var batteries []*Battery
  for _, entry := range entries {
    path := filepath.Join(PowerSupplyPath, entry.Name())
    
    // Follow symlinks - power_supply entries are often symlinks
    info, err := os.Stat(path)
    if err != nil || !info.IsDir() {
      continue
    }

    // Check if it's a battery
    batteryType, err := readString(filepath.Join(path, "type"))
    if err != nil || batteryType != "Battery" {
      continue
    }

    // Check if present
    present, err := readInt(filepath.Join(path, "present"))
    if err != nil || present != 1 {
      continue
    }

    // Filter out wireless peripherals (mice, keyboards)
    scope, err := readString(filepath.Join(path, "scope"))
    if err == nil && scope == "Device" {
      continue
    }

    batteries = append(batteries, &Battery{
      Path: path,
      Name: entry.Name(),
    })
  }

  // Cache the result
  cachedBatteries = batteries
  batteriesScanned = true

  if len(batteries) == 0 {
    return nil, fmt.Errorf("no batteries found")
  }

  return batteries, nil
}

// IsPresent checks if any battery exists
func IsPresent() bool {
  batteries, err := FindAllBatteries()
  return err == nil && len(batteries) > 0
}

// GetCombinedCapacity returns total battery percentage across all batteries
func GetCombinedCapacity() (int, error) {
  batteries, err := FindAllBatteries()
  if err != nil {
    return 0, err
  }

  // For single battery, use simple sequential reads
  if len(batteries) == 1 {
    bat := batteries[0]
    energyNow, err := readInt(filepath.Join(bat.Path, "energy_now"))
    if err != nil {
      energyNow, err = readInt(filepath.Join(bat.Path, "charge_now"))
      if err != nil {
        return 0, err
      }
    }

    energyFull, err := readInt(filepath.Join(bat.Path, "energy_full"))
    if err != nil {
      energyFull, err = readInt(filepath.Join(bat.Path, "charge_full"))
      if err != nil {
        return 0, err
      }
    }

    if energyFull == 0 {
      return 0, fmt.Errorf("no valid battery capacity data")
    }

    percentage := (energyNow*100 + energyFull/2) / energyFull
    return percentage, nil
  }

  // For multiple batteries, read in parallel with goroutines
  type batteryData struct {
    energyNow  int
    energyFull int
  }

  ch := make(chan batteryData, len(batteries))

  for _, bat := range batteries {
    go func(b *Battery) {
      var data batteryData

      // Try energy_now, fallback to charge_now
      data.energyNow, _ = readInt(filepath.Join(b.Path, "energy_now"))
      if data.energyNow == 0 {
        data.energyNow, _ = readInt(filepath.Join(b.Path, "charge_now"))
      }

      // Try energy_full, fallback to charge_full
      data.energyFull, _ = readInt(filepath.Join(b.Path, "energy_full"))
      if data.energyFull == 0 {
        data.energyFull, _ = readInt(filepath.Join(b.Path, "charge_full"))
      }

      ch <- data
    }(bat)
  }

  var totalNow, totalFull int
  for range batteries {
    data := <-ch
    totalNow += data.energyNow
    totalFull += data.energyFull
  }

  if totalFull == 0 {
    return 0, fmt.Errorf("no valid battery capacity data")
  }

  // Round to nearest integer
  percentage := (totalNow*100 + totalFull/2) / totalFull
  return percentage, nil
}

// GetStatus returns battery status (Charging, Discharging, Full, etc.)
func (b *Battery) GetStatus() (string, error) {
  return readString(filepath.Join(b.Path, "status"))
}

// Helper functions

// readString reads a sysfs file and returns trimmed string
func readString(path string) (string, error) {
  data, err := os.ReadFile(path)
  if err != nil {
    return "", err
  }
  return strings.TrimSpace(string(data)), nil
}

// readInt reads a sysfs file and returns integer value
func readInt(path string) (int, error) {
  data, err := os.ReadFile(path)
  if err != nil {
    return 0, err
  }
  
  // Manually trim and parse to avoid string allocation
  str := string(data)
  start := 0
  end := len(str)
  
  // Trim leading/trailing whitespace
  for start < end && (str[start] == ' ' || str[start] == '\t' || str[start] == '\n' || str[start] == '\r') {
    start++
  }
  for end > start && (str[end-1] == ' ' || str[end-1] == '\t' || str[end-1] == '\n' || str[end-1] == '\r') {
    end--
  }
  
  return strconv.Atoi(str[start:end])
}
