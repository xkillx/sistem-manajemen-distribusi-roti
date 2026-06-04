# Product Requirements Document (PRD)

# Sistem Manajemen Distribusi Roti (SMDR)

Version: MVP 1.0

Status: Finalized for MVP implementation

---

# 1. Executive Summary

SMDR adalah aplikasi untuk membantu pengusaha roti mengelola distribusi produk ke warung-warung yang menggunakan sistem titip jual.

Saat ini pencatatan dilakukan secara manual sehingga pemilik usaha sulit mengetahui:

* Berapa stok yang dititipkan ke setiap warung
* Berapa produk yang terjual
* Berapa produk yang diretur
* Berapa uang yang harus disetor
* Aktivitas sales di lapangan

Sistem ini memungkinkan sales mencatat aktivitas distribusi melalui smartphone dan owner memonitor aktivitas dengan data yang langsung tersedia setelah tersimpan. Untuk MVP, dashboard dapat memakai refresh manual atau polling ringan; WebSocket/SSE tidak wajib.

---

# 2. Goals

## Business Goals

* Mengurangi pencatatan manual
* Mempercepat proses pelaporan
* Mengurangi kehilangan data
* Mengetahui performa setiap warung
* Mengetahui performa setiap sales

## Success Metrics

* 100% kunjungan sales tercatat dalam sistem
* Owner dapat melihat laporan harian tanpa menunggu laporan manual
* Mengurangi kesalahan pencatatan stok minimal 80%

---

# 3. User Roles

## Owner

Memiliki akses penuh terhadap seluruh data.

MVP hanya mendukung satu usaha/single-tenant. Owner melihat seluruh data Sales, Warung, Produk, Kunjungan, dan laporan tanpa pembatasan cabang atau tenant.

### Kebutuhan

* Melihat dashboard
* Melihat laporan penjualan
* Melihat aktivitas sales
* Mengelola master data

---

## Sales

Petugas yang mengantarkan roti ke warung.

### Kebutuhan

* Melihat daftar warung
* Input stok titipan
* Input barang terjual
* Input retur
* Input pembayaran
* Check-in lokasi
* Melihat riwayat kunjungan sendiri

UI Sales diprioritaskan mobile-first untuk penggunaan smartphone di lapangan.

---

# 4. Scope MVP

## In Scope

### Master Data

#### Produk

* Nama produk
* SKU
* Harga jual
* Status aktif

#### Warung

* Nama warung
* Pemilik
* Nomor telepon
* Alamat
* Titik GPS

MVP tidak membutuhkan peta embedded. Titik GPS warung dan lokasi check-in cukup ditampilkan sebagai koordinat/jarak dan link buka di Google Maps.

#### Sales

* Nama
* Nomor telepon
* Username
* Password

---

### Distribusi

Sales dapat mencatat:

MVP hanya mengelola stok titipan di warung. Sistem tidak mencatat stok gudang, stok produksi, atau stok yang dibawa Sales, dan tidak memvalidasi jumlah titipan terhadap stok internal usaha.

#### Pengiriman

* Tanggal
* Warung
* Produk
* Jumlah dititipkan

#### Penjualan

* Produk
* Jumlah terjual

#### Retur

* Produk
* Jumlah retur

#### Pembayaran

* Nominal diterima
* Metode pembayaran
* Catatan

Catatan pembayaran wajib jika metode pembayaran adalah lainnya; untuk tunai dan transfer, catatan opsional.

---

### GPS Check-in

Saat sales mengunjungi warung:

* Ambil lokasi GPS
* Simpan waktu kunjungan

---

### Ringkasan Transaksi Warung

Menampilkan posisi terkini per warung:

* Stok titipan per produk
* Total nilai penjualan
* Total pembayaran diterima
* Selisih pembayaran

Ringkasan ini digunakan oleh Sales saat rekonsiliasi kunjungan pada warung yang sedang dikunjungi dan oleh Owner sebagai pintasan sebelum membuka laporan warung. Sales boleh melihat selisih pembayaran berjalan untuk warung yang sedang dikunjungi, tetapi tidak memiliki akses laporan global semua warung.

