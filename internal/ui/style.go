package ui

import (
	_ "embed"

	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

//go:embed style.css
var css string

// buildHeader creates the header bar with a close button.
func buildHeader(win *gtk.ApplicationWindow) *gtk.Box {
	header := gtk.NewBox(gtk.OrientationHorizontal, 0)

	header.AddCSSClass("header")

	title := gtk.NewLabel("Launchpad")

	title.AddCSSClass("title")
	title.SetHAlign(gtk.AlignStart)

	header.Append(title)

	spacer := gtk.NewLabel("")

	spacer.SetHExpand(true)

	header.Append(spacer)

	closeBtn := gtk.NewButtonFromIconName("window-close-symbolic")

	closeBtn.AddCSSClass("close-button")
	closeBtn.SetTooltipText("Close (Esc)")
	closeBtn.ConnectClicked(func() {
		win.Close()
	})

	header.Append(closeBtn)

	return header
}

// applyCSS applies custom styling to the window.
func applyCSS() {
	provider := gtk.NewCSSProvider()

	provider.LoadFromString(css)

	gtk.StyleContextAddProviderForDisplay(
		gdk.DisplayGetDefault(),
		provider,
		gtk.STYLE_PROVIDER_PRIORITY_APPLICATION,
	)
}
