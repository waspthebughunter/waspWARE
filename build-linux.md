# WASPWARE - Linux Derleme Rehberi

## 📋 Genel Bakış
Bu rehber, WASPWARE'i çeşitli Linux platformlarında derlemek için kapsamlı talimatlar sağlar, farklı mimariler için çapraz derleme dahil.

---

## 🔧 Gereksinimler

### 1. Go Dili Kurulumu

#### Ubuntu/Debian
```bash
# Go deposunu ekle
wget https://go.dev/release/go1.21.3.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.3.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc

# Doğrula
go version
```

#### CentOS/RHEL/Fedora
```bash
sudo yum install -y go
# VEYA
sudo dnf install -y go
```

#### Arch Linux
```bash
sudo pacman -S go
```

### 2. Derleme Araçları Kurulumu (Opsiyonel)
```bash
# Yerel çapraz derleme için
sudo apt-get install gcc musl-dev   # Debian/Ubuntu
sudo yum install gcc musl-devel      # CentOS/RHEL
```

---

## 🏗️ Temel Derleme Komutları

### Basit Derleme (Mevcut Platform)
```bash
cd waspWARE
go build -o waspWARE ./waspWARE/waspWARE.go
```

### Optimizasyon Bayrakları ile Derleme
```bash
go build -ldflags="-s -w" -o waspWARE ./waspWARE/waspWARE.go

# Opsiyonel: Daha küçük binary için UPX sıkıştırma (~%57 boyut azalması)
upx -9 --best waspWARE

# Doğrula
ls -lh waspWARE
# Çıktı: ~851 KB (2.0 MB'den)
```

### Versiyon Bilgisi ile Derleme
```bash
go build -ldflags="-s -w -X 'main.Version=1.0.0'" \
  -o waspWARE ./waspWARE/waspWARE.go

# UPX sıkıştırma
upx -9 --best waspWARE
```

---

## 🎯 Farklı Linux Platformları için Çapraz Derleme

### amd64/x86_64 (Standart Masaüstü)
```bash
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o waspWARE-amd64 ./waspWARE/waspWARE.go

# Opsiyonel: UPX sıkıştırma
upx -9 --best waspWARE-amd64
```

### ARM64 (Apple Silicon, Raspberry Pi 4/5, AWS Graviton)
```bash
GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o waspWARE-arm64 ./waspWARE/waspWARE.go

# Opsiyonel: UPX sıkıştırma
upx -9 --best waspWARE-arm64
```

### ARMv7 (Raspberry Pi 3/4, eski ARM cihazlar)
```bash
GOOS=linux GOARCH=arm GOARM=7 go build -ldflags="-s -w" -o waspWARE-armhf ./waspWARE/waspWARE.go

# Opsiyonel: UPX sıkıştırma
upx -9 --best waspWARE-armhf
```

### ARMv8 (Raspberry Pi 4/5, yeni ARM cihazlar)
```bash
GOOS=linux GOARCH=arm GOARM=8 go build -ldflags="-s -w" -o waspWARE-armv8 ./waspWARE/waspWARE.go

# Opsiyonel: UPX sıkıştırma
upx -9 --best waspWARE-armv8
```

### 32-bit x86 (Eski Sistemler)
```bash
GOOS=linux GOARCH=386 go build -ldflags="-s -w" -o waspWARE-386 ./waspWARE/waspWARE.go

# Opsiyonel: UPX sıkıştırma
upx -9 --best waspWARE-386
```

---

## 🪟 Linux'tan Windows (.exe) için Çapraz Derleme

### Linux'ta Windows Çalıştırılabilir Dosyası Derleme

GOOS ortam değişkenini kullanarak WASPWARE'i Windows için çapraz derleyebilirsiniz:

#### Temel Windows Derlemesi (amd64)
```bash
# Windows amd64 için derle (standart Windows 10/11)
GOOS=windows GOARCH=amd64 go build -o waspWARE.exe ./waspWARE/waspWARE.go

# Çalıştırılabilir dosyayı doğrula
file waspWARE.exe
# Çıktı: waspWARE.exe: PE32+ executable (console) x86-64, for MS Windows
```

#### Farklı Windows Mimarileri için Derleme

##### Windows on ARM (Surface Pro X, yeni cihazlar)
```bash
GOOS=windows GOARCH=arm64 go build -o waspWARE-arm64.exe ./waspWARE/waspWARE.go
```

