package input

import (
	"fmt"
	"log"
)

func GetYesOrNoFromUser() bool {
	var res string
	_, err := fmt.Scanln(&res)
	if err != nil {
		log.Fatalf("\nError getting response from user: %s", err.Error())
	}
	if res == "y" {
		return true
	}
	return false
}
