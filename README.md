# Sistem Manajemen Distribusi Roti (SMDR)

Aplikasi web untuk mengelola distribusi roti dengan sistem titip jual antara pemilik usaha, sales, dan warung. MVP single-tenant: Sales mencatat kunjungan dari smartphone, Owner memantau data master, stok titipan, penjualan, retur, pembayaran, dashboard, dan laporan secara real-time.

## Tech Stack

- **Frontend:** Angular 18 SPA
- **Backend:** Go dengan Gin
- **API:** REST
- **Database:** PostgreSQL 16
- **Akses Database:** SQL-first dengan sqlc + pgx, tanpa ORM
- **Auth:** JWT dalam httpOnly secure cookie, bcrypt password hashing
- **Deploy:** Docker Compose di VPS Ubuntu

## Quick Start

```bash
cp .env.example .env
# Edit .env dengan secret Anda
docker compose up -d
```

Akun Owner di-seed dari `.env` saat pertama kali dijalankan. Login di `http://localhost:4200`.

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

### Migrasi Database
```bash
docker compose exec postgres psql -U smdr -d smdr -f /migrations/001_init.sql
```

## Konteks Domain

Lihat `CONTEXT.md` untuk glosarium domain dan `docs/PRD.md` untuk dokumen kebutuhan produk lengkap.
