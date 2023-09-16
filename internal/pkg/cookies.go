package pkg

import (
	"net/http"

	"github.com/zellyn/kooky"
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
