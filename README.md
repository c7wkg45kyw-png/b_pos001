# BPOS001 Backend

POS module backend built with Clean Architecture, Gin, GORM, PostgreSQL, zap, JWT, and Swagger UI.

## Run

    go run ./cmd/server

Swagger UI: http://localhost:8082/docs/swagger.html

All /api/v1/* endpoints require Authorization: Bearer <jwt> with merchant_id, branch_id, aud=pos_module, and scopes such as pos:read, pos:create, pos:update, pos:delete, or pos:*.

## BAUT001 Token Compatibility

BPOS001 accepts BAUT001 access tokens when both services share the same `JWT_SECRET` and `JWT_ISSUER`. The middleware follows BAUT001 behavior: it validates signature and issuer, but does not enforce `aud`.

Required token/header context:

- Token claim: `merchant_id`
- Token claim: `scopes` with `pos:read`, `pos:create`, `pos:update`, `pos:delete`, or `pos:*`
- Optional token claim: `branch_id`
- If the token does not include `branch_id`, send header `X-Branch-ID: <branch_id>`

Example decoded payload:

```json
{
  "iss": "https://global-commerce.com",
  "sub": "<user-id>",
  "aud": ["sku_module,pos_module,inventory_module,report_module"],
  "merchant_id": "mch_555666777",
  "client_id": "vendor_portal_web",
  "type": "USER",
  "scopes": ["sku:read", "sku:create", "sku:update", "sku:delete", "pos:create", "pos:read"],
  "iat": 1783932400,
  "exp": 1784018800
}
```

Login through BAUT001:

```bash
curl -X POST http://localhost:8081/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"admin1234"}'
```

Use that token against BPOS001:

```bash
curl http://localhost:8082/api/v1/posclients \
  -H 'Authorization: Bearer <access_token>' \
  -H 'X-Branch-ID: BRANCH-HQ-001'
```
# b_pos001
