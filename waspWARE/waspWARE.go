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
	"sync"
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

func decryptAES(encryptedData []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if len(encryptedData) < nonceSize {
		return nil, fmt.Errorf("şifreli veri çok kısa")
	}
	nonce := encryptedData[:nonceSize]
	data := encryptedData[nonceSize:]
	decryptedData, err := gcm.Open(nil, nonce, data, nil)
	if err != nil {
		return nil, err
	}
	return decryptedData, nil
}

// FileProcessor - İşlemci yapılandırması
type FileProcessor struct {
	key         []byte
	decryptMode bool
	totalFiles  int
	processed   int
	mu          sync.Mutex
}

// processFile - Tek bir dosyayı şifrele/şifreyi çöz
func (fp *FileProcessor) processFile(filePath string, binaryName, binaryPattern string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("okuma hatası: %v", err)
	}

	var processedContent []byte
	if fp.decryptMode {
		processedContent, err = decryptAES(content, fp.key)
	} else {
		processedContent, err = encryptAES(content, fp.key)
	}
	if err != nil {
		return fmt.Errorf("işlem hatası: %v", err)
	}

	err = os.WriteFile(filePath, processedContent, 0644)
	if err != nil {
		return fmt.Errorf("yazma hatası: %v", err)
	}

	fp.mu.Lock()
	fp.processed++
	progress := float64(fp.processed) / float64(fp.totalFiles) * 100
	fmt.Printf("\r  🔄 İlerleme: %.1f%% (%d/%d dosya)", progress, fp.processed, fp.totalFiles)
	fp.mu.Unlock()

	return nil
}

// makeEncrypt - Çoklu iş parçacığı ile şifrele/şifreyi çöz
func makeEncrypt(path string, key []byte, decryptMode bool) error {
	binaryName := filepath.Base(os.Args[0])
	binaryPattern := binaryName + ".wasp"

	// Dosyaları topla
	var files []string
	filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filePath == binaryName || filePath == binaryPattern {
			fmt.Printf("  ℹ️  Kendi dosyası atlandı: %s\n", filePath)
			return nil
		}
		// Decrypt modunda sadece .wasp dosyalarını işle
		if decryptMode && !strings.HasSuffix(filePath, ".wasp") {
			return nil
		}
		files = append(files, filePath)
		return nil
	})

	fp := &FileProcessor{
		key:         key,
		decryptMode: decryptMode,
		totalFiles:  len(files),
	}

	if fp.totalFiles == 0 {
		fmt.Println("\r  ℹ️  İşlenecek dosya bulunamadı")
		return nil
	}

	// Worker pool - 4 iş parçacığı
	const numWorkers = 4
	jobChan := make(chan string, len(files))
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for filePath := range jobChan {
				fp.processFile(filePath, binaryName, binaryPattern)
			}
		}()
	}

	// İşleri dağıt
	for _, f := range files {
		jobChan <- f
	}
	close(jobChan)
	wg.Wait()

	fmt.Println("\r  ✓ Tüm dosyalar işlendi")

	// Terminal açık kalacak - çıkış yapmıyoruz
	return nil
}

// changeExtensions - Şifreleme için .wasp ekle, şifreyi çözmek için kaldır
func changeExtensions(path string, decryptMode bool) error {
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		var newPath string
		if decryptMode {
			// .wasp uzantısını kaldır
			if strings.HasSuffix(filePath, ".wasp") {
				newPath = strings.TrimSuffix(filePath, ".wasp")
				fmt.Printf("  🔄 %s -> %s\n", filePath, newPath)
			}
		} else {
			// .wasp uzantısını ekle
			newPath = filePath + ".wasp"
			fmt.Printf("  🔄 %s -> %s\n", filePath, newPath)
		}

		if newPath != "" && newPath != filePath {
			err := os.Rename(filePath, newPath)
			if err != nil {
				fmt.Printf("  ⚠️  Uzantı değiştirme hatası (%s): %v\n", filePath, err)
				return nil
			}
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

func runEncryption(key string, directory string, decryptMode bool) {
	modeText := "ŞİFRELEME"
	if decryptMode {
		modeText = "ŞİFRESINI ÇÖZME"
	}

	fmt.Println("\n" + repeatChar('=', 60))
	fmt.Printf("🔐 WASPWARE - %s MODU\n", modeText)
	fmt.Println(repeatChar('=', 60))
	fmt.Printf("\n🔑 Anahtar: %s\n", key)
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

	fmt.Printf("\n🚀 '%s' dizinindeki dosyalar işleniyor...\n\n", directory)

	err = makeEncrypt(directory, keyBytes, decryptMode)
	if err != nil {
		fmt.Printf("❌ İşlem sırasında hata: %v\n", err)
		return
	}

	fmt.Println("\n🔄 Dosya uzantıları güncelleniyor...")
	err = changeExtensions(directory, decryptMode)
	if err != nil {
		fmt.Printf("⚠️  Uzantı değiştirme sırasında bazı hatalar: %v\n", err)
	}

	fmt.Println("\n" + repeatChar('=', 60))
	if decryptMode {
		fmt.Println("✅ ŞİFRESINI ÇÖZME TAMAMLANDI!")
	} else {
		fmt.Println("✅ ŞİFRELEME TAMAMLANDI!")
	}
	fmt.Println(repeatChar('=', 60))
	fmt.Printf("\n📁 İşlenen dosyalar: %s/\n", directory)
	if decryptMode {
		fmt.Println("💾 Dosyalar AES-256-GCM ile şifresini çözüldü.")
	} else {
		fmt.Println("💾 Dosyalar AES-256-GCM ile şifrelenmiştir.")
	}
	fmt.Println("🔑 Anahtarı güvenli bir yerde saklayınız!")
	fmt.Println("\n💡 Terminal açık kalıyor. İşlem tamamlandı.")
	fmt.Println(repeatChar('=', 60))
}

func main() {
	keyPtr := flag.String("key", "", "Anahtar")
	dizinPtr := flag.String("dizin", "", "Hedef dizin")
	decryptPtr := flag.Bool("decrypt", false, "Şifresini çözme modu")

	flag.Parse()

	if *keyPtr != "" && *dizinPtr != "" {
		fmt.Printf("🔑 Komut satırı anahtarı kullanılıyor: %s\n", *keyPtr)
		runEncryption(*keyPtr, *dizinPtr, *decryptPtr)
		return
	}

	fmt.Println("\n" + repeatChar('=', 60))
	fmt.Println("🔐 WASPWARE - Dosya Şifreleme Aracı")
	fmt.Println(repeatChar('=', 60))
	fmt.Println("\n📝 AES-256-GCM şifreleme kullanır.")
	fmt.Println("📝 Çoklu iş parçacığı (4 worker) ile hızlı işlem.")
	fmt.Println("📝 İlerleme göstergesi ile takip edilebilir.\n")

	fmt.Print("🔄 İşlem modu seçiniz [1=Şifrele, 2=Şifresini Çöz]: ")
	modeInput, _ := askForKey()
	modeInput = trimSpace(modeInput)

	var decryptMode bool
	if modeInput == "1" {
		decryptMode = false
	} else if modeInput == "2" {
		decryptMode = true
	} else {
		fmt.Println("❌ Geçersiz seçim. Şifreleme modu kullanılıyor.")
		decryptMode = false
	}

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

	runEncryption(userKey, directory, decryptMode)
        fmt.Println("\nÇıkmak için Enter'a basın...")
        bufio.NewReader(os.Stdin).ReadBytes('\n')
}
