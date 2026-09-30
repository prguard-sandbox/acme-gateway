# acme-gateway

HTTP gateway in front of the Acme services. Routes `/store/*` and `/billing/*` to their backends.

## Configuration

- `PORT`: listen port (default 8080)
- `STORE_URL`, `BILLING_URL`: backend base URLs
