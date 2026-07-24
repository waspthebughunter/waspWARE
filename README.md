# WASPWARE - Dosya Şifreleme Aracı

![WASPWARE Logo](https://github.com/waspthebughunter/waspWARE/assets/100480448/d562bb91-c0be-4ce8-89b4-1b3aef901f13)
![WASPWARE Screenshot](https://github.com/waspthebughunter/waspWARE/assets/100480448/b1351478-cc41-4813-bc94-e921ae9cfab2)

> **Uyarı:** Bu araç eğitim amaçlıdır. Zararlı eylemler için sorumluluk kabul edilmez!

---

## 📋 Genel Bakış

WASPWARE, dosyaları ve dizinleri AES-256-GCM şifreleme algoritması kullanarak güvenli bir şekilde şifreleyen güçlü bir Go tabanlı dosya şifreleme aracına sahiptir. Şifrelenen tüm dosyalar kolay tanımlanması için otomatik olarak `.wasp` uzantısına sahip olacak şekilde yeniden adlandırılır.

### Ana Özellikler

- 🔐 **AES-256-GCM Şifreleme** - Endüstri standardı şifreleme algoritması
- 📁 **Dizin Genişletmeli Şifreleme** - Tüm klasörleri yeniden adlandırma ile şifreleme
- 🔑 **Güvenli Anahtar Yönetimi** - Opsiyonel şifre tabanlı şifreleme anahtarları
- 🔄 **Otomatik Uzantı Değişimi** - Şifrelenen dosyalar `.wasp` uzantısı alır
- 💾 **Tek Çalıştırılabilir Dosya** - Dış bağımlılıklar gerektirmez
- 🎯 **Çapraz Platform** - Linux, Windows ve çoklu mimariler için derleme

---

## 🚀 Hızlı Başlangıç

### Temel Kullanım

```bash
# Aracı derle (optimizasyon ile)
go build -ldflags="-s -w" -o waspWARE ./waspWARE/waspWARE.go

# UPX sıkıştırma (daha küçük binary için önerilir)
upx -9 --best waspWARE

# İnteraktif olarak çalıştır
./waspWARE

# Veya komut satırı bayrakları ile
./waspWARE -key "gizli-anahtarınız" -dizin "/şifrelemek-istediğiniz/dizin"
```

### UPX Sıkıştırma (Opsiyonel ama Önerilir)

```bash
# Binary'yi UPX ile sıkıştırın (~%57 boyut azalması)
upx -9 --best waspWARE

# Sıkıştırmayı doğrula
ls -lh waspWARE
# Çıktı: ~851 KB (2.0 MB'den)
```

### İnteraktif Mod

WASPWARE'yi argüman olmadan çalıştırdığınızda:
1. Şifreleme anahtarı için sizden isteyecek (veya rastgele bir tane oluşturacak)
2. Onay isteyecek
3. Hedef dizini talep edecek
4. O dizindeki tüm dosyaları şifreleyecek

---

## 🔧 Komut Satırı Seçenekleri

```bash
./waspWARE -key "anahtarınız" -dizin "/hedef/dizin"
```

- `-key`: Şifreleme anahtarı (opsiyonel - rastgele oluşturulabilir)
- `-dizin`: Şifrelemek istediğiniz dizin (opsiyonel - varsayılan olarak mevcut dizin)

---

## 📖 Dokümantasyon

### Derleme Rehberleri

- **[Linux Derleme Rehberi](./build-linux.md)** - Farklı mimariler için çapraz derleme dahil olmak üzere Linux platformlarında derlemek için kapsamlı talimatlar.
- **[Windows Derleme Rehberi](./build-windows.md)** - Windows'ta ve Linux'tan çapraz derleme için adım adım talimatlar.

### Desteklenen Platformlar

| Platform | Mimari | Binary Boyutu (Ham) | Binary Boyutu (UPX) |
|----------|--------|---------------------|---------------------|
| Linux | amd64 (x86_64) | ~2.0 MB | **~851 KB** ✅ |
| Linux | arm64 (ARM 64-bit) | ~2.0 MB | **~851 KB** ✅ |
| Linux | armv7 (ARM 32-bit) | ~2.0 MB | **~851 KB** ✅ |
| Linux | armv8 (ARM 64-bit) | ~2.0 MB | **~851 KB** ✅ |
| Linux | 386 (x86 32-bit) | ~2.0 MB | **~851 KB** ✅ |
| Windows | amd64 | ~2.0 MB | **~851 KB** ✅ |
| Windows | arm64 | ~2.0 MB | **~851 KB** ✅ |
| Windows | 386 | ~2.0 MB | **~851 KB** ✅ |

### Uyumlu Dağıtımlar

- Ubuntu, Debian, Mint
- CentOS, RHEL, Fedora
- Arch Linux, Manjaro
- Raspberry Pi OS
- Alpine Linux
- Go yüklü herhangi bir Linux dağıtımı

---

## 🔐 Nasıl Çalışır

1. **Anahtar Oluşturma**: Bir şifreleme anahtarı sağlarsınız veya WASPWARE rastgele bir tane oluşturur
2. **Dizin Tarama**: Araç hedef dizini yeniden adlandırma ile tarar
3. **AES Şifreleme**: Her dosya AES-256-GCM kullanarak şifrelenir
4. **Uzantı Değişimi**: Şifrelenen dosyalar `.wasp` uzantısına sahip olacak şekilde yeniden adlandırılır
5. **Tamamlama**: Tüm dosyalar şifrelenir ve güvenli depolmaya hazır hale gelir

---

## 📝 Kullanım Örnekleri

### Belirli Bir Dizini Şifreleme

```bash
./waspWARE -key "GizliAnahtar123" -dizin "/kullanici/dokumanlar"
```

### Rastgele Anahtar Oluşturma

```bash
./waspWARE  # Anahtar girmeden Enter'a basarak rastgele anahtar oluşturur
```

### Mevcut Dizini Şifreleme

```bash
./waspWARE -key "anahtarim"  # Yol verilmezse mevcut dizini şifreler
```

---

## 🔒 Güvenlik Dikkatleri

- **AES-256-GCM** onaylanmış şifreleme ile ilişkili veri (AEAD) sağlar
- **Rastgele Nonce**: Her şifreleme kriptografik olarak güvenli rastgele nonce kullanır
- **Anahtar Yönetimi**: Şifreleme anahtarlarınızı güvenli bir şekilde saklayın - kaybedildikleri geri kazanılamaz
- **Dosya İzinleri**: Şifrelenen dosyalar orijinal izinlerini korur

---

## 🛠️ Derleme Talimatları

### Yerel Derleme (Mevcut Platform)

```bash
# Optimizasyon bayrakları ile derle
go build -ldflags="-s -w" -o waspWARE ./waspWARE/waspWARE.go

# Opsiyonel: Daha küçük binary için UPX sıkıştırma (~%57 boyut azalması)
upx -9 --best waspWARE

# Doğrula
ls -lh waspWARE
# Çıktı: ~851 KB (2.0 MB'den)
```

### Çapraz Derleme Örnekleri

#### Linux için ARM64
```bash
GOOS=linux GOARCH=arm64 go build -o waspWARE-arm64 ./waspWARE/waspWARE.go
```

#### Windows Çalıştırılabilir Dosyası Linux'tan
```bash
GOOS=windows GOARCH=amd64 go build -o waspWARE.exe ./waspWARE/waspWARE.go
```

### Docker Derleme

```bash
docker run --rm -v $(pwd):/app golang:1.21-alpine \
  go build -o /app/waspWARE ./waspWARE/waspWARE.go
```

---

## 📚 Ek Kaynaklar

- [Go Çapraz Derleme Dokümantasyonu](https://go.dev/doc/install/source#environment)
- [AES Şifreleme En İyi Uygulamaları](https://csrc.nist.gov/publications/detail/sp/800-38a/final)
- [Go Proje Düzeni](https://github.com/golang-standards/project-layout)

---

## ⚠️ Önemli Notlar

- **Anahtarları Yedekleyin**: Şifreleme anahtarları kaybedildikleri geri kazanılamaz
- **Üretimden Önce Test Edin**: Her zaman şifrelemeyi küçük bir alt küme üzerinde test edin
- **Anahtar Depolama**: Anahtarları güvenli konumlarda saklayın (şifre yöneticileri, şifreli dosyalar)
- **Dosya Bütünlüğü**: Şifrelenen dosyalar orijinal izinlerini ve meta verilerini korur

---

## 📄 Lisans

Bu proje eğitim amaçları için sağlanmaktadır. Sorumlu ve yasal olarak kullanın.

---

## 🙏 Teşekkürler

- Go Language Team
- AES Şifreleme Standartları (NIST SP 800-38A)
- Açık kaynak topluluğu

---

<div align="center">
  <strong>WASPWARE - Dosya Şifrelemeyi Basit Hale Getiriyor</strong>
</div>
