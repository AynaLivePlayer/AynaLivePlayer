package lyrics

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type ActiveLyricPosition int

const (
	// ActiveLyricPositionMiddle positions the active lyric line in the middle of the widget
	ActiveLyricPositionMiddle ActiveLyricPosition = iota

	// ActiveLyricPositionUpperMiddle positions the active lyric line
	// in the upper-middle of the widget, roughly 1/3 of the way down
	ActiveLyricPositionUpperMiddle
)

// LyricsViewer is a widget for displaying song lyrics.
// It supports synced and unsynced mode. In synced mode, the active line
// is highlighted and the widget can advance to the next line
// with an animated scroll. In unsynced mode all lyrics are shown
// in the active color and the user is allowed to scroll freely.
type LyricsViewer struct {
	widget.BaseWidget

	// Alignment controls the text alignment of the lyric lines
	Alignment fyne.TextAlign

	// TextSizeName is the theme size name that controls the size of the lyric lines.
	// Defaults to theme.SizeNameSubHeadingText.
	TextSizeName fyne.ThemeSizeName

	// ActiveLyricColorName is the theme color name that the currently active
	// lyric line will be drawn in synced mode, or all lyrics in non-synced mode.
	// Defaults to theme.ColorNameForeground.
	ActiveLyricColorName fyne.ThemeColorName

	// InactiveLyricColorName is the theme color name that the inactive lyric lines
	// will be drawn in synced mode. Defaults to theme.ColorNameDisabled.
	InactiveLyricColorName fyne.ThemeColorName

	// HoveredLyricColorName is the theme color name that hovered lyric lines
	// will be drawn in synced mode when an OnLyricTapped callback is set.
	// Defaults to theme.ColorNameHover.
	HoveredLyricColorName fyne.ThemeColorName

	// ActiveLyricPosition sets the vertical positioning of the active lyric line
	// in synced mode.
	ActiveLyricPosition ActiveLyricPosition

	// OnLyricTapped sets a callback function that is invoked when a
	// synced lyric line is tapped. The line number is *one-indexed*.
	// Typically used to seek to the timecode of the given lyric.
	// When showing unsynced lyrics, or if this callback is unset,
	// the visual styling of the widget will not indicate interactivity.
	OnLyricTapped func(lineNum int)

	lines  []string
	synced bool

	// one-indexed - 0 means before the first line
	// during an animation, currentLine is the line
	// that will be scrolled when the animation is finished
	currentLine int

	prototypeLyricLineSize fyne.Size

	scroll *container.Scroll
	vbox   *fyne.Container

	// nil when an animation is not currently running
	anim            *fyne.Animation
	animStartOffset float32
}

// NewLyricsViewer returns a new lyrics viewer.
func NewLyricsViewer() *LyricsViewer {
	s := &LyricsViewer{}
	s.ExtendBaseWidget(s)
	s.prototypeLyricLineSize = s.newLyricLine("Hello...", 0, false).MinSize()
	return s
}

// SetLyrics sets the lyrics and also resets the current line to 0 if synced.
func (l *LyricsViewer) SetLyrics(lines []string, synced bool) {
	l.lines = lines
	l.synced = synced
	l.currentLine = 0
	if l.scroll != nil {
		if synced {
			l.scroll.Direction = container.ScrollNone
		} else {
			l.scroll.Direction = container.ScrollVerticalOnly
		}
	}
	l.updateContent()
}

// SetCurrentLine sets the current line that the lyric viewer is scrolled to.
// Argument is *one-indexed* - SetCurrentLine(0) means setting the scroll to be
// before the first line. In unsynced mode this is a no-op. This function is
// typically called when the user has seeked the playing song to a new position.
func (l *LyricsViewer) SetCurrentLine(line int) {
	if line < 0 || line > len(l.lines) {
		// do not panic, just ignore invalid input
		return
	}
	if l.vbox == nil || !l.synced {
		l.currentLine = line
		return // renderer not created yet or unsynced mode
	}
	inactiveColor := l.inactiveLyricColor()
	if l.checkStopAnimation() && l.currentLine > 1 {
		// we were in the middle of animation
		// make sure prev line is right color
		l.setLineColor(l.vbox.Objects[l.currentLine-1].(*lyricLine), inactiveColor, true)
	}
	if l.currentLine != 0 {
		l.setLineColor(l.vbox.Objects[l.currentLine].(*lyricLine), inactiveColor, true)
	}
	l.currentLine = line
	if l.currentLine != 0 {
		l.setLineColor(l.vbox.Objects[l.currentLine].(*lyricLine), l.activeLyricColor(), true)
	}
	l.scroll.Offset.Y = l.offsetForLine(l.currentLine)
	l.scroll.Refresh()
}

