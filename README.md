# Launchpad

A GTK4-based popup launcher for Ubuntu Mate, written in Go. Displays a configurable list of application entries with
icons and commands, supports separators, and can be triggered by a global hotkey (default: <kbd>Super</kbd>+<kbd>
Esc</kbd>).

---

## Features

- **Icon + command entries** — Each entry shows a themed icon and name; supports both system theme icons and custom icon
  file paths.
- **Separators** — Organize entries into groups with visual separators.
- **Auto-close** — Window closes when you click outside it, or press <kbd>Esc</kbd>.
- **Shift-to-keep-open** — Hold <kbd>Shift</kbd> while clicking (or press <kbd>Shift</kbd>+<kbd>Enter</kbd>) to launch
  an entry without closing the launcher.
- **Keyboard navigation** — Navigate entries with arrow keys, activate with <kbd>Enter</kbd> or <kbd>Space</kbd>, close
  with <kbd>Esc</kbd>.
- **Single instance** — Only one launchpad window can be open at a time.
- **Global hotkey** — Register <kbd>Super</kbd>+<kbd>Esc</kbd> (or custom) via the `register` command using Mate's Marco
  window manager.
- **Skeleton config** — Generate a starter config with `init`.
- **Icon fallback** — Unresolvable icons gracefully fall back to `application-x-executable`.

---

## Installation

### Prerequisites

- Go 1.26+
- GTK4 development libraries (`libgtk-4-dev`)
- libadwaita development libraries (`libadwaita-1-dev`)
- `golangci-lint` (optional, for development)

On Debian/Ubuntu Mate:

```bash
sudo apt install golang-go libgtk-4-dev libadwaita-1-dev golangci-lint
```

### Build from source

```bash
git clone https://github.com/abunjevac/launchpad.git
cd launchpad
task build
```

Or without Taskfile:

```bash
go build -o launchpad ./cmd/launchpad/
```

### Install

```bash
sudo cp launchpad /usr/local/bin/
```

Or use `task install` / `go install`:

```bash
go install ./cmd/launchpad/
```

---

## Configuration

### Default location

`~/.launchpad.yaml`

### Skeleton config

Generate a starter config:

```bash
launchpad init
```

Or specify a custom path:

```bash
launchpad init -c /path/to/config.yaml
```

### Full config reference

```yaml
# Global shortcut key combination (used by `launchpad register`)
# Default: "<Super>Escape"
shortcut: "<Super>Escape"

entries:
  # Regular entry with themed icon
  - name: "Terminal"
    icon: "utilities-terminal"
    command: "mate-terminal"

  # Separator between groups
  - separator: true

  # Entry with custom icon file path
  - name: "Custom Script"
    icon: "/home/user/icons/myapp.png"
    command: "/home/user/.local/bin/myapp"

  # Entry with desktop-file-style arguments
  - name: "Sublime Text"
    icon: "sublime-text"
    command: "/snap/bin/sublime-text.subl %F"

  - name: "Files"
    icon: "system-file-manager"
    command: "caja %U"

  - name: "Evolution"
    icon: "evolution"
    command: "evolution %U"
```

### Icon resolution

1. **File path** — If the icon starts with `/`, `./`, or `~/`, it's treated as a filesystem path.
2. **Theme icon** — Otherwise, it's looked up in the current icon theme (Adwaita, Humanity, hicolor, etc.).
3. **Fallback** — If neither resolves, `application-x-executable` is used.

---

## Usage

### Run the launcher

```bash
launchpad run
```

Or simply:

```bash
launchpad
```

To use a custom config path:

```bash
launchpad -c ~/.config/launchpad/launchpad.yaml
```

### Global hotkey

#### Register

```bash
launchpad register
```

This uses Mate's Marco gsettings to bind the shortcut. Picks a free slot (`command-10` / `run-command-10`).

To use a different shortcut:

```bash
launchpad register --shortcut "<Control><Alt>L"
```

Or set it in your config:

```yaml
shortcut: "<Control><Alt>L"
```

Then `launchpad register` will read it from the config.

#### Unregister

```bash
launchpad unregister
```

#### Manual setup (if `register` doesn't work)

1. Open **System → Preferences → Hardware → Keyboard Shortcuts**
2. Add a new shortcut:
    - **Name**: Launchpad
    - **Command**: `/usr/local/bin/launchpad` (or wherever you installed it)
    - **Shortcut**: <kbd>Super</kbd>+<kbd>Esc</kbd>

Or via CLI:

```bash
gsettings set org.mate.Marco.keybinding-commands.command-10 "/usr/local/bin/launchpad"
gsettings set org.mate.Marco.global-keybindings.run-command-10 "<Super>Escape"
```

---

## CLI Reference

| Command                | Description                                  |
|------------------------|----------------------------------------------|
| `launchpad run`        | Run the launchpad UI popup (default command) |
| `launchpad init`       | Generate a skeleton config file              |
| `launchpad register`   | Register the global shortcut key             |
| `launchpad unregister` | Remove the global shortcut key               |

### Flags

| Flag         | Alias | Description                                                 |
|--------------|-------|-------------------------------------------------------------|
| `--config`   | `-c`  | Path to config file (default: `~/.launchpad.yaml`)          |
| `--shortcut` | `-s`  | Shortcut key combination (for `register`, overrides config) |

---

## Keyboard Shortcuts (within launcher)

| Key                                 | Action                 |
|-------------------------------------|------------------------|
| <kbd>↑</kbd> / <kbd>↓</kbd>         | Navigate entries       |
| <kbd>Enter</kbd> / <kbd>Space</kbd> | Launch selected entry  |
| <kbd>Shift</kbd>+<kbd>Enter</kbd>   | Launch without closing |
| <kbd>Shift</kbd>+click              | Launch without closing |
| <kbd>Esc</kbd>                      | Close launcher         |

---

## Development

### Project structure

```
launchpad/
├── cmd/
│   └── launchpad/
│       └── main.go            # Entry point, CLI subcommands
├── internal/
│   ├── cmd/
│   │   ├── keybinding.go      # register/unregister via gsettings
│   │   └── init.go            # skeleton config generator
│   ├── config/
│   │   ├── config.go          # YAML config structs + loader
│   │   └── config_test.go     # config unit tests
│   └── ui/
│       ├── ui.go              # App lifecycle, window management, centering
│       ├── entries.go         # Entry list, keyboard/mouse handling, launching
│       ├── style.go           # CSS loading, header bar
│       └── style.css          # Embedded GTK stylesheet
├── README.md
├── .golangci.yml
├── Taskfile.yml
├── go.mod
└── go.sum
```

### Dependencies

- [`gotk4`](https://github.com/diamondburned/gotk4) — GTK4 Go bindings
- [`urfave/cli`](https://github.com/urfave/cli) — CLI framework
- `gopkg.in/yaml.v3` — YAML parsing

### Building

```bash
task build      # build binary
task run        # build + run
task lint       # golangci-lint + go vet
task check      # full quality gate
task tidy       # go mod tidy
```

### Linting

```bash
golangci-lint run
go vet ./...
```

---

## Known limitations

- **Window manager**: Designed and tested on Ubuntu Mate with the Marco window manager. May need adjustments for other
  WMs/compositors.

---

## License

MIT
