// IPC (Inter-Process Communication) for Hoopoo components
//
// This package provides communication between hoopoo processes, allowing
// the ytdlp-host to signal the overlay, notification system to trigger
// overlays, and other component coordination.

package ipc

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
)

// MessageType represents the type of IPC message
type MessageType string

const (
	MessageShowOverlay    MessageType = "show_overlay"
	MessageHideOverlay    MessageType = "hide_overlay"
	MessageUpdateProgress MessageType = "update_progress"
	MessageComplete       MessageType = "complete"
	MessageError          MessageType = "error"
)

// Message represents an IPC message
type Message struct {
	Type    MessageType            `json:"type"`
	Payload map[string]interface{} `json:"payload"`
}

// Server represents an IPC server (e.g., overlay daemon listening)
type Server struct {
	socketPath string
	listener   net.Listener
	handlers   map[MessageType]func(map[string]interface{}) error
}

// NewServer creates a new IPC server
func NewServer(name string) (*Server, error) {
	socketPath := getSocketPath(name)
	
	// Remove old socket if exists
	os.Remove(socketPath)
	
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create socket: %w", err)
	}
	
	s := &Server{
		socketPath: socketPath,
		listener:   listener,
		handlers:   make(map[MessageType]func(map[string]interface{}) error),
	}
	
	return s, nil
}

// Handle registers a handler for a message type
func (s *Server) Handle(msgType MessageType, handler func(map[string]interface{}) error) {
	s.handlers[msgType] = handler
}

// Listen starts listening for IPC messages
func (s *Server) Listen() error {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return err
		}
		
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	
	decoder := json.NewDecoder(conn)
	var msg Message
	
	if err := decoder.Decode(&msg); err != nil {
		return
	}
	
	if handler, ok := s.handlers[msg.Type]; ok {
		handler(msg.Payload)
	}
}

// Close closes the IPC server
func (s *Server) Close() error {
	os.Remove(s.socketPath)
	return s.listener.Close()
}

// Client represents an IPC client (e.g., ytdlp-host sending messages)
type Client struct {
	socketPath string
}

// NewClient creates a new IPC client
func NewClient(name string) *Client {
	return &Client{
		socketPath: getSocketPath(name),
	}
}

// Send sends a message to the IPC server
func (c *Client) Send(msgType MessageType, payload map[string]interface{}) error {
	conn, err := net.Dial("unix", c.socketPath)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer conn.Close()
	
	msg := Message{
		Type:    msgType,
		Payload: payload,
	}
	
	encoder := json.NewEncoder(conn)
	return encoder.Encode(msg)
}

// IsRunning checks if the IPC server is running
func (c *Client) IsRunning() bool {
	_, err := os.Stat(c.socketPath)
	return err == nil
}

func getSocketPath(name string) string {
	runtimeDir := os.Getenv("XDG_RUNTIME_DIR")
	if runtimeDir == "" {
		runtimeDir = "/tmp"
	}
	
	return filepath.Join(runtimeDir, fmt.Sprintf("hoopoo-%s.sock", name))
}
