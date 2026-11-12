// hoopoo:summary=Track keyboard backlight changes via polling (universal fallback).
// hoopoo:hidden=true

package main

import (
	"fmt"
	"os"
	"time"

	"hoopoo/pkg/backlight/tracker/polling"
	"hoopoo/pkg/system/paths"
)

func main() {
	// Find brightness file (cached for performance)
	brightnessFile, err := polling.FindBrightnessFile()
	if err != nil || brightnessFile == "" {
		os.Exit(0) // No keyboard backlight, exit gracefully
	}
	
	trackerFile := paths.Runtime() + "/hoopoo/kbd-brightness-tracker"
	
	// Save initial brightness
	if err := polling.SaveInitialBrightness(brightnessFile, trackerFile); err != nil {
		// Non-fatal
	}
	
	// Read last saved value
	lastValue := 0
	if data, err := os.ReadFile(trackerFile); err == nil {
		fmt.Sscanf(string(data), "%d", &lastValue)
	}
	
	// Create single ticker (reuse instead of sleep loop)
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	
	// Pre-allocate buffer for reading (avoid repeated allocations)
	for range ticker.C {
		current, err := polling.ReadBrightness(brightnessFile)
		if err != nil {
			continue
		}
		
		// Skip if unchanged
		if current == lastValue {
			continue
		}
		
		// Save new value
		if err := polling.SaveBrightness(trackerFile, current); err == nil {
			lastValue = current
		}
	}
}
