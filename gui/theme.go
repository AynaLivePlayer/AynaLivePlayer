package gui

import (
	"AynaLivePlayer/resource"
	"fyne.io/fyne/v2"
	"image/color"

	xtheme "fyne.io/x/fyne/theme"
)

type myTheme struct{}

var _ fyne.Theme = (*myTheme)(nil)

// return bundled font resource
func (*myTheme) Font(s fyne.TextStyle) fyne.Resource {
	if s.Monospace {
		return resource.FontMSYaHei
	}
	if s.Bold {
		if s.Italic {
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
	return xtheme.AdwaitaTheme().Color(n, v)
}

func (*myTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return xtheme.AdwaitaTheme().Icon(n)
}

func (*myTheme) Size(n fyne.ThemeSizeName) float32 {
	return xtheme.AdwaitaTheme().Size(n)
}
