package steam

import "net/http"

type steam struct {
	Steam
	cookies []*http.Cookie
}

type Steam interface {
	Login(string, string) error
	LoginCaptcha(string, string, Captcha) error
	LoginGuard(string, string, string, string) error
	Cookie() []*http.Cookie
}

func NewSteam() Steam {
	return new(steam)
}

func (s *steam) Cookie() []*http.Cookie {
	return s.cookies
}
