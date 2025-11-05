//go:build darwin

package sysmediacontrol

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework MediaPlayer -framework AppKit

#import <Foundation/Foundation.h>
#import <MediaPlayer/MediaPlayer.h>
#import <AppKit/AppKit.h>

// Forward declaration for Go export function
extern void handleCommand(int);

// Command handler
static MPRemoteCommandHandlerStatus commandHandler(MPRemoteCommandEvent *event, int command) {
    handleCommand(command);
    return MPRemoteCommandHandlerStatusSuccess;
}

// Initialize media player controls
static void initMediaPlayer() {
    @autoreleasepool {
        MPRemoteCommandCenter *commandCenter = [MPRemoteCommandCenter sharedCommandCenter];

        [[commandCenter playCommand] addTargetWithHandler:^MPRemoteCommandHandlerStatus(MPRemoteCommandEvent * _Nonnull event) {
            return commandHandler(event, 0); // 0 = play
        }];

        [[commandCenter pauseCommand] addTargetWithHandler:^MPRemoteCommandHandlerStatus(MPRemoteCommandEvent * _Nonnull event) {
            return commandHandler(event, 1); // 1 = pause
        }];

        [[commandCenter nextTrackCommand] addTargetWithHandler:^MPRemoteCommandHandlerStatus(MPRemoteCommandEvent * _Nonnull event) {
            return commandHandler(event, 2); // 2 = next
        }];

        [[commandCenter previousTrackCommand] addTargetWithHandler:^MPRemoteCommandHandlerStatus(MPRemoteCommandEvent * _Nonnull event) {
            return commandHandler(event, 3); // 3 = previous
        }];

        // Enable commands
        [commandCenter playCommand].enabled = YES;
        [commandCenter pauseCommand].enabled = YES;
        [commandCenter nextTrackCommand].enabled = YES;
        [commandCenter previousTrackCommand].enabled = YES;
    }
}

// Update now playing info
static void updateNowPlaying(const char *title, const char *artist, const char *album,
                     double duration, double position, int isPlaying) {
    @autoreleasepool {
        MPNowPlayingInfoCenter *center = [MPNowPlayingInfoCenter defaultCenter];
        NSMutableDictionary *nowPlayingInfo = [center.nowPlayingInfo mutableCopy];
        if (nowPlayingInfo == nil) {
            nowPlayingInfo = [NSMutableDictionary dictionary];
        }

        if (title != NULL) {
            [nowPlayingInfo setObject:[NSString stringWithUTF8String:title]
                              forKey:MPMediaItemPropertyTitle];
        }

        if (artist != NULL) {
            [nowPlayingInfo setObject:[NSString stringWithUTF8String:artist]
                              forKey:MPMediaItemPropertyArtist];
        }

        if (album != NULL) {
            [nowPlayingInfo setObject:[NSString stringWithUTF8String:album]
                              forKey:MPMediaItemPropertyAlbumTitle];
        }

        if (duration > 0) {
            [nowPlayingInfo setObject:[NSNumber numberWithDouble:duration]
                              forKey:MPMediaItemPropertyPlaybackDuration];
        }

        if (position >= 0) {
            [nowPlayingInfo setObject:[NSNumber numberWithDouble:position]
                              forKey:MPNowPlayingInfoPropertyElapsedPlaybackTime];
        }

        [nowPlayingInfo setObject:[NSNumber numberWithDouble:(isPlaying ? 1.0 : 0.0)]
                          forKey:MPNowPlayingInfoPropertyPlaybackRate];

        center.nowPlayingInfo = nowPlayingInfo;
    }
}

// Update artwork from URL
static void updateArtworkFromURL(const char *urlString) {
    @autoreleasepool {
        if (urlString == NULL) return;

        NSString *urlStr = [NSString stringWithUTF8String:urlString];
        NSURL *url = [NSURL URLWithString:urlStr];
        if (url == NULL) return;

        dispatch_async(dispatch_get_global_queue(DISPATCH_QUEUE_PRIORITY_DEFAULT, 0), ^{
            NSData *imageData = [NSData dataWithContentsOfURL:url];
            if (imageData) {
                NSImage *image = [[NSImage alloc] initWithData:imageData];
                if (image) {
                    MPMediaItemArtwork *artwork = [[MPMediaItemArtwork alloc]
                        initWithBoundsSize:image.size
                        requestHandler:^NSImage * _Nonnull(CGSize size) {
                            return image;
                        }];

                    dispatch_async(dispatch_get_main_queue(), ^{
                        MPNowPlayingInfoCenter *center = [MPNowPlayingInfoCenter defaultCenter];
                        NSMutableDictionary *nowPlayingInfo = [center.nowPlayingInfo mutableCopy];
                        if (nowPlayingInfo == nil) {
                            nowPlayingInfo = [NSMutableDictionary dictionary];
                        }
                        [nowPlayingInfo setObject:artwork forKey:MPMediaItemPropertyArtwork];
                        center.nowPlayingInfo = nowPlayingInfo;
                    });
                }
            }
        });
    }
}

