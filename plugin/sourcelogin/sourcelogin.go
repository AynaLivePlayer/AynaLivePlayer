package sourcelogin

import (
	"AynaLivePlayer/core/events"
	"AynaLivePlayer/global"
	"AynaLivePlayer/gui/component"
	config2 "AynaLivePlayer/gui/views/config"
	"AynaLivePlayer/pkg/config"
	"AynaLivePlayer/pkg/i18n"
	"AynaLivePlayer/pkg/logger"
	"AynaLivePlayer/resource"
	"bytes"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/AynaLivePlayer/miaosic"
	"github.com/skip2/go-qrcode"
)

const MODULE_PLGUIN_NETEASELOGIN = "plugin.neteaselogin"

type SourceLogin struct {
	SessionPath string `json:"session_path"`
	sessions    map[string]string
	log         logger.ILogger
	panel       fyne.CanvasObject
}

func (w *SourceLogin) OnLoad() {
	_ = config.LoadJson(w.SessionPath, &w.sessions)
}

func (w *SourceLogin) OnSave() {
	_ = config.SaveJson(w.SessionPath, &w.sessions)
}

func NewSourceLogin() *SourceLogin {
	return &SourceLogin{
		SessionPath: "config/source_session.json",
		log:         global.Logger.WithPrefix("plugin.sourcelogin"),
		sessions:    make(map[string]string),
	}
}

func (w *SourceLogin) Name() string {
	return "SourceLogin"
}

func (w *SourceLogin) Enable() error {
	config.LoadConfig(w)
	config2.AddConfigLayout(w)
	return nil
}

func (w *SourceLogin) Disable() error {
	w.log.Info("save session for all provider")
	providers := miaosic.ListAvailableProviders()
	for _, pname := range providers {
		if p, ok := miaosic.GetProvider(pname); ok {
			pl, ok2 := p.(miaosic.Loginable)
			if ok2 {
				w.log.Info("save session for %s", pname)
				w.sessions[pname] = pl.SaveSession()
			}
		}
	}
	return nil
}

func (w *SourceLogin) Title() string {
	return i18n.T("plugin.sourcelogin.title")
}

func (w *SourceLogin) Description() string {
	return i18n.T("plugin.sourcelogin.description")
}

func (w *SourceLogin) CreatePanel() fyne.CanvasObject {
	if w.panel != nil {
		return w.panel
	}
	currentUser := widget.NewLabel(i18n.T("plugin.sourcelogin.current_user.notlogin"))
	currentStatus := container.NewHBox(
		widget.NewLabel(i18n.T("plugin.sourcelogin.current_user")),
		currentUser)

	providers := miaosic.ListAvailableProviders()
	loginableProviders := make([]string, 0)
	loginables := make(map[string]miaosic.MediaProvider)
	for _, pname := range providers {
		if p, ok := miaosic.GetProvider(pname); ok {
			pl, ok2 := p.(miaosic.Loginable)
			if ok2 {
				loginableProviders = append(loginableProviders, pname)
				loginables[pname] = p
				if session, ok3 := w.sessions[pname]; ok3 {
					err := pl.RestoreSession(session)
					if err != nil {
						w.log.Error("failed to restore session for ", pname)
					}
				}
			}
		}
	}
	providerChoice := widget.NewSelect(loginableProviders, func(s string) {
		w.log.Info("switching provider to ", s)
		if s != "" {
			pvdr, _ := miaosic.GetProvider(s)
			provider := pvdr.(miaosic.Loginable)
			if provider.IsLogin() {
				currentUser.SetText(i18n.T("plugin.sourcelogin.current_user.loggedin"))
			} else {
				currentUser.SetText(i18n.T("plugin.sourcelogin.current_user.notlogin"))
			}
		}
	})

	sourcePanel := container.NewGridWithColumns(2,
		providerChoice, currentStatus)

	logoutBtn := component.NewAsyncButton(
		i18n.T("plugin.sourcelogin.logout"),
		func() {
			err := loginables[providerChoice.Selected].(miaosic.Loginable).Logout()
			if err != nil {
				_ = global.EventBus.Publish(events.ErrorUpdate,
					events.ErrorUpdateEvent{Error: err})
				return
			}
			currentUser.SetText(i18n.T("plugin.sourcelogin.current_user.notlogin"))
			w.sessions[providerChoice.Selected] = ""
		},
	)
	qrcodeImg := canvas.NewImageFromResource(resource.ImageEmptyQrCode)
	qrcodeImg.SetMinSize(fyne.NewSize(200, 200))
	qrcodeImg.FillMode = canvas.ImageFillContain
	var currentLoginSession *miaosic.QrLoginSession
	//var key string
	qrStatus := widget.NewLabel("AAAAAAAA")
	qrStatus.SetText("")
	newQrBtn := component.NewAsyncButton(
		i18n.T("plugin.sourcelogin.qr.new"),
		func() {
			var err error
			if providerChoice.Selected == "" {
				return
			}
			qrStatus.SetText("")
			w.log.Info("getting a new qr code for login")
			pvdr, _ := miaosic.GetProvider(providerChoice.Selected)
			provider := pvdr.(miaosic.Loginable)
			currentLoginSession, err = provider.QrLogin()
			if err != nil {
				_ = global.EventBus.Publish(events.ErrorUpdate,
					events.ErrorUpdateEvent{Error: err})
				return
			}
			w.log.Debugf("trying encode url %s to qrcode", currentLoginSession.Url)
			data, err := qrcode.Encode(currentLoginSession.Url, qrcode.Medium, 256)
			if err != nil {
				_ = global.EventBus.Publish(events.ErrorUpdate,
					events.ErrorUpdateEvent{Error: err})
				return
			}
			//w.log.Debug("create img from raw data")
			pic := canvas.NewImageFromReader(bytes.NewReader(data), "qrcode")
			qrcodeImg.Resource = pic.Resource
			qrcodeImg.Refresh()
		},
	)
	finishQrBtn := component.NewAsyncButton(
		i18n.T("plugin.sourcelogin.qr.finish"),
		func() {
			if currentLoginSession == nil {
				return
			}
			currentProvider := providerChoice.Selected
			if currentProvider == "" {
				return
			}
			pvdr, _ := miaosic.GetProvider(currentProvider)
			provider := pvdr.(miaosic.Loginable)
			w.log.Info("checking qr status")
			result, err := provider.QrLoginVerify(currentLoginSession)
			if err != nil {
				_ = global.EventBus.Publish(events.ErrorUpdate,
					events.ErrorUpdateEvent{Error: err})
				return
			}
			qrStatus.SetText(result.Message)
			if result.Success {
				currentLoginSession = nil
				qrcodeImg.Resource = resource.ImageEmptyQrCode
				qrcodeImg.Refresh()
				providerChoice.OnChanged(currentProvider)
				w.sessions[currentProvider] = provider.SaveSession()
			}
		},
	)
	controlBox := container.NewHBox(newQrBtn, finishQrBtn, logoutBtn)
	qrImagePanel := container.NewCenter(
		container.NewVBox(qrcodeImg, qrStatus),
	)
	w.panel = container.NewVBox(sourcePanel, controlBox, qrImagePanel)
	return w.panel
}
