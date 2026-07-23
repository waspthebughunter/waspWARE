# WASPWARE - Windows Build Guide

## 📋 Overview
This guide provides step-by-step instructions for building WASPWARE on Windows platforms, including cross-compilation from Linux systems.

---

## 🔧 Prerequisites

### Option 1: Install Go on Windows (Native Build)
```powershell
# Download and install Go from official site
winget install Golang.Golang

# OR download from: https://go.dev/dl/
```

### Option 2: Cross-Compile from Linux
You can build Windows executables from a Linux system using the `GOOS` environment variable. See [build-linux.md](./build-linux.md) for detailed instructions.

### 1. Verify Go Installation (Windows)
```powershell
go version
# Expected output: go version go1.21.x windows/amd64
```

### 2. Set Environment Variables (Optional)
```powershell
# Add Go to PATH if not already present
$env:Path += ";C:\Go\bin"

# Verify
echo $env:PATH
```

---

## 🏗️ Basic Build Commands

### Simple Build
```powershell
cd waspWARE
go build -o waspWARE.exe ./waspWARE/waspWARE.go
```

### Build with Release Flags (Optimized)
```powershell
go build -ldflags="-s -w" -o waspWARE.exe ./waspWARE/waspWARE.go
```

### Build with Version Info
```powershell
go build -ldflags="-s -w -X 'main.Version=1.0.0'" \
  -o waspWARE.exe ./waspWARE/waspWARE.go
```

---

## 🪟 Cross-Compilation from Linux to Windows

### Build Windows Executable on Linux

From a Linux system, you can cross-compile WASPWARE for Windows:

```bash
# Build for Windows amd64 (standard Windows 10/11)
GOOS=windows GOARCH=amd64 go build -o waspWARE.exe ./waspWARE/waspWARE.go

# Verify the executable
file waspWARE.exe
# Output: waspWARE.exe: PE32+ executable (console) x86-64, for MS Windows
```

### Build for Different Windows Architectures from Linux

##### Windows on ARM (Surface Pro X, newer devices)
```bash
GOOS=windows GOARCH=arm64 go build -o waspWARE-arm64.exe ./waspWARE/waspWARE.go
```

##### 32-bit Windows (Legacy Systems)
```bash
GOOS=windows GOARCH=386 go build -o waspWARE-386.exe ./waspWARE/waspWARE.go
```

---

## 🐳 Docker Build for Windows

### Using Multi-Stage Dockerfile
```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY waspWARE/waspWARE.go .

# Build arguments for cross-compilation
ARG TARGETOS=windows
ARG TARGETARCH=amd64

RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -ldflags="-s -w" -o waspWARE.exe .

FROM alpine:latest
RUN apk add --no-cache winpty
COPY --from=builder /app/waspWARE.exe /waspWARE.exe
ENTRYPOINT ["winpty", "/waspWARE.exe"]
```

Build and run:
```bash
docker build -t waspware-windows .
docker run --rm -v $(pwd):/app waspware-windows
```

---

## 🎯 Cross-Compilation for Different Windows Platforms

### Build for ARM64 (Windows on ARM, Surface Pro X)
```powershell
GOOS=windows GOARCH=arm64 go build -o waspWARE-arm64.exe ./waspWARE/waspWARE.go
```

### Build for 32-bit Windows (Legacy Systems)
```powershell
GOOS=windows GOARCH=386 go build -o waspWARE-386.exe ./waspWARE/waspWARE.go
```

---

## 🐳 Docker Build (Cross-Platform)

### Using Docker Desktop
```powershell
# Create Dockerfile
New-Item -Path "Dockerfile" -Force

# Add content to Dockerfile
@"
FROM golang:1.21-alpine
WORKDIR /app
COPY waspWARE/waspWARE.go .
RUN GOOS=windows GOARCH=amd64 go build -o waspWARE.exe .
"@ | Out-File -FilePath "Dockerfile"

# Build
docker build -t waspware-builder .

# Run
docker run --rm -v "${PWD}":/app waspware-builder
```

---

## ✅ Verification Steps

### 1. Check Binary Exists
```powershell
Get-Item waspWARE.exe | Select-Object Name, Length, LastWriteTime
```

### 2. Test Execution
```powershell
.\waspWARE.exe -help
# OR run interactively
.\waspWARE.exe
```

### 3. Verify File Properties
```powershell
Get-Item waspWARE.exe | Format-List Name, Length, LastWriteTime, Mode
```

### 4. Verify Windows Executable Format
```powershell
# Check if it's a valid PE executable
$file = Get-Item waspWARE.exe
$file | Select-Object Name, @{Name='Size(MB)';Expression={[math]::Round($file.Length/1MB,2)}}

# Using System.Diagnostics.FileVersionInfo
$versionInfo = [System.Diagnostics.FileVersionInfo]::GetVersionInfo("$PWD/waspWARE.exe")
$versionInfo | Select-Object FileName, FileVersion, FileType
```

---

## 🪟 Cross-Compilation from Linux Verification

### Verify Windows Executable Built on Linux
```bash
# Check file type
file waspWARE.exe
# Expected: PE32+ executable (console) x86-64, for MS Windows

# Check binary size
ls -lh waspWARE.exe

# Verify PE header
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

## 📦 Distribution Package

### Create ZIP Archive
```powershell
Compress-Archive -Path waspWARE.exe,README.md -DestinationPath waspWARE-windows.zip
```

### Create Installer (Optional)
```powershell
# Using Inno Setup (download from https://jrsoftware.org/isinfo.php)
# Create setup.iss script and compile with iscc
```

---

## 🔧 Troubleshooting

### Issue: "go: cannot find module"
**Solution:** Ensure you're in the correct directory
```powershell
cd waspWARE
go build -o waspWARE.exe ./waspWARE/waspWARE.go
```

### Issue: Build fails with permission error
**Solution:** Run PowerShell as Administrator
```powershell
Start-Process powershell -ArgumentList "cd waspWARE; go build -o waspWARE.exe" -Verb RunAs
```

### Issue: Missing Go installation
**Solution:** Install Go from https://go.dev/dl/
```powershell
# Verify installation
go version
```

---

## 📚 Additional Resources

- [Go Documentation](https://golang.org/doc/)
- [Cross-compilation Guide](https://go.dev/doc/install/source#environment)
- [Windows Build Issues](https://github.com/golang/go/wiki/WindowsBuildIssues)

---

## 🎯 Quick Start Script

Save as `build-windows.ps1`:
```powershell
#!/usr/bin/env pwsh
# WASPWARE Windows Build Script

Write-Host "🔨 Building WASPWARE for Windows..." -ForegroundColor Cyan

# Navigate to project
cd waspWARE

# Build with optimization
go build -ldflags="-s -w" -o waspWARE.exe ./waspWARE/waspWARE.go

if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ Build successful!" -ForegroundColor Green
    Write-Host "📦 Binary: waspWARE.exe" -ForegroundColor Yellow
    Write-Host "📊 Size: $(Get-Item waspWARE.exe).Length bytes" -ForegroundColor Yellow
} else {
    Write-Host "❌ Build failed!" -ForegroundColor Red
    exit 1
}
```

**Run:**
```powershell
.\build-windows.ps1
```

---

## 📝 Notes

- WASPWARE is a **single executable** with no external dependencies
- Binary size: ~2.8 MB (optimized) to ~3.0 MB (standard)
- Compatible with Windows 10/11 and Windows Server
- No additional runtime required (Go runtime embedded in binary)
