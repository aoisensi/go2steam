package steam

import "net/url"

const (
	urlCaptcha = "https://steamcommunity.com/public/captcha.php"
)

type captcha struct {
	gid, answer string
}
type Captcha interface {
	GetURL() string
	SetAnswer(string)
	GetAnswer() string
	GetGID() string
}

func (c *captcha) GetURL() string {
	return urlCaptcha + "?" + url.Values{"gid": {c.gid}}.Encode()
}

func (c *captcha) SetAnswer(answer string) {
	c.answer = answer
}

func (c *captcha) GetAnswer() string {
	return c.answer
}
func (c *captcha) GetGID() string {
	return c.gid
}

func newCaptcha(gid string) Captcha {
	r := new(captcha)
	r.gid = gid
	return r
}
