# Admin User Service

Back-office microservice for **end users**, **roles**, and **permissions**. Admin gateway: `/api/v1/users`, `/roles`, `/permissions`. This service exposes `/users`, `/roles`, `/permissions` directly.

Cloned from the **user/user-service** project pattern (Gin + Cobra + `turahe/pkg`). Note: the `user/user-service` submodule currently contains merchant-service code; this admin service implements the user-domain API expected by the gateway.

## API (protected, admin JWT)

| Resource | Method | Path |
|----------|--------|------|
| Users | GET | `/users` |
| Users | POST | `/users` |
| Users | GET | `/users/:id` |
| Users | PUT | `/users/:id` |
| Users | DELETE | `/users/:id` |
| Users | PATCH | `/users/:id/status` |
| Roles | GET | `/roles` |
| Roles | GET | `/roles/:id` |
| Permissions | GET | `/permissions` |
| Permissions | GET | `/permissions/:id` |
| Health | GET | `/health` |

## Quick start

```bash
cp .env.example .env
# Ensure users + roles tables exist in DATABASE_DBNAME (base)
make migrate-up   # permissions table only
make run
```

## Docker

`docker compose up -d --build admin-user-service` (port **8104**).