---

### Dashboard Owner

Menampilkan:

* Total warung
* Total sales
* Total produk terjual dalam pcs
* Total retur dalam pcs
* Total pemasukan sebagai pembayaran diterima dalam rupiah
* Aktivitas sales hari ini

Aktivitas sales hari ini ditampilkan sebagai daftar kunjungan terbaru yang memuat waktu check-in, sales, warung, jarak check-in, status kunjungan, nilai penjualan, pembayaran diterima, dan catatan jika kunjungan tanpa transaksi.

Dashboard Owner menyediakan refresh manual dan dapat melakukan polling ringan setiap 60 detik.

UI Owner diprioritaskan untuk desktop/laptop dan tetap responsive untuk perangkat lain.

---

### Reporting

#### Laporan Penjualan

Filter:

* Harian
* Mingguan
* Bulanan

Filter laporan MVP menggunakan rentang tanggal `date_from` dan `date_to` berdasarkan Hari Bisnis Asia/Jakarta. UI boleh menyediakan shortcut seperti hari ini, 7 hari terakhir, dan bulan ini.

Laporan Penjualan menampilkan ringkasan total dan tabel detail per tanggal kunjungan, warung, sales, dan produk agar Owner dapat menelusuri sumber angka.

#### Laporan Retur

Filter:

* Periode
* Produk
* Warung
* Sales

Menampilkan:

* Jumlah retur per produk
* Daftar kunjungan terkait

#### Laporan Warung

Menampilkan:

* Total titipan
* Total terjual
* Total retur
* Nilai penjualan
* Pembayaran diterima
* Selisih pembayaran

#### Laporan Sales

Menampilkan:

* Jumlah kunjungan
* Jumlah kunjungan dengan transaksi
* Produk terjual dalam pcs
* Nilai penjualan
* Pembayaran diterima
* Retur dalam pcs

---

# 5. Out of Scope (Future Release)

Tidak termasuk dalam MVP:

* Mobile native Android
* Mobile iOS
* Offline mode
* Route optimization
* Accounting
* Piutang
* Import master data dari Excel/CSV
* Multi gudang
* Multi cabang
* Merge warung duplikat
* Tracking kendaraan
* WhatsApp integration
* Invoice PDF otomatis
* Barcode scanner
* Foto bukti kunjungan, retur, atau pembayaran
* Offsite backup database
* Account lockout otomatis setelah login gagal
* Export Excel
* AI forecasting

---

# 6. User Flow

## Sales Mengisi Stok

Login

↓

Pilih atau daftarkan Warung

↓

GPS Check-in

↓

Isi form Kunjungan terpadu

↓

Input Produk dan Jumlah Dititipkan

↓

Simpan

---

## Sales Mencatat Penjualan

Pilih Warung

↓

GPS Check-in

↓

Isi form Kunjungan terpadu

↓

Input Produk Terjual

↓

Input Retur

↓

Input Pembayaran

↓

Simpan

Untuk MVP, input titipan, penjualan, retur, pembayaran, dan catatan dilakukan dalam satu form Kunjungan terpadu setelah GPS check-in.

Form Kunjungan menampilkan harga jual produk dan menghitung nilai penjualan otomatis. Sales tidak dapat mengubah harga penjualan secara manual; harga penjualan diambil dari harga jual produk saat kunjungan diselesaikan.

Daftar produk untuk Titipan baru hanya berisi produk aktif. Daftar produk untuk Penjualan dan Retur berisi produk dengan stok awal kunjungan lebih dari nol di warung tersebut, termasuk produk nonaktif yang masih memiliki stok titipan.

Kunjungan tidak dihapus permanen melalui UI MVP. Kunjungan yang salah dibatalkan; efeknya keluar dari stok dan laporan operasional, tetapi riwayat pembatalannya tetap tersimpan.

