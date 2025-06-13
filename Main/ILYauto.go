package main

import (
	"ILY/Data"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"net/smtp"
	"strings"
	"time"
	"unicode"

	"golang.org/x/crypto/argon2"
)

func init() {
	if err := Data.InitDB(); err != nil {
		panic(fmt.Sprintf("Ошибка инициализации БД: %v", err))
	}
}

var (
	lastEmailSentTime time.Time
	verificationCode  string
	currentUserEmail  string
)

func generateVerificationCode() string {
	const digits = "0123456789"
	code := make([]byte, 6)
	for i := range code {
		num, _ := rand.Int(rand.Reader, big.NewInt(10))
		code[i] = digits[num.Int64()]
	}
	return string(code)
}

func sendVerificationEmail(email, code string) error {
	from := "lutykaf@gmail.com"
	password := "dfiu gohs mnud ioxy"

	to := []string{email}
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	subject := "ILYPO - Код подтверждения"
	body := fmt.Sprintf("Ваш код подтверждения: %s\nДанный код действителен в течение 10 минут.", code)
	message := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n%s", from, strings.Join(to, ","), subject, body)

	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, []byte(message))
	return err
}

func registr() {
	Clear()
	var email, inputCode string

	for {
		To(29, 2)
		fmt.Print("╔════════════════════════╗")
		To(29, 3)
		fmt.Print("║  Введите 1 для выхода  ║")
		To(29, 4)
		fmt.Print("╚════════════════════════╝")

		To(9, 6)
		fmt.Print("╔════════════════════════╗")
		To(9, 7)
		fmt.Print("║   Введите ваш email:   ║")
		To(9, 8)
		fmt.Print("╚════════════════════════╝")
		To(36, 7)
		fmt.Print("-> ")
		Vod(&email)
		if email == "1" {
			Logoreg()
		}

		if exists, err := Data.UserExists(email); err == nil && exists {
			To(9, 13)
			fmt.Print("╔═══════════════════════════════════╗")
			To(9, 14)
			fmt.Print("║  Этот email уже зарегистрирован!  ║")
			To(9, 15)
			fmt.Print("╚═══════════════════════════════════╝")
			time.Sleep(2 * time.Second)
			Clear()
			continue
		}

		if _, created, attempts, err := Data.GetVerificationCode(email); err == nil {
			if attempts >= 3 {
				if time.Since(created) < 10*time.Minute {
					waitTime := 10*time.Minute - time.Since(created)
					To(9, 14)
					fmt.Printf(" Превышены попытки! Ждите %.0f минут ", waitTime.Minutes())
					time.Sleep(2 * time.Second)
					Clear()
					continue
				} else {
					Data.DeleteVerificationCode(email)
				}
			}
		}

		verificationCode = generateVerificationCode()
		if err := Data.SaveVerificationCode(email, verificationCode); err != nil {
			To(9, 13)
			fmt.Print("╔═══════════════════════════════════════════╗")
			To(9, 14)
			fmt.Print("║  Ошибка сохранения кода! Попробуйте снова  ║")
			To(9, 15)
			fmt.Print("╚═══════════════════════════════════════════╝")
			time.Sleep(2 * time.Second)
			Clear()
			continue
		}

		if err := sendVerificationEmail(email, verificationCode); err != nil {
			To(9, 13)
			fmt.Print("╔═══════════════════════════════════════════╗")
			To(9, 14)
			fmt.Print("║  Ошибка отправки письма! Проверьте email  ║")
			To(9, 15)
			fmt.Print("╚═══════════════════════════════════════════╝")
			time.Sleep(2 * time.Second)
			Clear()
			continue
		}

		for attempts := 3; attempts > 0; attempts-- {
			Clear()
			To(9, 13)
			fmt.Print("╔════════════════════════════════════════════════════╗")
			To(9, 14)
			fmt.Print("║  Письмо с кодом отправлено! Проверьте почту(СПАМ)  ║")
			To(9, 15)
			fmt.Print("╚════════════════════════════════════════════════════╝")

			To(29, 2)
			fmt.Print("╔════════════════════════╗")
			To(29, 3)
			fmt.Print("║  Введите 1 для выхода  ║")
			To(29, 4)
			fmt.Print("╚════════════════════════╝")

			To(9, 6)
			fmt.Print("╔════════════════════════╗")
			To(9, 7)
			fmt.Print("║ Введите код из письма: ║")
			To(9, 8)
			fmt.Print("╚════════════════════════╝")
			To(36, 7)
			fmt.Print("-> ")
			Vod(&inputCode)
			if inputCode == "1" {
				Logoreg()
			}

			dbCode, created, dbAttempts, err := Data.GetVerificationCode(email)
			if err != nil {
				To(9, 17)
				fmt.Print("╔═══════════════════════════════════════════╗")
				To(9, 18)
				fmt.Print("║  Ошибка проверки кода! Попробуйте снова.  ║")
				To(9, 19)
				fmt.Print("╚═══════════════════════════════════════════╝")
				time.Sleep(2 * time.Second)
				break
			}

			if time.Since(created.UTC()) > 10*time.Minute {
				To(9, 17)
				fmt.Print("╔═══════════════════════════════════════════╗")
				To(9, 18)
				fmt.Print("║  Код устарел! Запросите новый код.       ║")
				To(9, 19)
				fmt.Print("╚═══════════════════════════════════════════╝")
				Data.DeleteVerificationCode(email)
				time.Sleep(2 * time.Second)
				break
			}

			if dbAttempts >= 3 {
				To(9, 17)
				fmt.Print("╔═══════════════════════════════════════════╗")
				To(9, 18)
				fmt.Print("║  Превышены попытки! Запросите новый код.  ║")
				To(9, 19)
				fmt.Print("╚═══════════════════════════════════════════╝")
				Data.DeleteVerificationCode(email)
				time.Sleep(2 * time.Second)
				break
			}

			if inputCode != dbCode {
				Data.IncrementAttempts(email)
				To(19, 2)
				fmt.Print("╔═════════════════════════════════════╗")
				To(19, 3)
				fmt.Printf("║  Неверный код! Осталось попыток: %d  ║", attempts-1)
				To(19, 4)
				fmt.Print("╚═════════════════════════════════════╝")
				time.Sleep(1 * time.Second)
				continue
			}

			Data.DeleteVerificationCode(email)
			currentUserEmail = email

			To(24, 10)
			fmt.Print("╔══════════════════════════════╗")
			To(24, 11)
			fmt.Print("║  Регистрация Email окончена  ║")
			To(24, 12)
			fmt.Print("╚══════════════════════════════╝")
			time.Sleep(2 * time.Second)
			Clear()
			registrPass()
			return
		}
		Clear()
		To(13, 9)
		fmt.Print("╔══════════════════════════════════════════════════╗")
		To(13, 10)
		fmt.Print("║  Превышено количество попыток. Начинаем заново.  ║")
		To(13, 11)
		fmt.Print("╚══════════════════════════════════════════════════╝")
		time.Sleep(2 * time.Second)
		Clear()
	}
}

