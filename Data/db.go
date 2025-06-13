package Data

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/crypto/argon2"
)

var DB *sql.DB

func InitDB() error {
	workingDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("не удалось получить рабочую директорию: %v", err)
	}

	dataDir := filepath.Join(workingDir, "Data")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return fmt.Errorf("не удалось создать папку Data: %v", err)
	}

	dbPath := filepath.Join(dataDir, "passwords.db")
	fmt.Printf("База данных будет создана по пути: %s\n", dbPath)

	DB, err = sql.Open("sqlite3", dbPath+"?_journal_mode=WAL&_foreign_keys=1")
	if err != nil {
		return fmt.Errorf("ошибка подключения к БД: %v", err)
	}

	if err := DB.Ping(); err != nil {
		return fmt.Errorf("не удалось подключиться к БД: %v", err)
	}

	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT UNIQUE NOT NULL CHECK(length(email) > 0),
			master_key_hash TEXT NOT NULL CHECK(length(master_key_hash) > 0),
			salt TEXT NOT NULL CHECK(length(salt) > 0),
			created_at TEXT DEFAULT CURRENT_TIMESTAMP
		)
	`)

	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS verification_codes (
        email TEXT PRIMARY KEY,
        code TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        attempts INTEGER DEFAULT 0
    )`)

	if err != nil {
		return fmt.Errorf("ошибка создания таблицы: %v", err)
	}

	return nil
}

func SaveUser(email, hash, salt string) error {
	if DB == nil {
		return fmt.Errorf("база данных не инициализирована")
	}

	tx, err := DB.Begin()
	if err != nil {
		return fmt.Errorf("ошибка начала транзакции: %v", err)
	}
	defer tx.Rollback()

	var exists bool
	err = tx.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)",
		email,
	).Scan(&exists)

	if err != nil {
		return fmt.Errorf("ошибка проверки email: %v", err)
	}

	if exists {
		return fmt.Errorf("email %s уже зарегистрирован", email)
	}

	_, err = tx.Exec(
		"INSERT INTO users (email, master_key_hash, salt) VALUES (?, ?, ?)",
		email, hash, salt,
	)

	if err != nil {
		return fmt.Errorf("ошибка сохранения данных: %v", err)
	}

	return tx.Commit()
}

func UserExists(email string) (bool, error) {
	if DB == nil {
		return false, fmt.Errorf("база данных не инициализирована")
	}

	var exists bool
	err := DB.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)",
		email,
	).Scan(&exists)

	return exists, err
}

func GetFirstUserEmail() (string, error) {
	if DB == nil {
		return "", fmt.Errorf("база данных не инициализирована")
	}

	var email string
	err := DB.QueryRow("SELECT email FROM users LIMIT 1").Scan(&email)
	return email, err
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

func ValidateMasterPassword(password string) (bool, error) {
	if DB == nil {
		return false, fmt.Errorf("база данных не инициализирована")
	}

	var storedHash, saltHex string
	err := DB.QueryRow(
		"SELECT master_key_hash, salt FROM users LIMIT 1",
	).Scan(&storedHash, &saltHex)

	if err != nil {
		return false, fmt.Errorf("ошибка получения данных: %v", err)
	}

	salt, err := hex.DecodeString(saltHex)
	if err != nil {
		return false, fmt.Errorf("ошибка декодирования соли: %v", err)
	}

	inputHash := fmt.Sprintf("%x", argon2.IDKey(
		[]byte(password),
		salt,
		3,
		64*1024,
		4,
		32,
	))

	return inputHash == storedHash, nil
}

func UserExistsByID(userID int) (bool, error) {
	if DB == nil {
		return false, fmt.Errorf("база данных не инициализирована")
	}

	var exists bool
	err := DB.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)",
		userID,
	).Scan(&exists)

	return exists, err
}

func SaveVerificationCode(email, code string) error {
	_, err := DB.Exec(`
        INSERT OR REPLACE INTO verification_codes 
        (email, code, created_at, attempts) 
        VALUES (?, ?, datetime('now'), 0)`,
		email, code)
	return err
}

func GetVerificationCode(email string) (code string, createdAt time.Time, attempts int, err error) {
	row := DB.QueryRow("SELECT code, created_at, attempts FROM verification_codes WHERE email = ?", email)
	err = row.Scan(&code, &createdAt, &attempts)
	return
}

func IncrementAttempts(email string) error {
	_, err := DB.Exec(`
        UPDATE verification_codes 
        SET attempts = attempts + 1 
        WHERE email = ?`,
		email)
	return err
}

func DeleteVerificationCode(email string) error {
	_, err := DB.Exec("DELETE FROM verification_codes WHERE email = ?", email)
	return err
}
