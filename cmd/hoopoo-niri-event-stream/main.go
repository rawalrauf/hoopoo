// hoopoo:summary=Central niri event stream dispatcher with fan-out to multiple handlers

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"hoopoo/pkg/niri"
)

// EventHandler processes niri events.
// Stateless handlers return immediately after processing.
// Stateful handlers can consume additional events from the stream channel.
type EventHandler interface {
	Handle(event map[string]interface{}, stream <-chan map[string]interface{})
}

// EventDispatcher fans out events to multiple independent handlers
type EventDispatcher struct {
	handlers []EventHandler
	channels []chan map[string]interface{}
	wg       sync.WaitGroup
}

func NewDispatcher() *EventDispatcher {
	return &EventDispatcher{}
}

func (d *EventDispatcher) Register(handler EventHandler) {
	d.handlers = append(d.handlers, handler)
	ch := make(chan map[string]interface{}, 100)
	d.channels = append(d.channels, ch)
}

func (d *EventDispatcher) Start() error {
	scanner, cmd, err := niri.EventStream()
	if err != nil {
		return fmt.Errorf("failed to start event stream: %w", err)
	}
	defer cmd.Process.Kill()

	// Start handler goroutines
	for i, handler := range d.handlers {
		d.wg.Add(1)
		go d.runHandler(handler, d.channels[i])
	}

	// Read events and broadcast to all handlers
	for scanner.Scan() {
		line := scanner.Text()

		var event map[string]interface{}
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			continue
		}

		// Broadcast to all handler channels (non-blocking)
		for _, ch := range d.channels {
			select {
			case ch <- event:
			default:
				// Channel full, skip event for slow handler
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("event stream error: %w", err)
	}

	for _, ch := range d.channels {
		close(ch)
	}

	d.wg.Wait()
	return nil
}

func (d *EventDispatcher) runHandler(handler EventHandler, stream <-chan map[string]interface{}) {
	defer d.wg.Done()

	for event := range stream {
		handler.Handle(event, stream)
	}
}

func main() {
	dispatcher := NewDispatcher()

	// Register handlers
	dispatcher.Register(&MonitorRecoveryHandler{})
	dispatcher.Register(NewWaybarTitleHandler())

	if err := dispatcher.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Dispatcher error: %v\n", err)
		os.Exit(1)
	}
}
