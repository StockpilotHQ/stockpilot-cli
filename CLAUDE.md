# Stockpilot CLI

## Project overview

`stockpilot-cli` is the official command-line interface for [Stockpilot](https://stockpilot.com) — an inventory management platform. The CLI lets users and AI agents interact with the Stockpilot API from the terminal.

**Goals:**
- Single binary, zero dependencies, installed via `curl`
- Entity-oriented commands (not one-to-one endpoint wrappers)
- AI agent-friendly: every command supports `--json` output
- Shipped with `SKILL.md` so agents know how to use it without docs

## Brand

- Always write **Stockpilot** — capital S, lowercase p. Never "StockPilot".
- CLI binary name: `stockpilot`

## Architecture

**Language:** Go  
**CLI framework:** [cobra](https://github.com/spf13/cobra)  
**HTTP client:** `net/http` (stdlib, no extra deps)  
**Output formatting:** [tablewriter](https://github.com/olekukonko/tablewriter) for tables, stdlib JSON for `--json`  
**Credentials:** stored in `~/.config/stockpilot/config.json`

## Project structure

```
stockpilot-cli/
├── cmd/stockpilot/main.go      # Entry point
├── internal/
│   ├── api/client.go           # HTTP client, auth headers, error handling
│   ├── config/config.go        # Read/write ~/.config/stockpilot/config.json
│   └── output/formatter.go     # Table + JSON output helpers
├── commands/
│   ├── root.go                 # Root cobra command, --json flag, version
│   ├── login.go                # stockpilot login
│   ├── whoami.go               # stockpilot whoami
│   ├── inventory.go            # stockpilot inventory *
│   ├── orders.go               # stockpilot orders *
│   ├── products.go             # stockpilot products *
│   ├── customers.go            # stockpilot customers *
│   ├── analytics.go            # stockpilot analytics *
│   ├── shipping.go             # stockpilot shipping *
│   └── status.go               # stockpilot status (smart compound command)
├── skills/
│   └── SKILL.md                # Agent skill file — how to use this CLI with AI agents
├── scripts/
│   └── install.sh              # curl install script
├── .github/workflows/
│   └── release.yml             # GoReleaser cross-platform builds
├── CLAUDE.md
├── go.mod
├── go.sum
└── .gitignore
```

## API

- **Base URL:** `https://api.stockpilot.dev`
- **Auth:** `X-CLIENT-ID` and `X-CLIENT-SECRET` headers on every request
- **Credentials file:** `~/.config/stockpilot/config.json`
  ```json
  { "client_id": "...", "client_secret": "..." }
  ```
- **Pagination:** `page` + `page_size` query params (default 100, max 1000)
- **Product lookup:** always exactly one of `id`, `sku`, or `barcode`
- **Location format:** hierarchical bin, e.g. `A1-001-01`
- **Threshold format:** `5u` (units) or `33w` (weeks of stock)

## Commands (1.0 scope)

| Command | What it does |
|---|---|
| `stockpilot login` | Save client ID + secret to config file |
| `stockpilot whoami` | Verify credentials, show organization |
| `stockpilot status` | Pending orders count + low stock items (compound) |
| `stockpilot inventory list` | Paginated inventory list |
| `stockpilot inventory get` | Single item by `--sku`, `--id`, or `--barcode` |
| `stockpilot inventory update` | Update quantity/location/threshold |
| `stockpilot orders list` | List orders (supports `--status` filter) |
| `stockpilot orders get ORDER_ID` | Single order details |
| `stockpilot orders fulfil ORDER_ID` | Mark order as fulfilled |
| `stockpilot orders cancel ORDER_ID` | Cancel an order |
| `stockpilot products list` | List products |
| `stockpilot products get --id ID` | Single product |
| `stockpilot customers list` | List customers |
| `stockpilot analytics sales` | Sales per item (`--sku`, `--from`, `--to`) |
| `stockpilot analytics summary` | Overall sales summary |
| `stockpilot shipping label ORDER_ID` | Request shipping label |

## Output modes

Every command outputs a human-readable table by default.  
Pass `--json` to get raw JSON (for agents and scripts).

```bash
stockpilot inventory list --json | jq '.[] | select(.quantity < 5)'
```

## Development

```bash
go run ./cmd/stockpilot --help
go build -o stockpilot ./cmd/stockpilot
go test ./...
```

## Release

Releases are built with [GoReleaser](https://goreleaser.com) via GitHub Actions.  
Artifacts: `stockpilot_darwin_arm64`, `stockpilot_darwin_amd64`, `stockpilot_linux_amd64`, `stockpilot_windows_amd64.exe`

Install script at `scripts/install.sh` detects OS/arch and downloads the right binary.
