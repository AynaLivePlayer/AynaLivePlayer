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
	Now   *LyricLine
	Index int
	Total int
	Prev  []*LyricLine
	Next  []*LyricLine
}

type Lyric struct {
	Lyrics []*LyricLine
}

func LoadLyric(lyric string) *Lyric {
	tmp := make(map[float64]*LyricLine)
	times := make([]float64, 0)
	for _, line := range strings.Split(lyric, "\n") {
		lrc := timeTagRegex.ReplaceAllString(line, "")
		if len(lrc) > 0 && lrc[len(lrc)-1] == '\r' {
			lrc = lrc[:len(lrc)-1]
		}
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

func (l *Lyric) findIndexV1(time float64) int {
	for i := 0; i < len(l.Lyrics)-1; i++ {
		if l.Lyrics[i].Time <= time && time < l.Lyrics[i+1].Time {
			return i
		}
	}
	return -1
}

func (l *Lyric) findIndex(time float64) int {
	start := 0
	end := len(l.Lyrics) - 1
	mid := (start + end) / 2
	for start < end {
		if l.Lyrics[mid].Time <= time && time < l.Lyrics[mid+1].Time {
			return mid
		}
		if l.Lyrics[mid].Time > time {
			end = mid
		} else {
			start = mid
		}
		mid = (start + end) / 2
	}
	return -1
}

func (l *Lyric) Find(time float64) *LyricLine {
	idx := l.findIndex(time)
	if idx == -1 {
		return nil
	}
	return l.Lyrics[idx]
}

func (l *Lyric) FindContext(time float64, prev int, next int) *LyricContext {
	prev = -prev
	idx := l.findIndex(time)
	if idx == -1 {
		return nil
	}
	if (idx + prev) < 0 {
		prev = -idx
	}
	if (idx + 1 + next) > len(l.Lyrics) {
		next = len(l.Lyrics) - idx - 1
	}
	return &LyricContext{
		Now:   l.Lyrics[idx],
		Index: idx,
		Total: len(l.Lyrics),
		Prev:  l.Lyrics[idx+prev : idx],
		Next:  l.Lyrics[idx+1 : idx+1+next],
	}
}
