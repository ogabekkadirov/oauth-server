package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

const (
	saltLength  = 32
	encryptCost = 13
)

// HashPassword ...
func HashPassword(saltedPass string) string {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(saltedPass), encryptCost)
	if err != nil {
		log.Fatal(err)
	}
	return string(hashedPass)
}

// Combine ...
func Combine(salt string, rawPassword string) string {
	arr := []string{salt, rawPassword}
	saltedPassword := strings.Join(arr, "")
	return saltedPassword
}

// GenerateSalt ...
func GenerateSalt() string {
	data := make([]byte, saltLength)
	_, err := rand.Read(data)
	if err != nil {
		log.Fatal(err)
	}

	salt := base64.URLEncoding.EncodeToString(data)
	return salt
}

// PasswordMatch ..
func PasswordMatch(rawPassword, hash string) bool {
	// saltedPassword := Combine(salt, rawPassword)
	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(rawPassword)) != nil {
		return false
	}
	return true
}
