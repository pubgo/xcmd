package xcmd_mdr

import (
	"github.com/MichaelMure/go-term-markdown"
	"github.com/awesome-gocui/gocui"
)

const padding = 4
const renderView = "render"

type ui struct {
	raw string
	// current width of the view
	width   int
	XOffset int
	YOffset int

	// number of lines in the rendered markdown
	lines int
}

func newUi(g *gocui.Gui) (*ui, error) {
	result := &ui{
		width: -1,
	}

	g.SetManagerFunc(result.layout)

	// Quit
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, result.quit); err != nil {
		return nil, err
	}
	if err := g.SetKeybinding(renderView, 'q', gocui.ModNone, result.quit); err != nil {
		return nil, err
	}

	// Up
	if err := g.SetKeybinding(renderView, 'k', gocui.ModNone, result.up); err != nil {
		return nil, err
	}
	if err := g.SetKeybinding(renderView, gocui.KeyCtrlP, gocui.ModNone, result.up); err != nil {
		return nil, err
	}
	if err := g.SetKeybinding(renderView, gocui.KeyArrowUp, gocui.ModNone, result.up); err != nil {
		return nil, err
	}
	// Down
	if err := g.SetKeybinding(renderView, 'j', gocui.ModNone, result.down); err != nil {
		return nil, err
	}
	if err := g.SetKeybinding(renderView, gocui.KeyCtrlN, gocui.ModNone, result.down); err != nil {
		return nil, err
	}
	if err := g.SetKeybinding(renderView, gocui.KeyArrowDown, gocui.ModNone, result.down); err != nil {
		return nil, err
	}

	// PageUp
	if err := g.SetKeybinding(renderView, gocui.KeyPgup, gocui.ModNone, result.pageUp); err != nil {
		return nil, err
	}
	// PageDown
	if err := g.SetKeybinding(renderView, gocui.KeyPgdn, gocui.ModNone, result.pageDown); err != nil {
		return nil, err
	}
	if err := g.SetKeybinding(renderView, gocui.KeySpace, gocui.ModNone, result.pageDown); err != nil {
		return nil, err
	}

	return result, nil
}

func (ui *ui) setContent(content []byte) {
	ui.raw = string(content)
	ui.width = -1
}

func (ui *ui) layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	v, err := g.SetView(renderView, ui.XOffset, -ui.YOffset, maxX, maxY, 0)
	if err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}

		v.Frame = false
		v.Wrap = false
	}

	if len(ui.raw) > 0 && ui.width != maxX {
		ui.width = maxX
		v.Clear()
		_, _ = v.Write(ui.render(g))
	}

	_, err = g.SetCurrentView(renderView)
	if err != nil {
		return err
	}

	return nil
}

func (ui *ui) render(g *gocui.Gui) []byte {
	maxX, _ := g.Size()

	opts := []markdown.Options{
		// needed when going through gocui
		markdown.WithImageDithering(markdown.DitheringWithBlocks),
	}

	rendered := markdown.Render(ui.raw, maxX-1-padding, padding, opts...)
	ui.lines = 0
	for _, b := range rendered {
		if b == '\n' {
			ui.lines++
		}
	}
	return rendered
}

func (ui *ui) quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func (ui *ui) up(g *gocui.Gui, v *gocui.View) error {
	ui.YOffset -= 1
	ui.YOffset = max(ui.YOffset, 0)
	return nil
}

func (ui *ui) down(g *gocui.Gui, v *gocui.View) error {
	_, maxY := g.Size()
	ui.YOffset += 1
	ui.YOffset = min(ui.YOffset, ui.lines-maxY+1)
	ui.YOffset = max(ui.YOffset, 0)
	return nil
}

func (ui *ui) pageUp(g *gocui.Gui, v *gocui.View) error {
	_, maxY := g.Size()
	ui.YOffset -= maxY / 2
	ui.YOffset = max(ui.YOffset, 0)
	return nil
}

func (ui *ui) pageDown(g *gocui.Gui, v *gocui.View) error {
	_, maxY := g.Size()
	ui.YOffset += maxY / 2
	ui.YOffset = min(ui.YOffset, ui.lines-maxY+1)
	ui.YOffset = max(ui.YOffset, 0)
	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}