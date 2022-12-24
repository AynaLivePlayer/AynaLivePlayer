package provider

import (
	"github.com/go-resty/resty/v2"
	"time"
)

func httpGet(url string, header map[string]string) (*resty.Response, error) {
	resp, err := resty.New().
		SetTimeout(time.Second * 3).R().
		SetHeaders(header).
		Get(url)
	return resp, err
}
func httpGetString(url string, header map[string]string) string {
	resp, err := httpGet(url, header)
	if err != nil {
		return ""
	}
	return resp.String()
}

func httpHead(url string, header map[string]string) (*resty.Response, error) {
	resp, err := resty.New().
		SetTimeout(time.Second * 3).R().
		SetHeaders(header).
		Head(url)
	return resp, err
}
