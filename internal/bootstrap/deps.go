package bootstrap

import (
	"fmt"
	"os/exec"
)

// Ensure checks if a command is available in the system PATH
func Ensure(cmd string) error {
	if _, err := exec.LookPath(cmd); err == nil {
		return nil
	}
	
	switch cmd {
	case "yt-dlp":
		return fmt.Errorf("yt-dlp not found. Please install it:\n" +
			"  macOS: brew install yt-dlp\n" +
			"  Linux: Use your package manager (apt, dnf, pacman) or pip install yt-dlp\n" +
			"  Windows: Download from https://github.com/yt-dlp/yt-dlp/releases")
	case "ffmpeg":
		return fmt.Errorf("ffmpeg not found. Please install it:\n" +
			"  macOS: brew install ffmpeg\n" +
			"  Linux: Use your package manager (apt, dnf, pacman)\n" +
			"  Windows: Download from https://ffmpeg.org/download.html")
	default:
		return fmt.Errorf("%s not found in PATH", cmd)
	}
}

// GetCommandPath returns the full path to a command
func GetCommandPath(cmd string) (string, error) {
	path, err := exec.LookPath(cmd)
	if err != nil {
		return "", fmt.Errorf("%s not found in PATH", cmd)
	}
	return path, nil
}