##### 32-bit Windows (Eski Sistemler)
```bash
GOOS=windows GOARCH=386 go build -o waspWARE-386.exe ./waspWARE/waspWARE.go
```

##### Windows ARM (Windows on ARM cihazlar)
```bash
GOOS=windows GOARCH=arm go build -o waspWARE-arm.exe ./waspWARE/waspWARE.go
```

#### Optimizasyonlu Windows Derlemesi
```bash
# Daha küçük binary için optimizasyon bayrakları ile derle
GOOS=windows GOARCH=amd64 go build \
  -ldflags="-s -w" \
  -o waspWARE.exe ./waspWARE/waspWARE.go

# Boyut ve özellikleri doğrula
ls -lh waspWARE.exe
file waspWARE.exe
```

---

## 🐳 Docker Derleme (Çapraz Derleme için Önerilir)

### Resmi Go Görselini Kullanma
```bash
# Göreyi çek
docker pull golang:1.21-alpine

# Mevcut platform için derle
docker run --rm -v $(pwd):/app golang:1.21-alpine \
  go build -o /app/waspWARE ./waspWARE/waspWARE.go

# ARM64 için derle
docker run --rm -v $(pwd):/app golang:1.21-alpine \
  GOOS=linux GOARCH=arm64 go build -o /app/waspWARE-arm64 ./waspWARE/waspWARE.go

# ARMv7 için derle
docker run --rm -v $(pwd):/app golang:1.21-alpine \
  GOOS=linux GOARCH=arm GOARM=7 go build -o /app/waspWARE-armhf ./waspWARE/waspWARE.go
```

### Windows için Çoklu Aşamalı Dockerfile Kullanma

`Dockerfile` oluştur:
```dockerfile
# Aşama 1: Linux için derleme
FROM golang:1.21-alpine AS builder-linux

WORKDIR /app
COPY waspWARE/waspWARE.go .

RUN go build -ldflags="-s -w" -o waspWARE ./waspWARE/waspWARE.go

# Aşama 2: Windows için derleme
FROM golang:1.21-alpine AS builder-windows

WORKDIR /app
COPY waspWARE/waspWARE.go .

RUN GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o waspWARE.exe .

# Aşama 3: Windows binary'si olan son görsel
FROM alpine:latest
RUN apk add --no-cache winpty

COPY --from=builder-windows /app/waspWARE.exe /waspWARE.exe

ENTRYPOINT ["winpty", "/waspWARE.exe"]
```

Derle ve çalıştır:
```bash
docker build -t waspware:latest .

# Linux'ta çalıştır (winpty aracılığıyla çalışacak)
docker run --rm -v $(pwd):/app waspware:latest \
  -key "mykey" -dizin "/target/directory"

# Veya Windows binary'sini başka bir sisteme kopyala
docker run --rm -v $(pwd)/bin:/output waspware:latest \
  cp /waspWARE.exe /output/waspWARE-windows-amd64.exe
```

---

### Çoklu Mimarili Docker Derlemesi

`build-all-platforms.sh` oluştur:
```bash
#!/bin/bash
# WASPWARE Çoklu-Platform Derleme Senaryosu (Linux + Windows)

set -e

BINARY_DIR="bin"
mkdir -p "$BINARY_DIR"

echo "🔨 WASPWARE için tüm platformlarda derleniyor..."

# Linux Platformları
echo ""
echo "🐧 Linux platformları için derleniyor..."

GOOS=linux GOARCH=amd64 go build -o "$BINARY_DIR/waspWARE-amd64" ./waspWARE/waspWARE.go
GOOS=linux GOARCH=arm64 go build -o "$BINARY_DIR/waspWARE-arm64" ./waspWARE/waspWARE.go
GOOS=linux GOARCH=arm GOARM=7 go build -o "$BINARY_DIR/waspWARE-armhf" ./waspWARE/waspWARE.go
GOOS=linux GOARCH=386 go build -o "$BINARY_DIR/waspWARE-386" ./waspWARE/waspWARE.go

# Windows Platformları
echo ""
echo "🪟 Windows platformları için derleniyor..."

GOOS=windows GOARCH=amd64 go build -o "$BINARY_DIR/waspWARE-windows-amd64.exe" ./waspWARE/waspWARE.go
GOOS=windows GOARCH=arm64 go build -o "$BINARY_DIR/waspWARE-windows-arm64.exe" ./waspWARE/waspWARE.go

echo ""
echo "✅ Tüm derlemeler tamamlandı!"
ls -lh "$BINARY_DIR/"
```

