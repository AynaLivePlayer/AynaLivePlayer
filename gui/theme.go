package gui

import (
	"AynaLivePlayer/resource"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"image/color"
)

type myTheme struct{}

var _ fyne.Theme = (*myTheme)(nil)

// return bundled font resource
func (*myTheme) Font(s fyne.TextStyle) fyne.Resource {
	l().Debugf("12313123")
	if s.Monospace {
		return resource.FontMSYaHei
	}
	if s.Bold {
		if s.Italic {
			//return theme.DefaultTheme().Font(s)
			return resource.FontMSYaHeiBold
		}
		return resource.FontMSYaHei
	}
	if s.Italic {
		return resource.FontMSYaHei
	}
	return resource.FontMSYaHei
}

func (*myTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(n, v)
}

func (*myTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
}

func (*myTheme) Size(n fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(n)
}
