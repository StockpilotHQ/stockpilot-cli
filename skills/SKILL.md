# Stockpilot CLI — Agent Skill

You have access to the `stockpilot` CLI, which gives you full access to the Stockpilot inventory management platform.

## Authentication

Credentials are stored in `~/.config/stockpilot/config.json` after running `stockpilot login`.
Always verify credentials before proceeding: `stockpilot whoami --json`

## Output

Every command supports `--json` for machine-readable output.
Always use `--json` when processing output programmatically.

## Quick reference

| Goal | Command |
|---|---|
| Verify credentials | `stockpilot whoami --json` |
| Overall status | `stockpilot status --json` |
| List all inventory | `stockpilot inventory list --json` |
| Get item by SKU | `stockpilot inventory get --sku WIDGET-01 --json` |
| Get item by ID | `stockpilot inventory get --id 123 --json` |
| Update stock level | `stockpilot inventory update --sku WIDGET-01 --quantity 50` |
| Update bin location | `stockpilot inventory update --sku WIDGET-01 --location A1-001-01` |
| Set reorder threshold | `stockpilot inventory update --sku WIDGET-01 --threshold 5u` |
| List all orders | `stockpilot orders list --json` |
| List pending orders | `stockpilot orders list --status pending --json` |
| Get single order | `stockpilot orders get ORDER_ID --json` |
| Fulfil an order | `stockpilot orders fulfil ORDER_ID` |
| Cancel an order | `stockpilot orders cancel ORDER_ID` |
| List products | `stockpilot products list --json` |
| Get product | `stockpilot products get --id 123 --json` |
| List customers | `stockpilot customers --json` |
| Sales for item | `stockpilot analytics sales --id 123 --from 2026-01-01 --json` |
| Sales summary | `stockpilot analytics summary --from 2026-01-01 --json` |
| List shipping integrations | `stockpilot shipping integrations --json` |
| Request shipping label | `stockpilot shipping label ORDER_ID --integration INTEGRATION_ID` |

## Common workflows

### Check if there's enough stock for pending orders
```bash
stockpilot orders list --status pending --json
stockpilot inventory list --json
```

### Find slow/fast moving items
```bash
stockpilot analytics summary --from 2026-01-01 --json
```

### Fulfil an order end-to-end
```bash
# 1. Get order details
stockpilot orders get ORDER_ID --json
# 2. Request shipping label
stockpilot shipping label ORDER_ID --integration INTEGRATION_ID
# 3. Mark as fulfilled
stockpilot orders fulfil ORDER_ID
```

### Replenish low stock
```bash
# 1. Find low stock items
stockpilot inventory list --json | jq '.[] | select(.quantity < 10)'
# 2. Update quantities after receiving stock
stockpilot inventory update --sku WIDGET-01 --quantity 100
```

## Important notes

- Product lookup always requires exactly one of: `--id`, `--sku`, or `--barcode`
- Location format: hierarchical bin code, e.g. `A1-001-01`
- Threshold format: `5u` (units) or `33w` (weeks of stock)
- Dates: ISO format `YYYY-MM-DD`
- Pagination: use `--page` and `--page-size` (max 1000) for large datasets
- All list commands default to 100 items per page
