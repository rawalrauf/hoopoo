// Niri event stream

package niri

import (
  "bufio"
  "os/exec"
)

// EventStream returns a scanner for the niri event stream
func EventStream() (*bufio.Scanner, *exec.Cmd, error) {
  cmd := exec.Command("niri", "msg", "--json", "event-stream")
  stdout, err := cmd.StdoutPipe()
  if err != nil {
    return nil, nil, err
  }

  if err := cmd.Start(); err != nil {
    return nil, nil, err
  }

  scanner := bufio.NewScanner(stdout)
  return scanner, cmd, nil
}
