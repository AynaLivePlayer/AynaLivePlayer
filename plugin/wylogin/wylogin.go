package wylogin

import (
	"AynaLivePlayer/config"
	"AynaLivePlayer/gui"
	"AynaLivePlayer/i18n"
	"AynaLivePlayer/logger"
	"AynaLivePlayer/provider"
	"bytes"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	qrcode "github.com/skip2/go-qrcode"
	"net/http"
)

const MODULE_PLGUIN_NETEASELOGIN = "plugin.neteaselogin"

var lg = logger.Logger.WithField("Module", MODULE_PLGUIN_NETEASELOGIN)

type WYLogin struct {
	MusicU string
	CSRF   string
	panel  fyne.CanvasObject
}

func NewWYLogin() *WYLogin {
	return &WYLogin{
		MusicU: "MUSIC_U=;",
		CSRF:   "__csrf=;",
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
		lg.Info("updating netease status")
		provider.NeteaseAPI.UpdateStatus()
		lg.Info("finish updating netease status")
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

	refreshBtn := gui.NewAsyncButton(
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
	logoutBtn := gui.NewAsyncButton(
		i18n.T("plugin.neteaselogin.logout"),
		func() {
			provider.NeteaseAPI.Logout()
			currentUser.SetText(i18n.T("plugin.neteaselogin.current_user.notlogin"))
		},
	)
	controlBtns := container.NewHBox(refreshBtn, logoutBtn)
	qrcodeImg := canvas.NewImageFromResource(gui.ResEmptyImage)
	qrcodeImg.SetMinSize(fyne.NewSize(200, 200))
	qrcodeImg.FillMode = canvas.ImageFillContain
	var key string
	qrStatus := widget.NewLabel("AAAAAAAA")
	qrStatus.SetText("")
	newQrBtn := gui.NewAsyncButton(
		i18n.T("plugin.neteaselogin.qr.new"),
		func() {
			qrStatus.SetText("")
			lg.Info("getting a new qr code for login")
			key = provider.NeteaseAPI.GetQrLoginKey()
			if key == "" {
				lg.Warn("fail to get qr code key")
				return
			}
			lg.Debugf("trying encode url %s to qrcode", provider.NeteaseAPI.GetQrLoginUrl(key))
			data, err := qrcode.Encode(provider.NeteaseAPI.GetQrLoginUrl(key), qrcode.Medium, 256)
			if err != nil {
				lg.Warnf("generate qr code failed: %s", err)
				return
			}
			lg.Debug("create img from raw data")
			pic := canvas.NewImageFromReader(bytes.NewReader(data), "qrcode")
			qrcodeImg.Resource = pic.Resource
			qrcodeImg.Refresh()
		},
	)
	finishQrBtn := gui.NewAsyncButton(
		i18n.T("plugin.neteaselogin.qr.finish"),
		func() {
			if key == "" {
				return
			}
			lg.Info("checking qr status")
			ok, msg := provider.NeteaseAPI.CheckQrLogin(key)
			qrStatus.SetText(msg)
			if ok {
				key = ""
				qrcodeImg.Resource = gui.ResEmptyImage
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
