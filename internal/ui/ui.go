package ui

import (
	"os"
	"sync/atomic"

	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"

	"github.com/abunjevac/launchpad/internal/config"
)

// keepOpen suppresses ConnectLeave from closing the window.
// Used by Shift+Enter to keep Launchpad open after launching an app.
var keepOpen atomic.Bool

// Run starts the launchpad GTK application.
func Run(cfg *config.Config) {
	app := gtk.NewApplication("io.github.abunjevac.launchpad", gio.ApplicationDefaultFlags)

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
