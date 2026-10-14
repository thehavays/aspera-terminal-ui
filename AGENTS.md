# Aspera Terminal UI — Agent & Developer Guide

This document is the primary onboarding and operational guide for AI agents and developers working in the `aspera-terminal-ui` repository.

---

## 📖 Project Overview

Aspera Terminal UI is an interactive, keyboard-driven terminal user interface for IBM Aspera P2P/MySpace systems. It simplifies listing, sharing, tab-completing, and downloading files via Aspera directly from the terminal without leaving the shell environment.

### 🛠️ Core Technology Stack
- **Language / Runtime**: Go 1.22+ (Module: `aspera-terminal-ui`)
- **CLI Framework**: `github.com/spf13/cobra`
- **TUI Framework**: `github.com/charmbracelet/bubbletea`, `bubbles`, `lipgloss`
- **Table Formatting**: `github.com/jedib0t/go-pretty/v6`
- **Security / Keyring**: `github.com/zalando/go-keyring` (requires `libsecret-1-dev` on Linux)

---

## 📁 Repository Structure

```
aspera-terminal-ui/
├── .agents/                 # Agent rules, workflows, and standards
│   ├── GIT_WORKFLOW.md      # Detailed Git branching & GitHub MCP rules
│   ├── CI.md                # Continuous Integration standards & quality gates
│   ├── CD.md                # Continuous Deployment & Snapcraft architecture
│   └── skills/              # Agent skills (aspera-terminal-ui skill)
├── .github/
│   ├── workflows/
│   │   ├── ci.yml           # Automated Go CI test workflow
│   │   └── snapcraft-release.yml # Snapcraft build & publish workflow
│   ├── ISSUE_TEMPLATE/      # GitHub issue forms (bug_report.yml, feature_request.yml)
│   └── pull_request_template.md
├── cmd/                     # Cobra CLI command definitions (login, list, download, share)
├── doc/                     # Project documentation and assets
├── internal/                # Internal application packages (api, auth, tui, storage)
├── snap/                    # Canonical Snapcraft packaging configuration and metadata
├── .editorconfig            # Standard editor formatting rules (tabs for Go)
├── .env.example             # Documented environment configuration
├── .gitignore               # Ignored binaries, caches, and secrets
├── go.mod                   # Go module definition
├── go.sum                   # Go module checksums
├── LICENSE                  # MIT License
├── main.go                  # Application entry point
└── README.md                # User & contributor documentation
```

---

## ⚡ Canonical Commands

Always use the following commands for development and verification:

| Action | Command | Purpose |
| :--- | :--- | :--- |
| **Build Binary** | `go build -o aspera-terminal-ui .` | Compiles application executable |
| **Run TUI** | `go run main.go` | Launches application directly |
| **Run Tests** | `go test -v ./...` | Executes Go test suite |
| **Static Analysis** | `go vet ./...` | Runs Go compiler static analysis |
| **Format Code** | `gofmt -s -w .` | Applies standard Go code formatting |

---

## 🌿 Git & Version Control Workflow

For complete version control guidelines, see [.agents/GIT_WORKFLOW.md](./.agents/GIT_WORKFLOW.md).

### Summary Rules:
1. **Protected Branches**: Direct pushes to `main` and `dev` are strictly forbidden.
2. **Branch Roles**:
   - `dev`: Primary integration branch. All feature and fix branches must originate from and target `dev`.
   - `main`: Production release branch. Updated only by merging `dev` via Release Pull Request.
3. **Branch Naming**: `<type>/<kebab-case-description>` (e.g. `feat/transfer-progress`, `fix/login-refresh`).
4. **Commit Format**: Conventional Commits (`feat: ...`, `fix: ...`, `chore: ...`).
5. **Pre-flight Verification**: Always run `go test ./...` and `go vet ./...` before opening PRs.
6. **User Confirmation Gates (MANDATORY)**:
   - **Before Opening PR**: ALWAYS ask user approval before calling `create_pull_request`.
   - **Before Merging PR**: ALWAYS ask user approval before calling `merge_pull_request`.
7. **Issue Closure Policy**: Issues are created first with `issue_write`, referenced in dev PRs (`ref #...`), and closed ONLY when merged to `main` via Release PR (`closes #...`).

---

## 🔄 CI/CD Guidelines
- **CI**: See [.agents/CI.md](./.agents/CI.md) for quality gates.
- **CD**: See [.agents/CD.md](./.agents/CD.md) for Snapcraft packaging and release flow.
