# WASPWARE - Linux Build Guide

## 📋 Overview
This guide provides comprehensive instructions for building WASPWARE on various Linux platforms, including cross-compilation for different architectures.

---

## 🔧 Prerequisites

### 1. Install Go Language

#### Ubuntu/Debian
```bash
# Add Go repository
wget https://go.dev/release/go1.21.3.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.3.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc

# Verify
go version
```

#### CentOS/RHEL/Fedora
```bash
sudo yum install -y go
# OR
sudo dnf install -y go
```

#### Arch Linux
```bash
sudo pacman -S go
```

### 2. Install Build Tools (Optional)
```bash
# For native cross-compilation
sudo apt-get install gcc musl-dev   # Debian/Ubuntu
sudo yum install gcc musl-devel      # CentOS/RHEL
```

---

## 🏗️ Basic Build Commands

### Simple Build (Current Platform)
```bash
cd waspWARE
go build -o waspWARE ./waspWARE/waspWARE.go
```

### Build with Release Flags (Optimized)
```bash
go build -ldflags="-s -w" -o waspWARE ./waspWARE/waspWARE.go
```

### Build with Version Info
```bash
go build -ldflags="-s -w -X 'main.Version=1.0.0'" \
  -o waspWARE ./waspWARE/waspWARE.go
```

---

## 🎯 Cross-Compilation for Different Linux Platforms

### amd64/x86_64 (Standard Desktop)
```bash
GOOS=linux GOARCH=amd64 go build -o waspWARE-amd64 ./waspWARE/waspWARE.go
```

### ARM64 (Apple Silicon, Raspberry Pi 4/5, AWS Graviton)
```bash
GOOS=linux GOARCH=arm64 go build -o waspWARE-arm64 ./waspWARE/waspWARE.go
```

### ARMv7 (Raspberry Pi 3/4, older ARM devices)
```bash
GOOS=linux GOARCH=arm GOARM=7 go build -o waspWARE-armhf ./waspWARE/waspWARE.go
```

### ARMv8 (Raspberry Pi 4/5, newer ARM devices)
```bash
GOOS=linux GOARCH=arm GOARM=8 go build -o waspWARE-armv8 ./waspWARE/waspWARE.go
```

### 32-bit x86 (Legacy Systems)
```bash
GOOS=linux GOARCH=386 go build -o waspWARE-386 ./waspWARE/waspWARE.go
```

---

## 🪟 Cross-Compilation to Windows (.exe) from Linux

### Build Windows Executable on Linux

You can cross-compile WASPWARE for Windows from a Linux system using the `GOOS` environment variable:

#### Basic Windows Build (amd64)
```bash
# Build for Windows amd64 (standard Windows 10/11)
GOOS=windows GOARCH=amd64 go build -o waspWARE.exe ./waspWARE/waspWARE.go

# Verify the executable
file waspWARE.exe
# Output: waspWARE.exe: PE32+ executable (console) x86-64, for MS Windows
```

#### Build for Different Windows Architectures

##### Windows on ARM (Surface Pro X, newer devices)
```bash
GOOS=windows GOARCH=arm64 go build -o waspWARE-arm64.exe ./waspWARE/waspWARE.go
```

##### 32-bit Windows (Legacy Systems)
```bash
GOOS=windows GOARCH=386 go build -o waspWARE-386.exe ./waspWARE/waspWARE.go
```

##### Windows ARM (Windows on ARM devices)
```bash
GOOS=windows GOARCH=arm go build -o waspWARE-arm.exe ./waspWARE/waspWARE.go
```

#### Optimized Windows Build
```bash
# Build with release flags for smaller binary
GOOS=windows GOARCH=amd64 go build \
  -ldflags="-s -w" \
  -o waspWARE.exe ./waspWARE/waspWARE.go

# Verify size and properties
ls -lh waspWARE.exe
file waspWARE.exe
```

---

## 🐳 Docker Build (Recommended for Cross-Compilation)

