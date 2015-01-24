package main

import (
	"fmt"

	"github.com/k0kubun/pp"

	go2steam ".."
)

func main() {
	steam := go2steam.NewSteam()
	var user, pass string

	var errl error
	for {
		if errl != nil {
			fmt.Println(errl)
		}
		switch errt := errl.(type) {
		case *go2steam.ErrorLoginCapchaAuth:
			fmt.Println("Please open url and input Captcha.")
			cap := errt.Captcha()
			fmt.Println(cap.GetURL())
			fmt.Print("Answer > ")
			var ans string
			fmt.Scanln(&ans)
			cap.SetAnswer(ans)
			errl = steam.LoginCaptcha(user, pass, cap)
		case *go2steam.ErrorLoginEMailAuth:
			fmt.Printf("Valve sent email to \"%s\" domain.\n", errt.Domain)
			var code, device string
			fmt.Print("Special-Access-Code > ")
			fmt.Scanln(&code)
			fmt.Print("Device Name > ")
			fmt.Scanln(&device)
			errl = steam.LoginGuard(user, pass, code, device)
		default:
			fmt.Print("Username > ")
			fmt.Scanln(&user)
			fmt.Print("Password > ")
			fmt.Scanln(&pass)
			errl = steam.Login(user, pass)
		}

		if errl == nil {
			break
		}

	}
	pp.Println(steam.Cookie())
}
