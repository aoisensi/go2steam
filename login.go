package steam

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
	"strconv"
)

const (
	urlLoginGetRSAKey = "https://steamcommunity.com/login/getrsakey"
	urlLoginDoLogin   = "https://steamcommunity.com/login/dologin"
)

var (
	ErrFailedLoginGetRSAKey = errors.New("Failed to get RSA Key")
)

type ErrorLoginEMailAuth struct {
	Domain  string
	SteamID string
}

func (e *ErrorLoginEMailAuth) Error() string {
	return "You need to verify SteamGuard"
}

type ErrorLoginCapchaAuth struct {
	CaptchaGID string
}

func (e *ErrorLoginCapchaAuth) Error() string {
	return "You need to verify capcha"
}

func (e *ErrorLoginCapchaAuth) Captcha() Captcha {
	return newCaptcha(e.CaptchaGID)
}

type jsonLoginGetRSAKey struct {
	Success      bool
	PublicKeyMod string `json:"publickey_mod"`
	PublicKeyExp string `json:"publickey_exp"`
	Timestamp    string
	TokenGID     string `json:"token_gid"`
}

type jsonLoginDoLogin struct {
	Success bool
	Message string

	BadCaptcha    bool   `json:"bad_captcha"`
	CaptchaNeeded bool   `json:"captcha_needed"`
	CaptchaGID    string `json:"captcha_gid"`

	EMailAuthNeeded bool `json:"emailauth_needed"`
	EMailDomain     string
	EMailSteamID    string

	TransferURL        string                 `json:"transfer_url"`
	TransferParameters map[string]interface{} `json:"transfer_parameters"`

	RequiresTwofactor bool `json:"requires_twofactor"`
}

func (s *steam) Login(username, password string) error {
	v := url.Values{
		"username": {username},
		"password": {password},
	}
	return s.login(v)
}

func (s *steam) LoginCaptcha(username, password string, captcha Captcha) error {
	v := url.Values{
		"username":     {username},
		"password":     {password},
		"captchagid":   {captcha.GetGID()},
		"captcha_text": {captcha.GetAnswer()},
	}
	return s.login(v)
}

func (s *steam) LoginGuard(username, password, code, name string) error {
	v := url.Values{
		"username":          {username},
		"password":          {password},
		"emailauth":         {code},
		"loginfriendlyname": {name},
	}
	return s.login(v)
}

func (s *steam) login(v url.Values) error {
	password, ts, err := loginGetRSA(v.Get("username"), v.Get("password"))
	if err != nil {
		return err
	}
	v.Set("password", password)
	v.Set("rsatimestamp", ts)
	v.Set("remember_login", "true")
	cookies, err := loginDoLogin(v)
	if err != nil {
		return err
	}
	s.cookies = cookies
	return nil
}

func loginDoLogin(v url.Values) ([]*http.Cookie, error) {

	resp, err := http.PostForm(urlLoginDoLogin, v)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	login := new(jsonLoginDoLogin)
	if err := json.Unmarshal(body, login); err != nil {
		return nil, err
	}
	if login.Success {
		cookie, err := login.transfer()
		if err != nil {
			return nil, err
		}

		return cookie, nil
	}
	switch {
	case login.EMailAuthNeeded:
		return nil, &ErrorLoginEMailAuth{
			Domain:  login.EMailDomain,
			SteamID: login.EMailSteamID,
		}
	case login.CaptchaNeeded:
		return nil, &ErrorLoginCapchaAuth{
			CaptchaGID: login.CaptchaGID,
		}
	default:
		return nil, errors.New(login.Message)
	}
}

func loginGetRSA(username, password string) (string, string, error) {
	u := urlLoginGetRSAKey
	resp, err := http.PostForm(u, url.Values{"username": {username}})
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}
	key := new(jsonLoginGetRSAKey)
	json.Unmarshal(body, key)

	if !key.Success {
		return "", "", ErrFailedLoginGetRSAKey
	}

	pubkey := key.getPubKey()

	res, err := rsa.EncryptPKCS1v15(rand.Reader, pubkey, []byte(password))
	if err != nil {
		return "", "", err
	}
	rp := base64.StdEncoding.EncodeToString(res)
	return rp, key.Timestamp, err
}

func (r *jsonLoginGetRSAKey) getPubKey() *rsa.PublicKey {
	mod := new(big.Int)
	modb, _ := hex.DecodeString(r.PublicKeyMod)
	mod.SetBytes(modb)
	exp, _ := strconv.ParseInt(r.PublicKeyExp, 16, 32)
	return &rsa.PublicKey{N: mod, E: int(exp)}
}

func (l *jsonLoginDoLogin) transfer() ([]*http.Cookie, error) {
	p := url.Values{}
	for k, v := range l.TransferParameters {
		switch vt := v.(type) {
		case string:
			p.Add(k, vt)
		case bool:
			if vt {
				p.Add(k, "true")
			} else {
				p.Add(k, "false")
			}
		}
	}
	resp, err := http.PostForm(l.TransferURL, p)
	if err != nil {
		panic(err)
		return nil, err
	}
	defer resp.Body.Close()
	return resp.Cookies(), nil
}
