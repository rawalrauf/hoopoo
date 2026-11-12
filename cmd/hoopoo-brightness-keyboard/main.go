// hoopoo:summary=Adjust keyboard backlight brightness using available steps.
// hoopoo:args=<up|down|cycle|off|restore>

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"hoopoo/pkg/backlight/keyboard"
)

func main() {
	if len(os.Args) < 2 {
		os.Args = append(os.Args, "up")
	}
	
	direction := os.Args[1]
	
	device, err := keyboard.FindKeyboardDevice()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	
	stateFile := keyboard.StateFile()
	trackerFile := keyboard.TrackerFile()
	
	// Handle off command
	if direction == "off" {
		lock, err := keyboard.LockFile(stateFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer lock.Close()
		
		// Read from tracker's runtime state
		if trackerValue, err := keyboard.ReadStateFile(trackerFile); err == nil {
			if err := keyboard.WriteStateFile(stateFile, trackerValue); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}
		
		if err := setBrightness(device, 0); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}
	
	// Handle restore command
	if direction == "restore" {
		lock, err := keyboard.LockFile(stateFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer lock.Close()
		
		// Priority: state file, then tracker file
		var saved int
		var found bool
		
		if val, err := keyboard.ReadStateFile(stateFile); err == nil {
			saved = val
			found = true
		} else if val, err := keyboard.ReadStateFile(trackerFile); err == nil {
			saved = val
			found = true
		}
		
		if found {
			if err := setBrightness(device, saved); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}
		return
	}
	
	// Handle up/down/cycle commands
	// Read current and max brightness in parallel
	type result struct {
		value int
		err   error
	}
	
	currentCh := make(chan result, 1)
	maxCh := make(chan result, 1)
	
	go func() {
		val, err := keyboard.GetBrightness(device)
		currentCh <- result{val, err}
	}()
	
	go func() {
		val, err := keyboard.GetMaxBrightness(device)
		maxCh <- result{val, err}
	}()
	
	currentResult := <-currentCh
	maxResult := <-maxCh
	
	if currentResult.err != nil {
		fmt.Fprintln(os.Stderr, currentResult.err)
		os.Exit(1)
	}
	
	if maxResult.err != nil {
		fmt.Fprintln(os.Stderr, maxResult.err)
		os.Exit(1)
	}
	
	current := currentResult.value
	max := maxResult.value
	
	step := keyboard.CalculateStep(max)
	newBrightness := keyboard.CalculateNewBrightness(current, max, step, direction)
	
	if err := setBrightness(device, newBrightness); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	
	// Update tracker file only if inotify service is installed
	// (fast file check instead of systemctl call)
	if isInotifyInstalled() {
		keyboard.WriteStateFile(trackerFile, newBrightness)
	}
	
	// Display OSD in background (detached process)
	percent := (newBrightness * 100) / max
	cmd := exec.Command("hoopoo", "swayosd", "kbd-brightness", fmt.Sprintf("%d", percent))
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	cmd.Start()
}

func setBrightness(device string, level int) error {
	cmd := exec.Command("brightnessctl", "-d", device, "set", fmt.Sprintf("%d", level))
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}

func isInotifyInstalled() bool {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false
	}
	_, err = os.Stat(filepath.Join(homeDir, ".config/systemd/user/kbd-track-inotify.service"))
	return err == nil
}
