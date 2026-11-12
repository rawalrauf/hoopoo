// Keyboard backlight brightness control.
package keyboard

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"hoopoo/pkg/system/paths"
)

// FindKeyboardDevice searches for keyboard backlight device.
// Checks both /sys/class/leds and /sys/class/backlight.
// Caches result for performance (device doesn't change during runtime).
var cachedDevice string
var deviceFound bool

func FindKeyboardDevice() (string, error) {
	if deviceFound {
		return cachedDevice, nil
	}
	
	for _, class := range []string{"leds", "backlight"} {
		pattern := filepath.Join("/sys/class", class, "*kbd_backlight*")
		matches, err := filepath.Glob(pattern)
		if err != nil {
			continue
		}
		for _, match := range matches {
			if _, err := os.Stat(match); err == nil {
				cachedDevice = filepath.Base(match)
				deviceFound = true
				return cachedDevice, nil
			}
		}
	}
	return "", fmt.Errorf("no keyboard backlight device found")
}

// GetBrightness reads current brightness level directly from sysfs.
// Cache device path to avoid repeated lookups.
var cachedBrightnessPath string

func GetBrightness(device string) (int, error) {
	if cachedBrightnessPath == "" {
		if _, err := os.Stat(filepath.Join("/sys/class/leds", device, "brightness")); err == nil {
			cachedBrightnessPath = filepath.Join("/sys/class/leds", device, "brightness")
		} else {
			cachedBrightnessPath = filepath.Join("/sys/class/backlight", device, "brightness")
		}
	}
	
	data, err := os.ReadFile(cachedBrightnessPath)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.TrimSpace(string(data)))
}

// GetMaxBrightness reads maximum brightness level.
// Cache device path to avoid repeated lookups.
var cachedMaxPath string

func GetMaxBrightness(device string) (int, error) {
	if cachedMaxPath == "" {
		if _, err := os.Stat(filepath.Join("/sys/class/leds", device, "max_brightness")); err == nil {
			cachedMaxPath = filepath.Join("/sys/class/leds", device, "max_brightness")
		} else {
			cachedMaxPath = filepath.Join("/sys/class/backlight", device, "max_brightness")
		}
	}
	
	data, err := os.ReadFile(cachedMaxPath)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.TrimSpace(string(data)))
}

// CalculateStep returns appropriate step size (10% of max, minimum 1).
func CalculateStep(maxBrightness int) int {
	step := maxBrightness / 10
	if step < 1 {
		return 1
	}
	return step
}

// CalculateNewBrightness computes new brightness based on direction.
func CalculateNewBrightness(current, max, step int, direction string) int {
	switch direction {
	case "up":
		new := current + step
		if new > max {
			return max
		}
		return new
	case "down":
		new := current - step
		if new < 0 {
			return 0
		}
		return new
	case "cycle":
		new := current + step
		if new > max {
			return 0
		}
		return new
	default:
		return current
	}
}

// StateFile returns path to persistent state file.
func StateFile() string {
	return filepath.Join(paths.State(), "kbd-brightness")
}

// TrackerFile returns path to runtime tracker file.
func TrackerFile() string {
	return filepath.Join(paths.Runtime(), "hoopoo", "kbd-brightness-tracker")
}

// ReadStateFile reads brightness value from state file with validation.
func ReadStateFile(path string) (int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	value := strings.TrimSpace(string(data))
	level, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid brightness value: %s", value)
	}
	return level, nil
}

// WriteStateFile writes brightness value to state file with directory creation.
// Cache directory to avoid repeated os.MkdirAll calls.
var cachedStateDir string

func WriteStateFile(path string, level int) error {
	if cachedStateDir == "" {
		cachedStateDir = filepath.Dir(path)
		if err := os.MkdirAll(cachedStateDir, 0755); err != nil {
			return err
		}
	}
	return os.WriteFile(path, []byte(strconv.Itoa(level)), 0644)
}

// LockFile creates an exclusive lock on a file (like flock).
// Cache lock directory to avoid repeated os.MkdirAll calls.
var cachedLockDir string

func LockFile(path string) (*os.File, error) {
	lockPath := path + ".lock"
	
	if cachedLockDir == "" {
		cachedLockDir = filepath.Dir(lockPath)
		if err := os.MkdirAll(cachedLockDir, 0755); err != nil {
			return nil, err
		}
	}
	
	f, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		f.Close()
		return nil, err
	}
	
	return f, nil
}
