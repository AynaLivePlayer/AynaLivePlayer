package wylogin

import (
	"AynaLivePlayer/adapters/provider"
	"AynaLivePlayer/core/adapter"
	"AynaLivePlayer/gui"
	"AynaLivePlayer/gui/component"
	"AynaLivePlayer/pkg/config"
	"AynaLivePlayer/pkg/i18n"
	"AynaLivePlayer/resource"
	"bytes"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/skip2/go-qrcode"
	"net/http"
)

const MODULE_PLGUIN_NETEASELOGIN = "plugin.neteaselogin"

type WYLogin struct {
	config.BaseConfig
	MusicU string
	CSRF   string
	panel  fyne.CanvasObject
	cb     adapter.IControlBridge
	log    adapter.ILogger
}

func NewWYLogin(cb adapter.IControlBridge) *WYLogin {
	return &WYLogin{
		MusicU: "MUSIC_U=;",
		CSRF:   "__csrf=;",
		cb:     cb,
		log:    cb.Logger().WithModule(MODULE_PLGUIN_NETEASELOGIN),
	}
}

func (w *WYLogin) Name() string {
	return "NeteaseLogin"
}

func (w *WYLogin) Enable() error {
	config.LoadConfig(w)
	w.loadCookie()
	gui.AddConfigLayout(w)
	go func() {
		w.log.Info("updating netease status")
		provider.NeteaseAPI.UpdateStatus()
		w.log.Info("finish updating netease status")
	}()
	return nil
}

func (w *WYLogin) Disable() error {
	w.saveCookie()
	return nil
}

func (w *WYLogin) loadCookie() {
	provider.NeteaseAPI.ReqData.Cookies = (&http.Response{
		Header: map[string][]string{
			"Set-Cookie": []string{w.MusicU, w.CSRF},
		},
	}).Cookies()
}

func (w *WYLogin) saveCookie() {
	for _, c := range provider.NeteaseAPI.ReqData.Cookies {
		if c.Name == "MUSIC_U" {
			w.MusicU = c.String()
		}
		if c.Name == "__csrf" {
			w.CSRF = c.String()
		}
	}
}

func (w *WYLogin) Title() string {
	return i18n.T("plugin.neteaselogin.title")
}

func (w *WYLogin) Description() string {
	return i18n.T("plugin.neteaselogin.description")
}

func (w *WYLogin) CreatePanel() fyne.CanvasObject {
	if w.panel != nil {
		return w.panel
	}
	currentUser := widget.NewLabel(i18n.T("plugin.neteaselogin.current_user.notlogin"))
	currentStatus := container.NewHBox(
		widget.NewLabel(i18n.T("plugin.neteaselogin.current_user")),
		currentUser)

	refreshBtn := component.NewAsyncButton(
		i18n.T("plugin.neteaselogin.refresh"),
		func() {
			provider.NeteaseAPI.UpdateStatus()
			if provider.NeteaseAPI.IsLogin() {
				currentUser.SetText(provider.NeteaseAPI.Nickname())
			} else {
				currentUser.SetText(i18n.T("plugin.neteaselogin.current_user.notlogin"))
			}

		},
	)
	logoutBtn := component.NewAsyncButton(
		i18n.T("plugin.neteaselogin.logout"),
		func() {
			provider.NeteaseAPI.Logout()
			currentUser.SetText(i18n.T("plugin.neteaselogin.current_user.notlogin"))
		},
	)
	controlBtns := container.NewHBox(refreshBtn, logoutBtn)
	qrcodeImg := canvas.NewImageFromResource(resource.ImageEmpty)
	qrcodeImg.SetMinSize(fyne.NewSize(200, 200))
	qrcodeImg.FillMode = canvas.ImageFillContain
	var key string
	qrStatus := widget.NewLabel("AAAAAAAA")
	qrStatus.SetText("")
	newQrBtn := component.NewAsyncButton(
		i18n.T("plugin.neteaselogin.qr.new"),
		func() {
			qrStatus.SetText("")
			w.log.Info("getting a new qr code for login")
			key = provider.NeteaseAPI.GetQrLoginKey()
			if key == "" {
				w.log.Warn("fail to get qr code key")
				return
			}
			w.log.Debugf("trying encode url %s to qrcode", provider.NeteaseAPI.GetQrLoginUrl(key))
			data, err := qrcode.Encode(provider.NeteaseAPI.GetQrLoginUrl(key), qrcode.Medium, 256)
			if err != nil {
				w.log.Warnf("generate qr code failed: %s", err)
				return
			}
			w.log.Debug("create img from raw data")
			pic := canvas.NewImageFromReader(bytes.NewReader(data), "qrcode")
			qrcodeImg.Resource = pic.Resource
			qrcodeImg.Refresh()
		},
	)
	finishQrBtn := component.NewAsyncButton(
		i18n.T("plugin.neteaselogin.qr.finish"),
		func() {
			if key == "" {
				return
			}
			w.log.Info("checking qr status")
			ok, msg := provider.NeteaseAPI.CheckQrLogin(key)
			qrStatus.SetText(msg)
			if ok {
				key = ""
				qrcodeImg.Resource = resource.ImageEmpty
				qrcodeImg.Refresh()
			}
		},
	)
	loginPanel := container.NewCenter(
		container.NewVBox(
			qrcodeImg,
			container.NewHBox(newQrBtn, finishQrBtn, qrStatus),
		),
	)
	w.panel = container.NewVBox(controlBtns, currentStatus, loginPanel)
	return w.panel
}
