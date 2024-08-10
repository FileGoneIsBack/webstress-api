# Simple Go HTTPS Handler

## Overview

This is a simple webstresser in golang using my https handle and serve (https://github.com/FileGoneIsBack/Golang-HTTPS-Handler)
this uses api and raw server connections (i will not be dropping the client files uncompiled thats something your gonna need to trust or remake)

## Prerequisites

1. **Go Programming Language**: You need Go installed to build and run the application.
2. **GCC Compiler**: Required for building Go applications with cgo.
3. **Cloudflare Account**: Required to obtain SSL/TLS certificates for your domain.

## Installation

### 1. Install Go

To install Go, follow the instructions for your operating system:

- **Windows**: [Download Go Installer](https://golang.org/dl/) and follow the installation instructions.
- **macOS**: You can use Homebrew:
```
brew install go
```
- **Linux**
```
wget https://golang.org/dl/go1.x.x.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.x.x.linux-amd64.tar.gz
```

### 2. Install GCC
- **Windows**: Install MinGW or MSYS2.
- **macOS**: Install Xcode Command Line Tools:
```
xcode-select --install
```
- **Linux**: Install GCC using your package manager:
```
sudo apt-get install build-essential
```
### 3. Obtain Domain Certificates and Keys
- Log in to Cloudflare: Access your Cloudflare account or sign up if you don’t have one.
- Add Your Domain: Follow Cloudflare’s instructions to add your domain.
- Obtain Certificates: Navigate to the SSL/TLS section and get the domain certificate and key. Download them to your local machine.

### 4. Clone the repo
```
git clone https://github.com/your-username/your-repo-name.git
cd your-repo-name
```

### 5. Edit Files
edit the config.json and the cert/key files in the assets folder!

### 6. build the src
```
go run .
```
- sometimes you might need to enable cgo 
```
export CGO_ENABLED=1
```
- compile it fully 
```
go build api
./api
```
