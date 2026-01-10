package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Abhi1264/vidforge/internal/bootstrap"
	"github.com/Abhi1264/vidforge/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

var version = "dev"

func main() {
	showVersion := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Println(version)
		return
	}

	if os.Getenv("CI") == "true" {
		fmt.Println("VidForge CI mode: startup OK")
		return
	}

	if err := bootstrap.Ensure("yt-dlp"); err != nil {
		log.Fatal(err)
	}
	if err := bootstrap.Ensure("ffmpeg"); err != nil {
		log.Fatal(err)
	}

	p := tea.NewProgram(
		ui.NewModel(),
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
