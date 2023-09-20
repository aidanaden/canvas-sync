package utils

import "time"

const (
	OUTPUT_FORMAT = "2006-01-02T15:04:05.000Z"
)

func TimestampToJavaScriptISO(t time.Time) string {
	return t.UTC().Format(OUTPUT_FORMAT)
}

func SGTfromUTC(t time.Time) (*time.Time, error) {
	location, err := time.LoadLocation("Asia/Singapore")
	if err != nil {
		return nil, err
	}
	local := t.In(location)
	return &local, nil
}

func FormatEventDate(t time.Time) string {
	return t.Format("Mon, 02 Jan 06 03:04 PM")
}
