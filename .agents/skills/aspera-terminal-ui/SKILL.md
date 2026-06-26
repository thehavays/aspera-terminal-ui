---
name: aspera-terminal-ui-assistant
description: Reference and instructions for running and using the aspera-terminal-ui CLI application.
---

# Aspera Terminal UI Assistant

This skill provides reference instructions for running and using the Aspera Terminal UI (`aspera-terminal-ui`) application.

## CLI Usage Reference

Antigravity can run the built `aspera-terminal-ui` binary or use `go run main.go` to invoke various subcommands. Below is the command reference.

### 1. Build and Prepare
To build the binary:
```bash
go build -o aspera-terminal-ui main.go
```

To install the embedded Aspera `ascp` client tools and configure shell completions (bash/zsh):
```bash
./aspera-terminal-ui install
```

### 2. Login & Logout
Authentication is required before making transfer list or download requests.

* **Login (Interactive)**:
  ```bash
  ./aspera-terminal-ui login
  ```
  *Note: This command runs interactively and prompts the user on `stdin` for: Endpoint URL, Client ID, Client Secret, Username, and Password. When executing this command, prepare to feed input or wait for the user.*

* **Logout**:
  ```bash
  ./aspera-terminal-ui logout
  ```
  *Note: Clears all active tokens and credentials.*

### 3. List Transfers
To query or inspect transfers:
```bash
./aspera-terminal-ui list [flags]
```
* **Flags**:
  - `-r, --received` (default): Show received transfers.
  - `-s, --shared`: Show shared transfers.
  *Note: This starts an interactive Bubble Tea terminal user interface. Do not run this in automated headless scripts unless the user is watching or interacting.*

### 4. Download Files
To download files associated with a transfer request:
```bash
./aspera-terminal-ui download [REQUEST_ID] [flags]
```
* **Parameters**:
  - `REQUEST_ID`: The unique ID of the transfer request. If omitted, the command starts an interactive TUI for selection.
* **Flags**:
  - `-a, --all`: Download all files in the request without prompting for file selection.
  - `-p, --path <directory>`: Specify the destination directory (default is a subdirectory named after the request ID in the current working directory).

### 5. Share Files
To upload and share files with recipients:
```bash
./aspera-terminal-ui share [FILES...] --to <emails> --subject <subject> [flags]
```
* **Flags**:
  - `-t, --to <emails>` (required, comma-separated): Recipient email addresses.
  - `-s, --subject <subject>` (required): The subject of the transfer.
  - `-b, --body <body>`: A message description body.
  - `--cc <emails>` (comma-separated): CC email addresses.
  - `--bcc <emails>` (comma-separated): BCC email addresses.
