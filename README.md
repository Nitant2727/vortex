# Vortex - YouTube Search CLI

A beautiful terminal-based YouTube search application built with Go. Search YouTube videos directly from your terminal with a modern, interactive interface.

![Vortex Demo](demo.gif)

## Features

- 🔍 Real-time YouTube search
- 🎨 Beautiful terminal UI with borders and colors
- ⌨️ Interactive keyboard controls
- 📋 Tabulated search results
- 🚀 Fast and lightweight
- 🔗 Direct video opening in browser

## Prerequisites

- Go 1.16 or higher
- A YouTube Data API key ([Get one here](https://console.developers.google.com/))

## Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/vortex.git
cd vortex
```

2. Install dependencies:
```bash
go mod download
```

3. Set up your YouTube API key:
```bash
export YOUTUBE_API_KEY="your-api-key-here"
```

## Usage

Run the application:
```bash
go run main.go
```

### Controls
- Type your search query and press `Enter` to search
- Use `↑/↓` or `k/j` keys to navigate through results
- Press `o` to open the selected video in your browser
- Press `?` to toggle help menu
- Press `Esc` or `Ctrl+C` to quit

## Project Structure

```
.
├── main.go              # Main application entry point
├── pkg/
│   ├── models/          # Data structures and types
│   ├── ui/             # UI styles and components
│   ├── utils/          # Helper functions
│   └── youtube/        # YouTube API integration
├── go.mod              # Go module file
└── README.md           # This file
```

## Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - Terminal UI framework
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Style definitions
- [go-pretty](https://github.com/jedib0t/go-pretty) - Table formatting
- [Google API Go Client](https://github.com/googleapis/google-api-go-client) - YouTube Data API client

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

