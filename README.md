# WASPWARE - File Encryption Tool

![WASPWARE Logo](https://github.com/waspthebughunter/waspWARE/assets/100480448/d562bb91-c0be-4ce8-89b4-1b3aef901f13)
![WASPWARE Screenshot](https://github.com/waspthebughunter/waspWARE/assets/100480448/b1351478-cc41-4813-bc94-e921ae9cfab2)

> **Disclaimer:** This tool is for educational purposes only. I am not responsible for your harmful actions!

---

## 📋 Overview

WASPWARE is a powerful file encryption utility written in Go that uses AES-256-GCM encryption to securely encrypt files and directories. All encrypted files are automatically renamed with a `.wasp` extension for easy identification.

### Key Features

- 🔐 **AES-256-GCM Encryption** - Industry-standard encryption algorithm
- 📁 **Directory-wide Encryption** - Encrypt entire folders recursively
- 🔑 **Secure Key Management** - Optional password-based encryption keys
- 🔄 **Automatic Extension Change** - Encrypted files get `.wasp` extension
- 💾 **Single Executable** - No external dependencies required
- 🎯 **Cross-Platform** - Build for Linux, Windows, and multiple architectures

---

## 🚀 Quick Start

### Basic Usage

```bash
# Compile the tool
go build -o waspWARE ./waspWARE/waspWARE.go

# Run interactively
./waspWARE

# Or use command-line flags
./waspWARE -key "your-secret-key" -dizin "/path/to/encrypt"
```

### Interactive Mode

When you run WASPWARE without arguments, it will:
1. Prompt you for an encryption key (or generate a random one)
2. Ask for confirmation
3. Request the target directory
4. Encrypt all files in that directory

---

## 🔧 Command-Line Options

```bash
./waspWARE -key "your-key" -dizin "/target/directory"
```

- `-key`: Encryption key (optional - can be generated randomly)
- `-dizin`: Target directory to encrypt (optional - defaults to current directory)

---

## 📖 Documentation

### Build Guides

- **[Linux Build Guide](./build-linux.md)** - Comprehensive instructions for building on Linux platforms, including cross-compilation for different architectures.
- **[Windows Build Guide](./build-windows.md)** - Step-by-step instructions for building on Windows and cross-compiling from Linux.

### Supported Platforms

| Platform | Architecture | Binary Size |
|----------|-------------|-------------|
| Linux | amd64 (x86_64) | ~3.0 MB |
| Linux | arm64 (ARM 64-bit) | ~2.8 MB |
| Linux | armv7 (ARM 32-bit) | ~2.8 MB |
| Linux | armv8 (ARM 64-bit) | ~2.8 MB |
| Linux | 386 (x86 32-bit) | ~2.9 MB |
| Windows | amd64 | ~3.0 MB |
| Windows | arm64 | ~2.8 MB |
| Windows | 386 | ~2.9 MB |

### Compatible Distributions

- Ubuntu, Debian, Mint
- CentOS, RHEL, Fedora
- Arch Linux, Manjaro
- Raspberry Pi OS
- Alpine Linux
- Any Linux distribution with Go installed

---

## 🔐 How It Works

1. **Key Generation**: You provide an encryption key, or WASPWARE generates a random one
2. **Directory Scan**: The tool recursively scans the target directory
3. **AES Encryption**: Each file is encrypted using AES-256-GCM
4. **Extension Change**: Encrypted files are renamed with `.wasp` extension
5. **Completion**: All files are encrypted and ready for secure storage

---

## 📝 Usage Examples

### Encrypt a Specific Directory

```bash
./waspWARE -key "MySecretKey123" -dizin "/home/user/documents"
```

### Generate Random Key

```bash
./waspWARE  # Press Enter without entering a key for random generation
```

### Encrypt Current Directory

```bash
./waspWARE -key "mykey"  # Will encrypt current directory if no path provided
```

---

## 🔒 Security Considerations

- **AES-256-GCM** provides authenticated encryption with associated data (AEAD)
- **Random Nonce**: Each encryption uses a cryptographically secure random nonce
- **Key Management**: Store your encryption keys securely - they cannot be recovered if lost
- **File Permissions**: Encrypted files maintain their original permissions

---

## 🛠️ Build Instructions

### Native Build (Current Platform)

```bash
go build -ldflags="-s -w" -o waspWARE ./waspWARE/waspWARE.go
```

### Cross-Compilation Examples

#### Linux for ARM64
```bash
GOOS=linux GOARCH=arm64 go build -o waspWARE-arm64 ./waspWARE/waspWARE.go
```

#### Windows Executable from Linux
```bash
GOOS=windows GOARCH=amd64 go build -o waspWARE.exe ./waspWARE/waspWARE.go
```

### Docker Build

```bash
docker run --rm -v $(pwd):/app golang:1.21-alpine \
  go build -o /app/waspWARE ./waspWARE/waspWARE.go
```

---

## 📚 Additional Resources

- [Go Cross-compilation Documentation](https://go.dev/doc/install/source#environment)
- [AES Encryption Best Practices](https://csrc.nist.gov/publications/detail/sp/800-38a/final)
- [Go Project Layout](https://github.com/golang-standards/project-layout)

---

## ⚠️ Important Notes

- **Backup Your Keys**: Encryption keys cannot be recovered if lost
- **Test Before Production**: Always test encryption on a small subset first
- **Key Storage**: Store keys in secure locations (password managers, encrypted files)
- **File Integrity**: Encrypted files maintain their original permissions and metadata

---

## 📄 License

This project is provided for educational purposes. Use responsibly and legally.

---

## 🙏 Acknowledgments

- Go Language Team
- AES Encryption Standards (NIST SP 800-38A)
- The open-source community

---

<div align="center">
  <strong>WASPWARE - Secure File Encryption Made Simple</strong>
</div>
