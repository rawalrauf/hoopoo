// Path resolution for hoopoo directories

package paths

import (
  "os"
  "path/filepath"
)

func UserData() string {
  dataHome := os.Getenv("XDG_DATA_HOME")
  if dataHome == "" {
    home, _ := os.UserHomeDir()
    dataHome = filepath.Join(home, ".local/share")
  }
  return dataHome
}

func Hoopoo() string {
  if path := os.Getenv("HOOPOO_PATH"); path != "" {
    return path
  }
  return filepath.Join(UserData(), "hoopoo")
}

func State() string {
  stateHome := os.Getenv("XDG_STATE_HOME")
  if stateHome == "" {
    home, _ := os.UserHomeDir()
    stateHome = filepath.Join(home, ".local/state")
  }
  return filepath.Join(stateHome, "hoopoo")
}

func Runtime() string {
  if runtime := os.Getenv("XDG_RUNTIME_DIR"); runtime != "" {
    return runtime
  }
  return "/tmp"
}

func UserConfig() string {
  configHome := os.Getenv("XDG_CONFIG_HOME")
  if configHome == "" {
    home, _ := os.UserHomeDir()
    configHome = filepath.Join(home, ".config")
  }
  return configHome
}

func HoopooConfig() string {
  return filepath.Join(UserConfig(), "hoopoo")
}

func Bin() string {
  return filepath.Join(Hoopoo(), "bin")
}

func UserApplications() string {
  return filepath.Join(UserData(), "applications")
}
