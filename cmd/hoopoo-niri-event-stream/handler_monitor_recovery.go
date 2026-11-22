// Stateless handler: processes WorkspacesChanged events and returns immediately

package main

import (
	"os/exec"
	"time"

	"hoopoo/pkg/system/process"
)

type MonitorRecoveryHandler struct {
	lastTrigger int64
}

const debounceSeconds = 3

func (h *MonitorRecoveryHandler) shouldTrigger() bool {
	now := time.Now().Unix()
	
	if now-h.lastTrigger < debounceSeconds {
		return false
	}
	
	h.lastTrigger = now
	return true
}

func (h *MonitorRecoveryHandler) hasExternalMonitors() bool {
	return exec.Command("hoopoo", "hw", "external", "monitors").Run() == nil
}

func (h *MonitorRecoveryHandler) Handle(event map[string]interface{}, stream <-chan map[string]interface{}) {
	if _, ok := event["WorkspacesChanged"]; !ok {
		return
	}
	
	if !h.shouldTrigger() {
		return
	}
	
	exec.Command("hoopoo", "niri", "monitor", "internal", "recover").Run()
	
	if process.IsRunning("wl-mirror") && !h.hasExternalMonitors() {
		exec.Command("hoopoo", "niri", "monitor", "internal", "mirror", "off").Run()
	}
}
