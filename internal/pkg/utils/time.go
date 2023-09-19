package utils

import "time"

const (
	OUTPUT_FORMAT = "2006-01-02T15:04:05.000Z"
)

func TimestampToJavaScriptISO(t time.Time) string {
	return t.UTC().Format(OUTPUT_FORMAT)
}
