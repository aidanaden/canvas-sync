package utils

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"

	"github.com/pkg/browser"
	"github.com/pterm/pterm"
	"github.com/zellyn/kooky"
	_ "github.com/zellyn/kooky/browser/all" // register cookie store finders!
)

func ExtractCookies(domainSuffix string) []*http.Cookie {
	var rawCookies []*http.Cookie
	if cookies := kooky.ReadCookies(kooky.DomainHasSuffix(domainSuffix)); cookies != nil {
		for _, cookie := range cookies {
			rawCookie := http.Cookie{Name: cookie.Name, Value: cookie.Value}
			rawCookies = append(rawCookies, &rawCookie)
		}
	}
	return rawCookies
}

// extracts cookies from canvas url e.g. http://canvas.nus.edu.sg
func ExtractCanvasBrowserCookies(rawUrl string) http.CookieJar {
	url, err := url.Parse(rawUrl)
	if err != nil {
		pterm.Error.Printfln("Error parsing canvas url: %s", err.Error())
		os.Exit(1)
	}
	var rawCookies []*http.Cookie
	rawCookies = ExtractCookies(url.Host)
	if len(rawCookies) == 0 {
		if err := browser.OpenURL(rawUrl); err != nil {
			pterm.Error.Printfln("Error opening url in browser: %s", err.Error())
			os.Exit(1)
		}
		for {
			rawCookies = ExtractCookies(url.Host)
			if len(rawCookies) > 0 {
				break
			}
		}
	}
	jar, err := cookiejar.New(nil)
	if err != nil {
		pterm.Error.Printfln("Error creating cookie jar: %s", err.Error())
		os.Exit(1)
	}

	jar.SetCookies(url, rawCookies)
	return jar
}
