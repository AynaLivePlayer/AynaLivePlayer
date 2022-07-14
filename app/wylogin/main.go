package main

import (
	"fmt"
	neteaseApi "github.com/XiaoMengXinX/Music163Api-Go/api"
	"github.com/skip2/go-qrcode"
	"net/http"

	//neteaseTypes "github.com/XiaoMengXinX/Music163Api-Go/types"
	neteaseUtil "github.com/XiaoMengXinX/Music163Api-Go/utils"
)

var reqData = neteaseUtil.RequestData{
	Cookies: []*http.Cookie{},
	Headers: neteaseUtil.Headers{
		{
			"X-Real-IP",
			"118.88.88.88",
		},
	},
}

func IsLogin() {

}

func main() {
	status, err := neteaseApi.GetLoginStatus(reqData)
	if err != nil {
		return
	}
	fmt.Println(status.Profile.UserId)
	fmt.Println(status.Account.Id)
	unikey, err := neteaseApi.GetQrUnikey(reqData)
	if err != nil {
		return
	}
	fmt.Println(unikey.Unikey)
	qrcode.WriteFile(fmt.Sprintf("https://music.163.com/login?codekey=%s", unikey.Unikey), qrcode.Medium, 256, "qrcode.png")
	//wait user input
	var input string
	_, _ = fmt.Scanln(&input)
	login, h, err := neteaseApi.CheckQrLogin(reqData, unikey.Unikey)
	if err != nil {
		return
	}
	for _, c := range (&http.Response{Header: h}).Cookies() {
		fmt.Println(c)
	}
	fmt.Println(login.Nickname, login.Message, login.Code)
	reqData.Cookies = (&http.Response{Header: h}).Cookies()
	status, err = neteaseApi.GetLoginStatus(reqData)
	if err != nil {
		return
	}
	fmt.Println(status.Profile.UserId)
	fmt.Println(status.Account.Id)
}
