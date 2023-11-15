package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// KEY İÇİN RANDOM BASE64 STRİNG ÜRETEN FONKSİYON
func generateRandomBase64String(length int) (string, error) {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(randomBytes), nil
}

// ÜRETİLEN RANDOM BASE64 STRİNG İLE KEY ÜRETEN FONKSİYON
func createKey() string {
	keyLength := 16
	randomBase64String, err := generateRandomBase64String(keyLength)
	if err != nil {
		fmt.Println("Error generating random key:", err)
		return ""
	}
	return randomBase64String
}

func main() {
	fmt.Println("KULLANILAN KEY : ", createKey(), "\nBU KEYİ SAKLAYINIZ")
}