// NextLine advances the lyric viewer to the next line with an animated scroll.
// In unsynced mode this is a no-op.
func (l *LyricsViewer) NextLine() {
	if l.vbox == nil || !l.synced {
		return // no renderer yet, or unsynced lyrics (no-op)
	}

	if l.currentLine == len(l.lines) {
		return // already at last line
	}
	if l.checkStopAnimation() {
		// we were in the middle of animation - short-circuit it to completed
		// make sure prev and current lines are right color and scrolled to the end
		if l.currentLine > 1 {
			l.setLineColor(l.vbox.Objects[l.currentLine-1].(*lyricLine), l.inactiveLyricColor(), true)
		}
		l.setLineColor(l.vbox.Objects[l.currentLine].(*lyricLine), l.activeLyricColor(), true)
		l.scroll.Offset.Y = l.offsetForLine(l.currentLine)
	}
	l.currentLine++

	var prevLine, nextLine *lyricLine
	if l.currentLine > 1 {
		prevLine = l.vbox.Objects[l.currentLine-1].(*lyricLine)
	}
	if l.currentLine <= len(l.lines) {
		nextLine = l.vbox.Objects[l.currentLine].(*lyricLine)
	}

	l.setupScrollAnimation(prevLine, nextLine)
	l.anim.Start()
}

func (l *LyricsViewer) Refresh() {
	l.updateContent()
}

func (l *LyricsViewer) MinSize() fyne.Size {
	// overridden because NoScroll will have minSize encompass the full lyrics
	// note also that leaving this to the renderer MinSize, based on the
	// VBox with RichText lines inside Scroll, may lead to race conditions
	// (https://github.com/fyne-io/fyne/issues/4890)
	minHeight := l.prototypeLyricLineSize.Height*3 + theme.Padding()*2
	return fyne.NewSize(l.prototypeLyricLineSize.Width, minHeight)
}

func (l *LyricsViewer) Resize(size fyne.Size) {
	l.updateSpacerSize(size)
	l.BaseWidget.Resize(size)
	if l.vbox == nil {
		return // renderer not created yet
	}
	if l.anim == nil {
		l.scroll.Offset = fyne.NewPos(0, l.offsetForLine(l.currentLine))
		l.scroll.Refresh()
	} else {
		// animation is running - update its reference scroll pos
		l.animStartOffset = l.offsetForLine(l.currentLine - 1)
	}
}

func (l *LyricsViewer) updateSpacerSize(size fyne.Size) {
	if l.vbox == nil {
		return // renderer not created yet
	}

	ht := size.Height / 2
	if l.ActiveLyricPosition == ActiveLyricPositionUpperMiddle {
		ht = size.Height / 3
	}

	var topSpaceHeight, bottomSpaceHeight float32
	if l.synced {
		topSpaceHeight = ht + l.prototypeLyricLineSize.Height/2
		// end spacer only needs to be big enough - can't be too big
		// so use a very simple height calculation
		bottomSpaceHeight = size.Height
	}
	l.vbox.Objects[0].(*vSpace).Height = topSpaceHeight
	l.vbox.Objects[len(l.vbox.Objects)-1].(*vSpace).Height = bottomSpaceHeight
}

func (l *LyricsViewer) updateContent() {
	if l.vbox == nil {
		return // renderer not created yet
	}
	l.checkStopAnimation()

	lnObj := len(l.vbox.Objects)
	if lnObj == 0 {
		l.vbox.Objects = append(l.vbox.Objects, NewVSpace(0), NewVSpace(0))
		lnObj = 2
	}
	l.updateSpacerSize(l.Size())
	endSpacer := l.vbox.Objects[lnObj-1]
	for i, line := range l.lines {
		lineNum := i + 1 // one-indexed
		useActiveColor := !l.synced || l.currentLine == lineNum
		if lineNum < lnObj-1 {
			rt := l.vbox.Objects[lineNum].(*lyricLine)
			if useActiveColor {
				l.setLineColor(rt, l.activeLyricColor(), false)
			} else {
				l.setLineColor(rt, l.inactiveLyricColor(), false)
			}
			l.setLineTextAndProperties(rt, line, lineNum, true)
		} else if lineNum < lnObj {
			// replacing end spacer (last element in Objects) with a new richtext
			l.vbox.Objects[lineNum] = l.newLyricLine(line, lineNum, useActiveColor)
		} else {
			// extending the Objects slice
			l.vbox.Objects = append(l.vbox.Objects, l.newLyricLine(line, lineNum, useActiveColor))
		}
	}
	for i := len(l.lines) + 1; i < lnObj; i++ {
		l.vbox.Objects[i] = nil
	}
	l.vbox.Objects = l.vbox.Objects[:len(l.lines)+1]
	l.vbox.Objects = append(l.vbox.Objects, endSpacer)
	l.vbox.Refresh()
	l.scroll.Offset.Y = l.offsetForLine(l.currentLine)
	l.scroll.Refresh()
}

