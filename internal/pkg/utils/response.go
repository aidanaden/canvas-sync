package utils

import (
	"io"
	"net/http"
	"os"

	"github.com/pterm/pterm"
)

func SetQueryAccessToken(req *http.Request, accessToken string) {
	if accessToken == "" {
		return
	}
	q := req.URL.Query()
	q.Set("access_token", accessToken)
	req.URL.RawQuery = q.Encode()
}

func ExtractResponseToString(res *http.Response) string {
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		pterm.Error.Printfln("Error: failed to extract response: %s", err.Error())
		os.Exit(1)
	}
	json := string(body)
	return json
}
