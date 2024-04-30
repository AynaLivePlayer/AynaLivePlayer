package events

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUnmarshalEventData(t *testing.T) {
	eventData := LiveRoomAddCmdEvent{
		Title:    "test",
		Provider: "asdfasd",
		RoomKey:  "asdfasdf",
	}
	data, err := json.Marshal(eventData)
	require.NoError(t, err)
	val, err := UnmarshalEventData(LiveRoomAddCmd, data)
	require.NoError(t, err)
	resultData, ok := val.(LiveRoomAddCmdEvent)
	require.True(t, ok)
	require.Equal(t, eventData, resultData)
}
