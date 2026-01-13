# VidForge Development Guide

This document contains technical information for developers working on VidForge.

## Architecture

### Project Structure

```
vidforge/
├── cmd/
│   └── vidforge/
│       └── main.go          # Application entry point
├── internal/
│   ├── bootstrap/
│   │   ├── deps.go          # Dependency verification and download
│   │   └── os.go            # OS detection utilities
│   ├── config/
│   │   └── config.go        # Configuration management
│   ├── downloader/
│   │   ├── formats.go       # Format listing functionality
│   │   ├── job.go           # Download job execution
│   │   ├── manager.go       # Concurrent download manager
│   │   ├── profile.go       # Quality profiles
│   │   └── url.go           # URL validation utilities
│   └── ui/
│       ├── formats.go       # Format list item implementation
│       ├── help.go          # Help text
│       ├── model.go         # Bubbletea model (state)
│       ├── styles.go        # UI styling
│       ├── update.go        # Event handling (update logic)
│       └── view.go          # UI rendering
├── go.mod
├── go.sum
└── README.md
```

### Core Components

#### 1. Bootstrap Package (`internal/bootstrap`)

Handles dependency verification and optional download:

- **`deps.go`**: Dependency management
  - Checks system PATH for `yt-dlp` and `ffmpeg`
  - Offers to download yt-dlp if not found (interactive mode)
  - Stores downloaded binaries in `~/.vidforge/bin/`
  - Returns full path to executables for use by downloader

- **`os.go`**: OS and architecture detection utilities
  - `DetectOS()`: Returns current OS and architecture
  - `IsMacOS()`, `IsLinux()`, `IsWindows()`: Platform checks
  - `IsIntel()`, `IsARM()`: Architecture checks

#### 2. Downloader Package (`internal/downloader`)

Core download functionality:

- **`manager.go`**: Concurrent download manager
  - Manages a queue of download jobs
  - Supports multiple concurrent downloads (configurable workers)
  - Provides progress updates via channels
  - Supports job cancellation

- **`job.go`**: Individual download job execution
  - Runs `yt-dlp` processes
  - Parses progress output using regex
  - Handles context cancellation
  - Supports resume functionality

- **`profile.go`**: Download quality profiles
  - Pre-defined quality profiles with yt-dlp flags
  - Profile selection by quality string
  - Default profile selection

- **`formats.go`**: Format listing utility
  - Parses `yt-dlp -F` output
  - Returns available formats for a URL

- **`url.go`**: URL utilities
  - YouTube URL detection
  - Platform-specific handling

#### 3. UI Package (`internal/ui`)

Terminal user interface built with Bubbletea:

- **`model.go`**: Application state
  - Text input for URLs
  - Job tracking (map of job states)
  - Profile selection state
  - SponsorBlock toggle state
  - Help/profile menu visibility

- **`update.go`**: Event handling
  - Keyboard input processing
  - Progress updates from download manager
  - State transitions
  - Job submission logic

- **`view.go`**: UI rendering
  - Main view rendering
  - Progress bar rendering
  - Profile selection view
  - Help screen rendering
  - Job list rendering

- **`styles.go`**: UI styling
  - Lipgloss styles for different UI elements
  - Color schemes and formatting

- **`help.go`**: Help text content

#### 4. Config Package (`internal/config`)

Configuration management:

- Default download paths per platform
- Configuration file support (planned)
- User preference storage

### Data Flow

1. **Initialization**:
   - `main.go` verifies dependencies via `bootstrap.Ensure()`
   - Offers to download yt-dlp if not found (interactive mode)
   - Exits with error if ffmpeg not found
   - Creates a new Bubbletea program with `ui.NewModel()`
   - Initializes download manager with 3 worker goroutines

2. **Download Submission**:
   - User enters URL and presses Enter
   - Model validates URL and creates a new job
   - Job is submitted to the download manager queue
   - Manager assigns job to an available worker

3. **Download Execution**:
   - Worker goroutine uses `bootstrap.GetCommandPath()` to locate `yt-dlp`
   - Executes `yt-dlp` command from system PATH or local directory
   - Progress is parsed from stdout using regex
   - Progress updates are sent via channel to the UI model