### Using Official Go Image
```bash
# Pull image
docker pull golang:1.21-alpine

# Build for current platform
docker run --rm -v $(pwd):/app golang:1.21-alpine \
  go build -o /app/waspWARE ./waspWARE/waspWARE.go

# Build for ARM64
docker run --rm -v $(pwd):/app golang:1.21-alpine \
  GOOS=linux GOARCH=arm64 go build -o /app/waspWARE-arm64 ./waspWARE/waspWARE.go

# Build for ARMv7
docker run --rm -v $(pwd):/app golang:1.21-alpine \
  GOOS=linux GOARCH=arm GOARM=7 go build -o /app/waspWARE-armhf ./waspWARE/waspWARE.go
```

### Using Multi-Stage Dockerfile for Windows

Create `Dockerfile`:
```dockerfile
# Stage 1: Build for Linux
FROM golang:1.21-alpine AS builder-linux

WORKDIR /app
COPY waspWARE/waspWARE.go .

RUN go build -ldflags="-s -w" -o waspWARE ./waspWARE/waspWARE.go

# Stage 2: Build for Windows
FROM golang:1.21-alpine AS builder-windows

WORKDIR /app
COPY waspWARE/waspWARE.go .

RUN GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o waspWARE.exe .

# Stage 3: Final image with Windows binary
FROM alpine:latest
RUN apk add --no-cache winpty

COPY --from=builder-windows /app/waspWARE.exe /waspWARE.exe

ENTRYPOINT ["winpty", "/waspWARE.exe"]
```

Build and run:
```bash
docker build -t waspware:latest .

# Run on Linux (will execute via winpty)
docker run --rm -v $(pwd):/app waspware:latest \
  -key "mykey" -dizin "/target/directory"

# Or copy Windows binary to another system
docker run --rm -v $(pwd)/bin:/output waspware:latest \
  cp /waspWARE.exe /output/waspWARE-windows-amd64.exe
```

---

### Docker Multi-Architecture Build

Create `build-all-platforms.sh`:
```bash
#!/bin/bash
# WASPWARE Multi-Platform Build Script (Linux + Windows)

set -e

BINARY_DIR="bin"
mkdir -p "$BINARY_DIR"

echo "🔨 Building WASPWARE for all platforms..."

# Linux Platforms
echo ""
echo "🐧 Building for Linux platforms..."

GOOS=linux GOARCH=amd64 go build -o "$BINARY_DIR/waspWARE-amd64" ./waspWARE/waspWARE.go
GOOS=linux GOARCH=arm64 go build -o "$BINARY_DIR/waspWARE-arm64" ./waspWARE/waspWARE.go
GOOS=linux GOARCH=arm GOARM=7 go build -o "$BINARY_DIR/waspWARE-armhf" ./waspWARE/waspWARE.go
GOOS=linux GOARCH=386 go build -o "$BINARY_DIR/waspWARE-386" ./waspWARE/waspWARE.go

# Windows Platforms
echo ""
echo "🪟 Building for Windows platforms..."

GOOS=windows GOARCH=amd64 go build -o "$BINARY_DIR/waspWARE-windows-amd64.exe" ./waspWARE/waspWARE.go
GOOS=windows GOARCH=arm64 go build -o "$BINARY_DIR/waspWARE-windows-arm64.exe" ./waspWARE/waspWARE.go

echo ""
echo "✅ All builds completed!"
ls -lh "$BINARY_DIR/"
```

Make executable and run:
```bash
chmod +x build-all-platforms.sh
./build-all-platforms.sh
```

---

### Docker Multi-Stage for All Platforms

Create `Dockerfile.multi`:
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

Build for specific platform:
```bash
# Linux amd64
docker build --target linux -t waspware-linux .

# Windows amd64
docker build --target windows -t waspware-windows .
```

---

## 📦 Build All Platforms Script

Create `build-all.sh`:
```bash
#!/bin/bash
# WASPWARE Multi-Platform Build Script

set -e

BINARY_DIR="bin"
mkdir -p "$BINARY_DIR"

echo "🔨 Building WASPWARE for all platforms..."

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

echo "✅ All builds completed!"
ls -lh "$BINARY_DIR/"
```

Make executable and run:
```bash
chmod +x build-all.sh
./build-all.sh
```

---

## ✅ Verification Steps

### 1. Check Binary Exists
```bash
ls -lh waspWARE*
```

### 2. Test Execution
```bash
./waspWARE -help
# OR run interactively
./waspWARE
```

### 3. Verify File Properties
```bash
file waspWARE
stat waspWARE
```

