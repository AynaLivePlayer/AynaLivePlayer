package gctx

import (
	_logger "AynaLivePlayer/pkg/logger"
	"fyne.io/fyne/v2"
)

// gui context

const (
	EventChannel = "gui"
)

var Logger _logger.ILogger = nil
var Context *GuiContext = nil

type GuiContext struct {
	App                 fyne.App    // application
	Window              fyne.Window // main window
	EventChannel        string
	onMainWindowClosing []func()
}

func NewGuiContext(app fyne.App, mainWindow fyne.Window) *GuiContext {
	return &GuiContext{
		App:                 app,
		Window:              mainWindow,
		EventChannel:        EventChannel,
		onMainWindowClosing: make([]func(), 0),
	}
}

func (c *GuiContext) Init() {
	c.Window.SetOnClosed(func() {
		for idx, f := range c.onMainWindowClosing {
			Logger.Debugf("runing gui closing handler #%d", idx)
			f()
		}
	})
}

func (c *GuiContext) OnMainWindowClosing(f func()) {
	c.onMainWindowClosing = append(c.onMainWindowClosing, f)
}
