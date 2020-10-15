package crypt

import (
	"errors"
	"fmt"

	"github.com/howeyc/gopass"
)

func ScanPassword(showTip string) (string, error) {

	fmt.Println(showTip)

	password, err := gopass.GetPasswd()

	if err != nil {
		return "", err
	}

	if string(password) == "" {
		return "", errors.New(showTip + " empty")
	}

	return string(password), nil
}
