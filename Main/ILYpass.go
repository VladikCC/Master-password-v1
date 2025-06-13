package main

import (
	"ILY/Data"
	"fmt"
	"time"
)

func login() {
	var pass string
	Clear()

	for attempts := 3; attempts > 0; attempts-- {
		To(28, 4)
		fmt.Print("Добро пожаловать в ILY")
		To(8, 5)
		fmt.Print("!Просьба прочитать файл (README) созданный в главной папке main!")
		To(5, 11)
		fmt.Print("╔════════════════════════╗")
		To(5, 12)
		fmt.Print("║Введите пароль для входа║")
		To(5, 13)
		fmt.Print("╚════════════════════════╝")
		To(32, 12)
		fmt.Print("-> ")
		Vod(&pass)

		valid, err := Data.ValidateMasterPassword(pass)
		if err != nil {
			To(9, 10)
			fmt.Print("Ошибка проверки пароля")
			time.Sleep(2 * time.Second)
			continue
		}

		if valid {
			Clear()
			To(24, 10)
			fmt.Print("╔══════════════════════════════╗")
			To(24, 11)
			fmt.Print("║     Вход выполнен успешно    ║")
			To(24, 12)
			fmt.Print("╚══════════════════════════════╝")
			time.Sleep(2 * time.Second)
			Clear()
			menu()
			break
		} else {
			To(9, 16)
			fmt.Print("╔════════════════════════════════════════╗")
			To(9, 17)
			fmt.Printf("║  Неверный пароль! Осталось попыток: %d  ║ ", attempts-1)
			To(9, 18)
			fmt.Print("╚════════════════════════════════════════╝")
			time.Sleep(2 * time.Second)
			Clear()
		}
	}
	To(24, 9)
	fmt.Print("╔══════════════════════════════╗")
	To(24, 10)
	fmt.Print("║!Превышено количество попыток!║")
	To(24, 11)
	fmt.Print("╚══════════════════════════════╝")
	fmt.Scanln()
	time.Sleep(2 * time.Second)
}