Sales dapat melihat riwayat kunjungannya sendiri dengan filter sederhana hari ini dan 7 hari terakhir. Sales dapat mengoreksi atau membatalkan kunjungan miliknya hanya pada Hari Bisnis yang sama; riwayat hari sebelumnya bersifat read-only untuk Sales. Sales tidak dapat melihat riwayat kunjungan Sales lain.

Backend memvalidasi stok titipan dalam transaksi database saat finalisasi atau koreksi kunjungan. Jika stok berubah karena kunjungan lain dan request akan membuat stok negatif, request ditolak dan UI meminta Sales memuat ulang ringkasan warung.

Kunjungan hanya dapat diselesaikan jika memiliki minimal salah satu dari Titipan, Penjualan, Retur, Pembayaran, atau catatan alasan Kunjungan Tanpa Transaksi. Draft boleh kosong setelah GPS check-in, tetapi tidak boleh diselesaikan kosong tanpa catatan.

Catatan umum Kunjungan bersifat opsional jika Kunjungan memiliki Titipan, Penjualan, Retur, atau Pembayaran. Catatan alasan wajib hanya untuk Kunjungan Tanpa Transaksi.

---

## Owner Melihat Laporan

Login

↓

Dashboard

↓

Pilih Periode

↓

Lihat Laporan

---

# 7. Functional Requirements

## FR-001 Login

User dapat login menggunakan username dan password.

Username unik secara global untuk semua akun Owner dan Sales, dengan perbandingan case-insensitive.

Owner pertama dibuat melalui konfigurasi deployment awal. MVP hanya mendukung satu akun Owner dan tidak menyediakan UI untuk menambah, mengedit, atau menonaktifkan akun Owner. Owner dapat mengganti password akunnya sendiri. Jika Owner lupa password dan tidak bisa login, reset dilakukan melalui command admin/deployment yang membuat password sementara Owner; Owner wajib menggantinya setelah login. Tidak ada registrasi publik atau forgot-password via email pada MVP.

Deployment production MVP hanya melakukan seed Owner pertama. Produk, warung, dan sales dibuat melalui UI; data contoh hanya boleh digunakan untuk development/test.

---

## FR-002 Kelola Produk

Owner dapat:

* Tambah produk
* Edit produk
* Nonaktifkan produk

Produk tidak dihapus permanen melalui UI MVP agar histori laporan tetap stabil.

Produk hanya dapat dibuat, diedit, atau dinonaktifkan oleh Owner. Sales tidak dapat membuat Produk baru di lapangan.

---

## FR-003 Kelola Warung

Owner dapat:

* Tambah warung
* Edit warung
* Simpan lokasi GPS

Sales dapat membuat warung baru saat akuisisi lapangan dengan data minimum nama warung dan GPS check-in pertama. Warung hasil akuisisi dapat langsung menerima Titipan pada Kunjungan pertama meskipun pemilik, nomor telepon, dan alamat belum diisi. Saat membuat warung baru, sistem memberi peringatan ringan jika ada warung dengan nama mirip atau lokasi dekat, tetapi tidak memblokir pembuatan warung.

Sales dapat mencari warung berdasarkan nama. Daftar warung untuk Sales menampilkan nama bersama alamat/catatan lokasi dan informasi lokasi yang tersedia agar warung bernama sama tetap dapat dibedakan. Pengurutan berdasarkan jarak dari lokasi Sales bersifat opsional dan bukan blocker MVP; default dapat menggunakan nama atau terakhir dikunjungi.

Sales dapat mengedit terbatas nama, kontak, alamat/catatan lokasi, dan Titik GPS Warung aktif, termasuk Warung yang dibuat Sales lain. Perubahan Titik GPS Warung harus dilakukan secara eksplisit dan masuk audit log. Sales tidak dapat menonaktifkan Warung; nonaktifkan Warung hanya Owner.

Warung tidak dihapus permanen melalui UI MVP agar histori laporan tetap stabil.

---

## FR-004 Kelola Sales

Owner dapat:

* Tambah sales
* Reset password
* Nonaktifkan akun