func registrPass() {

	var confirmPassword, password string

	for {
		To(29, 2)
		fmt.Print("╔════════════════════════╗")
		To(29, 3)
		fmt.Print("║  Введите 1 для выхода  ║")
		To(29, 4)
		fmt.Print("╚════════════════════════╝")
		To(3, 19)
		fmt.Print("╔═════════════════════════════════════════════════════════════════════════╗")
		To(3, 20)
		fmt.Print("║ Пароль должен иметь 24 символа, специальный знак, латинские буквы (A,a) ║")
		To(3, 21)
		fmt.Print("╚═════════════════════════════════════════════════════════════════════════╝")
		To(9, 6)
		fmt.Print("╔══════════════════════════════╗")
		To(9, 7)
		fmt.Print("║     Введите ILYpassword:     ║")
		To(9, 8)
		fmt.Print("╚══════════════════════════════╝")
		To(9, 10)
		fmt.Print("╔══════════════════════════════╗")
		To(9, 11)
		fmt.Print("║   Подтвердите ILYpassword:   ║")
		To(9, 12)
		fmt.Print("╚══════════════════════════════╝")
		To(42, 11)
		fmt.Print("-> ")
		To(42, 7)
		fmt.Print("-> ")
		Vod(&password)
		if password == "1" {
			Logoreg()
		}

		if err := validatePassword(password); err != nil {
			To(9, 15)
			fmt.Printf("Ошибка: %v", err)
			time.Sleep(3 * time.Second)
			Clear()
			continue
		} else {
			break
		}

	}

	for {
		To(29, 2)
		fmt.Print("╔════════════════════════╗")
		To(29, 3)
		fmt.Print("║  Введите 1 для выхода  ║")
		To(29, 4)
		fmt.Print("╚════════════════════════╝")
		To(3, 19)
		fmt.Print("╔═════════════════════════════════════════════════════════════════════════╗")
		To(3, 20)
		fmt.Print("║ Пароль должен иметь 24 символа, специальный знак, латинские буквы (A,a) ║")
		To(3, 21)
		fmt.Print("╚═════════════════════════════════════════════════════════════════════════╝")
		To(9, 6)
		fmt.Print("╔══════════════════════════════╗")
		To(9, 7)
		fmt.Print("║     Введите ILYpassword:     ║")
		To(9, 8)
		fmt.Print("╚══════════════════════════════╝")
		To(42, 7)
		fmt.Print("-> ", password)
		To(9, 10)
		fmt.Print("╔══════════════════════════════╗")
		To(9, 11)
		fmt.Print("║   Подтвердите ILYpassword:   ║")
		To(9, 12)
		fmt.Print("╚══════════════════════════════╝")
		To(42, 11)
		fmt.Print("-> ")
		Vod(&confirmPassword)
		if confirmPassword == "1" {
			Logoreg()
		}

		if password != confirmPassword {
			To(9, 14)
			fmt.Print("╔════════════════════════╗")
			To(9, 15)
			fmt.Print("║  Пароли не совпадают!  ║")
			To(9, 16)
			fmt.Print("╚════════════════════════╝")
			time.Sleep(2 * time.Second)
			Clear()
			To(29, 2)
			fmt.Print("╔════════════════════════╗")
			To(29, 3)
			fmt.Print("║  Введите 1 для выхода  ║")
			To(29, 4)
			fmt.Print("╚════════════════════════╝")
			To(3, 19)
			fmt.Print("╔═════════════════════════════════════════════════════════════════════════╗")
			To(3, 20)
			fmt.Print("║ Пароль должен иметь 24 символа, специальный знак, латинские буквы (A,a) ║")
			To(3, 21)
			fmt.Print("╚═════════════════════════════════════════════════════════════════════════╝")
			To(9, 6)
			fmt.Print("╔══════════════════════════════╗")
			To(9, 7)
			fmt.Print("║     Введите ILYpassword:     ║")
			To(9, 8)
			fmt.Print("╚══════════════════════════════╝")
			To(42, 7)
			fmt.Print("-> ", password)
			if password == "1" {
				Logoreg()
			}
			continue
		} else {

			hash, salt := generatePasswordHash(password)

			if err := Data.SaveUser(currentUserEmail, hash, salt); err != nil {
				To(9, 20)
				fmt.Print("Ошибка сохранения данных!")
				time.Sleep(2 * time.Second)
				continue
			}

			Clear()
			To(24, 10)
			fmt.Print("╔══════════════════════════════╗")
			To(24, 11)
			fmt.Print("║     Регистрация окончена     ║")
			To(24, 12)
			fmt.Print("╚══════════════════════════════╝")
			time.Sleep(2 * time.Second)
			Clear()

			testPass := password
			valid, err := Data.ValidateMasterPassword(testPass)
			fmt.Printf("[TEST] Проверка пароля '%s': %v, ошибка: %v\n", testPass, valid, err)

			exists, _ := Data.UserExists(currentUserEmail)
			fmt.Printf("\n[DEBUG] Проверка записи: %v\n", exists)
			menu()
			break
		}

	}
}

