package utils

import (
	"encoding/json"
	"errors"
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
		pterm.Error.Printfln("Failed to extract response: %s", err.Error())
		os.Exit(1)
	}
	json := string(body)
	return json
}

type GithubReleaseInfo struct {
	TagName    string `json:"tag_name"`
	CommitHash string `json:"target_commitish"`
}

func GetCavasSyncLatestVersionHash() (*GithubReleaseInfo, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/repos/aidanaden/canvas-sync/releases/latest", nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	jsonStr := ExtractResponseToString(res)
	var releaseInfo GithubReleaseInfo
	json.Unmarshal([]byte(jsonStr), &releaseInfo)
	if releaseInfo.CommitHash == "" {
		return nil, errors.New("error getting latest canvas-sync release version")
	}
	return &releaseInfo, nil
}
