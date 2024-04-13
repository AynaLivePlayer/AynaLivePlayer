package events

const MediaProviderUpdate = "update.media.provider.update"

type MediaProviderUpdateEvent struct {
	Providers []string
}
