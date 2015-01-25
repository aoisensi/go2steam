package steam

import (
	"encoding/json"
	"net/http"
	"net/url"
)

var (
	cookieURL *url.URL
)

func (s *steam) Cookies() []byte {
	data := s.service.Jar.Cookies(cookieURL)
	result, _ := json.Marshal(data)
	return result
}

func (s *steam) SetCookies(data []byte) {
	var cookies []*http.Cookie
	json.Unmarshal(data, &cookies)
	s.service.Jar.SetCookies(cookieURL, cookies)
}

func (s *steam) sessionId() (string, error) {
	cs := s.service.Jar.Cookies(cookieURL)
	for _, c := range cs {
		if c.Name == "sessionid" {
			return url.QueryUnescape(c.Value)
		}
	}
	return "", nil
}

func initCookie() {
	cookieURL, _ = url.Parse("https://steamcommunity.com/")
}
