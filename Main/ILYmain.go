package main

import (
	"ILY/Data"
	"bufio"
	"fmt"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func Clear() {
	fmt.Print("\033c")
}

func Vod(target *string) {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	*target = strings.TrimSpace(input)
}

func To(x int, y int) {
	fmt.Printf("\033[%d;%dH", y, x)
}

func Logoreg() {
	var void string

	Clear()

	for {
		To(28, 4)
		fmt.Print("Добро пожаловать в ILY")
		To(8, 5)
		fmt.Print("!Просьба прочитать файл (README) созданный в главной папке main!")
		To(27, 7)
		fmt.Print("╔══════════════════════╗")
		To(27, 8)
		fmt.Print("║ Нажмите любую кнопку ║")
		To(27, 9)
		fmt.Print("╚══════════════════════╝")

		To(5, 14)
		fmt.Print("╔════════════════════════╗")
		To(5, 15)
		fmt.Print("║    Для Регистрации     ║")
		To(5, 16)
		fmt.Print("╚════════════════════════╝")

		To(32, 15)
		fmt.Print("-> ")
		Vod(&void)
		registr()
	}
}

func main() {
	createREADME()

	if err := Data.InitDB(); err != nil {
		fmt.Printf("FATAL: %v\n", err)
		return
	}

	fmt.Println("\n[Проверка подключения к БД]")
	checkDatabaseTables()
	rows, err := Data.DB.Query("SELECT name FROM sqlite_master WHERE type='table'")
	if err != nil {
		fmt.Printf("Ошибка запроса: %v\n", err)
	} else {
		defer rows.Close()
		fmt.Println("Таблицы в базе данных:")
		for rows.Next() {
			var name string
			rows.Scan(&name)
			fmt.Println("-", name)
		}
	}

	AuthFlow()
}

func createREADME() {
	file, err := os.Create("(!README.RUS!)")
	if err != nil {
		fmt.Println("Unable to create file:", err)
		os.Exit(1)
	}
	defer file.Close()
	file.WriteString(" !Добро пожаловать! Это подробная инструкция и информация для пользователя! \n")
	file.WriteString(" [1] О данном ПО? \n")
	file.WriteString("->Данное ПО создавалось для вашего личного использования, а именно \n")
	file.WriteString("->Сервер хранения данных развертывается прямо на вашем ПК на базе данных типа SQLite. \n")
	file.WriteString("->(!А ТАКЖЕ ИСХОДНЫЙ КОД В ОТКРЫТОМ ДОСТУПЕ!) \n")
	file.WriteString("->Это значит, что никто не следит за вашими данными и не имеет к ним доступа, кроме вас. \n")
	file.WriteString("->Также ваши пароли, которые вы будете вводить, автоматически будут добавляться сюда \n")
	file.WriteString("->И при этом надежно шифроваться, так что все ваши пароли вы сможете сохранить под надежным сейфом данных. \n")
	file.WriteString("->И если хакеры или кто-то другой попытается получить с вашего ПК доступ к вашей базе данных, где всё хранится, \n")
	file.WriteString("->Кроме вашего email, который нужен для двухфакторной аутентификации, они ничего не увидят, так как пароли в базе данных шифруются. \n")
	file.WriteString(" [2] Для чего создавалось это ПО? \n")
	file.WriteString("->Данное ПО создано в целях вашей собственной безопасности, а именно: \n")
	file.WriteString("->Когда вы регистрируетесь где-либо (например, возьмём всем известный Telegram), \n")
	file.WriteString("->Да, там есть двухфакторная аутентификация, как и в нашем ПО, но никакая двухфакторная аутентификация не обеспечит полную безопасность. \n")
	file.WriteString("->Если вы при регистрации в Telegram ввели, к примеру, пароль: ваш номер, дату рождения, имя собаки, вашего родного города или девушки, \n")
	file.WriteString("->Такой пароль можно взломать за 3-10 минут, используя вредоносное ПО (метод подбора), или даже просто зная один из ваших других паролей. \n")
	file.WriteString("->Так как многие люди особо не любят менять пароли в разных приложениях, ваша защита будет ничтожна, если человек знает основы языка C. \n")
	file.WriteString("->Ему не составит проблем написать ПО для слежки или просто использовать базовый метод подбора с обходом защиты. \n")
	file.WriteString("->И что вы можете сделать??? \n")
	file.WriteString("->Использовать НАШЕ ПРИЛОЖЕНИЕ (ПО)! А именно: при регистрации где-либо (допустим, в том же Telegram) активировать наше ПО, \n")
	file.WriteString("->Авторизоваться, введя email, и самое ГЛАВНОЕ! \n")
	file.WriteString("->!ПРИ ВВОДЕ ПАРОЛЯ ДЛЯ НАШЕГО ПО ВВОДИТЕ ЕГО МАКСИМАЛЬНО СЛОЖНЫМ И ЗАПИШИТЕ НА ЛИСТИК ИЛИ ЗАПОМНИТЕ! \n")
	file.WriteString("->Например, это может быть что-то вроде: 0!456900101LubaTo778 или то, что связано с вами: Iloveyou2007Nastya2008Muhamed \n")
	file.WriteString("->Это очень важный момент, так как мы обеспечиваем защиту ваших данных и удобство, но если кто-то узнает пароль от данного ПО, \n")
	file.WriteString("->То при его вводе вы можете потерять все данные. Да, это кажется ненадежным, но мы обеспечиваем двухфакторную аутентификацию и поддержку других паролей. \n")
	file.WriteString("->Искренне просим записать пароль и никому его не сообщать. В принципе, наша защита такая же, как на криптобиржах. \n")
	file.WriteString("->Так что после того, как вы ввели пароль и зарегистрировались, можете со спокойной душой зайти в ваш Telegram или Google-аккаунт, \n")
	file.WriteString("->Ввести любой пароль (например, 98274324twu7jcvns98-1fja8), и программа сама запомнит его в нашем ПО. \n")
	file.WriteString("->Люди, которые захотят взломать ваш пароль, точно не смогут его подобрать или использовать какие-либо другие методы. \n")
	file.WriteString("->И если вам нужно будет ввести пароль, просто зайдите в наше ПО, введите мастер-пароль (чтобы подтвердить, что это вы), скопируйте нужный пароль и вставьте его. \n")
	file.WriteString("->Это работает не только для Telegram, но и для всех других площадок и приложений, которые вам нужны. \n")
	file.WriteString("\n")
	file.WriteString(" [Разработчик Дилигул В.А.] \n")
	file.WriteString(" [Если хотите помочь разрабочику - Тбанк карта мир 2200 7017 0262 4278]")
}

func checkDatabaseTables() {
	rows, err := Data.DB.Query("SELECT name FROM sqlite_master WHERE type='table'")
	if err != nil {
		fmt.Printf("Ошибка проверки таблиц: %v\n", err)
		return
	}
	defer rows.Close()

	fmt.Println("Существующие таблицы:")
	for rows.Next() {
		var tableName string
		rows.Scan(&tableName)
		fmt.Println("-", tableName)
	}
}
