// Path resolution for hoopoo directories

package paths

import (
  "os"
  "path/filepath"
)

func Hoopoo() string {
  if path := os.Getenv("HOOPOO_PATH"); path != "" {
    return path
  }
  return filepath.Join(os.Getenv("HOME"), ".local/share/hoopoo")
}

func State() string {
  stateHome := os.Getenv("XDG_STATE_HOME")
  if stateHome == "" {
    stateHome = filepath.Join(os.Getenv("HOME"), ".local/state")
  }
  return filepath.Join(stateHome, "hoopoo")
}

func Runtime() string {
  if runtime := os.Getenv("XDG_RUNTIME_DIR"); runtime != "" {
    return runtime
  }
  return "/tmp"
}

func Config() string {
  configHome := os.Getenv("XDG_CONFIG_HOME")
  if configHome == "" {
    configHome = filepath.Join(os.Getenv("HOME"), ".config")
  }
  return filepath.Join(configHome, "hoopoo")
}

func Bin() string {
  return filepath.Join(Hoopoo(), "bin")
}
