# VidForge — Universal Video Downloader

![Go Version](https://img.shields.io/badge/go-1.25.5-blue.svg)
![GitHub Release](https://img.shields.io/github/v/release/Abhi1264/vidforge)
![License](https://img.shields.io/badge/license-MIT-green.svg)
![Homebrew](https://img.shields.io/homebrew/v/vidforge)

VidForge is a beautiful terminal-based video downloader that supports 1000+ platforms including YouTube, Vimeo, Twitter, TikTok, and many more. Built with Go and powered by `yt-dlp`.

![VidForge Demo](https://via.placeholder.com/800x400?text=VidForge+Demo)

## ✨ Features

- 🎬 **Universal Support**: Download from YouTube, Vimeo, Twitter, TikTok, and 1000+ platforms
- 🎨 **Beautiful Interface**: Modern, intuitive terminal UI
- 📊 **Concurrent Downloads**: Download multiple videos simultaneously
- 🎚️ **Quality Profiles**: Pre-configured profiles for different use cases
- 🚫 **SponsorBlock**: Automatically skip sponsor segments on YouTube
- 🔄 **Smart Resume**: Automatically resume interrupted downloads
- 📁 **Custom Locations**: Choose where to save your downloads
- 🖥️ **Cross-Platform**: Works on macOS, Linux, and Windows

## 📦 Installation

### Quick Install (Recommended)

**macOS / Linux (Homebrew):**
```bash
brew tap Abhi1264/homebrew-tap
brew install vidforge
```

**Windows (Scoop):**
```powershell
scoop bucket add vidforge https://github.com/Abhi1264/vidforge-scoop
scoop install vidforge
```

### Download from GitHub Releases

1. Download the latest release for your platform from the [Releases](https://github.com/Abhi1264/vidforge/releases) page
2. Extract the archive
3. Run vidforge

**Note:** VidForge will automatically offer to download `yt-dlp` on first run if it's not installed. You'll need to install `ffmpeg` separately:

**macOS:**
```bash
brew install ffmpeg
```

**Linux (Debian/Ubuntu):**
```bash
sudo apt install ffmpeg
```

**Linux (Fedora):**
```bash
sudo dnf install ffmpeg
```

**Linux (Arch):**
```bash
sudo pacman -S ffmpeg
```

**Windows:**
Download from [ffmpeg.org](https://ffmpeg.org/download.html) and add to PATH

## 🚀 Quick Start

1. Launch VidForge:
   ```bash
   vidforge
   ```

2. Paste a video URL and press Enter

3. Watch your download progress in real-time!

## ⌨️ Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `Enter` | Start download |
| `↑` / `↓` or `j` / `k` | Navigate downloads |
| `p` | Pause/cancel selected download |
| `s` | Toggle SponsorBlock (YouTube) |
| `f` | Change quality profile |
| `d` | Set download location |
| `?` | Show help |
| `q` / `Ctrl+C` | Quit |

## 🎚️ Quality Profiles

VidForge includes 6 pre-configured quality profiles:

1. **Best Quality** - 4K/8K, VP9/AV1, MKV format
2. **High Quality** - 1080p, H.264, MP4 format
3. **Balanced** - 720p, H.264, MP4 format (default)
4. **Mobile Saver** - ≤720p, low bitrate, MP4 format
5. **Audio Only** - Extract audio as MP3
6. **Archive** - Best quality + metadata + subtitles

Press `f` to switch between profiles, or use number keys 1-6 for quick selection.

## 🚫 SponsorBlock Integration

VidForge automatically removes sponsor segments from YouTube videos:
- Sponsor segments
- Self-promotion
- Interaction requests (like/subscribe)
- Intro and outro segments

Toggle with the `s` key. Automatically enabled for YouTube URLs.

## 📁 Custom Download Location

Press `d` to set a custom download location. Supports:
- Absolute paths: `/path/to/downloads`
- Relative paths: `./videos`
- Home directory: `~/Downloads`

The location is saved for future downloads.

## 🔧 Troubleshooting

### "yt-dlp not found"

VidForge will offer to download yt-dlp automatically on first run. If this fails, install manually:

**macOS:**
```bash
brew install yt-dlp
```

**Linux:**
```bash
# Ubuntu/Debian
sudo apt install yt-dlp

# Or use pip
pip install yt-dlp
```

**Windows:**
Download from [yt-dlp releases](https://github.com/yt-dlp/yt-dlp/releases) and add to PATH

### "ffmpeg not found"

ffmpeg must be installed manually. See the [Installation](#-installation) section above.

### Download Failed

- Check your internet connection
- Verify the URL is valid
- Some platforms may have region restrictions
- Try updating yt-dlp: `pip install -U yt-dlp` or `brew upgrade yt-dlp`

### Slow Downloads

- Reduce concurrent downloads (configurable in code, default is 3)
- Check your internet speed
- Some platforms may throttle download speeds

## 🤝 Contributing

Contributions are welcome! See [DEVELOPMENT.md](DEVELOPMENT.md) for technical details and development guidelines.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built with [Bubbletea](https://github.com/charmbracelet/bubbletea)
- Powered by [yt-dlp](https://github.com/yt-dlp/yt-dlp)
- SponsorBlock integration via [SponsorBlock API](https://sponsor.ajay.app/)

## 🔗 Links

- [Report a Bug](https://github.com/Abhi1264/vidforge/issues)
- [Request a Feature](https://github.com/Abhi1264/vidforge/issues)
- [Development Guide](DEVELOPMENT.md)
- [yt-dlp Documentation](https://github.com/yt-dlp/yt-dlp#readme)

---

Made with ❤️ by [Abhi1264](https://github.com/Abhi1264)
