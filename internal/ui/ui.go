package ui

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"

	"github.com/abunjevac/launchpad/internal/config"
)

// keepOpen suppresses ConnectLeave from closing the window.
// Used by Shift+Enter to keep Launchpad open after launching an app.
var keepOpen atomic.Bool

// Run starts the launchpad GTK application.
func Run(cfg *config.Config) {
	app := gtk.NewApplication("io.github.abunjevac.launchpad", gio.ApplicationFlagsNone)

	app.ConnectActivate(func() {
		newWindow(app, cfg)
	})

	if code := app.Run([]string{}); code != 0 {
		os.Exit(code)
	}
}

func newWindow(app *gtk.Application, cfg *config.Config) {
	win := gtk.NewApplicationWindow(app)

	win.SetDecorated(false)
	win.SetResizable(false)
	win.SetModal(false)
	win.SetTitle("io.github.abunjevac.launchpad")

	vbox := gtk.NewBox(gtk.OrientationVertical, 0)

	vbox.AddCSSClass("launchpad")

	header := buildHeader(win)

	vbox.Append(header)

	scrolled := gtk.NewScrolledWindow()

	scrolled.SetPolicy(gtk.PolicyNever, gtk.PolicyAutomatic)
	scrolled.SetPropagateNaturalWidth(true)
	scrolled.SetPropagateNaturalHeight(true)
	scrolled.AddCSSClass("entries-scroll")

	list := gtk.NewListBox()

	list.SetSelectionMode(gtk.SelectionSingle)
	list.SetActivateOnSingleClick(true)
	list.AddCSSClass("entries")

	populateEntries(list, cfg)

	scrolled.SetChild(list)
	vbox.Append(scrolled)
	win.SetChild(vbox)

	applyCSS()
	connectKeyboard(win, list, cfg)
	connectAutoclose(win)
	connectRowActivated(list, win, cfg)

	win.SetDefaultSize(320, 0)
	win.Present()

	centerWindow(win)
}

// connectAutoclose closes the window when it loses focus,
// unless keepOpen is set.
func connectAutoclose(win *gtk.ApplicationWindow) {
	focusCtrl := gtk.NewEventControllerFocus()

	focusCtrl.ConnectLeave(func() {
		if !keepOpen.Load() {
			win.Close()
		}
	})

	win.AddController(focusCtrl)
}

// centerWindow positions the window at the center of the monitor
// where the surface currently resides.
func centerWindow(win *gtk.ApplicationWindow) {
	display := gdk.DisplayGetDefault()

	if display == nil {
		return
	}

	// go via Widget.Native to disambiguate from embedded coreglib.Object.Native
	native := win.Widget.Native()

	if native == nil {
		return
	}

	surface := gdk.BaseSurface(native.Surface())

	var once sync.Once

	surface.ConnectLayout(func(width, height int) {
		if width < 50 || height < 50 {
			return
		}

		once.Do(func() {
			go func() {
				monitor := display.MonitorAtSurface(surface)

				if monitor == nil {
					return
				}

				geo := monitor.Geometry()

				centerX := max(geo.X()+(geo.Width()-width)/2, 0)
				centerY := max(geo.Y()+(geo.Height()-height)/2, 0)

				for range 10 {
					time.Sleep(10 * time.Millisecond)

					cmd := exec.Command(
						"xdotool", "search", "--onlyvisible",
						"--name", "io.github.abunjevac.launchpad",
						"windowmove", strconv.Itoa(centerX), strconv.Itoa(centerY),
					)

					if err := cmd.Run(); err == nil {
						return
					}
				}

				_, _ = fmt.Fprintf(os.Stderr, "centering failed after 10 retries\n")
			}()
		})
	})
}
