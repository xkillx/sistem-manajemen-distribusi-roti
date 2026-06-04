# Sistem Manajemen Distribusi Roti (SMDR)

A bread consignment distribution management system for Indonesian small businesses. MVP is a single-tenant web application: Sales record field visits from a smartphone, and the Owner monitors master data, stock, sales, returns, payments, dashboard, and reports.

## Tech Stack

- **Frontend:** Angular 18 SPA
- **Backend:** Go with Gin
- **API:** REST
- **Database:** PostgreSQL 16
- **ORM:** SQL-first with sqlc + pgx
- **Auth:** JWT in httpOnly secure cookie, bcrypt passwords
- **Deploy:** Docker Compose on Ubuntu VPS

## Quick Start

```bash
cp .env.example .env
# Edit .env with your secrets
docker compose up -d
```

The Owner account is seeded from `.env` on first run. Login at `http://localhost:4200`.

## Development

### Frontend
```bash
cd frontend
npm install
npm start
```

### Backend
```bash
cd backend
go run ./cmd/server
```

### Database Migrations
```bash
docker compose exec postgres psql -U smdr -d smdr -f /migrations/001_init.sql
```

## Domain Context

See `CONTEXT.md` for the domain glossary and `docs/PRD.md` for the full product requirements.
