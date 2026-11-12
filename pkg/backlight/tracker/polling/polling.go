// Polling-based keyboard backlight tracker.
package polling

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// FindBrightnessFile searches for keyboard backlight brightness file.
// Checks both /sys/class/leds and /sys/class/backlight.
func FindBrightnessFile() (string, error) {
	for _, class := range []string{"leds", "backlight"} {
		pattern := filepath.Join("/sys/class", class, "*kbd_backlight*")
		matches, err := filepath.Glob(pattern)
		if err != nil {
			continue
		}
		for _, match := range matches {
			brightnessFile := filepath.Join(match, "brightness")
			if _, err := os.Stat(brightnessFile); err == nil {
				return brightnessFile, nil
			}
		}
	}
	return "", nil // Not found, exit gracefully
}

// ReadBrightness reads brightness from file.
func ReadBrightness(brightnessFile string) (int, error) {
	data, err := os.ReadFile(brightnessFile)
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
func SaveInitialBrightness(brightnessFile, trackerFile string) error {
	current, err := ReadBrightness(brightnessFile)
	if err != nil {
		return err
	}
	// Save all values including zero (no filtering)
	return SaveBrightness(trackerFile, current)
}