func validatePassword(password string) error {
	if len(password) < 24 {
		return errors.New("пароль должен быть не короче 24 символов")
	}

	var (
		hasUpper   = false
		hasLower   = false
		hasDigit   = false
		hasSpecial = false
	)

	allowedSpecials := "=~`$%^!_-()[]&#@*?%+*-/:\"'.,\\/"
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case strings.ContainsRune(allowedSpecials, char):
			hasSpecial = true
		default:
			return fmt.Errorf("недопустимый символ: %c", char)
		}
	}

	errorMessages := []string{}
	if !hasUpper {
		errorMessages = append(errorMessages, "минимум 1 заглавная буква")
	}
	if !hasLower {
		errorMessages = append(errorMessages, "минимум 1 строчная буква")
	}
	if !hasDigit {
		errorMessages = append(errorMessages, "минимум 1 цифра")
	}
	if !hasSpecial {
		errorMessages = append(errorMessages, "минимум 1 спецсимвол (=~`$%^!_-()[]&#@*?%+*-/:\"'.,\\/)")
	}

	if len(errorMessages) > 0 {
		return errors.New("требуется: " + strings.Join(errorMessages, ", "))
	}

	return nil
}

func generatePasswordHash(password string) (string, string) {
	salt := make([]byte, 16)
	rand.Read(salt)

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		3,
		64*1024,
		4,
		32,
	)

	return fmt.Sprintf("%x", hash), fmt.Sprintf("%x", salt)
}
