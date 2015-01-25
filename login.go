package steam

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
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

func (s *steam) Login(username, password string, opt LoginOption) error {
	password, ts, err := s.loginGetRSA(username, password)
	if err != nil {
		return err
	}
	v := url.Values{}
	if opt != nil {
		for vk, vv := range opt {
			v.Set(vk, vv)
		}
	}
	v.Set("username", username)
	v.Set("password", password)
	v.Set("rsatimestamp", ts)
	v.Set("remember_login", "true")
	_, err = s.loginDoLogin(v)
	if err != nil {
		return err
	}
	return nil
}

func (s *steam) loginDoLogin(v url.Values) ([]*http.Cookie, error) {

	resp, err := s.service.PostForm(urlLoginDoLogin, v)
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
		cookie, err := s.loginTransfer(login)
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
		return nil, fmt.Errorf("failed login %v", login)
	}
}

func (s *steam) loginGetRSA(username, password string) (string, string, error) {
	u := urlLoginGetRSAKey
	resp, err := s.service.PostForm(u, url.Values{"username": {username}})
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
	return rp, key.Timestamp, nil
}

func (r *jsonLoginGetRSAKey) getPubKey() *rsa.PublicKey {
	mod := new(big.Int)
	modb, _ := hex.DecodeString(r.PublicKeyMod)
	mod.SetBytes(modb)
	exp, _ := strconv.ParseInt(r.PublicKeyExp, 16, 32)
	return &rsa.PublicKey{N: mod, E: int(exp)}
}

func (s *steam) loginTransfer(l *jsonLoginDoLogin) ([]*http.Cookie, error) {
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
	resp, err := s.service.PostForm(l.TransferURL, p)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return resp.Cookies(), nil
}

type LoginOption map[string]string

func (o LoginOption) SetGuard(code, name string) {
	o["emailauth"] = code
	o["loginfriendlyname"] = name
}

func (o LoginOption) SetCaptcha(captcha Captcha) {
	if captcha == nil {
		return
	}
	o["captchagid"] = captcha.GetGID()
	o["captcha_text"] = captcha.GetAnswer()
}
