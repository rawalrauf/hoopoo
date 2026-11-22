// XDG Desktop Portal integration via D-Bus

package portal

import (
  "fmt"
  "os"
  "strings"
  "time"
  
  "github.com/godbus/dbus/v5"
)

const (
  portalBus       = "org.freedesktop.portal.Desktop"
  portalPath      = "/org/freedesktop/portal/desktop"
  answerTimeout   = 600 * time.Second
)

// Connection manages D-Bus connection to desktop portal
type Connection struct {
  conn *dbus.Conn
}

// NewConnection creates a new portal connection
func NewConnection() (*Connection, error) {
  conn, err := dbus.SessionBus()
  if err != nil {
    return nil, fmt.Errorf("failed to connect to session bus: %w", err)
  }
  
  return &Connection{conn: conn}, nil
}

// Close closes the portal connection
func (c *Connection) Close() error {
  return c.conn.Close()
}

// generateToken creates a unique token for this process
func (c *Connection) generateToken(prefix string) string {
  return fmt.Sprintf("%s%d", prefix, os.Getpid())
}

// subscribePath subscribes to portal response signals on a given path
func (c *Connection) subscribePath(path string) error {
  matchRule := fmt.Sprintf(
    "type='signal',interface='org.freedesktop.portal.Request',member='Response',path='%s'",
    path,
  )
  return c.conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0, matchRule).Err
}

// unsubscribePath unsubscribes from portal response signals
func (c *Connection) unsubscribePath(path string) {
  matchRule := fmt.Sprintf(
    "type='signal',interface='org.freedesktop.portal.Request',member='Response',path='%s'",
    path,
  )
  c.conn.BusObject().Call("org.freedesktop.DBus.RemoveMatch", 0, matchRule)
}

// predictPath generates the expected response path based on token
func (c *Connection) predictPath(token string) string {
  sender := strings.ReplaceAll(c.conn.Names()[0][1:], ".", "_")
  return fmt.Sprintf("/org/freedesktop/portal/desktop/request/%s/%s", sender, token)
}
