# Continuous Deployment (CD) Guidelines for Aspera Terminal UI

---

## ⚙️ Environment & CD Context

* **HOSTING_PLATFORM**: Snapcraft (Canonical Snap Store) & GitHub Releases
* **GITHUB_USERNAME**: `thehavays`
* **REPOSITORY**: `aspera-terminal-ui`
* **DEPLOYMENT_TYPE**: CLI / TUI Package Distribution
* **RUNTIME**: Linux (x86_64 / arm64)

---

## 🌐 Delivery & Distribution Architecture

`aspera-terminal-ui` is distributed as a standalone binary and packaged for the Canonical Snap Store:

### Distribution Channels:
1. **Canonical Snap Store**:
   - Package name: `aspera-terminal-ui`
   - Installation:
     ```bash
     sudo snap install aspera-terminal-ui
     sudo snap connect aspera-terminal-ui:password-manager-service
     ```
   - Automated CD Workflow: [`.github/workflows/snapcraft-release.yml`](../.github/workflows/snapcraft-release.yml)
   - Trigger: Pushes of semantic version tags (`v*`).

2. **Source Build & Local Installation**:
   ```bash
   go build -o aspera-terminal-ui
   ./aspera-terminal-ui install
   ```

---

## 🌿 Environment & Branch Mappings

| Branch / Tag | Stage | Purpose |
| :--- | :--- | :--- |
| **`main`** | **Production Release** | Stable release branch. |
| **`dev`** | **Integration / Staging** | Active development branch for validated feature PRs. |
| **`feat/*`** | **Feature Branches** | Feature development branches targeting `dev`. |
| **`v*` Tags** | **Snap Store Publish** | Automatically triggers the Snapcraft build & publish workflow to the `stable` channel. |

---

## 🚀 Production Release Workflow

1. **Verify Integration**: Ensure all features, tests, and vet checks pass on `dev`.
2. **Open Release PR**: Ask user confirmation, then create Release PR from `dev` to `main` via GitHub MCP `create_pull_request`.
3. **Merge Release PR**: Ask user confirmation, then merge via GitHub MCP `merge_pull_request` (auto-closes resolved issues).
4. **Publish Release & Tag**:
   - Tag the release commit on `main`:
     ```bash
     git tag -a v1.0.12 -m "Release v1.0.12"
     git push origin v1.0.12
     ```
   - GitHub Actions automatically executes `snapcraft-release.yml` to build the snap and publish to the Snap Store.
