package core

import (
	"AynaLivePlayer/config"
	"AynaLivePlayer/controller"
	"AynaLivePlayer/model"
	"AynaLivePlayer/repo/provider"
)

type ProviderController struct {
	config.BaseConfig
	Priority []string
	LocalDir string
}

func (pc *ProviderController) Name() string {
	return "Provider"
}

func NewProviderController() controller.IProviderController {
	p := &ProviderController{
		Priority: []string{"netease", "kuwo", "bilibili", "local", "bilibili-video"},
		LocalDir: "./music",
	}
	config.LoadConfig(p)
	provider.NewLocal(p.LocalDir)
	return p
}

func (pc *ProviderController) GetPriority() []string {
	return pc.Priority
}

func (pc *ProviderController) PrepareMedia(media *model.Media) error {
	var err error
	if media.Title == "" || !media.Cover.Exists() {
		lg.Trace("fetching media info")
		if err = provider.UpdateMedia(media); err != nil {
			lg.Warn("fail to prepare media when fetch info", err)
			return err
		}
	}
	if media.Url == "" {
		lg.Trace("fetching media url")
		if err = provider.UpdateMediaUrl(media); err != nil {
			lg.Warn("fail to prepare media when url", err)
			return err
		}
	}
	if media.Lyric == "" {
		lg.Trace("fetching media lyric")
		if err = provider.UpdateMediaLyric(media); err != nil {
			lg.Warn("fail to prepare media when lyric", err)
		}
	}
	return nil
}

func (pc *ProviderController) MediaMatch(keyword string) *model.Media {
	lg.Infof("Match media for %s", keyword)
	for _, p := range pc.Priority {
		if pr, ok := provider.Providers[p]; ok {
			m := pr.MatchMedia(keyword)
			if m == nil {
				continue
			}
			if err := provider.UpdateMedia(m); err == nil {
				return m
			}
		} else {
			lg.Warnf("Provider %s not exist", p)
		}
	}
	return nil
}

func (pc *ProviderController) Search(keyword string) ([]*model.Media, error) {
	lg.Infof("Search for %s", keyword)
	for _, p := range pc.Priority {
		if pr, ok := provider.Providers[p]; ok {
			r, err := pr.Search(keyword)
			if err != nil {
				lg.Warn("Provider %s return err", err)
				continue
			}
			return r, err
		} else {
			lg.Warnf("Provider %s not exist", p)
		}
	}
	return nil, provider.ErrorNoSuchProvider
}

func (pc *ProviderController) SearchWithProvider(keyword string, p string) ([]*model.Media, error) {
	lg.Infof("Search for %s using %s", keyword, p)
	if pr, ok := provider.Providers[p]; ok {
		r, err := pr.Search(keyword)
		return r, err
	}
	lg.Warnf("Provider %s not exist", p)
	return nil, provider.ErrorNoSuchProvider
}

func (pc *ProviderController) PreparePlaylist(playlist controller.IPlaylist) error {
	lg.Debug("Prepare playlist ", playlist.Name())
	medias, err := provider.GetPlaylist(&playlist.Model().Meta)
	if err != nil {
		lg.Warn("prepare playlist failed ", err)
		return err
	}
	controller.ApplyUser(medias, controller.SystemUser)
	playlist.Replace(medias)
	return nil
}
