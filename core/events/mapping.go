package events

import (
	"AynaLivePlayer/core/model"
	"encoding/json"
	"errors"
	"reflect"
)

var EventsMapping = map[string]any{
	CmdLiveRoomAdd:                      CmdLiveRoomAddData{},
	LiveRoomProviderUpdate:              LiveRoomProviderUpdateEvent{},
	CmdLiveRoomRemove:                   CmdLiveRoomRemoveData{},
	UpdateLiveRoomRooms:                 UpdateLiveRoomRoomsData{},
	UpdateLiveRoomStatus:                UpdateLiveRoomStatusData{},
	CmdLiveRoomConfigChange:             CmdLiveRoomConfigChangeData{},
	CmdLiveRoomOperation:                CmdLiveRoomOperationData{},
	PlayerVolumeChangeCmd:               PlayerVolumeChangeCmdEvent{},
	PlayerPlayCmd:                       PlayerPlayCmdEvent{},
	PlayerPlayErrorUpdate:               PlayerPlayErrorUpdateEvent{},
	PlayerSeekCmd:                       PlayerSeekCmdEvent{},
	PlayerToggleCmd:                     PlayerToggleCmdEvent{},
	PlayerSetPauseCmd:                   PlayerSetPauseCmdEvent{},
	PlayerPlayNextCmd:                   PlayerPlayNextCmdEvent{},
	CmdGetCurrentLyric:                  CmdGetCurrentLyricData{},
	UpdateCurrentLyric:                  UpdateCurrentLyricData{},
	PlayerLyricPosUpdate:                PlayerLyricPosUpdateEvent{},
	PlayerPlayingUpdate:                 PlayerPlayingUpdateEvent{},
	PlayerPropertyPauseUpdate:           PlayerPropertyPauseUpdateEvent{},
	PlayerPropertyPercentPosUpdate:      PlayerPropertyPercentPosUpdateEvent{},
	PlayerPropertyStateUpdate:           PlayerPropertyStateUpdateEvent{},
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
	CmdMiaosicSearch:                    CmdMiaosicSearchData{},
	ReplyMiaosicSearch:                  ReplyMiaosicSearchData{},
	GUISetPlayerWindowOpenCmd:           GUISetPlayerWindowOpenCmdEvent{},
}

func init() {
	for _, v := range []model.PlaylistID{model.PlaylistIDSystem, model.PlaylistIDPlayer} {
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

func UnmarshalEventData(eventId string, data []byte) (any, error) {
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
