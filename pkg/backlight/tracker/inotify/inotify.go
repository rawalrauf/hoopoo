// Event-driven keyboard backlight tracker using inotify.
package inotify

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// FindHwChangedFile searches for brightness_hw_changed file.
// Tries multiple vendor-specific paths.
func FindHwChangedFile() (string, error) {
	patterns := []string{
		"/sys/devices/platform/*/leds/*kbd_backlight*/brightness_hw_changed",
		"/sys/devices/platform/*/backlight/*kbd_backlight*/brightness_hw_changed",
		"/sys/class/leds/*kbd_backlight*/brightness_hw_changed",
		"/sys/class/backlight/*kbd_backlight*/brightness_hw_changed",
	}
	
	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			continue
		}
		for _, match := range matches {
			if _, err := os.Stat(match); err == nil {
				return match, nil
			}
		}
	}
	
	return "", fmt.Errorf("no brightness_hw_changed file found")
}

// ReadBrightness reads brightness from device sysfs path.
// Cache device path to avoid repeated lookups.
var cachedDevicePath string

func ReadBrightness(device string) (int, error) {
	if cachedDevicePath == "" {
		if _, err := os.Stat(filepath.Join("/sys/class/leds", device, "brightness")); err == nil {
			cachedDevicePath = filepath.Join("/sys/class/leds", device, "brightness")
		} else {
			cachedDevicePath = filepath.Join("/sys/class/backlight", device, "brightness")
		}
	}
	
	data, err := os.ReadFile(cachedDevicePath)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.TrimSpace(string(data)))
}

// SaveBrightness writes brightness to tracker file.
// Cache directory to avoid repeated os.MkdirAll calls.
var cachedTrackerDir string

func SaveBrightness(trackerFile string, value int) error {
	if cachedTrackerDir == "" {
		cachedTrackerDir = filepath.Dir(trackerFile)
		if err := os.MkdirAll(cachedTrackerDir, 0755); err != nil {
			return err
		}
	}
	return os.WriteFile(trackerFile, []byte(strconv.Itoa(value)), 0644)
}

// SaveInitialBrightness saves current brightness on service start.
func SaveInitialBrightness(device, trackerFile string) error {
	current, err := ReadBrightness(device)
	if err != nil {
		return err
	}
	// Save all values including zero (no filtering)
	return SaveBrightness(trackerFile, current)
}
