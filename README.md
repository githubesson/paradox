![panel.png](https://files.catbox.moe/ebb0jf.png)

# Paradox - Golang macOS Stealer PoC

A proof-of-concept implementation demonstrating how macOS stealers function. This project serves to illustrate common techniques and behaviors employed by macOS malware for research and defensive purposes.

## Legal Notice

This software is provided for educational and research purposes only. The author does not endorse or encourage any malicious use of this code. Users are responsible for complying with all applicable laws and regulations. 

## ⚠️ Disclaimer

This is a **proof-of-concept only**. It was created for educational and research purposes to help understand malware behavior and improve defensive measures. Do not use this code for malicious purposes. The author assumes no liability for any misuse of this code.

## Educational Purpose

The code in this repository helps security researchers and defenders:
- Understand common macOS stealer techniques
- Study malware behavior patterns
- Develop better detection and prevention methods
- Learn about macOS security mechanisms

## Requirements

### Server
- Go 1.24.1 or higher
- `go.mod` (dependencies will be automatically installed when building)
- Python 3.x

### Client
- Go 1.24.1 or higher
- `go.mod` (dependencies will be automatically installed when building)

### Frontend
- nodejs + npm (author used node v23.10.0 and npm 10.9.2)

## Installation & Setup
1. Clone the repository:
   ```bash
   git clone https://github.com/githubesson/paradox
   cd paradox
   ```

2. Build the server:
   ```bash
   cd server
   go build -o paradox-server .
   ```

3. Setup the frontend panel:
   ```bash
   cd ../frontend
   npm i
   ```

## Usage

1. Start the server:
   ```bash
   cd server
   ./paradox-server
   ```

2. The server will listen on 127.0.0.1:8080.

3. Start the web panel:
   ```bash
   cd ../frontend
   npm run dev
   ```

4. Register an account in the web panel (127.0.0.1:3000/register)

5. Build a payload via the web panel

6. Run the client

## Project Structure

```
.
├── built/      # Clients built by the server
├── server/     # Go server implementation
├── frontend/   # Frontend panel
└── payload/    # Clientside code
```

Note: This is a proof-of-concept implementation. Some features may be incomplete or require additional configuration.

## Contact

- Twitter - https://x.com/7N7
- Others - https://misleadi.ng