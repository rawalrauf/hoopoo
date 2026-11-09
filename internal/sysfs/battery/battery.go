// Battery information from sysfs

package battery

import (
  "fmt"
  "os"
  "path/filepath"
  "strconv"
  "strings"
)

var PowerSupplyPath = getPowerSupplyPath()

func getPowerSupplyPath() string {
  if path := os.Getenv("HOOPOO_POWER_SUPPLY_PATH"); path != "" {
    return path
  }
  return "/sys/class/power_supply"
}

type Battery struct {
  Path string
  Name string
}

var cachedBatteries []*Battery
var batteriesScanned bool

func FindAll() ([]*Battery, error) {
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
    
    info, err := os.Stat(path)
    if err != nil || !info.IsDir() {
      continue
    }

    batteryType, err := readString(filepath.Join(path, "type"))
    if err != nil || batteryType != "Battery" {
      continue
    }

    present, err := readInt(filepath.Join(path, "present"))
    if err != nil || present != 1 {
      continue
    }

    scope, err := readString(filepath.Join(path, "scope"))
    if err == nil && scope == "Device" {
      continue
    }

    batteries = append(batteries, &Battery{
      Path: path,
      Name: entry.Name(),
    })
  }

  cachedBatteries = batteries
  batteriesScanned = true

  if len(batteries) == 0 {
    return nil, fmt.Errorf("no batteries found")
  }

  return batteries, nil
}

func IsPresent() bool {
  batteries, err := FindAll()
  return err == nil && len(batteries) > 0
}

func CombinedCapacity() (int, error) {
  batteries, err := FindAll()
  if err != nil {
    return 0, err
  }

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

  type batteryData struct {
    energyNow  int
    energyFull int
  }

  ch := make(chan batteryData, len(batteries))

  for _, bat := range batteries {
    go func(b *Battery) {
      var data batteryData

      data.energyNow, _ = readInt(filepath.Join(b.Path, "energy_now"))
      if data.energyNow == 0 {
        data.energyNow, _ = readInt(filepath.Join(b.Path, "charge_now"))
      }

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

  percentage := (totalNow*100 + totalFull/2) / totalFull
  return percentage, nil
}

func (b *Battery) Status() (string, error) {
  return readString(filepath.Join(b.Path, "status"))
}

func readString(path string) (string, error) {
  data, err := os.ReadFile(path)
  if err != nil {
    return "", err
  }
  return strings.TrimSpace(string(data)), nil
}

func readInt(path string) (int, error) {
  data, err := os.ReadFile(path)
  if err != nil {
    return 0, err
  }
  
  str := string(data)
  start := 0
  end := len(str)
  
  for start < end && (str[start] == ' ' || str[start] == '\t' || str[start] == '\n' || str[start] == '\r') {
    start++
  }
  for end > start && (str[end-1] == ' ' || str[end-1] == '\t' || str[end-1] == '\n' || str[end-1] == '\r') {
    end--
  }
  
  return strconv.Atoi(str[start:end])
}
