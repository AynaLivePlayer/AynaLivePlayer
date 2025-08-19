package model

type AudioDevice struct {
	Name        string
	Description string
}

type PlayerState int

const (
	PlayerStatePlaying PlayerState = iota
	PlayerStateLoading
	PlayerStateIdle
)

func (s PlayerState) NextState(next PlayerState) PlayerState {
	if s == PlayerStatePlaying {
		return next
	}
	if s == PlayerStateIdle {
		return next
	}
	if s == PlayerStateLoading {
		if next != PlayerStatePlaying {
			return PlayerStateLoading
		}
		return next
	}
	return next
}
