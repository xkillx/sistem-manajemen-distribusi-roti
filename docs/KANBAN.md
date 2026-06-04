# Kanban Task List - SMDR MVP 1.0

> Sumber: PRD MVP SMDR v1.0

---

# Foundation

| ID | Task | Priority | Dependencies | Status |
|------|------|------|------|------|
| SMDR-001 | Finalisasi Requirement MVP | High | - | Done |
| SMDR-002 | Buat UI Wireframe MVP | High | SMDR-001 | Ready |
| SMDR-003 | Setup Project Frontend Angular | High | SMDR-001 | Done |
| SMDR-004 | Setup Project Backend Go Gin | High | SMDR-001 | Done |
| SMDR-005 | Setup Database PostgreSQL | High | SMDR-001 | Done |
| SMDR-006 | Setup Deployment MVP | High | SMDR-003, SMDR-004, SMDR-005 | Done |
| SMDR-007 | Implementasi Login | High | SMDR-003, SMDR-004, SMDR-005 | Done |
| SMDR-008 | Implementasi Role Access Owner & Sales | High | SMDR-007 | Done |
| SMDR-009 | Implementasi Logout & Session Handling | Medium | SMDR-007 | Done |

---

# Master Data

| ID | Task | Priority | Dependencies |
|------|------|------|------|
| SMDR-010 | Schema & API Produk | High | SMDR-005, SMDR-008 |
| SMDR-011 | UI Kelola Produk | High | SMDR-010 |
| SMDR-012 | Schema & API Warung | High | SMDR-005, SMDR-008 |
| SMDR-013 | UI Kelola Warung | High | SMDR-012 |
| SMDR-014 | UI Daftar Warung untuk Sales | High | SMDR-012 |
| SMDR-015 | Schema & API Sales | High | SMDR-005, SMDR-008 |
| SMDR-016 | UI Kelola Sales | High | SMDR-015 |

---

# GPS Check-in

| ID | Task | Priority | Dependencies |
|------|------|------|------|
| SMDR-017 | Schema & API Check-in Kunjungan | High | SMDR-005, SMDR-008, SMDR-012 |
| SMDR-018 | UI GPS Check-in Sales | High | SMDR-014, SMDR-017 |

---

# Distribusi

| ID | Task | Priority | Dependencies |
|------|------|------|------|
| SMDR-019 | Schema & API Pengiriman Titipan | High | SMDR-010, SMDR-012, SMDR-017 |
| SMDR-020 | UI Input Stok Titipan | High | SMDR-014, SMDR-018, SMDR-019 |

---

# Penjualan, Retur, dan Pembayaran

| ID | Task | Priority | Dependencies |
|------|------|------|------|
| SMDR-021 | Schema & API Penjualan | High | SMDR-010, SMDR-012, SMDR-019 |
| SMDR-022 | Schema & API Retur | High | SMDR-010, SMDR-012, SMDR-019 |
| SMDR-023 | Schema & API Pembayaran | High | SMDR-012, SMDR-021, SMDR-022 |
| SMDR-024 | UI Input Penjualan, Retur & Pembayaran | High | SMDR-021, SMDR-022, SMDR-023 |
| SMDR-025 | Ringkasan Transaksi Warung: stok titipan, nilai penjualan, pembayaran, selisih | High | SMDR-019, SMDR-021, SMDR-022, SMDR-023 |

---

# Dashboard

| ID | Task | Priority | Dependencies |
|------|------|------|------|
| SMDR-026 | API Ringkasan Dashboard Owner | High | SMDR-010, SMDR-012, SMDR-015, SMDR-019, SMDR-021, SMDR-022, SMDR-023 |
| SMDR-027 | UI Dashboard Owner | High | SMDR-026 |

---

# Reporting

| ID | Task | Priority | Dependencies |
|------|------|------|------|
| SMDR-028 | API Laporan Penjualan | High | SMDR-021, SMDR-023 |
| SMDR-029 | API Laporan Retur | Medium | SMDR-022 |
| SMDR-030 | API Laporan Warung | High | SMDR-025 |
| SMDR-031 | API Laporan Sales | High | SMDR-017, SMDR-021, SMDR-023 |
| SMDR-032 | UI Reporting Owner | High | SMDR-028, SMDR-029, SMDR-030, SMDR-031 |

