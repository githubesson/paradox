# Paradox - Golang macOS Stealer PoC

A proof-of-concept implementation demonstrating how macOS stealers function, developed in 48 hours as an educational example. This project serves to illustrate common techniques and behaviors employed by macOS malware for research and defensive purposes.

## ⚠️ Disclaimer

This is a **proof-of-concept only**. It was created for educational and research purposes to help understand malware behavior and improve defensive measures. Do not use this code for malicious purposes. The author assumes no liability for any misuse of this code.

## Overview

This project was rapidly developed over a 48-hour period to demonstrate basic concepts of how macOS stealers typically operate. The implementation is intentionally rough around the edges to serve as a learning tool rather than a polished product.

## Educational Purpose

The code in this repository helps security researchers and defenders:
- Understand common macOS stealer techniques
- Study malware behavior patterns
- Develop better detection and prevention methods
- Learn about macOS security mechanisms

## Development Context

- Built in: Under 48 hours
- Platform: macOS
- Status: Proof of Concept / Educational Example

## Legal Notice

This software is provided for educational and research purposes only. The author does not endorse or encourage any malicious use of this code. Users are responsible for complying with all applicable laws and regulations. 

## Requirements

### Server
- Go 1.24.1 or higher
- `go.mod` (dependencies will be automatically installed when building)
- Python 3.x

### Client
- Go 1.24.1 or higher
- `go.mod` (dependencies will be automatically installed when building)

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

## Usage

1. Start the server:
   ```bash
   cd server
   ./paradox-server
   ```

2. The server will listen on 127.0.0.1:8080.

3. Build the client:
   ```bash
   curl http://localhost:8080/build
   ```
   If you set everything up correctly, the response should be along the lines of: {"build_id":"randombuildid (you will need this for the next step)","filename":"randomfilename","message":"Payload built successfully"}. Self built clients will not work due to build id validation on client upload.

4. Download the client:
   ```bash
   curl http://localhost:8080/download/build/:buildidfromjsonresponse:
   ```

5. Run the client

## Project Structure

```
.
├── built/      # Clients built by the server
├── server/     # Go server implementation
└── payload/    # Clientside code
```

Note: This is a proof-of-concept implementation. Some features may be incomplete or require additional configuration.

## Contact

Twitter - https://x.com/7N7
Others - https://misleadi.ng