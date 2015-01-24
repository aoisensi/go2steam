package steam

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

type steam struct {
	Steam
	service http.Client
	lgnv    url.Values
}

type Steam interface {
	Login() error
	SetLogin(string, string)
	SetLoginCaptcha(Captcha)
	SetLoginGuard(string, string)
	Cookies() []byte
	SetCookies([]byte)

	LoadGroupFromCustom(string) (Group, error)
}

func NewSteam() Steam {
	steam := new(steam)
	steam.service = http.Client{}
	jar, _ := cookiejar.New(nil)
	steam.service.Jar = jar
	steam.lgnv = url.Values{}
	return steam
}

func init() {
	initCookie()
}
