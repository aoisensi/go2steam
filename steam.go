package steam

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

type steam struct {
	Steam
	service     http.Client
	loginValues url.Values
}

type Steam interface {
	Login(string, string, LoginOption) error
	Cookies() []byte
	SetCookies([]byte)

	LoadGroupFromCustom(string) (Group, error)
}

func NewSteam() Steam {
	steam := new(steam)
	steam.service = http.Client{}
	jar, _ := cookiejar.New(nil)
	steam.service.Jar = jar
	steam.loginValues = url.Values{}
	return steam
}

func init() {
	initCookie()
}