Akun Sales dibuat atau direset menggunakan password sementara. Sales wajib mengganti password sementara sebelum dapat membuat kunjungan.

Sales dapat mengganti password sendiri setelah login. Jika Sales lupa password, Owner dapat melakukan reset password dan memberikan password sementara baru.

Sales tidak dihapus permanen melalui UI MVP agar histori laporan tetap stabil.

Username Sales tidak dapat diubah setelah akun dibuat. Owner dapat mengubah nama dan nomor telepon, mereset password, atau menonaktifkan akun.

---

## FR-005 Input Distribusi

Sales dapat mencatat stok yang dititipkan.

---

## FR-006 Input Penjualan

Sales dapat mencatat produk yang terjual.

---

## FR-007 Input Retur

Sales dapat mencatat produk yang dikembalikan.

---

## FR-008 Input Pembayaran

Sales dapat mencatat pembayaran dari warung.

---

## FR-009 GPS Check-in

Sistem menyimpan:

* Latitude
* Longitude
* Waktu kunjungan

---

## FR-010 Dashboard

Owner dapat melihat ringkasan bisnis.

---

## FR-011 Reporting

Owner dapat menghasilkan laporan:

* Penjualan
* Retur
* Kunjungan sales
* Performa warung

---

# 8. Non Functional Requirements

## Acceptance Testing

Flow Sales utama harus diuji end-to-end: login, pilih atau daftarkan Warung, GPS check-in, simpan Kunjungan berisi Titipan, buat Kunjungan berikutnya berisi Penjualan/Retur/Pembayaran, dan lihat riwayat kunjungan sendiri.

Flow Owner utama harus diuji end-to-end: login, kelola Produk/Sales/Warung, lihat dashboard, lihat laporan Penjualan/Retur/Warung/Sales, serta koreksi atau batalkan Kunjungan.

## Performance

* Response time < 2 detik

## Security

* Password hash bcrypt
* Password minimal 8 karakter
* JWT authentication via httpOnly secure cookie
* HTTPS
* Rate limit sederhana pada endpoint login
* Audit log minimal untuk aksi penting

Audit log MVP mencatat actor, action, entity, entity ID, waktu, dan ringkasan perubahan untuk create/update/cancel kunjungan, reset password sales, nonaktifkan sales/warung/produk, dan perubahan harga produk. Audit login success/failure tidak wajib masuk audit log MVP; login cukup ditangani melalui log aplikasi dan rate limit.

## Availability

* 99% uptime

## Connectivity

* MVP berjalan online-only.
* Sales harus memiliki koneksi internet untuk login, mengambil data, GPS check-in, dan menyimpan kunjungan.
* Jika request gagal karena jaringan, sistem menampilkan error dan Sales mencoba menyimpan ulang; tidak ada antrean transaksi offline.

## Browser Support

* Chrome
* Edge
* Safari

Flow Sales dengan GPS check-in diprioritaskan untuk Chrome Android. Edge dan Safari didukung untuk penggunaan dasar, dashboard, dan reporting; perbedaan perilaku izin GPS di Safari tidak menjadi blocker utama MVP kecuali kebutuhan operasional berubah.

---

# 9. Recommended Tech Stack

Frontend:

* Angular SPA

Backend:

* Go (Gin)

API:

* REST

Database:

* PostgreSQL

Database Access:

* SQL-first with sqlc + pgx, without full ORM

Authentication:

* JWT

Deployment:

* VPS Ubuntu with Docker Compose
* Reverse proxy HTTPS
* Scheduled daily local PostgreSQL backup with 7-day retention

Storage:

* httpOnly secure cookie for authentication session
* Browser localStorage only for non-sensitive UI preferences if needed

---

# 10. MVP Timeline

Week 1

* Requirement finalization
* UI Wireframe

Week 2

* Master Data

Week 3

* Distribusi

Week 4

* Penjualan dan Retur

Week 5

* Dashboard

Week 6

* Reporting

Week 7

* Testing

Week 8

* Production Release

---
