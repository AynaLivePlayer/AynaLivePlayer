package events

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUnmarshalEventData(t *testing.T) {
	eventData := CmdLiveRoomAddData{
		Title:    "test",
		Provider: "asdfasd",
		RoomKey:  "asdfasdf",
	}
	data, err := json.Marshal(eventData)
	require.NoError(t, err)
	val, err := UnmarshalEventData(CmdLiveRoomAdd, data)
	require.NoError(t, err)
	resultData, ok := val.(CmdLiveRoomAddData)
	require.True(t, ok)
	require.Equal(t, eventData, resultData)
}
