package events

import (
	"AynaLivePlayer/core/model"
	"AynaLivePlayer/pkg/event"
	"encoding/json"
	"errors"
	"reflect"
)

var EventsMapping = map[event.EventId]any{
	LiveRoomAddCmd:                      LiveRoomAddCmdEvent{},
	LiveRoomProviderUpdate:              LiveRoomProviderUpdateEvent{},
	LiveRoomRemoveCmd:                   LiveRoomRemoveCmdEvent{},
	LiveRoomRoomsUpdate:                 LiveRoomRoomsUpdateEvent{},
	LiveRoomStatusUpdate:                LiveRoomStatusUpdateEvent{},
	LiveRoomConfigChangeCmd:             LiveRoomConfigChangeCmdEvent{},
	LiveRoomOperationCmd:                LiveRoomOperationCmdEvent{},
	PlayerVolumeChangeCmd:               PlayerVolumeChangeCmdEvent{},
	PlayerPlayCmd:                       PlayerPlayCmdEvent{},
	PlayerPlayErrorUpdate:               PlayerPlayErrorUpdateEvent{},
	PlayerSeekCmd:                       PlayerSeekCmdEvent{},
	PlayerToggleCmd:                     PlayerToggleCmdEvent{},
	PlayerSetPauseCmd:                   PlayerSetPauseCmdEvent{},
	PlayerPlayNextCmd:                   PlayerPlayNextCmdEvent{},
	PlayerLyricRequestCmd:               PlayerLyricRequestCmdEvent{},
	PlayerLyricReload:                   PlayerLyricReloadEvent{},
	PlayerLyricPosUpdate:                PlayerLyricPosUpdateEvent{},
	PlayerPlayingUpdate:                 PlayerPlayingUpdateEvent{},
	PlayerPropertyPauseUpdate:           PlayerPropertyPauseUpdateEvent{},
	PlayerPropertyPercentPosUpdate:      PlayerPropertyPercentPosUpdateEvent{},
	PlayerPropertyIdleActiveUpdate:      PlayerPropertyIdleActiveUpdateEvent{},
	PlayerPropertyTimePosUpdate:         PlayerPropertyTimePosUpdateEvent{},
	PlayerPropertyDurationUpdate:        PlayerPropertyDurationUpdateEvent{},
	PlayerPropertyVolumeUpdate:          PlayerPropertyVolumeUpdateEvent{},
	PlayerVideoPlayerSetWindowHandleCmd: PlayerVideoPlayerSetWindowHandleCmdEvent{},
	PlayerSetAudioDeviceCmd:             PlayerSetAudioDeviceCmdEvent{},
	PlayerAudioDeviceUpdate:             PlayerAudioDeviceUpdateEvent{},
	PlaylistManagerSetSystemCmd:         PlaylistManagerSetSystemCmdEvent{},
	PlaylistManagerSystemUpdate:         PlaylistManagerSystemUpdateEvent{},
	PlaylistManagerRefreshCurrentCmd:    PlaylistManagerRefreshCurrentCmdEvent{},
	PlaylistManagerGetCurrentCmd:        PlaylistManagerGetCurrentCmdEvent{},
	PlaylistManagerCurrentUpdate:        PlaylistManagerCurrentUpdateEvent{},
	PlaylistManagerInfoUpdate:           PlaylistManagerInfoUpdateEvent{},
	PlaylistManagerAddPlaylistCmd:       PlaylistManagerAddPlaylistCmdEvent{},
	PlaylistManagerRemovePlaylistCmd:    PlaylistManagerRemovePlaylistCmdEvent{},
	MediaProviderUpdate:                 MediaProviderUpdateEvent{},
	SearchCmd:                           SearchCmdEvent{},
	SearchResultUpdate:                  SearchResultUpdateEvent{},
}

func init() {
	for _, v := range []model.PlaylistID{model.PlaylistIDSystem, model.PlaylistIDPlayer, model.PlaylistIDHistory} {
		EventsMapping[PlaylistDetailUpdate(v)] = PlaylistDetailUpdateEvent{}
		EventsMapping[PlaylistMoveCmd(v)] = PlaylistMoveCmdEvent{}
		EventsMapping[PlaylistSetIndexCmd(v)] = PlaylistSetIndexCmdEvent{}
		EventsMapping[PlaylistDeleteCmd(v)] = PlaylistDeleteCmdEvent{}
		EventsMapping[PlaylistInsertCmd(v)] = PlaylistInsertCmdEvent{}
		EventsMapping[PlaylistInsertUpdate(v)] = PlaylistInsertUpdateEvent{}
		EventsMapping[PlaylistNextCmd(v)] = PlaylistNextCmdEvent{}
		EventsMapping[PlaylistNextUpdate(v)] = PlaylistNextUpdateEvent{}
		EventsMapping[PlaylistModeChangeCmd(v)] = PlaylistModeChangeCmdEvent{}
		EventsMapping[PlaylistModeChangeUpdate(v)] = PlaylistModeChangeUpdateEvent{}
	}
}

func UnmarshalEventData(eventId event.EventId, data []byte) (any, error) {
	val, ok := EventsMapping[eventId]
	if !ok {
		return nil, errors.New("event id not found")
	}
	newVal := reflect.New(reflect.TypeOf(val))
	err := json.Unmarshal(data, newVal.Interface())
	if err != nil {
		return nil, err
	}
	return newVal.Elem().Interface(), nil
}
