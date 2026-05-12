# Aspera Terminal UI (atui)

A powerful, interactive terminal user interface for Aspera P2P/MySpace systems. `atui` simplifies the process of listing, sharing, and downloading files via Aspera without leaving your terminal.

## Features

- **Interactive TUI**: Navigate through received and shared requests with a sleek, keyboard-driven interface.
- **Smart Download**: Select specific files from a request or download everything with a single command.
- **Dynamic Completion**: Tab-completion support for commands and Request IDs (fetches your latest requests automatically).
- **Auto-Token Refresh**: Never worry about "Session Expired" errors. `atui` handles token renewal and re-authentication in the background.
- **Embedded Aspera Client**: No need to manually install Aspera; `atui` can extract and set up the necessary binaries for you.

## Installation

### Prerequisites
- Go 1.26 or higher (for building from source)
- Linux (currently optimized for Linux environments)

### Building from Source
```bash
go build -o atui
```

### Initial Setup
After building or installing, run the following command to extract embedded Aspera binaries and set up shell completion:
```bash
./atui install
```

## Usage

### 1. Login
```bash
./atui login
```
You will be prompted for your endpoint, client credentials, and user account. Your password is saved locally and securely to enable automatic token refresh.

### 2. List Requests
```bash
./atui list
```
View your received transfers in an interactive list. Use `Tab` to switch between **Received** and **Shared** lists.

### 3. Download Files
You can download files interactively or by providing a Request ID:

**Interactive Mode:**
```bash
./atui download
```

**Direct Mode:**
```bash
./atui download REQ123456789
```

### 4. Share Files
```bash
./atui share --to user@example.com --subject "Important Files" /path/to/file1 /path/to/file2
```

## Shell Completion
To enable Tab completion for your shell:

**Bash:**
```bash
echo 'source <(atui completion bash)' >> ~/.bashrc
source ~/.bashrc
```

**Zsh:**
```zsh
echo 'source <(atui completion zsh)' >> ~/.zshrc
source ~/.zshrc
```

## Configuration
Configuration is stored in `~/.config/atui/config.json`.

## License
MIT
