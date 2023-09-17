package utils

import (
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/pkg/browser"
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
		panic(err)
	}
	var rawCookies []*http.Cookie
	rawCookies = ExtractCookies(url.Host)
	if rawCookies == nil || len(rawCookies) == 0 {
		if err := browser.OpenURL(rawUrl); err != nil {
			panic(err)
		}
		for {
			rawCookies = ExtractCookies(url.Host)
			if rawCookies != nil && len(rawCookies) > 0 {
				break
			}
		}
	}
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalf("Got error while creating cookie jar %s", err.Error())
	}

	jar.SetCookies(url, rawCookies)
	return jar
}
