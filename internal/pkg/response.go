package pkg

import (
	"io"
	"net/http"
)

func ExtractResponseToString(res *http.Response) string {
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	json := string(body)
	return json
}
