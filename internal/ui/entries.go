package ui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"

	"github.com/abunjevac/launchpad/internal/config"
)

// populateEntries fills the list box with entries from config.
func populateEntries(list *gtk.ListBox, cfg *config.Config) {
	for _, entry := range cfg.Entries {
		if entry.Separator {
			row := gtk.NewListBoxRow()

			row.SetSelectable(false)
			row.SetFocusable(false)

			sep := gtk.NewSeparator(gtk.OrientationHorizontal)

			sep.SetMarginStart(8)
			sep.SetMarginEnd(8)
			sep.SetMarginTop(4)
			sep.SetMarginBottom(4)

			row.SetChild(sep)
			list.Append(row)

			continue
		}

		row := buildEntryRow(entry)

		list.Append(row)
	}
}

// buildEntryRow creates a list box row for a single entry.
func buildEntryRow(entry config.Entry) *gtk.ListBoxRow {
	row := gtk.NewListBoxRow()

	row.AddCSSClass("entry-row")

	img := createIcon(entry.Icon)

	if img == nil {
		img = gtk.NewImageFromIconName("application-x-executable")
	}

	img.SetPixelSize(32)

	content := gtk.NewBox(gtk.OrientationHorizontal, 8)

	content.SetMarginStart(12)
	content.SetMarginEnd(12)
	content.SetMarginTop(8)
	content.SetMarginBottom(8)

	content.Append(img)

	name := gtk.NewLabel(entry.Name)

	name.SetHAlign(gtk.AlignStart)
	name.SetVAlign(gtk.AlignCenter)
	name.AddCSSClass("entry-name")

	content.Append(name)

	row.SetTooltipText(entry.Name)
	row.SetChild(content)

	return row
}

// createIcon resolves an icon string and returns a Gtk Image widget.
// Supports file paths (starts with /, ./ or ~/) and named theme icons.
func createIcon(iconStr string) *gtk.Image {
	if iconStr == "" {
		return nil
	}

	var path string

	switch {
	case strings.HasPrefix(iconStr, "~/"):
		home, err := os.UserHomeDir()
		if err != nil {
			break
		}

		path = filepath.Join(home, iconStr[2:])

	case strings.HasPrefix(iconStr, "/") || strings.HasPrefix(iconStr, "./"):
		path = iconStr
	}

	if path != "" {
		if _, err := os.Stat(path); err == nil {
			return gtk.NewImageFromFile(path)
		}
	}

	return gtk.NewImageFromIconName(iconStr)
}

// rowToEntry finds the config entry corresponding to a list row.
func rowToEntry(row *gtk.ListBoxRow, cfg *config.Config) *config.Entry {
	targetIdx := row.Index()
	listIdx := 0

	for i := range cfg.Entries {
		entry := &cfg.Entries[i]

		if entry.Separator {
			listIdx++

			continue
		}

		if listIdx == targetIdx {
			return entry
		}

		listIdx++
	}

	return nil
}

// moveSelection moves the list selection by delta positions,
// skipping non-selectable rows (separators).
func moveSelection(list *gtk.ListBox, delta int) {
	if list == nil {
		return
	}

	selected := list.SelectedRow()
	start := 0

	if selected != nil {
		start = max(int(selected.Index())+delta, 0)
	}

	for {
		row := list.RowAtIndex(start)

		if row == nil {
			return
		}

		if row.Selectable() {
			list.SelectRow(row)
			row.GrabFocus()

			return
		}

		start += delta

		if start < 0 {
			return
		}
	}
}

// launchEntry executes an entry's command.
func launchEntry(entry *config.Entry, shiftHeld bool, win *gtk.ApplicationWindow) {
	if entry.Command == "" {
		return
	}

	cmd := exec.Command("sh", "-c", entry.Command)

	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "launching %s: %v\n", entry.Name, err)

		return
	}

	go func() {
		_ = cmd.Wait()
	}()

	if !shiftHeld {
		win.Close()

		return
	}

	// Suppress autoclose so the launched app stealing focus doesn't
	// close Launchpad. Re-enable after 2s so Launchpad doesn't linger
	// indefinitely.
	keepOpen.Store(true)

	time.AfterFunc(2*time.Second, func() {
		keepOpen.Store(false)
	})
}

// shiftHeld queries the current keyboard modifier state to check if
// Shift is held. This works regardless of when Shift was pressed
// (even before window mapping) because it queries the actual device.
func shiftHeld() bool {
	display := gdk.DisplayGetDefault()

	if display == nil {
		return false
	}

	seater := display.DefaultSeat()

	if seater == nil {
		return false
	}

	seat := gdk.BaseSeat(seater)
	keyboard := seat.Keyboard()

	if keyboard == nil {
		return false
	}

	return gdk.BaseDevice(keyboard).ModifierState().Has(gdk.ShiftMask)
}

// connectKeyboard handles keyboard navigation.
// Arrow keys select entries. Enter and Space activate the selected row.
func connectKeyboard(win *gtk.ApplicationWindow, list *gtk.ListBox, cfg *config.Config) {
	keyCtrl := gtk.NewEventControllerKey()

	keyCtrl.SetPropagationPhase(gtk.PhaseCapture)
	keyCtrl.ConnectKeyPressed(func(keyValue uint, keyCode uint, state gdk.ModifierType) bool {
		switch keyValue {
		case gdk.KEY_Escape:
			win.Close()

			return true

		case gdk.KEY_Up, gdk.KEY_Left:
			moveSelection(list, -1)

			return true

		case gdk.KEY_Down, gdk.KEY_Right:
			moveSelection(list, 1)

			return true

		case gdk.KEY_Return, gdk.KEY_KP_Enter, gdk.KEY_space:
			row := list.SelectedRow()

			if row == nil {
				return true
			}

			entry := rowToEntry(row, cfg)

			if entry == nil {
				return true
			}

			launchEntry(entry, shiftHeld(), win)

			return true
		}

		return false
	})

	win.AddController(keyCtrl)
}

// connectRowActivated handles row activation from mouse clicks.
func connectRowActivated(list *gtk.ListBox, win *gtk.ApplicationWindow, cfg *config.Config) {
	list.ConnectRowActivated(func(row *gtk.ListBoxRow) {
		entry := rowToEntry(row, cfg)

		if entry == nil {
			return
		}

		launchEntry(entry, shiftHeld(), win)
	})
}
