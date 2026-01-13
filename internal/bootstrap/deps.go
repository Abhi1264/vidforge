package bootstrap

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// getBinaryDir returns the directory where downloaded binaries are stored
func getBinaryDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".vidforge", "bin"), nil
}

// getLocalBinaryPath returns the path to a locally downloaded binary
func getLocalBinaryPath(cmd string) (string, error) {
	binDir, err := getBinaryDir()
	if err != nil {
		return "", err
	}

	binaryName := cmd
	if runtime.GOOS == "windows" {
		binaryName = cmd + ".exe"
	}

	return filepath.Join(binDir, binaryName), nil
}

// downloadBinary downloads a binary for the current platform
func downloadBinary(cmd string) error {
	binDir, err := getBinaryDir()
	if err != nil {
		return fmt.Errorf("failed to get binary directory: %w", err)
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return fmt.Errorf("failed to create binary directory: %w", err)
	}

	var url string
	var destPath string

	switch cmd {
	case "yt-dlp":
		destPath, _ = getLocalBinaryPath("yt-dlp")
		switch runtime.GOOS {
		case "darwin":
			url = "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_macos"
		case "linux":
			if runtime.GOARCH == "arm64" {
				url = "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_linux_aarch64"
			} else {
				url = "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_linux"
			}
		case "windows":
			url = "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp.exe"
		default:
			return fmt.Errorf("unsupported platform: %s/%s", runtime.GOOS, runtime.GOARCH)
		}

	case "ffmpeg":
		destPath, _ = getLocalBinaryPath("ffmpeg")
		// For ffmpeg, we need to download from different sources
		return fmt.Errorf("automatic ffmpeg download not yet implemented. Please install manually:\n" +
			"  macOS: brew install ffmpeg\n" +
			"  Linux: sudo apt install ffmpeg (or your package manager)\n" +
			"  Windows: Download from https://ffmpeg.org/download.html")

	default:
		return fmt.Errorf("unknown command: %s", cmd)
	}

	fmt.Printf("Downloading %s...\n", cmd)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status: %s", resp.Status)
	}

	out, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	// Make executable on Unix systems
	if runtime.GOOS != "windows" {
		if err := os.Chmod(destPath, 0755); err != nil {
			return fmt.Errorf("failed to make executable: %w", err)
		}
	}

	fmt.Printf("✓ Downloaded %s to %s\n", cmd, destPath)
	return nil
}

// Ensure checks if a command is available and offers to download if not
func Ensure(cmd string) error {
	// 1. Check if command is in PATH
	if _, err := exec.LookPath(cmd); err == nil {
		return nil
	}

	// 2. Check if we have it locally downloaded
	localPath, err := getLocalBinaryPath(cmd)
	if err == nil {
		if _, err := os.Stat(localPath); err == nil {
			return nil
		}
	}

	// 3. Command not found - offer to download
	fmt.Printf("\n%s not found in system PATH.\n", cmd)

	// Check if we're in a non-interactive environment
	if os.Getenv("CI") == "true" || !isTerminal() {
		return fmt.Errorf("%s not found. Please install it:\n"+
			"  macOS: brew install %s\n"+
			"  Linux: Use your package manager (apt, dnf, pacman)\n"+
			"  Windows: Download from official website", cmd, cmd)
	}

	// Ask user if they want to download
	if cmd == "ffmpeg" {
		// ffmpeg is more complex, just show install instructions
		return fmt.Errorf("ffmpeg not found. Please install it:\n" +
			"  macOS: brew install ffmpeg\n" +
			"  Linux: sudo apt install ffmpeg (or your package manager)\n" +
			"  Windows: Download from https://ffmpeg.org/download.html")
	}

	fmt.Printf("Would you like to download %s automatically? [Y/n]: ", cmd)
	var response string
	fmt.Scanln(&response)

	response = strings.ToLower(strings.TrimSpace(response))
	if response == "" || response == "y" || response == "yes" {
		if err := downloadBinary(cmd); err != nil {
			return fmt.Errorf("failed to download %s: %w\n"+
				"Please install manually:\n"+
				"  macOS: brew install %s\n"+
				"  Linux: Use your package manager\n"+
				"  Windows: Download from official website", cmd, err, cmd)
		}
		return nil
	}

	return fmt.Errorf("%s required but not installed", cmd)
}

// GetCommandPath returns the full path to a command
func GetCommandPath(cmd string) (string, error) {
	// Check PATH first
	if path, err := exec.LookPath(cmd); err == nil {
		return path, nil
	}

	// Check local directory
	localPath, err := getLocalBinaryPath(cmd)
	if err == nil {
		if _, err := os.Stat(localPath); err == nil {
			return localPath, nil
		}
	}

	return "", fmt.Errorf("%s not found", cmd)
}

// isTerminal checks if we're running in an interactive terminal
func isTerminal() bool {
	fileInfo, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}
