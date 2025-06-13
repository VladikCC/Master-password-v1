package main

import (
	"ILY/Data"
	"fmt"
)

func CheckAuth() bool {
	exists, err := Data.UserExistsByID(1)
	if err != nil {
		fmt.Printf("Ошибка проверки авторизации: %v\n", err)
		return false
	}
	return exists
}

func AuthFlow() {
	if CheckAuth() {
		login()
	} else {
		Logoreg()
	}
}
