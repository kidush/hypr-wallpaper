# hypr-wallpaper

A TUI file browser for selecting and setting wallpapers on Hyprland using hyprpaper.

## Features

- Browse folders and images with vim-style navigation (hjkl)
- Live image preview using überzug++
- Set wallpaper with Enter key
- Automatically updates `~/.config/hypr/hyprpaper.conf` for persistence

## Dependencies

### Runtime

- [hyprpaper](https://github.com/hyprwm/hyprpaper) - Hyprland wallpaper utility
- [ueberzugpp](https://github.com/jstkdng/ueberzugpp) - Terminal image previews

### Build

- Go 1.21+

## Installation

### From source

```bash
go install github.com/thiagoflins/hypr-wallpaper@latest
```

### Manual build

```bash
git clone https://github.com/thiagoflins/hypr-wallpaper.git
cd hypr-wallpaper
go build -o hypr-wallpaper .
sudo cp hypr-wallpaper /usr/local/bin/
```

## Usage

```bash
# Start from home directory
hypr-wallpaper

# Start from specific directory
hypr-wallpaper ~/Pictures/wallpapers
```

## Keybindings

| Key | Action |
|-----|--------|
| `j` / `↓` | Move down |
| `k` / `↑` | Move up |
| `l` / `→` / `Enter` | Enter folder / Set wallpaper |
| `h` / `←` / `Backspace` | Go to parent folder |
| `q` / `Esc` | Quit |

## License

MIT
