package canvas

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/playwright-community/playwright-go"
	"github.com/pterm/pterm"
)

func LoginToCanvas(page playwright.Page, canvasUrl *url.URL) (playwright.Page, error) {
	loginUrl := url.URL{
		Host:   canvasUrl.Host,
		Scheme: canvasUrl.Scheme,
		Path:   "/login/saml/105",
	}
	loginSuccessUrl := url.URL{
		Host:   canvasUrl.Host,
		Scheme: canvasUrl.Scheme,
		Path:   "/",
		RawQuery: url.Values{
			"login_success": {"1"},
		}.Encode(),
	}

	if _, err := page.Goto(loginUrl.String()); err != nil {
		return nil, fmt.Errorf("failed to navigate to login url %s: %s", loginUrl.String(), err.Error())
	}

	// login if not logged in yet (prompt for credentials)
	if strings.Contains(page.URL(), "https://vafs.u.nus.edu/adfs/ls/?SAMLRequest=") {
		rawUsername, err := pterm.DefaultInteractiveTextInput.Show("Please enter your canvas username")
		if err != nil {
			return nil, fmt.Errorf("failed to get username input: %s", err.Error())
		}
		if rawUsername == "" {
			return nil, fmt.Errorf("username cannot be empty")
		}
		rawPassword, err := pterm.DefaultInteractiveTextInput.WithMask("*").Show("Please enter your canvas password")
		if err != nil {
			return nil, fmt.Errorf("failed to get password input: %s", err.Error())
		}
		if rawPassword == "" {
			return nil, fmt.Errorf("password cannot be empty")
		}
		usernameInputLoc := page.Locator("#userNameInput")
		if err := usernameInputLoc.Fill(rawUsername); err != nil {
			return nil, fmt.Errorf("failed to enter username on login page: %s", err.Error())
		}
		usernameValue, err := usernameInputLoc.InputValue()
		if err != nil {
			return nil, fmt.Errorf("failed to get username input value: %s", err.Error())
		} else if usernameValue == "" {
			return nil, fmt.Errorf("failed to fill username into form")
		}
		passwordInputLoc := page.Locator("#passwordInput")
		if err := passwordInputLoc.Fill(rawPassword); err != nil {
			pterm.Error.Printfln("Error entering password on login page: %s", err.Error())
			os.Exit(1)
		}
		passwordValue, err := passwordInputLoc.InputValue()
		if err != nil {
			return nil, fmt.Errorf("failed to get password input value: %s", err.Error())
		} else if passwordValue == "" {
			return nil, fmt.Errorf("failed to fill password into form")
		}
		if err := passwordInputLoc.Click(); err != nil {
			return nil, fmt.Errorf("failed to get click into password input: %s", err.Error())
		}
		if err := page.Keyboard().Down("Enter"); err != nil {
			return nil, fmt.Errorf("failed to sign in: %s", err.Error())
		}
		if err := page.WaitForURL(loginSuccessUrl.String()); err != nil {
			if page.URL() != loginSuccessUrl.String() {
				return nil, fmt.Errorf("login failed: current page is %s", page.URL())
			}
		}
	}
	return page, nil
}
