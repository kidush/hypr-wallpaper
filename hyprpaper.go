package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type Monitor struct {
	Name string `json:"name"`
}

func SetWallpaper(wallpaperPath string) (bool, string) {
	absPath, err := filepath.Abs(wallpaperPath)
	if err != nil {
		return false, fmt.Sprintf("Error: %v", err)
	}

	monitors, err := getMonitors()
	if err != nil {
		return false, fmt.Sprintf("Error getting monitors: %v", err)
	}

	if len(monitors) == 0 {
		return false, "No monitors found"
	}

	preloadCmd := exec.Command("hyprctl", "hyprpaper", "preload", absPath)
	if err := preloadCmd.Run(); err != nil {
		return false, fmt.Sprintf("Error preloading: %v", err)
	}

	for _, monitor := range monitors {
		wallpaperArg := fmt.Sprintf("%s,%s", monitor.Name, absPath)
		wallpaperCmd := exec.Command("hyprctl", "hyprpaper", "wallpaper", wallpaperArg)
		if err := wallpaperCmd.Run(); err != nil {
			return false, fmt.Sprintf("Error setting wallpaper on %s: %v", monitor.Name, err)
		}
	}

	if err := updateHyprpaperConfig(absPath, monitors); err != nil {
		return true, fmt.Sprintf("Set! (config save failed: %v)", err)
	}

	return true, "Wallpaper set and saved!"
}

func getMonitors() ([]Monitor, error) {
	cmd := exec.Command("hyprctl", "monitors", "-j")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var monitors []Monitor
	if err := json.Unmarshal(output, &monitors); err != nil {
		return nil, err
	}

	return monitors, nil
}

func updateHyprpaperConfig(wallpaperPath string, monitors []Monitor) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(homeDir, ".config", "hypr", "hyprpaper.conf")

	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	var config string
	config += fmt.Sprintf("preload = %s\n", wallpaperPath)

	for _, monitor := range monitors {
		config += fmt.Sprintf("wallpaper = %s,%s\n", monitor.Name, wallpaperPath)
	}

	config += "splash = false\n"
	config += "ipc = on\n"

	return os.WriteFile(configPath, []byte(config), 0644)
}
