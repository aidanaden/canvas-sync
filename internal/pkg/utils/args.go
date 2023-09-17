package utils

import "strings"

func GetCourseCodesFromArgs(rawArgs []string) []string {
	args := make([]string, 0)
	for _, raw := range rawArgs {
		if strings.Contains(raw, ",") {
			splits := strings.Split(raw, ",")
			for _, split := range splits {
				args = append(args, strings.ToLower(split))
			}
		} else {
			args = append(args, strings.ToLower(raw))
		}
	}
	return args
}
