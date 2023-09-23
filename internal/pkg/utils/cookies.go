package utils

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/playwright-community/playwright-go"
	"github.com/pterm/pterm"
)

var DEFAULT_TIMEOUT float64 = 60_000

// extracts cookies from canvas url e.g. http://canvas.nus.edu.sg
func ExtractAndStoreCanvasCookies(canvasUrl string, cookiesFilePath string) []*http.Cookie {
	parsedCanvasUrl, err := url.Parse(canvasUrl)
	if err != nil {
		pterm.Error.Printfln("Error parsing canvas url %s: %s", canvasUrl, err.Error())
		os.Exit(1)
	}

	loginUrl := url.URL{
		Scheme: parsedCanvasUrl.Scheme,
		Host:   parsedCanvasUrl.Host,
		Path:   "/login/saml/105",
	}

	var rawCookies []playwright.Cookie
	if len(rawCookies) == 0 {
		pw, err := playwright.Run()
		if err != nil {
			pterm.Error.Printfln("Failed to launch playwright, launching browser...")
			os.Exit(1)
		}
		browser, err := pw.Chromium.Launch()
		if err != nil {
			pterm.Error.Printfln("Error launching playwright browser: %s", err.Error())
			os.Exit(1)
		}
		page, err := browser.NewPage()
		if err != nil {
			pterm.Error.Printfln("Error launching playwright page: %s", err.Error())
			os.Exit(1)
		}
		if _, err = page.Goto(loginUrl.String(), playwright.PageGotoOptions{Timeout: &DEFAULT_TIMEOUT}); err != nil {
			pterm.Error.Printfln("Error navigating to login url via playright: %s", err.Error())
			os.Exit(1)
		}
		rawUsername, err := pterm.DefaultInteractiveTextInput.Show("Please enter your canvas username")
		if err != nil {
			pterm.Error.Printfln("Error getting username input: %s", err.Error())
			os.Exit(1)
		}
		rawPassword, err := pterm.DefaultInteractiveTextInput.Show("Please enter your canvas password")
		if err != nil {
			pterm.Error.Printfln("Error getting password input: %s", err.Error())
			os.Exit(1)
		}
		usernameInputLoc := page.Locator("#userNameInput")
		if err := usernameInputLoc.Fill(rawUsername); err != nil {
			pterm.Error.Printfln("Error entering username on login page: %s", err.Error())
			os.Exit(1)
		}
		usernameValue, err := usernameInputLoc.InputValue()
		if err != nil {
			pterm.Error.Printfln("Error getting username input value: %s", err.Error())
			os.Exit(1)
		} else if usernameValue == "" {
			pterm.Error.Printfln("Invalid password")
			os.Exit(1)
		}
		passwordInputLoc := page.Locator("#passwordInput")
		if err := passwordInputLoc.Fill(rawPassword); err != nil {
			pterm.Error.Printfln("Error entering password on login page: %s", err.Error())
			os.Exit(1)
		}
		passwordValue, err := passwordInputLoc.InputValue()
		if err != nil {
			pterm.Error.Printfln("Error getting password input value: %s", err.Error())
			os.Exit(1)
		} else if passwordValue == "" {
			pterm.Error.Printfln("Invalid password")
			os.Exit(1)
		}
		if err := passwordInputLoc.Click(); err != nil {
			pterm.Error.Printfln("Error clicking password input: %s", err.Error())
			os.Exit(1)
		}
		if err := page.Keyboard().Down("Enter"); err != nil {
			pterm.Error.Printfln("Error pressing enter: %s", err.Error())
			os.Exit(1)
		}
		loginSuccessUrl := url.URL{
			Scheme: parsedCanvasUrl.Scheme,
			Host:   parsedCanvasUrl.Host,
			Path:   parsedCanvasUrl.Path + "/",
			RawQuery: url.Values{
				"login_success": {"1"},
			}.Encode(),
		}
		if err := page.WaitForURL(loginSuccessUrl.String()); err != nil {
			pterm.Error.Printfln("Failed to log in: %s", err.Error())
			os.Exit(1)
		}
		if page.URL() == loginSuccessUrl.String() {
			pterm.Success.Printfln("Successfully logged in!")
		} else {
			pterm.Error.Printfln("Failed to log in!")
			os.Exit(1)
		}
		rawCookies, err = page.Context().Cookies()
		if err != nil {
			pterm.Error.Printfln("Error extracting cookies after successful login: %s", err.Error())
			os.Exit(1)
		}
	}

	if err := StoreCookiesToFile(rawCookies, cookiesFilePath); err != nil {
		pterm.Error.Printfln("Error storing browser cookies: %s", err.Error())
		os.Exit(1)
	}

	var httpCookies []*http.Cookie
	for _, c := range rawCookies {
		httpCookies = append(httpCookies, &http.Cookie{Name: c.Name, Value: c.Value, HttpOnly: c.HttpOnly, Secure: c.Secure, Domain: c.Domain, Path: c.Path})
	}
	return httpCookies
}

func StoreCookiesToFile(rawCookies []playwright.Cookie, cookiesFilePath string) error {
	cookiesDir := filepath.Dir(cookiesFilePath)
	if err := os.MkdirAll(cookiesDir, 0755); err != nil {
		return err
	}
	cookiesJson, err := json.Marshal(rawCookies)
	if err != nil {
		return err
	}
	if err := os.WriteFile(cookiesFilePath, cookiesJson, 0755); err != nil {
		return err
	}
	return nil
}

func ExtractCookiesFromFile(cookiesFilePath string) ([]playwright.Cookie, error) {
	b, err := os.ReadFile(cookiesFilePath)
	if err != nil {
		return nil, err
	}
	var rawCookies []playwright.Cookie
	if err := json.Unmarshal(b, &rawCookies); err != nil {
		return nil, err
	}
	return rawCookies, nil
}