// Clear now playing info
static void clearNowPlaying() {
    @autoreleasepool {
        MPNowPlayingInfoCenter *center = [MPNowPlayingInfoCenter defaultCenter];
        center.nowPlayingInfo = nil;
    }
}
*/
import "C"
import (
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/global"
	"AynaLivePlayer/pkg/eventbus"
	"AynaLivePlayer/pkg/logger"
	"unsafe"
)

var (
	log              logger.ILogger
	currentTitle     string
	currentArtist    string
	currentAlbum     string
	currentDuration  float64
	currentPosition  float64
	currentIsPlaying bool
)

//export handleCommand
func handleCommand(command C.int) {
	switch command {
	case 0: // Play
		_ = global.EventBus.Publish(
			events.PlayerSetPauseCmd, events.PlayerSetPauseCmdEvent{Pause: false})
	case 1: // Pause
		_ = global.EventBus.Publish(
			events.PlayerSetPauseCmd, events.PlayerSetPauseCmdEvent{Pause: true})
	case 2: // Next
		_ = global.EventBus.Publish(
			events.PlayerPlayNextCmd, events.PlayerPlayNextCmdEvent{})
	case 3: // Previous
		_ = global.EventBus.Publish(events.PlayerSeekCmd, events.PlayerSeekCmdEvent{
			Position: 0,
			Absolute: true,
		})
	}
}

func updateNowPlayingInfo() {
	titleC := C.CString(currentTitle)
	artistC := C.CString(currentArtist)
	albumC := C.CString(currentAlbum)
	defer C.free(unsafe.Pointer(titleC))
	defer C.free(unsafe.Pointer(artistC))
	defer C.free(unsafe.Pointer(albumC))

	isPlaying := 0
	if currentIsPlaying {
		isPlaying = 1
	}

	C.updateNowPlaying(
		titleC,
		artistC,
		albumC,
		C.double(currentDuration),
		C.double(currentPosition),
		C.int(isPlaying),
	)
}

func InitSystemMediaControl() {
	log = global.Logger.WithPrefix("SMTC-Darwin")

	// Initialize media player controls
	C.initMediaPlayer()

	// Subscribe to player playing update events
	global.EventBus.Subscribe("", events.PlayerPlayingUpdate, "sysmediacontrol.update_playing", func(event *eventbus.Event) {
		data := event.Data.(events.PlayerPlayingUpdateEvent)

		if data.Removed {
			C.clearNowPlaying()
			currentTitle = ""
			currentArtist = ""
			currentAlbum = ""
			currentDuration = 0
			currentPosition = 0
			return
		}

		currentTitle = data.Media.Info.Title
		currentArtist = data.Media.Info.Artist
		currentAlbum = data.Media.Info.Album

		updateNowPlayingInfo()

		// Update artwork if available
		if data.Media.Info.Cover.Url != "" {
			urlC := C.CString(data.Media.Info.Cover.Url)
			C.updateArtworkFromURL(urlC)
			C.free(unsafe.Pointer(urlC))
		}
	})

	// Subscribe to pause state updates
	global.EventBus.Subscribe("", events.PlayerPropertyPauseUpdate, "sysmediacontrol.update_paused", func(event *eventbus.Event) {
		data := event.Data.(events.PlayerPropertyPauseUpdateEvent)
		currentIsPlaying = !data.Paused
		updateNowPlayingInfo()
	})

	// Subscribe to duration updates
	global.EventBus.Subscribe("", events.PlayerPropertyDurationUpdate, "sysmediacontrol.properties.duration", func(event *eventbus.Event) {
		data := event.Data.(events.PlayerPropertyDurationUpdateEvent)
		currentDuration = data.Duration
		updateNowPlayingInfo()
	})

	// Subscribe to time position updates
	global.EventBus.Subscribe("", events.PlayerPropertyTimePosUpdate, "sysmediacontrol.properties.time_pos", func(event *eventbus.Event) {
		data := event.Data.(events.PlayerPropertyTimePosUpdateEvent)
		currentPosition = data.TimePos
		updateNowPlayingInfo()
	})

	log.Info("macOS System Media Control initialized")
}

func Destroy() {
	C.clearNowPlaying()

	// Unsubscribe from all events
	global.EventBus.Unsubscribe(events.PlayerPlayingUpdate, "sysmediacontrol.update_playing")
	global.EventBus.Unsubscribe(events.PlayerPropertyPauseUpdate, "sysmediacontrol.update_paused")
	global.EventBus.Unsubscribe(events.PlayerPropertyDurationUpdate, "sysmediacontrol.properties.duration")
	global.EventBus.Unsubscribe(events.PlayerPropertyTimePosUpdate, "sysmediacontrol.properties.time_pos")

	log.Info("macOS System Media Control destroyed")
}
