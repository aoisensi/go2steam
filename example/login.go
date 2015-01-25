package main

import (
	"fmt"
	"io/ioutil"
	"os"

	go2steam "./.."
)

const (
	confCookie = "cookie.conf"
)

func main() {
	steam := go2steam.NewSteam()

	cookie, err := os.Open(confCookie)
	if err == nil {
		data, _ := ioutil.ReadAll(cookie)
		steam.SetCookies(data)
		defer cookie.Close()
	} else {
		fmt.Println(err)
	}

	defer func() {
		f, _ := os.Create(confCookie)
		f.Write(steam.Cookies())
		defer f.Close()
	}()

	var user, pass string
	fmt.Print("Username > ")
	fmt.Scanln(&user)
	fmt.Print("Password > ")
	fmt.Scanln(&pass)

	var errl error
	opt := go2steam.LoginOption{}
	for {
		switch errt := errl.(type) {
		case *go2steam.ErrorLoginCapchaAuth:
			fmt.Println("Please open url and input Captcha.")
			cap := errt.Captcha()
			fmt.Println(cap.GetURL())
			fmt.Print("Answer > ")
			var ans string
			fmt.Scanln(&ans)
			cap.SetAnswer(ans)
			opt.SetCaptcha(cap)
		case *go2steam.ErrorLoginEMailAuth:
			fmt.Printf("Valve sent email to \"%s\" domain.\n", errt.Domain)
			var code, device string
			fmt.Print("Special-Access-Code > ")
			fmt.Scanln(&code)
			fmt.Print("Device Name > ")
			fmt.Scanln(&device)
			opt.SetGuard(code, device)
		default:
			if errl != nil {
				panic(errl)
			}
		}
		errl = steam.Login(user, pass, opt)
		if errl == nil {
			break
		}
	}

}
