package todo

import (
	"AynaLivePlayer/adapters/provider"
	"AynaLivePlayer/core/adapter"
	"AynaLivePlayer/core/model"
	"AynaLivePlayer/pkg/config"
)

type ProviderController struct {
	config.BaseConfig
	Priority []string
	LocalDir string
	log      adapter.ILogger
}

func (pc *ProviderController) Name() string {
	return "Provider"
}

func NewProviderController(
	log adapter.ILogger,
) adapter.IProviderController {
	p := &ProviderController{
		Priority: []string{"netease", "kuwo", "bilibili", "local", "bilibili-video"},
		LocalDir: "./music",
		log:      log,
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
		pc.log.Debug("fetching media info")
		if err = provider.UpdateMedia(media); err != nil {
			pc.log.Warn("fail to prepare media when fetch info ", err)
			return err
		}
	}
	if media.Url == "" {
		pc.log.Debug("fetching media url")
		if err = provider.UpdateMediaUrl(media); err != nil {
			pc.log.Warn("fail to prepare media when url ", err)
			return err
		}
	}
	if media.Lyric == "" {
		pc.log.Debug("fetching media lyric")
		if err = provider.UpdateMediaLyric(media); err != nil {
			pc.log.Warn("fail to prepare media when lyric", err)
		}
	}
	return nil
}

func (pc *ProviderController) MediaMatch(keyword string) *model.Media {
	pc.log.Infof("Match media for %s", keyword)
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
			pc.log.Warnf("Provider %s not exist", p)
		}
	}
	return nil
}

func (pc *ProviderController) Search(keyword string) ([]*model.Media, error) {
	pc.log.Infof("Search for %s", keyword)
	for _, p := range pc.Priority {
		if pr, ok := provider.Providers[p]; ok {
			r, err := pr.Search(keyword)
			if err != nil {
				pc.log.Warn("Provider %s return err", err)
				continue
			}
			return r, err
		} else {
			pc.log.Warnf("Provider %s not exist", p)
		}
	}
	return nil, provider.ErrorNoSuchProvider
}

func (pc *ProviderController) SearchWithProvider(keyword string, p string) ([]*model.Media, error) {
	pc.log.Infof("Search for %s using %s", keyword, p)
	if pr, ok := provider.Providers[p]; ok {
		r, err := pr.Search(keyword)
		return r, err
	}
	pc.log.Warnf("Provider %s not exist", p)
	return nil, provider.ErrorNoSuchProvider
}

func (pc *ProviderController) PreparePlaylist(playlist adapter.IPlaylist) error {
	pc.log.Debug("Prepare playlist ", playlist.Identifier())
	medias, err := provider.GetPlaylist(&playlist.Model().Meta)
	if err != nil {
		pc.log.Warn("prepare playlist failed ", err)
		return err
	}
	model.ApplyUser(medias, SystemUser)
	playlist.Replace(medias)
	return nil
}