Çalıştırılabilir yap ve çalıştır:
```bash
chmod +x build-all-platforms.sh
./build-all-platforms.sh
```

---

### Tüm Platformlar için Çoklu Aşamalı Docker

`Dockerfile.multi` oluştur:
```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY waspWARE/waspWARE.go .

ARG TARGETOS=linux
ARG TARGETARCH=amd64
ARG GOARM=

RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    ${GOARM:+GOARM=$GOARM} go build -ldflags="-s -w" -o waspWARE .

FROM alpine:latest AS linux
COPY --from=builder /app/waspWARE /waspWARE
ENTRYPOINT ["/waspWARE"]

FROM alpine:latest AS windows
RUN apk add --no-cache winpty
COPY --from=builder /app/waspWARE.exe /waspWARE.exe
ENTRYPOINT ["winpty", "/waspWARE.exe"]
```

Belirli platform için derleme:
```bash
# Linux amd64
docker build --target linux -t waspware-linux .

# Windows amd64
docker build --target windows -t waspware-windows .
```

---

## 📦 Tüm Platformlar Derleme Senaryosu

`build-all.sh` oluştur:
```bash
#!/bin/bash
# WASPWARE Çoklu-Platform Derleme Senaryosu

set -e

BINARY_DIR="bin"
mkdir -p "$BINARY_DIR"

echo "🔨 WASPWARE için tüm platformlarda derleniyor..."

# amd64
GOOS=linux GOARCH=amd64 go build -o "$BINARY_DIR/waspWARE-amd64" ./waspWARE/waspWARE.go

# arm64
GOOS=linux GOARCH=arm64 go build -o "$BINARY_DIR/waspWARE-arm64" ./waspWARE/waspWARE.go

# armv7
GOOS=linux GOARCH=arm GOARM=7 go build -o "$BINARY_DIR/waspWARE-armhf" ./waspWARE/waspWARE.go

# armv8
GOOS=linux GOARCH=arm GOARM=8 go build -o "$BINARY_DIR/waspWARE-armv8" ./waspWARE/waspWARE.go

# 386
GOOS=linux GOARCH=386 go build -o "$BINARY_DIR/waspWARE-386" ./waspWARE/waspWARE.go

echo "✅ Tüm derlemeler tamamlandı!"
ls -lh "$BINARY_DIR/"
```

Çalıştırılabilir yap ve çalıştır:
```bash
chmod +x build-all.sh
./build-all.sh
```

---

## ✅ Doğrulama Adımları

### 1. Binary'nin Var Olup Olmadığını Kontrol Et
```bash
ls -lh waspWARE*
```

### 2. Çalıştırma Testi
```bash
./waspWARE -help
# VEYA interaktif olarak çalıştır
./waspWARE
```

### 3. Dosya Özelliklerini Doğrula
```bash
file waspWARE
stat waspWARE
```

### 4. Şifreleme Testi
```bash
mkdir -p /tmp/test-encryption
echo "Test dosya içeriği" > /tmp/test-encryption/test.txt
./waspWARE -key "test123" -dizin /tmp/test-encryption
ls -la /tmp/test-encryption/
```

---

## 📊 Mimariye Göre Binary Boyutları

| Platform | Mimari | Binary Boyutu | Kullanım Alanı |
|----------|--------|---------------|----------------|
| Linux | amd64 | ~3.0 MB | Standart x86_64 Linux |
| Linux | arm64 | ~2.8 MB | Apple Silicon, Raspberry Pi 4/5 |
| Linux | armv7 | ~2.8 MB | Raspberry Pi 3/4 (32-bit) |
| Linux | armv8 | ~2.8 MB | Raspberry Pi 4/5 (64-bit) |
| Linux | 386 | ~2.9 MB | Eski 32-bit sistemler |
| Windows | amd64 | ~3.0 MB | Windows 10/11 x86_64 |
| Windows | arm64 | ~2.8 MB | Windows on ARM cihazlar |
| Windows | 386 | ~2.9 MB | Eski 32-bit Windows |

---

## 🪟 Windows Derleme Doğrulama

