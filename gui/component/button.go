package component

import (
	"AynaLivePlayer/gui/gutil"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type AsyncButton struct {
	widget.Button
}

func NewAsyncButton(label string, tapped func()) *AsyncButton {
	b := &AsyncButton{
		Button: widget.Button{
			Text: label,
		},
	}
	b.ExtendBaseWidget(b)
	b.SetOnTapped(tapped)
	return b
}

func NewAsyncButtonWithIcon(label string, icon fyne.Resource, tapped func()) *AsyncButton {
	b := &AsyncButton{
		Button: widget.Button{
			Text: label,
			Icon: icon,
		},
	}
	b.ExtendBaseWidget(b)
	b.SetOnTapped(tapped)
	return b
}

func (b *AsyncButton) SetOnTapped(f func()) {
	b.Button.OnTapped = func() {
		b.Disable()
		go func() {
			f()
			//time.Sleep(3 * time.Second)
			// todo use fyne.Do after upgrade to v2.6.x
			gutil.RunInFyneThread(b.Enable)
			//b.Enable()
			//fyne.Do(b.Enable)
		}()
	}
}
