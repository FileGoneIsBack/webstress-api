# Simple Go HTTPS Handler

## Overview

This is a simple webstresser in golang using my https handle and serve (https://github.com/FileGoneIsBack/Golang-HTTPS-Handler)
this uses api and raw server connections 


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
sudo apt install snap
```

```
sudo snap install go --channel=1.21/stable --classic
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
git clone https://github.com/FileGoneIsBack/webstress-api
cd webstress-api
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

# connecting a server (updateV3.1)

1. download client files from assets/client
2. install client files to server
3. install golang to server (follow parts of tut before...)
4. go build . 

5a. Client setup
- make sure key matches website server.json key value
- make sure methods are correct w the site 

5b. Website setup
- make sure key is secure password 
- white list servers to connect in servers.json

6. ./client and they should connect, the server should send the methdos to the website to confirm matching funnel


scroll for latest!


====
i cannot give a date VERS 3.2
also added more admin options, fixed panel, updated error logs/sending notis, and much more ;)

Note: I've never used this source; this is something I made in my free time. I don't partake in these shitty activities, but use at your own risk and legally on your own networks.


change logs simpler then git logs
====
4/1/25 vers 3.5
-fix auto pay changed to https://nowpayments.io simply make an account go to settings and get api key and add cryptos
--while adding cryptos since ive only added like 2 your going to need to replace 
```
	var coinNameMap = map[string]string{
		"btc": "bitcoin",
		"eth": "ethereum",
	}
``` 
in core/master/api/payments/transaction.go

-also added notis used with
--sessions.SetFlash(w, r, "message", "from")

to come...

-im looking to fix auth tokens and tg bot to log messages by users and log who uses what token to keep the site secure!
-finish admin panel
-fix L7 chart
i been lazy and update this every few months its already passed a year old might make a new src but also been working on other shit.
