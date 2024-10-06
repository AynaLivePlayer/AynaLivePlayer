//go:build windows

package sysmediacontrol

import (
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/global"
	"AynaLivePlayer/pkg/config"
	"AynaLivePlayer/pkg/event"
	"github.com/go-ole/go-ole"
	"github.com/saltosystems/winrt-go"
	"github.com/saltosystems/winrt-go/windows/foundation"
	"github.com/saltosystems/winrt-go/windows/media"
	"github.com/saltosystems/winrt-go/windows/media/playback"
	"github.com/saltosystems/winrt-go/windows/storage/streams"
	"syscall"
	"unsafe"
)

const (
	TicksPerMicrosecond int64 = 10
	TicksPerMillisecond       = TicksPerMicrosecond * 1000
	TicksPerSecond            = TicksPerMillisecond * 1000
)

var (
	shell32, _                                 = syscall.LoadLibrary("shell32.dll")
	SetCurrentProcessExplicitAppUserModelID, _ = syscall.GetProcAddress(shell32, "SetCurrentProcessExplicitAppUserModelID")
)

var (
	smtc                   *media.SystemMediaTransportControls
	_player                *playback.MediaPlayer // Note: Do not use it!!! useless player, just for get smtc
	buttonPressedEventGUID = winrt.ParameterizedInstanceGUID(
		foundation.GUIDTypedEventHandler,
		media.SignatureSystemMediaTransportControls,
		media.SignatureSystemMediaTransportControlsButtonPressedEventArgs,
	)
)

func must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}

func withDisplayUpdater(f func(updater *media.SystemMediaTransportControlsDisplayUpdater)) {
	updater := must(smtc.GetDisplayUpdater())
	f(updater)
	updater.Release()
}

func withMusicProperties(f func(updater *media.SystemMediaTransportControlsDisplayUpdater, properties *media.MusicDisplayProperties)) {
	updater := must(smtc.GetDisplayUpdater())
	properties := must(updater.GetMusicProperties())
	f(updater, properties)
	properties.Release()
	updater.Release()
}

func InitSystemMediaControl() {
	_ = ole.RoInitialize(1)

	sptr, _ := syscall.UTF16PtrFromString("Aynakeya." + config.ProgramName)
	syscall.SyscallN(SetCurrentProcessExplicitAppUserModelID, uintptr(unsafe.Pointer(sptr)))
	_player = must(playback.NewMediaPlayer())
	smtc = must(_player.GetSystemMediaTransportControls())
	cmdManager := must(_player.GetCommandManager())
	_ = cmdManager.SetIsEnabled(false)
	cmdManager.Release()
	_ = smtc.SetIsEnabled(true)
	_ = smtc.SetIsPauseEnabled(true)
	_ = smtc.SetIsPlayEnabled(true)
	_ = smtc.SetIsNextEnabled(true)
	_ = smtc.SetIsPreviousEnabled(true)
	_ = smtc.SetPlaybackStatus(media.MediaPlaybackStatusPlaying)

	withDisplayUpdater(func(updater *media.SystemMediaTransportControlsDisplayUpdater) {
		_ = updater.SetType(media.MediaPlaybackTypeMusic)
	})

	global.EventManager.RegisterA(events.PlayerPlayingUpdate, "sysmediacontrol.update_playing", func(event *event.Event) {
		data := event.Data.(events.PlayerPlayingUpdateEvent)
		withMusicProperties(func(updater *media.SystemMediaTransportControlsDisplayUpdater, properties *media.MusicDisplayProperties) {
			properties.SetArtist(data.Media.Info.Artist)
			properties.SetTitle(data.Media.Info.Title)
			properties.SetAlbumTitle(data.Media.Info.Album)
			if data.Media.Info.Cover.Url != "" {
				imgUri, _ := foundation.UriCreateUri(data.Media.Info.Cover.Url)
				defer imgUri.Release()
				stream, _ := streams.RandomAccessStreamReferenceCreateFromUri(imgUri)
				defer stream.Release()
				_ = updater.SetThumbnail(stream)
			} else {
				// todo: using cover data
			}
			_ = updater.Update()
		})
		if data.Removed {
			smtc.SetPlaybackStatus(media.MediaPlaybackStatusChanging)
		}
	})

	global.EventManager.RegisterA(events.PlayerPropertyPauseUpdate, "sysmediacontrol.update_paused", func(event *event.Event) {
		if event.Data.(events.PlayerPropertyPauseUpdateEvent).Paused {
			smtc.SetPlaybackStatus(media.MediaPlaybackStatusPaused)
		} else {
			smtc.SetPlaybackStatus(media.MediaPlaybackStatusPlaying)
		}
	})

	pressedHandler := foundation.NewTypedEventHandler(
		ole.NewGUID(buttonPressedEventGUID),
		func(_ *foundation.TypedEventHandler, _ unsafe.Pointer, args unsafe.Pointer) {
			eventArgs := (*media.SystemMediaTransportControlsButtonPressedEventArgs)(args)
			defer eventArgs.Release()
			switch val, _ := eventArgs.GetButton(); val {
			case media.SystemMediaTransportControlsButtonPlay:
				global.EventManager.CallA(
					events.PlayerSetPauseCmd, events.PlayerSetPauseCmdEvent{Pause: false})
			case media.SystemMediaTransportControlsButtonPause:
				global.EventManager.CallA(
					events.PlayerSetPauseCmd, events.PlayerSetPauseCmdEvent{Pause: true})
			case media.SystemMediaTransportControlsButtonNext:
				global.EventManager.CallA(
					events.PlayerPlayNextCmd, events.PlayerPlayNextCmdEvent{})
			case media.SystemMediaTransportControlsButtonPrevious:
				global.EventManager.CallA(events.PlayerSeekCmd, events.PlayerSeekCmdEvent{
					Position: 0,
					Absolute: true,
				})
			}
		},
	)
	_, _ = smtc.AddButtonPressed(pressedHandler)
	pressedHandler.Release()

	// todo: finish timeline properties
	// cuz win 11 are not display timeline properties now
	// i just ignore it
}

func Destroy() {
	smtc.Release()
	_player.Release()
}