### Windows Çalıştırılabilir Dosyasını Doğrula
```bash
# Dosya tipini kontrol et
file waspWARE.exe
# Çıktı: waspWARE.exe: PE32+ executable (console) x86-64, for MS Windows

# Binary boyutunu kontrol et
ls -lh waspWARE.exe

# Geçerli Windows PE çalıştırılabilir dosyası olduğunu doğrula
readelf -h waspWARE.exe 2>/dev/null || objdump -f waspWARE.exe
```

### Windows'a Taşıma
```bash
# SCP/SFTP aracılığıyla Windows makinesine kopyala
scp waspWARE.exe user@windows-machine:/path/to/

# VEYA USB sürücü kullan
cp waspWARE.exe /media/usb/waspWARE.exe

# VEYA binary'yi e-posta ile gönder
# (Windows çalıştırılabilir dosyası eklenti olarak eklenebilir ve gönderilebilir)
```

### Windows'ta Çalıştırma
Windows'a taşındıktan sonra:
```powershell
# PowerShell
.\waspWARE.exe -key "mykey" -dizin "C:\target\directory"

# Komut Satırı
waspWARE.exe -key "mykey" -dizin "C:\target\directory"
```

---

## 🔧 Sorun Giderme

### Sorun: "go: module bulunamadı"
**Çözüm:** Doğru dizinde olduğunuzdan emin olun
```bash
cd waspWARE
go build -o waspWARE ./waspWARE/waspWARE.go
```

### Sorun: Çapraz derleme "bilinmeyen mimari" ile başarısız oluyor
**Çözüm:** Gerekli başlıkları kurun
```bash
# ARM64 çapraz derleme için
sudo apt-get install gcc-aarch64-linux-gnu
export CC=aarch64-linux-gnu-gcc

GOOS=linux GOARCH=arm64 CC=aarch64-linux-gnu-gcc \
  go build -o waspWARE-arm64 ./waspWARE/waspWARE.go
```

### Sorun: İzni hatası ile derleme başarısız oluyor
**Çözüm:** sudo kullanın veya izinleri ayarlayın
```bash
sudo chown -R $USER:$USER waspWARE
go build -o waspWARE ./waspWARE/waspWARE.go
```

### Sorun: Go kurulumu eksik
**Çözüm:** Resmi kaynaktan Go'yu kurun
```bash
# Ubuntu/Debian
curl -LO https://go.dev/dl/go1.21.3.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.3.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
```

### Sorun: Windows derlemesi "bilinmeyen mimari" ile başarısız oluyor
**Çözüm:** Gerekli çapraz derleme araçlarını kurun
```bash
# Windows hedefi için GCC kurulumu (opsiyonel, Go bunu olmadan da çapraz derleyebilir)
sudo apt-get install gcc-multilib

# ARM64 Windows derlemeleri için
GOOS=windows GOARCH=arm64 go build -o waspWARE-arm64.exe ./waspWARE/waspWARE.go
```

### Sorun: Oluşturulan .exe Linux'ta tanınmıyor
**Çözüm:** Doğru platform için derlediğinizden emin olun
```bash
# Derleme hedefini doğrula
echo "Şu platform için derleniyor: GOOS=$GOOS GOARCH=$GOARCH"

# Temizle ve yeniden derle
rm -f *.exe
GOOS=windows GOARCH=amd64 go build -o waspWARE.exe ./waspWARE/waspWARE.go

# Doğrula
file waspWARE.exe
```

### Sorun: Linux'ta .exe çalıştırırken "İzin reddedildi"
**Çözüm:** Bu beklenir - Windows çalıştırılabilir dosyaları doğrudan Linux'ta çalışamaz. Windows sistemine taşıyın.
```bash
# Dosya tipini kontrol et
file waspWARE.exe
# Şunu göstermeli: PE32+ executable for MS Windows

# Windows'a taşıyıp orada çalıştırın
```

---

## 📝 Notlar

