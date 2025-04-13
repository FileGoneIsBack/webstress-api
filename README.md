# Project Setup Guide

This project includes a web-based dashboard, Telegram bot, and client app. The system manages servers, methods, and API communications with optional CNC and attack handling logic.

---

## Requirements

- **Go** 1.22+
- **MySQL** or **MariaDB**
- **SQLite3** (for Windows or development use)
- A **VPS server** (e.g., from [OVH](https://ovh.com))
- A **Domain name** (configured with [Cloudflare](https://cloudflare.com))

### Recommended VSCode Extensions

| Extension                                                                                         | Description            |
| ------------------------------------------------------------------------------------------------- | ---------------------- |
| [Remote SSH](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-ssh)     | SSH access to server   |
| [DB Client JDBC](https://marketplace.visualstudio.com/items?itemName=cweijan.dbclient-jdbc)       | Manage MySQL/MariaDB   |
| [SQLite Editor](https://marketplace.visualstudio.com/items?itemName=yy0931.vscode-sqlite3-editor) | Edit SQLite3 databases |

---

## Installation

```bash
# Install Go
sudo snap install go --classic   # or use https://go.dev/dl/

# Install MySQL/MariaDB (Ubuntu/Debian)
sudo apt update
sudo apt install mysql-server    # or mariadb-server

# (Optional for Windows)
# Install SQLite3 and use with VSCode extension
```

### Configuration

- Enable `cgo` for the web handler.
- Update `config.json`:
  - Set `"secure": true` to enable HTTPS (MySQL).
- Start the app:

```bash
go run .
```

---

## Creating an SSL Certificate

1. **Buy a domain** from Namecheap, Cloudflare, or OVH.
2. **Create a Cloudflare account.**
3. **Update your domain's nameservers** to point to Cloudflare.
4. **Update the DNS record**:
   - Add an `A` record pointing the domain to your VPS IP.
5. **Set the SSL/TLS encryption mode to Full**.
6. In Cloudflare, go to **SSL/TLS > Origin Server**:
   - Click **Create Certificate**.
   - Save the generated certificate and key as:
     - `assets/cert.pem`
     - `assets/key.pem`
7. Set `"secure": true` in `config.json`.

Your site is now using HTTPS.

---

## Telegram Bot Setup

- Run the bot on another server or outside the main project directory.
- Configure the bot using `config.json`.
- Use a live MySQL server (requires `"secure": true`).
- Start the bot normally.

---

## Client App Setup

1. Edit `server.json` in the website project:
   - Add server IPs under the `"allowed"` list.
   - Set the `"key"` for authentication.
2. Edit the client `config.json`:
   - Ensure the `key` matches the server.
   - Set the correct `master` IP and port.
3. Modify method command entries as needed.
4. Build the client:

```bash
go build .
```

5. Move the compiled binary to the method directory and run it.

---

## JSON Configuration Overview

This system communicates between the web server and client servers using structured JSON files.

### Website Configuration

#### `config.json`

| Key        | Description                               |
| ---------- | ----------------------------------------- |
| `secure`   | Enables HTTPS on port 443 or HTTP on 8080 |
| `domain`   | Used for CNC attack handling              |
| `cert/key` | SSL certificate and key from Cloudflare   |
| `autobuy`  | Sellix key/email for purchases            |
| `fake`     | Fake users or attack display              |
| `cnc`      | CNC port group configuration              |
| `methods`  | Method group and behavior definitions     |

#### `methods.json`

```json
"HOLD": {
  "type": 1,
  "name": "HOLD (Amplification)",
  "description": "",
  "subnet": 1,
  "mtype": 1,
  "vip": false
}
```

- `type`: Numeric identifier
- `subnet`: Group reference from `config.json`
- `mtype`: `1` for Layer 4, `2` for Layer 7
- `vip`: true/false flag

#### `server.json`

```json
{
  "listner": 1234,
  "allowed": ["1.2.3.4"],
  "key": "secretKeyHere"
}
```

### Client Configuration

#### `config.json`

```json
{
  "name": "L4-Server-01",
  "slots": 100,
  "type": 1,
  "key": "secretKeyHere",
  "master": "your.server.ip:1234",
  "maxthreads": 75
}
```

#### `methods.json`

```json
"HOLD": "./dns $target/32 $port lists/dns.list $threads $pps $time"
```

- All files and methods must be in the same directory.
- If you see `exit status 1`, the method did not screen properly, but the client will remain connected.

---

## Additional Tips

- Keep method names in uppercase in all JSON files.
- Ensure that the client and server `key` values match.
- Avoid nesting folders in the methods directory.
- Use a reverse proxy like NGINX for added control.

---

## Final Notes

You are now ready to deploy and manage your CNC-style setup. Keep your configuration files consistent, your client/server keys secure, and maintain good organization for your methods and commands.

---

> Built with Go, SQL, and secure networking in mind.

