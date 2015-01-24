package steam

import (
	"net/http"
	"net/http/cookiejar"
)

type steam struct {
	Steam
	service http.Client
}

type Steam interface {
	Login(string, string) error
	LoginCaptcha(string, string, Captcha) error
	LoginGuard(string, string, string, string) error
	Cookies() []byte
	SetCookies([]byte)

	LoadGroupFromCustom(string) (Group, error)
}

func NewSteam() Steam {
	steam := new(steam)
	steam.service = http.Client{}
	jar, _ := cookiejar.New(nil)
	steam.service.Jar = jar
	return steam
}

func init() {
	initCookie()
}
