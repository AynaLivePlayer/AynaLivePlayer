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

	sptr, _ := syscall.UTF16PtrFromString(config.ProgramName)
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
			_ = updater.Update()
		})
	})

}

func Destroy() {
	smtc.Release()
	_player.Release()
}
