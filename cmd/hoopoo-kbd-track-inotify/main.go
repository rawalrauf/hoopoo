// hoopoo:summary=Track keyboard backlight changes via inotify (event-driven).
// hoopoo:hidden=true

package main

import (
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
	"hoopoo/pkg/backlight/keyboard"
	"hoopoo/pkg/backlight/tracker/inotify"
	"hoopoo/pkg/system/paths"
)

func main() {
	// Find keyboard device
	device, err := keyboard.FindKeyboardDevice()
	if err != nil {
		os.Exit(1)
	}
	
	// Find hw_changed file
	hwChangedFile, err := inotify.FindHwChangedFile()
	if err != nil {
		os.Exit(1)
	}
	
	trackerFile := paths.Runtime() + "/hoopoo/kbd-brightness-tracker"
	
	// Save initial brightness
	if err := inotify.SaveInitialBrightness(device, trackerFile); err != nil {
		// Non-fatal
	}
	
	// Read last saved value
	lastValue := 0
	if data, err := os.ReadFile(trackerFile); err == nil {
		fmt.Sscanf(string(data), "%d", &lastValue)
	}
	
	// Create fsnotify watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		os.Exit(1)
	}
	defer watcher.Close()
	
	// Watch hw_changed file
	if err := watcher.Add(hwChangedFile); err != nil {
		os.Exit(1)
	}
	
	// Event loop
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				// Read current brightness
				current, err := inotify.ReadBrightness(device)
				if err != nil {
					continue
				}
				
				// Skip if unchanged
				if current == lastValue {
					continue
				}
				
				// Save new value
				if err := inotify.SaveBrightness(trackerFile, current); err == nil {
					lastValue = current
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			// Log error but continue
			_ = err
		}
	}
}