---

# Quality & Security

| ID | Task | Priority | Dependencies |
|------|------|------|------|
| SMDR-034 | Validasi Input & Error Handling | High | Core APIs |
| SMDR-035 | Audit Security MVP | High | SMDR-007, SMDR-008 |
| SMDR-036 | Optimasi Response Time | Medium | Reporting & Dashboard |
| SMDR-037 | Testing Flow Sales end-to-end | High | SMDR-020, SMDR-024 |
| SMDR-038 | Testing Flow Owner end-to-end | High | SMDR-011, SMDR-013, SMDR-016, SMDR-027, SMDR-032 |
| SMDR-039 | Browser Compatibility Check | Medium | SMDR-037, SMDR-038 |

---

# Release

| ID | Task | Priority | Dependencies |
|------|------|------|------|
| SMDR-040 | Production Release MVP | High | SMDR-006, SMDR-035, SMDR-037, SMDR-038, SMDR-039 |

---

# Sprint Plan

## Sprint 1 — Foundation

- SMDR-001 — Finalisasi Requirement MVP
- SMDR-002 — Buat UI Wireframe MVP
- SMDR-003 — Setup Project Frontend Angular
- SMDR-004 — Setup Project Backend Go Gin
- SMDR-005 — Setup Database PostgreSQL
- SMDR-007 — Implementasi Login
- SMDR-008 — Implementasi Role Access Owner & Sales
- SMDR-009 — Implementasi Logout & Session Handling

## Sprint 2 — Master Data

- SMDR-010 — Schema & API Produk
- SMDR-011 — UI Kelola Produk
- SMDR-012 — Schema & API Warung
- SMDR-013 — UI Kelola Warung
- SMDR-014 — UI Daftar Warung untuk Sales
- SMDR-015 — Schema & API Sales
- SMDR-016 — UI Kelola Sales

## Sprint 3 — Distribusi

- SMDR-017 — Schema & API Check-in Kunjungan
- SMDR-018 — UI GPS Check-in Sales
- SMDR-019 — Schema & API Pengiriman Titipan
- SMDR-020 — UI Input Stok Titipan

## Sprint 4 — Penjualan, Retur, dan Pembayaran

- SMDR-021 — Schema & API Penjualan
- SMDR-022 — Schema & API Retur
- SMDR-023 — Schema & API Pembayaran
- SMDR-024 — UI Input Penjualan, Retur & Pembayaran
- SMDR-025 — Ringkasan Transaksi Warung: stok titipan, nilai penjualan, pembayaran, selisih

## Sprint 5 — Dashboard & Reporting

- SMDR-026 — API Ringkasan Dashboard Owner
- SMDR-027 — UI Dashboard Owner
- SMDR-028 — API Laporan Penjualan
- SMDR-029 — API Laporan Retur
- SMDR-030 — API Laporan Warung
- SMDR-031 — API Laporan Sales
- SMDR-032 — UI Reporting Owner

## Sprint 6 — Hardening & Release

- SMDR-034 — Validasi Input & Error Handling
- SMDR-035 — Audit Security MVP
- SMDR-036 — Optimasi Response Time
- SMDR-037 — Testing Flow Sales end-to-end
- SMDR-038 — Testing Flow Owner end-to-end
- SMDR-039 — Browser Compatibility Check
- SMDR-040 — Production Release MVP

---

# Future Release

| ID | Task | Priority | Dependencies |
|------|------|------|------|
| SMDR-F001 | Export Excel | Medium | SMDR-032 |

---

# Kanban Workflow

## Backlog

Task sudah teridentifikasi namun belum siap dikerjakan.

## Ready

Requirement jelas dan siap masuk sprint.

## In Progress

Sedang dalam proses development.

## Review

Menunggu code review, QA, atau validasi product owner.

## Done

Lolos acceptance criteria dan siap digunakan.

---

# Definition of Done

Sebuah task dianggap selesai apabila:

- Acceptance Criteria terpenuhi.
- Validasi input sudah berjalan.
- Error handling sudah tersedia.
- Role access sudah sesuai.
- Data tersimpan dan dapat dibaca kembali.
- Unit test atau integration test relevan sudah dijalankan.
- Tidak menyebabkan regression pada fitur utama.
- Dokumentasi teknis diperbarui jika diperlukan.