- WASPWARE, dış bağımlılıklar olmadan **tek çalıştırılabilir** bir dosyadır
- Go runtime gerekli değildir (binary'de gömülü)
- Uyumlu:
  - Ubuntu, Debian, Mint
  - CentOS, RHEL, Fedora
  - Arch Linux, Manjaro
  - Raspberry Pi OS
  - Alpine Linux
  - Go yüklü herhangi bir Linux dağıtımı

---

## 🎯 Hızlı Başlangıç Senaryoları

### Linux Derleme Senaryosu
`build-linux.sh` olarak kaydedin:
```bash
#!/bin/bash
# WASPWARE Linux Derleme Senaryosu

echo "🔨 WASPWARE için Linux'ta derleniyor..."

cd waspWARE

# Optimizasyon ile derle
go build -ldflags="-s -w" -o waspWARE ./waspWARE/waspWARE.go

if [ $? -eq 0 ]; then
    echo "✅ Derleme başarılı!"
    echo "📦 Binary: waspWARE"
    echo "📊 Boyut: $(stat -c%s waspWARE) bytes"
    
    # Desteklenen mimarileri göster
    echo ""
    echo "🎯 Desteklenen mimariler:"
    echo "  - amd64 (x86_64)"
    echo "  - arm64 (ARM 64-bit)"
    echo "  - armv7 (ARM 32-bit)"
    echo "  - armv8 (ARM 64-bit varyantı)"
    echo "  - 386 (x86 32-bit)"
else
    echo "❌ Derleme başarısız!"
    exit 1
fi
```

**Çalıştır:**
```bash
chmod +x build-linux.sh
./build-linux.sh
```

### Linux'tan Windows Derleme Senaryosu
`build-windows-from-linux.sh` olarak kaydedin:
```bash
#!/bin/bash
# WASPWARE Windows Derleme Senaryosu (Linux'tan)

echo "🪟 WASPWARE için Windows'ta derleniyor..."

cd waspWARE

# Windows amd64 için derle
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" \
  -o waspWARE.exe ./waspWARE/waspWARE.go

if [ $? -eq 0 ]; then
    echo "✅ Windows derlemesi başarılı!"
    echo "📦 Binary: waspWARE.exe"
    echo "📊 Boyut: $(stat -c%s waspWARE.exe) bytes"
    
    # Windows çalıştırılabilir dosyası olduğunu doğrula
    echo ""
    echo "🔍 Dosya tipi:"
    file waspWARE.exe
    
    # Desteklenen Windows mimarilerini göster
    echo ""
    echo "🎯 Desteklenen Windows mimarileri:"
    echo "  - amd64 (Windows x86_64)"
    echo "  - arm64 (Windows on ARM)"
    echo "  - 386 (Eski Windows)"
else
    echo "❌ Derleme başarısız!"
    exit 1
fi
```

**Çalıştır:**
```bash
chmod +x build-windows-from-linux.sh
./build-windows-from-linux.sh
```

### Çoklu-Platform Derleme Senaryosu
`build-all-platforms.sh` olarak kaydedin:
```bash
#!/bin/bash
# WASPWARE Çoklu-Platform Derleme Senaryosu

set -e

BINARY_DIR="bin"
mkdir -p "$BINARY_DIR"

echo "🔨 WASPWARE için tüm platformlarda derleniyor..."

# Linux Platformları
echo ""
echo "🐧 Linux platformları için derleniyor..."

GOOS=linux GOARCH=amd64 go build -o "$BINARY_DIR/waspWARE-amd64" ./waspWARE/waspWARE.go
GOOS=linux GOARCH=arm64 go build -o "$BINARY_DIR/waspWARE-arm64" ./waspWARE/waspWARE.go
GOOS=linux GOARCH=arm GOARM=7 go build -o "$BINARY_DIR/waspWARE-armhf" ./waspWARE/waspWARE.go
GOOS=linux GOARCH=386 go build -o "$BINARY_DIR/waspWARE-386" ./waspWARE/waspWARE.go

# Windows Platformları
echo ""
echo "🪟 Windows platformları için derleniyor..."

GOOS=windows GOARCH=amd64 go build -o "$BINARY_DIR/waspWARE-windows-amd64.exe" ./waspWARE/waspWARE.go
GOOS=windows GOARCH=arm64 go build -o "$BINARY_DIR/waspWARE-windows-arm64.exe" ./waspWARE/waspWARE.go

echo ""
echo "✅ Tüm derlemeler tamamlandı!"
ls -lh "$BINARY_DIR/"
```

**Çalıştır:**
```bash
chmod +x build-all-platforms.sh
./build-all-platforms.sh
```

---

## 📚 Ek Kaynaklar

- [Go Çapraz Derleme](https://go.dev/doc/install/source#environment)
- [Linux Derleme Sorunları](https://github.com/golang/go/wiki/LinuxBuildIssues)
- [ARM Çapraz derleme](https://go.dev/doc/install/source#environment)
- [Docker ile Go](https://docs.docker.com/language/golang/)