### 4. Test Encryption
```bash
mkdir -p /tmp/test-encryption
echo "Test file content" > /tmp/test-encryption/test.txt
./waspWARE -key "test123" -dizin /tmp/test-encryption
ls -la /tmp/test-encryption/
```

---

## 📊 Binary Sizes by Architecture

| Platform | Architecture | Binary Size | Use Case |
|----------|-------------|-------------|----------|
| Linux | amd64 | ~3.0 MB | Standard x86_64 Linux |
| Linux | arm64 | ~2.8 MB | Apple Silicon, Raspberry Pi 4/5 |
| Linux | armv7 | ~2.8 MB | Raspberry Pi 3/4 (32-bit) |
| Linux | armv8 | ~2.8 MB | Raspberry Pi 4/5 (64-bit) |
| Linux | 386 | ~2.9 MB | Legacy 32-bit systems |
| Windows | amd64 | ~3.0 MB | Windows 10/11 x86_64 |
| Windows | arm64 | ~2.8 MB | Windows on ARM devices |
| Windows | 386 | ~2.9 MB | Legacy 32-bit Windows |

---

## 🪟 Windows Build Verification

### Verify Windows Executable
```bash
# Check file type
file waspWARE.exe
# Output: waspWARE.exe: PE32+ executable (console) x86-64, for MS Windows

# Check binary size
ls -lh waspWARE.exe

# Verify it's a valid Windows PE executable
readelf -h waspWARE.exe 2>/dev/null || objdump -f waspWARE.exe
```

### Transfer to Windows
```bash
# Copy to Windows machine via SCP/SFTP
scp waspWARE.exe user@windows-machine:/path/to/

# Or use USB drive
cp waspWARE.exe /media/usb/waspWARE.exe

# Or email the binary
# (Windows executable can be attached and sent)
```

### Run on Windows
Once transferred to Windows:
```powershell
# PowerShell
.\waspWARE.exe -key "mykey" -dizin "C:\target\directory"

# Command Prompt
waspWARE.exe -key "mykey" -dizin "C:\target\directory"
```

---

## 🔧 Troubleshooting

### Issue: "go: cannot find module"
**Solution:** Ensure you're in the correct directory
```bash
cd waspWARE
go build -o waspWARE ./waspWARE/waspWARE.go
```

### Issue: Cross-compilation fails with "unknown architecture"
**Solution:** Install required headers
```bash
# For ARM64 cross-compilation
sudo apt-get install gcc-aarch64-linux-gnu
export CC=aarch64-linux-gnu-gcc

GOOS=linux GOARCH=arm64 CC=aarch64-linux-gnu-gcc \
  go build -o waspWARE-arm64 ./waspWARE/waspWARE.go
```

### Issue: Build fails with permission error
**Solution:** Use sudo or adjust permissions
```bash
sudo chown -R $USER:$USER waspWARE
go build -o waspWARE ./waspWARE/waspWARE.go
```

### Issue: Missing Go installation
**Solution:** Install Go from official source
```bash
# Ubuntu/Debian
curl -LO https://go.dev/dl/go1.21.3.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.3.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
```

### Issue: Windows build fails with "unknown architecture"
**Solution:** Install required cross-compilation tools
```bash
# Install GCC for Windows target (optional, Go can cross-compile without it)
sudo apt-get install gcc-multilib

# For ARM64 Windows builds
GOOS=windows GOARCH=arm64 go build -o waspWARE-arm64.exe ./waspWARE/waspWARE.go
```

### Issue: Generated .exe is not recognized on Windows
**Solution:** Ensure you're building for the correct platform
```bash
# Verify build target
echo "Building for: GOOS=$GOOS GOARCH=$GOARCH"

# Clean and rebuild
rm -f *.exe
GOOS=windows GOARCH=amd64 go build -o waspWARE.exe ./waspWARE/waspWARE.go

# Verify
file waspWARE.exe
```

### Issue: Permission denied when running .exe on Linux
**Solution:** This is expected - Windows executables can't run on Linux directly. Transfer to Windows system.
```bash
# Check file type
file waspWARE.exe
# Should show: PE32+ executable for MS Windows

# Transfer to Windows and run there
```

---

## 📝 Notes

