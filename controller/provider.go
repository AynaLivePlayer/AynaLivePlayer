package controller

import (
	"AynaLivePlayer/config"
	"AynaLivePlayer/player"
	"AynaLivePlayer/provider"
)

func PrepareMedia(media *player.Media) error {
	var err error
	if media.Title == "" || !media.Cover.Exists() {
		l.Trace("fetching media info")
		if err = provider.UpdateMedia(media); err != nil {
			l.Warn("fail to prepare media when fetch info", err)
			return err
		}
	}
	if media.Url == "" {
		l.Trace("fetching media url")
		if err = provider.UpdateMediaUrl(media); err != nil {
			l.Warn("fail to prepare media when url", err)
			return err
		}
	}
	if media.Lyric == "" {
		l.Trace("fetching media lyric")
		if err = provider.UpdateMediaLyric(media); err != nil {
			l.Warn("fail to prepare media when lyric", err)
		}
	}
	return nil
}

func MediaMatch(keyword string) *player.Media {
	l.Infof("Match media for %s", keyword)
	for _, p := range config.Provider.Priority {
		if pr, ok := provider.Providers[p]; ok {
			m := pr.MatchMedia(keyword)
			if m == nil {
				continue
			}
			if err := provider.UpdateMedia(m); err == nil {
				return m
			}
		} else {
			l.Warnf("Provider %s not exist", p)
		}
	}
	return nil
}

func Search(keyword string) ([]*player.Media, error) {
	l.Infof("Search for %s", keyword)
	for _, p := range config.Provider.Priority {
		if pr, ok := provider.Providers[p]; ok {
			r, err := pr.Search(keyword)
			if err != nil {
				l.Warn("Provider %s return err", err)
				continue
			}
			return r, err
		} else {
			l.Warnf("Provider %s not exist", p)
		}
	}
	return nil, provider.ErrorNoSuchProvider
}

func SearchWithProvider(keyword string, p string) ([]*player.Media, error) {
	l.Infof("Search for %s using %s", keyword, p)
	if pr, ok := provider.Providers[p]; ok {
		r, err := pr.Search(keyword)
		return r, err
	}
	l.Warnf("Provider %s not exist", p)
	return nil, provider.ErrorNoSuchProvider
}

func ApplyUser(medias []*player.Media, user interface{}) {
	for _, m := range medias {
		m.User = user
	}
}

func PreparePlaylist(playlist *player.Playlist) error {
	l.Debug("Prepare playlist ", playlist.Meta.(provider.Meta))
	medias, err := provider.GetPlaylist(playlist.Meta.(provider.Meta))
	if err != nil {
		l.Warn("prepare playlist failed ", err)
		return err
	}
	ApplyUser(medias, player.SystemUser)
	playlist.Replace(medias)
	return nil
}
