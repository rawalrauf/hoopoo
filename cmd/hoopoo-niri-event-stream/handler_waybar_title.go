// Stateful handler: hides waybar title and consumes events until focus returns

package main

import (
	"os"
	"os/exec"
	"strings"
)

type WaybarTitleHandler struct {
	styleFile string
	hideRule  string
}

func NewWaybarTitleHandler() *WaybarTitleHandler {
	return &WaybarTitleHandler{
		styleFile: os.Getenv("HOME") + "/.config/waybar/style.css",
		hideRule:  "#window label { font-size: 0; }\n",
	}
}

func (h *WaybarTitleHandler) hideTitle() error {
	content, err := os.ReadFile(h.styleFile)
	if err != nil {
		return err
	}
	
	if strings.Contains(string(content), h.hideRule) {
		return nil
	}
	
	f, err := os.OpenFile(h.styleFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	
	_, err = f.WriteString(h.hideRule)
	return err
}

func (h *WaybarTitleHandler) showTitle() {
	content, err := os.ReadFile(h.styleFile)
	if err != nil {
		return
	}
	
	newContent := strings.ReplaceAll(string(content), h.hideRule, "")
	os.WriteFile(h.styleFile, []byte(newContent), 0644)
}

func (h *WaybarTitleHandler) isWalkerOpen() bool {
	cmd := exec.Command("niri", "msg", "--json", "layers")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), `"namespace":"walker"`)
}

func (h *WaybarTitleHandler) isOverviewOpen() bool {
	cmd := exec.Command("niri", "msg", "--json", "overview-state")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(output)) == `{"is_open":true}`
}

func (h *WaybarTitleHandler) Handle(event map[string]interface{}, stream <-chan map[string]interface{}) {
	focusChanged, ok := event["WindowFocusChanged"]
	if !ok {
		return
	}
	
	focusMap, ok := focusChanged.(map[string]interface{})
	if !ok {
		return
	}
	
	if focusMap["id"] == nil && (h.isWalkerOpen() || h.isOverviewOpen()) {
		h.hideTitle()
		
		for nextEvent := range stream {
			if h.shouldShowTitle(nextEvent) {
				h.showTitle()
				return
			}
		}
	}
}

func (h *WaybarTitleHandler) shouldShowTitle(event map[string]interface{}) bool {
	if focusChanged, ok := event["WindowFocusChanged"]; ok {
		if focusMap, ok := focusChanged.(map[string]interface{}); ok {
			if focusMap["id"] != nil {
				return true
			}
		}
	}
	
	if windowChanged, ok := event["WindowOpenedOrChanged"]; ok {
		if windowMap, ok := windowChanged.(map[string]interface{}); ok {
			if window, ok := windowMap["window"].(map[string]interface{}); ok {
				if focused, ok := window["is_focused"].(bool); ok && focused {
					return true
				}
			}
		}
	}
	
	if overviewToggled, ok := event["OverviewToggled"]; ok {
		if overviewMap, ok := overviewToggled.(map[string]interface{}); ok {
			if isOpen, ok := overviewMap["is_open"].(bool); ok && !isOpen {
				return true
			}
		}
	}
	
	return false
}