- WASPWARE is a **single executable** with no external dependencies
- No Go runtime required (embedded in binary)
- Compatible with:
  - Ubuntu, Debian, Mint
  - CentOS, RHEL, Fedora
  - Arch Linux, Manjaro
  - Raspberry Pi OS
  - Alpine Linux
  - Any Linux distribution with Go installed

---

## 🎯 Quick Start Scripts

### Linux Build Script
Save as `build-linux.sh`:
```bash
#!/bin/bash
# WASPWARE Linux Build Script

echo "🔨 Building WASPWARE for Linux..."

cd waspWARE

# Build with optimization
go build -ldflags="-s -w" -o waspWARE ./waspWARE/waspWARE.go

if [ $? -eq 0 ]; then
    echo "✅ Build successful!"
    echo "📦 Binary: waspWARE"
    echo "📊 Size: $(stat -c%s waspWARE) bytes"
    
    # Show supported architectures
    echo ""
    echo "🎯 Supported architectures:"
    echo "  - amd64 (x86_64)"
    echo "  - arm64 (ARM 64-bit)"
    echo "  - armv7 (ARM 32-bit)"
    echo "  - armv8 (ARM 64-bit variant)"
    echo "  - 386 (x86 32-bit)"
else
    echo "❌ Build failed!"
    exit 1
fi
```

**Run:**
```bash
chmod +x build-linux.sh
./build-linux.sh
```

### Windows Build Script from Linux
Save as `build-windows-from-linux.sh`:
```bash
#!/bin/bash
# WASPWARE Windows Build Script (from Linux)

echo "🪟 Building WASPWARE for Windows..."

cd waspWARE

# Build for Windows amd64
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" \
  -o waspWARE.exe ./waspWARE/waspWARE.go

if [ $? -eq 0 ]; then
    echo "✅ Windows build successful!"
    echo "📦 Binary: waspWARE.exe"
    echo "📊 Size: $(stat -c%s waspWARE.exe) bytes"
    
    # Verify it's a Windows executable
    echo ""
    echo "🔍 File type:"
    file waspWARE.exe
    
    # Show supported Windows architectures
    echo ""
    echo "🎯 Supported Windows architectures:"
    echo "  - amd64 (Windows x86_64)"
    echo "  - arm64 (Windows on ARM)"
    echo "  - 386 (Legacy Windows)"
else
    echo "❌ Build failed!"
    exit 1
fi
```

**Run:**
```bash
chmod +x build-windows-from-linux.sh
./build-windows-from-linux.sh
```

### Multi-Platform Build Script
Save as `build-all-platforms.sh`:
```bash
#!/bin/bash
# WASPWARE Multi-Platform Build Script

set -e

BINARY_DIR="bin"
mkdir -p "$BINARY_DIR"

echo "🔨 Building WASPWARE for all platforms..."

# Linux Platforms
echo ""
echo "🐧 Building for Linux platforms..."

GOOS=linux GOARCH=amd64 go build -o "$BINARY_DIR/waspWARE-amd64" ./waspWARE/waspWARE.go
GOOS=linux GOARCH=arm64 go build -o "$BINARY_DIR/waspWARE-arm64" ./waspWARE/waspWARE.go
GOOS=linux GOARCH=arm GOARM=7 go build -o "$BINARY_DIR/waspWARE-armhf" ./waspWARE/waspWARE.go
GOOS=linux GOARCH=386 go build -o "$BINARY_DIR/waspWARE-386" ./waspWARE/waspWARE.go

# Windows Platforms
echo ""
echo "🪟 Building for Windows platforms..."

GOOS=windows GOARCH=amd64 go build -o "$BINARY_DIR/waspWARE-windows-amd64.exe" ./waspWARE/waspWARE.go
GOOS=windows GOARCH=arm64 go build -o "$BINARY_DIR/waspWARE-windows-arm64.exe" ./waspWARE/waspWARE.go

echo ""
echo "✅ All builds completed!"
ls -lh "$BINARY_DIR/"
```

**Run:**
```bash
chmod +x build-all-platforms.sh
./build-all-platforms.sh
```

---

## 📚 Additional Resources

- [Go Cross-compilation](https://go.dev/doc/install/source#environment)
- [Linux Build Issues](https://github.com/golang/go/wiki/LinuxBuildIssues)
- [ARM Cross-compilation](https://go.dev/doc/install/source#environment)
- [Docker with Go](https://docs.docker.com/language/golang/)