4. **UI Updates**:
   - UI listens to progress channel
   - Model updates job states on progress messages
   - View re-renders with updated progress bars and status

### Concurrency Model

VidForge uses a worker pool pattern for concurrent downloads:

- **Download Manager**: Maintains a channel-based queue of jobs
- **Worker Goroutines**: 3 workers process jobs concurrently
- **Progress Channel**: Single channel for all progress updates
- **Context Cancellation**: Each job has a context for cancellation support

## Building

### Prerequisites

- Go 1.25.5 or later
- `yt-dlp` and `ffmpeg` (for testing)

### Build Commands

```bash
# Development build
go build -o vidforge ./cmd/vidforge

# Production build with version
go build -ldflags="-s -w -X main.version=v1.3.2" -o vidforge ./cmd/vidforge

# Cross-compile for different platforms
GOOS=linux GOARCH=amd64 go build -o vidforge-linux-amd64 ./cmd/vidforge
GOOS=darwin GOARCH=arm64 go build -o vidforge-darwin-arm64 ./cmd/vidforge
GOOS=windows GOARCH=amd64 go build -o vidforge-windows-amd64.exe ./cmd/vidforge
```

## Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with race detector
go test -race ./...

# Run tests with coverage
go test -cover ./...
```

### Manual Testing

```bash
# Test with CI environment (skips dependency checks)
CI=true go run ./cmd/vidforge

# Test version flag
go run ./cmd/vidforge --version
```

## Development Workflow

### Adding New Features

1. Create feature branch: `git checkout -b feature/my-feature`
2. Implement changes
3. Test locally
4. Update documentation
5. Submit pull request

### Adding New Quality Profiles

Edit `internal/downloader/profile.go` and add a new `Profile`:

```go
{
    Name:        "Custom Profile",
    Specs:       "Your specs description",
    Quality:     "custom",
    Description: "Full description",
    Flags: []string{
        "-f", "your-format-string",
        "--merge-output-format", "mp4",
    },
}
```

### Modifying Keyboard Shortcuts

Edit `internal/ui/update.go` in the `Update()` function to add or modify key bindings.

## Release Process

### Using GoReleaser

1. Update version in changelog
2. Commit all changes
3. Tag release: `git tag -a v1.3.2 -m "Release v1.3.2"`
4. Push tag: `git push origin v1.3.2`
5. Run GoReleaser: `goreleaser release --clean`

### Manual Release

1. Build for all platforms
2. Create archives
3. Generate checksums
4. Create GitHub release
5. Upload artifacts

## Configuration

VidForge uses sensible defaults:

- **Concurrent Downloads**: 3 workers (configurable in `ui/model.go`)
- **Default Profile**: Balanced (720p, H.264, MP4)
- **SponsorBlock**: Enabled by default for YouTube URLs
- **Resume**: Enabled by default for all downloads
- **Download Location**: User's Downloads folder or configurable path

## Dependencies

### Go Dependencies

- `github.com/charmbracelet/bubbletea`: TUI framework
- `github.com/charmbracelet/bubbles`: UI components (textinput)
- `github.com/charmbracelet/lipgloss`: Terminal styling

### External Dependencies

- `yt-dlp`: Video downloader backend
- `ffmpeg`: Media processing tool

## Known Limitations

1. **Format Selection**: Currently uses profiles only; manual format selection not available
2. **Error Handling**: Basic error display; detailed error recovery may need improvement
3. **Concurrent Workers**: Fixed at 3 workers (configurable in code)
4. **ffmpeg Download**: Not implemented; users must install manually

## Future Enhancements

- [ ] Configurable download directory via UI
- [ ] Custom format selection UI
- [ ] Download history persistence
- [ ] Playlist support with queue management
- [ ] Subtitle selection and management
- [ ] Config file for user preferences
- [ ] Better error messages and recovery
- [ ] Download speed indicators
- [ ] ETA calculation and display
- [ ] Auto-update checking
- [ ] Automatic ffmpeg download

## Contributing

1. Fork the repository
2. Create your feature branch
3. Make your changes
4. Add tests if applicable
5. Update documentation
6. Submit a pull request

## License

MIT License - see LICENSE file for details