func (l *LyricsViewer) setupScrollAnimation(currentLine, nextLine *lyricLine) {
	// calculate total scroll distance for the animation
	scrollDist := theme.Padding()
	if currentLine != nil {
		scrollDist += currentLine.Size().Height / 2
	} else {
		scrollDist += l.prototypeLyricLineSize.Height / 2
	}
	if nextLine != nil {
		scrollDist += nextLine.Size().Height / 2
	} else {
		scrollDist += l.prototypeLyricLineSize.Height / 2
	}

	l.animStartOffset = l.scroll.Offset.Y
	var alreadyUpdated bool
	l.anim = fyne.NewAnimation(140*time.Millisecond, func(f float32) {
		l.scroll.Offset.Y = l.animStartOffset + f*scrollDist
		l.scroll.Refresh()
		if !alreadyUpdated && f >= 0.5 {
			if nextLine != nil {
				l.setLineColor(nextLine, l.activeLyricColor(), true)
			}
			if currentLine != nil {
				l.setLineColor(currentLine, l.inactiveLyricColor(), true)
			}
			alreadyUpdated = true
		}
		if f == 1 /*end of animation*/ {
			l.anim = nil
		}
	})
	l.anim.Curve = fyne.AnimationEaseInOut
}

func (l *LyricsViewer) offsetForLine(lineNum int /*one-indexed*/) float32 {
	if lineNum == 0 {
		return 0
	}
	pad := theme.Padding()
	offset := pad + l.prototypeLyricLineSize.Height/2
	for i := 1; i <= lineNum; i++ {
		if i > 1 {
			offset += l.vbox.Objects[i-1].MinSize().Height/2 + pad
		}
		offset += l.vbox.Objects[i].MinSize().Height / 2
	}
	return offset
}

func (l *LyricsViewer) newLyricLine(text string, lineNum int, useActiveColor bool) *lyricLine {
	ll := newLyricLine(text, nil)
	l.setLineTextAndProperties(ll, text, lineNum, false)
	ll.HoveredColorName = l.hoveredLyricColor()
	if useActiveColor {
		ll.ColorName = l.activeLyricColor()
	} else {
		ll.ColorName = l.inactiveLyricColor()
	}

	return ll
}

func (l *LyricsViewer) setLineTextAndProperties(ll *lyricLine, text string, lineNum int, refresh bool) {
	ll.Text = text
	ll.SizeName = l.textSizeName()
	ll.Alignment = l.Alignment
	ll.Tappable = l.synced && l.OnLyricTapped != nil
	ll.onTapped = func() {
		if l.OnLyricTapped != nil {
			l.OnLyricTapped(lineNum)
		}
	}
	if refresh {
		ll.Refresh()
	}
}

func (l *LyricsViewer) setLineColor(ll *lyricLine, colorName fyne.ThemeColorName, refresh bool) {
	ll.ColorName = colorName
	ll.HoveredColorName = l.hoveredLyricColor()
	if refresh {
		ll.Refresh()
	}
}

func (l *LyricsViewer) activeLyricColor() fyne.ThemeColorName {
	if l.ActiveLyricColorName != "" {
		return l.ActiveLyricColorName
	}
	return theme.ColorNameForeground
}

func (l *LyricsViewer) inactiveLyricColor() fyne.ThemeColorName {
	if l.InactiveLyricColorName != "" {
		return l.InactiveLyricColorName
	}
	return theme.ColorNameDisabled
}

func (l *LyricsViewer) hoveredLyricColor() fyne.ThemeColorName {
	if l.HoveredLyricColorName != "" {
		return l.HoveredLyricColorName
	}
	return theme.ColorNameHyperlink
}

func (l *LyricsViewer) textSizeName() fyne.ThemeSizeName {
	if l.TextSizeName != "" {
		return l.TextSizeName
	}
	return theme.SizeNameSubHeadingText
}

func (l *LyricsViewer) checkStopAnimation() bool {
	if l.anim != nil {
		l.anim.Stop()
		l.anim = nil
		return true
	}
	return false
}

func (l *LyricsViewer) CreateRenderer() fyne.WidgetRenderer {
	l.vbox = container.NewVBox()
	l.scroll = container.NewScroll(l.vbox)
	if l.synced {
		l.scroll.Direction = container.ScrollNone
	} else {
		l.scroll.Direction = container.ScrollVerticalOnly
	}
	l.updateContent()
	return widget.NewSimpleRenderer(l.scroll)
}

type vSpace struct {
	widget.BaseWidget

	Height float32
}

func NewVSpace(height float32) *vSpace {
	v := &vSpace{Height: height}
	v.ExtendBaseWidget(v)
	return v
}

func (v *vSpace) MinSize() fyne.Size {
	return fyne.NewSize(0, v.Height)
}

func (v *vSpace) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(layout.NewSpacer())
}
