# Stockpilot CLI Documentation

The `stockpilot` CLI gives you full access to [Stockpilot](https://stockpilot.com) from your terminal. It works standalone for humans and as a tool for AI agents (Claude, Cursor, Codex, and any agent that can run shell commands).

---

## Installation

```bash
curl -fsSL https://stockpilot.com/install | bash
```

Or download a binary directly from [GitHub Releases](https://github.com/stockpilot/stockpilot-cli/releases):

| Platform | Binary |
|---|---|
| macOS (Apple Silicon) | `stockpilot_darwin_arm64` |
| macOS (Intel) | `stockpilot_darwin_amd64` |
| Linux | `stockpilot_linux_amd64` |
| Windows | `stockpilot_windows_amd64.exe` |

---

## Authentication

You need a Stockpilot API Client ID and Client Secret. Get these from your Stockpilot account under **Settings → API**.

```bash
stockpilot login
# Client ID: xxxxxxxx
# Client Secret: ••••••••
# Logged in successfully.
```

Credentials are stored in `~/.config/stockpilot/config.json` (permissions: 600 — only your user can read it).

Verify your credentials at any time:

```bash
stockpilot whoami
```

```
ID   Organization      Unique ID
1    Acme Warehouse    acme-warehouse
```

---

## Output modes

Every command outputs a human-readable table by default.

Pass `--json` to get raw JSON — for scripts, pipes, and AI agents:

```bash
stockpilot inventory list --json
stockpilot orders list --status pending --json | jq '.[].id'
```

---

## Commands

### `stockpilot status`

Your morning briefing. Shows pending orders and low stock items in one shot.

```bash
stockpilot status
stockpilot status --json
```

> **Note for backend team:** This command currently makes two separate API calls (`/orders` + `/inventory`). A dedicated `GET /status/summary` endpoint would make this instant and more accurate. See [BACKEND_REQUESTS.md](./BACKEND_REQUESTS.md).

---

### `stockpilot inventory`

#### List inventory

```bash
stockpilot inventory list
stockpilot inventory list --page 2 --page-size 50
stockpilot inventory list --json
```

Flags:
- `--page` (default: 1)
- `--page-size` (default: 100, max: 1000)

Output columns: ID, SKU, Name, Qty, Location

#### Get a single item

Lookup by SKU, internal ID, or barcode — exactly one required:

```bash
stockpilot inventory get --sku WIDGET-01
stockpilot inventory get --id 123
stockpilot inventory get --barcode 8712345678901
```

Always returns full JSON (used as the detail view).

#### Update an item

Update quantity, bin location, and/or reorder threshold in a single command:

```bash
# Update quantity after receiving stock
stockpilot inventory update --sku WIDGET-01 --quantity 200

# Move to a new bin location
stockpilot inventory update --sku WIDGET-01 --location B2-003-02

# Set a reorder threshold
stockpilot inventory update --sku WIDGET-01 --threshold 10u

# Combine in one call
stockpilot inventory update --sku WIDGET-01 --quantity 200 --location B2-003-02 --threshold 10u
```

Flags:
- `--sku` / `--id` / `--barcode` — exactly one required to identify the item
- `--quantity` — new absolute quantity (integer)
- `--location` — bin location in hierarchical format, e.g. `A1-001-01` (rack-shelf-bin)
- `--threshold` — reorder threshold: `5u` (5 units) or `4w` (4 weeks of stock)

---

### `stockpilot orders`

#### List orders

```bash
stockpilot orders list
stockpilot orders list --status pending
stockpilot orders list --status fulfilled --page 1 --page-size 25
stockpilot orders list --json
```

Flags:
- `--status` — filter by status (e.g. `pending`, `fulfilled`, `cancelled`)
- `--page` (default: 1)
- `--page-size` (default: 100)

Output columns: ID, Status, Customer, Created

#### Get a single order

```bash
stockpilot orders get ORDER_ID
stockpilot orders get ORDER_ID --json
```

Returns full order detail including line items, shipping address, and fulfillment info.

#### Fulfil an order

```bash
stockpilot orders fulfil ORDER_ID
```

Marks the order as fulfilled. All line items are fulfilled.

#### Cancel an order

```bash
stockpilot orders cancel ORDER_ID
```

---

### `stockpilot products`

#### List products

```bash
stockpilot products list
stockpilot products list --page 1 --page-size 50
stockpilot products list --json
```

Output columns: ID, SKU, Name, Price

#### Get a product

```bash
stockpilot products get --id 123
```

---

### `stockpilot customers`

```bash
stockpilot customers list
stockpilot customers list --json
```

Output columns: ID, Code, Business, Email

---

### `stockpilot analytics`

#### Sales for a specific item

```bash
stockpilot analytics sales --id 123
stockpilot analytics sales --id 123 --from 2026-01-01
stockpilot analytics sales --id 123 --from 2026-01-01 --to 2026-03-31 --interval weekly
stockpilot analytics sales --id 123 --json
```

Flags:
- `--id` — inventory item ID (required)
- `--from` — start date, ISO format `YYYY-MM-DD`
- `--to` — end date, ISO format `YYYY-MM-DD`
- `--interval` — `daily` (default), `weekly`, or `monthly`

#### Sales summary

```bash
stockpilot analytics summary
stockpilot analytics summary --from 2026-01-01
stockpilot analytics summary --json
```

Flags:
- `--from` — start date
- `--to` — end date

---

### `stockpilot shipping`

#### List shipping integrations

```bash
stockpilot shipping integrations
stockpilot shipping integrations --json
```

Shows available shipping carriers and their integration IDs. Use the integration ID with `shipping label`.

#### Request a shipping label

```bash
stockpilot shipping label ORDER_ID --integration INTEGRATION_ID
```

Flags:
- `--integration` — shipping integration ID (required, get from `shipping integrations`)

---

### `stockpilot whoami`

```bash
stockpilot whoami
stockpilot whoami --json
```

Returns organization name, ID, and enabled feature flags.

---

### `stockpilot login`

```bash
stockpilot login
```

Prompts for Client ID and Client Secret. Verifies credentials against the API before saving. Credentials are stored at `~/.config/stockpilot/config.json` with mode 600.

---

## Workflows

### Morning ops check

```bash
stockpilot status
```

### Receive a stock delivery

```bash
# Check current levels
stockpilot inventory get --sku WIDGET-01

# Update after counting received goods
stockpilot inventory update --sku WIDGET-01 --quantity 250 --location A1-002-01
```

### Fulfil a pending order end-to-end

```bash
# See what needs shipping
stockpilot orders list --status pending

# Get full order details
stockpilot orders get ORDER_ID

# Request shipping label (get integration ID first if unsure)
stockpilot shipping integrations
stockpilot shipping label ORDER_ID --integration INTEGRATION_ID

# Mark as fulfilled
stockpilot orders fulfil ORDER_ID
```

### Find your fastest moving products

```bash
stockpilot analytics summary --from 2026-01-01 --json | jq 'sort_by(.units_sold) | reverse | .[0:10]'
```

### Check stock coverage for pending orders

```bash
# Get all pending orders (what products are needed)
stockpilot orders list --status pending --json > pending.json

# Get current inventory
stockpilot inventory list --json > inventory.json

# Feed both to an LLM or script to compare
```

### Find items below 10 units (until native filter is available)

```bash
stockpilot inventory list --json | jq '[.[] | select(.quantity < 10)]'
```

---

## Using with AI agents

The Stockpilot CLI is designed to work with any AI agent that can run shell commands — Claude, Cursor, Codex, n8n, and others.

### Setup for agents

1. Authenticate once on the machine where the agent runs:
   ```bash
   stockpilot login
   ```

2. Give your agent the skill file:
   - Path: `skills/SKILL.md` (in this repo)
   - Or paste it directly into your agent's system prompt

3. The agent can now run any `stockpilot` command using `--json` output.

### Example agent prompts

> "Check my inventory and tell me which items are likely to run out within 2 weeks based on recent sales."

```bash
stockpilot inventory list --json
stockpilot analytics summary --from 2026-03-01 --json
```

> "Fulfil all pending orders that have been waiting more than 24 hours."

```bash
stockpilot orders list --status pending --json
# agent filters by created_at, then for each:
stockpilot orders fulfil ORDER_ID
```

> "What were my top 5 selling products last month?"

```bash
stockpilot analytics summary --from 2026-03-01 --to 2026-03-31 --json
```

### Output format for agents

All `--json` output is an array of objects (for list commands) or a single object (for get/update commands). Errors are written to stderr with exit code 1.

```json
// List example
[
  { "id": 1, "sku": "WIDGET-01", "name": "Widget", "quantity": 42, "location": "A1-001-01" },
  { "id": 2, "sku": "GADGET-02", "name": "Gadget", "quantity": 8, "location": "B3-002-01" }
]

// Error example (stderr, exit 1)
Error: API error 401: invalid credentials
```

---

## Error reference

| Exit code | Meaning |
|---|---|
| 0 | Success |
| 1 | Error (API error, auth failure, invalid flags) |

| Error message | Fix |
|---|---|
| `not logged in — run: stockpilot login` | Run `stockpilot login` |
| `incomplete credentials — run: stockpilot login` | Re-run `stockpilot login` |
| `API error 401: ...` | Invalid or expired credentials |
| `API error 404: ...` | Resource not found — check your ID/SKU |
| `one of --sku, --id, or --barcode is required` | You must pass exactly one identifier flag |

---

## Global flags

| Flag | Description |
|---|---|
| `--json` | Output as JSON instead of a table |
| `--version` | Print CLI version |
| `--help` | Show help for any command |

---

## Data formats

| Format | Example | Used for |
|---|---|---|
| Date | `2026-01-31` | Analytics date ranges |
| Bin location | `A1-001-01` | Inventory location (rack-shelf-bin) |
| Unit threshold | `10u` | Reorder when below 10 units |
| Week threshold | `4w` | Reorder when below 4 weeks of stock |
| ISO datetime | `2026-01-31T14:30:00Z` | Order timestamps |

---

## Version history

| Version | What changed |
|---|---|
| 1.0 | Initial release: inventory, orders, products, customers, analytics, shipping |
