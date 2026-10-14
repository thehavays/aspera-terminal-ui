# Continuous Integration (CI) Guidelines for Aspera Terminal UI

---

## ⚙️ Environment & CI Context

* **GITHUB_USERNAME**: `thehavays`
* **REPOSITORY**: `aspera-terminal-ui`
* **STACK**: Go (1.22+), Cobra CLI, Bubbletea TUI, Lipgloss, `libsecret`
* **CI_STRATEGY**: Mandatory local verification gates and automated GitHub Actions PR checks

---

## 🎯 CI Philosophy & Purpose

Continuous Integration guarantees that all code integrated into `dev` and `main`:
1. **Compiles Cleanly**: Go compiler generates binaries without warnings or type discrepancies.
2. **Passes Static Analysis**: `go vet ./...` checks for common mistakes, shadowed variables, and unreachable code.
3. **Passes Automated Tests**: Go test suites (`go test -v ./...`) run cleanly with full unit and mock verification.
4. **Ensures Zero Leaks**: Strict zero-secret policy via `.env.example`, `config.json`, and `.gitignore`.

---

## 🚦 Mandatory Local Verification Gates

Before opening a Pull Request or pushing changes to remote, run these gates in order:

### 1. Formatting Check
```bash
gofmt -s -l .
```
* **Acceptance Criteria**: Zero unformatted Go files.

### 2. Static Analysis Check
```bash
go vet ./...
```
* **Acceptance Criteria**: Exit code 0, no vet issues reported.

### 3. Test Suite Verification
```bash
go test -v ./...
```
* **Acceptance Criteria**: Exit code 0, all package tests pass.

### 4. Build Check
```bash
go build ./...
```
* **Acceptance Criteria**: Compilation succeeds with zero errors.

### 5. Diff Hygiene & Cleanliness
```bash
git status
git diff
```
* **Acceptance Criteria**:
  - No leftover debug print statements (`fmt.Println` debug dumps).
  - No untracked cache, secrets (`config.json`), or binary artifacts.
  - Working tree is clean and scoped.

---

## 🤖 Automated CI Strategy: GitHub Actions

Configured workflow: [`.github/workflows/ci.yml`](../.github/workflows/ci.yml)
- **Triggers**: On `push` and `pull_request` targeting `dev` and `main`, plus `workflow_call`.
- **Environment**: Ubuntu Latest with Go `1.22`.
- **System Dependencies**: Installs `libsecret-1-dev` and `pkg-config` for system keyring integration.
- **Steps**: Checkout code, Setup Go, install dependencies, and run `go test -v ./...`.
