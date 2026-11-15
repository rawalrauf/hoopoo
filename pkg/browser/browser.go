// Browser detection and configuration utilities

package browser

import (
  "bufio"
  "os"
  "os/exec"
  "path/filepath"
  "regexp"
  "strings"
)

// GetDefault returns the default browser desktop file name
func GetDefault() (string, error) {
  cmd := exec.Command("xdg-settings", "get", "default-web-browser")
  output, err := cmd.Output()
  if err != nil {
    return "", err
  }
  return strings.TrimSpace(string(output)), nil
}

// GetExecutable finds the executable path from a .desktop file
func GetExecutable(desktopFile string) (string, error) {
  searchPaths := []string{
    filepath.Join(os.Getenv("HOME"), ".local/share/applications"),
    filepath.Join(os.Getenv("HOME"), ".nix-profile/share/applications"),
    "/usr/share/applications",
  }

  for _, path := range searchPaths {
    fullPath := filepath.Join(path, desktopFile)
    file, err := os.Open(fullPath)
    if err != nil {
      continue
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    re := regexp.MustCompile(`^Exec=([^ ]+)`)
    for scanner.Scan() {
      matches := re.FindStringSubmatch(scanner.Text())
      if len(matches) > 1 {
        return matches[1], nil
      }
    }
  }

  return "", os.ErrNotExist
}

// IsChromiumBased checks if browser supports Chromium-based features (web apps)
func IsChromiumBased(desktopFile string) bool {
  prefixes := []string{"google-chrome", "brave-browser", "microsoft-edge", 
                       "opera", "vivaldi", "helium", "chromium"}
  
  for _, prefix := range prefixes {
    if strings.HasPrefix(desktopFile, prefix) {
      return true
    }
  }
  return false
}

// FindChromiumBased finds the first available Chromium-based browser
func FindChromiumBased() (string, error) {
  candidates := []string{"chromium", "chromium-bin", "brave-browser", "google-chrome",
                        "microsoft-edge", "vivaldi", "opera", "thorium-browser",
                        "ungoogled-chromium", "helium"}

  searchPaths := []string{
    filepath.Join(os.Getenv("HOME"), ".local/share/applications"),
    filepath.Join(os.Getenv("HOME"), ".nix-profile/share/applications"),
    "/usr/share/applications",
  }

  for _, candidate := range candidates {
    for _, path := range searchPaths {
      pattern := filepath.Join(path, candidate+"*.desktop")
      matches, _ := filepath.Glob(pattern)
      if len(matches) > 0 {
        return filepath.Base(matches[0]), nil
      }
    }
  }

  return "", os.ErrNotExist
}

// DetectPrivateFlag returns the private browsing flag for the browser
func DetectPrivateFlag(execPath string) string {
  // Check if it's Firefox (has MOZ_LOG in help)
  cmd := exec.Command(execPath, "--help")
  output, _ := cmd.CombinedOutput()
  
  if strings.Contains(string(output), "MOZ_LOG") {
    return "--private-window"
  } else if strings.Contains(execPath, "edge") {
    return "--inprivate"
  }
  
  return "--incognito" // Default for Chromium-based
}

// ExtractDomain extracts clean domain from URL
// e.g., "https://www.youtube.com/watch?v=123" -> "youtube.com"
func ExtractDomain(url string) string {
  cleanDomain := url
  
  // Remove protocol
  re := regexp.MustCompile(`^[a-zA-Z0-9+.-]+://`)
  cleanDomain = re.ReplaceAllString(cleanDomain, "")
  
  // Remove www.
  cleanDomain = strings.TrimPrefix(cleanDomain, "www.")
  
  // Remove path
  if idx := strings.Index(cleanDomain, "/"); idx != -1 {
    cleanDomain = cleanDomain[:idx]
  }
  
  return cleanDomain
}

// URLToAppIDBase converts URL to app ID base format
// e.g., "https://youtube.com/watch" -> "youtube.com__watch"
func URLToAppIDBase(url string) string {
  // Remove protocol
  stripped := url
  re := regexp.MustCompile(`^[a-zA-Z0-9+.-]+://`)
  stripped = re.ReplaceAllString(stripped, "")
  
  // Remove trailing slash
  stripped = strings.TrimSuffix(stripped, "/")
  
  if strings.Contains(stripped, "/") {
    parts := strings.SplitN(stripped, "/", 2)
    domain := parts[0]
    path := strings.TrimSuffix(parts[1], "/")
    path = strings.ReplaceAll(path, "/", "_")
    return domain + "__" + path
  }
  
  return stripped + "__"
}

// GetBrowserPrefix returns the browser prefix for app ID construction
func GetBrowserPrefix(desktopFile string) string {
  switch {
  case strings.HasPrefix(desktopFile, "google-chrome"), strings.HasPrefix(desktopFile, "chrome-"):
    return "chrome"
  case strings.HasPrefix(desktopFile, "brave"):
    return "brave"
  case strings.HasPrefix(desktopFile, "chromium"):
    return "chromium"
  case strings.HasPrefix(desktopFile, "microsoft-edge"), strings.HasPrefix(desktopFile, "msedge"):
    return "msedge"
  case strings.HasPrefix(desktopFile, "opera"):
    return "opera"
  case strings.HasPrefix(desktopFile, "vivaldi"):
    return "vivaldi"
  case strings.HasPrefix(desktopFile, "helium"):
    return "helium"
  case strings.HasPrefix(desktopFile, "thorium"):
    return "thorium"
  default:
    return "chrome" // fallback
  }
}

// BuildWebAppID constructs the app ID for a web app
// Format: prefix-domain__path-Default
func BuildWebAppID(prefix, domain string) string {
  return prefix + "-" + domain + "__-Default"
}

