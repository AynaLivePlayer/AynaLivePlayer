package provider

import (
	"fmt"
	neteaseApi "github.com/XiaoMengXinX/Music163Api-Go/api"
	"net/http"
)

// Netease other method

func (n *Netease) UpdateStatus() {
	status, _ := neteaseApi.GetLoginStatus(n.ReqData)
	n.loginStatus = status
}

// IsLogin check if current cookie is a login user
func (n *Netease) IsLogin() bool {
	return n.loginStatus.Profile.UserId != 0
}

func (n *Netease) Nickname() string {
	return n.loginStatus.Profile.Nickname
}

func (n *Netease) GetQrLoginKey() string {
	unikey, err := neteaseApi.GetQrUnikey(n.ReqData)
	if err != nil {
		return ""
	}
	return unikey.Unikey
}

func (n *Netease) GetQrLoginUrl(key string) string {
	return fmt.Sprintf("https://music.163.com/login?codekey=%s", key)
}

func (n *Netease) CheckQrLogin(key string) (bool, string) {
	login, h, err := neteaseApi.CheckQrLogin(n.ReqData, key)
	if err != nil {
		return false, ""
	}
	// if login.Code == 800 || login.Code == 803. login success
	if login.Code != 800 && login.Code != 803 {
		return false, login.Message
	}
	cookies := make([]*http.Cookie, 0)
	for _, c := range (&http.Response{Header: h}).Cookies() {
		if c.Name == "MUSIC_U" || c.Name == "__csrf" {
			cookies = append(cookies, c)
		}
	}
	n.ReqData.Cookies = cookies
	return true, login.Message
}

func (n *Netease) Logout() {
	n.ReqData.Cookies = []*http.Cookie{
		{Name: "MUSIC_U", Value: ""},
		{Name: "__csrf", Value: ""},
	}
	return
}
