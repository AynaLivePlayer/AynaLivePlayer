package model

import (
	"github.com/spf13/cast"
	"regexp"
	"sort"
	"strings"
)

var timeTagRegex = regexp.MustCompile("\\[[0-9]+:[0-9]+(\\.[0-9]+)?\\]")

type LyricLine struct {
	Time        float64 // in seconds
	Lyric       string
	Translation string
}

type LyricContext struct {
	Current *LyricLine
	Prev    []*LyricLine
	Next    []*LyricLine
}

type Lyric struct {
	Lyrics []*LyricLine
}

func LoadLyric(lyric string) *Lyric {
	tmp := make(map[float64]*LyricLine)
	times := make([]float64, 0)
	for _, line := range strings.Split(lyric, "\n") {
		lrc := timeTagRegex.ReplaceAllString(line, "")
		for _, time := range timeTagRegex.FindAllString(line, -1) {
			ts := strings.Split(time[1:len(time)-1], ":")
			t := cast.ToFloat64(ts[0])*60 + cast.ToFloat64(ts[1])
			times = append(times, t)
			tmp[t] = &LyricLine{
				Time:  t,
				Lyric: lrc,
			}
		}
	}
	sort.Float64s(times)
	lrcs := make([]*LyricLine, len(times))
	for index, time := range times {
		lrcs[index] = tmp[time]
	}
	if len(lrcs) == 0 {
		lrcs = append(lrcs, &LyricLine{Time: 0, Lyric: ""})
	}
	lrcs = append(lrcs, &LyricLine{
		Time:  99999999999,
		Lyric: "",
	})
	return &Lyric{Lyrics: lrcs}
}

func (l *Lyric) Find(time float64) *LyricLine {
	for i := 0; i < len(l.Lyrics)-1; i++ {
		if l.Lyrics[i].Time <= time && time < l.Lyrics[i+1].Time {
			return l.Lyrics[i]
		}
	}
	return nil
}

func (l *Lyric) FindContext(time float64, prev int, next int) *LyricContext {
	for i := 0; i < len(l.Lyrics)-1; i++ {
		if l.Lyrics[i].Time <= time && time < l.Lyrics[i+1].Time {
			if (i + prev) < 0 {
				prev = -i
			}
			if (i + 1 + next) > len(l.Lyrics) {
				next = len(l.Lyrics) - i - 1
			}
			return &LyricContext{
				Current: l.Lyrics[i],
				Prev:    l.Lyrics[i+prev : i],
				Next:    l.Lyrics[i+1 : i+1+next],
			}
		}
	}
	return nil
}
