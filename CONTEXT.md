# Sistem Manajemen Distribusi Roti

Konteks ini menjelaskan bahasa domain untuk aplikasi distribusi roti dengan sistem titip jual antara usaha roti, sales, dan warung.

## Language

**Kunjungan**:
Satu aktivitas **Sales** ke satu **Warung** pada satu waktu, wajib dimulai dengan **GPS Check-in**. Satu **Kunjungan** dapat mencakup **Titipan**, **Penjualan**, **Retur**, dan **Pembayaran** yang terjadi pada kunjungan tersebut; satu **Warung** dapat memiliki lebih dari satu **Kunjungan** pada **Hari Bisnis** yang sama.
_Avoid_: Transaksi, visit

**Status Kunjungan**:
Keadaan **Kunjungan** dalam siklus pencatatan. Untuk MVP, statusnya adalah draft, selesai, atau dibatalkan; hanya Kunjungan selesai yang memengaruhi stok dan laporan operasional. Draft dibuat setelah **GPS Check-in**, satu **Sales** hanya boleh memiliki satu draft aktif, dan draft otomatis menjadi dibatalkan jika tidak diselesaikan sampai akhir **Hari Bisnis**.
_Avoid_: Status transaksi

**Kunjungan Tanpa Transaksi**:
**Kunjungan** selesai yang hanya mencatat kehadiran **Sales** tanpa **Titipan**, **Penjualan**, **Retur**, atau **Pembayaran**. Kunjungan ini wajib memiliki catatan alasan, masuk laporan aktivitas **Sales**, tetapi tidak memengaruhi stok atau nilai laporan operasional.
_Avoid_: Kunjungan gagal

**Koreksi Kunjungan**:
Perubahan langsung atas data **Kunjungan** yang sudah tersimpan karena terjadi salah input. Untuk MVP, **Sales** hanya boleh melakukan **Koreksi Kunjungan** pada **Tanggal Kunjungan** yang sama tanpa mengubah **Harga Penjualan**, sedangkan koreksi setelah hari itu dan koreksi **Harga Penjualan** menjadi wewenang **Owner** dan dicatat dalam audit log.
_Avoid_: Revisi transaksi

**Pembatalan Kunjungan**:
Tindakan menandai **Kunjungan** sebagai tidak valid. **Pembatalan Kunjungan** mengeluarkan seluruh **Titipan**, **Penjualan**, **Retur**, dan **Pembayaran** di dalamnya dari **Stok Titipan**, dashboard, dan laporan operasional, tetapi riwayat pembatalannya tetap tersimpan.
_Avoid_: Hapus kunjungan

**GPS Check-in**:
Catatan lokasi dan waktu saat **Sales** memulai **Kunjungan** di **Warung**. Untuk MVP, **Kunjungan** tidak valid tanpa **GPS Check-in**, tetapi check-in tidak ditolak berdasarkan radius lokasi.
_Avoid_: Absensi jika yang dimaksud adalah bukti kunjungan ke warung

**Titik GPS Warung**:
Lokasi rujukan **Warung** untuk menghitung **Jarak Check-in**. Titik ini dapat kosong saat **Owner** membuat **Warung** tanpa mengetahui lokasinya, lalu dibuat dari **GPS Check-in** pertama ke **Warung** tersebut atau dari perubahan eksplisit pada data **Warung**, bukan diperbarui otomatis oleh check-in berikutnya.
_Avoid_: Lokasi check-in jika yang dimaksud adalah lokasi master warung

**Jarak Check-in**:
Jarak antara lokasi **GPS Check-in** dan **Titik GPS Warung**. Untuk MVP, jarak ini menjadi informasi audit untuk **Owner**, bukan syarat untuk menerima atau menolak **Kunjungan**; pada check-in pertama yang menginisialisasi **Titik GPS Warung**, jarak belum menjadi informasi audit.
_Avoid_: Geofence jika yang dimaksud hanya audit jarak

**Tanggal Kunjungan**:
Tanggal ketika **Sales** mencatat aktivitas di **Warung**. Untuk MVP, **Penjualan** dilaporkan berdasarkan **Tanggal Kunjungan**, bukan tanggal produk sebenarnya terjual kepada pembeli akhir.
_Avoid_: Tanggal jual aktual

**Hari Bisnis**:
Rentang tanggal operasional yang digunakan untuk dashboard dan laporan harian. Untuk MVP, **Hari Bisnis** mengikuti timezone Indonesia Barat, Asia/Jakarta.
_Avoid_: UTC day jika yang dimaksud adalah hari operasional lokal

**Sales**:
Petugas lapangan yang mengunjungi **Warung** dan mencatat aktivitas distribusi roti.
_Avoid_: Kurir, driver

**Sales Aktif**:
**Sales** yang masih boleh login dan membuat **Kunjungan** baru. **Sales** yang tidak aktif tetap muncul pada riwayat dan laporan.
_Avoid_: Sales tersedia jika yang dimaksud adalah boleh login atau bekerja

**Password Sementara**:
Password awal atau hasil reset untuk akun **Sales** yang harus diganti sebelum **Sales** dapat membuat **Kunjungan**.
_Avoid_: Password permanen dari Owner

**Owner**:
Pengguna yang memiliki akses penuh untuk mengelola master data, akun **Sales**, dashboard, dan laporan. Untuk MVP, **Owner** adalah satu-satunya role pengelola.
_Avoid_: Admin, supervisor

**Role Akun**:
Jenis akses pengguna dalam sistem. Untuk MVP, hanya ada dua **Role Akun**: **Owner** dan **Sales**.
_Avoid_: Permission granular, role cabang

**Warung**:
Tempat titip jual yang menerima produk roti untuk dijual kepada pembeli akhir. Untuk MVP, **Warung** dapat dibuat oleh **Owner** dari dashboard atau oleh **Sales** saat akuisisi lapangan, tidak memiliki **Sales** tetap, dapat dikunjungi oleh **Sales** mana pun jika aktif, dan nama **Warung** tidak harus unik.
_Avoid_: Toko, outlet, customer

**Akuisisi Warung**:
Proses ketika **Sales** mendaftarkan **Warung** baru yang setuju menerima titip jual. Pada MVP, **Warung** hasil akuisisi langsung aktif; nama warung dan **GPS Check-in** pertama wajib, sedangkan alamat, nama pemilik, dan nomor telepon bersifat opsional.
_Avoid_: Registrasi customer jika yang dimaksud adalah pendaftaran warung lapangan

**Warung Aktif**:
**Warung** yang masih boleh dipilih untuk **Kunjungan** baru. **Warung** hasil **Akuisisi Warung** langsung aktif; **Warung** yang tidak aktif tetap muncul pada riwayat dan laporan, tetapi tidak menerima **Kunjungan** baru.
_Avoid_: Warung tersedia jika yang dimaksud adalah boleh dikunjungi

**Produk**:
Jenis roti yang dapat dititipkan ke **Warung** dan dicatat dalam **Penjualan** atau **Retur**. Nama **Produk** aktif harus unik agar tidak membingungkan saat dipilih di lapangan.
_Avoid_: Barang, item

**SKU Produk**:
Kode unik untuk membedakan **Produk**. Untuk MVP, **SKU Produk** wajib unik tetapi dapat dibuat otomatis jika **Owner** tidak mengisinya.
_Avoid_: Kode barang jika yang dimaksud adalah identitas produk sistem

**Jumlah Produk**:
Banyaknya **Produk** dalam satuan pcs. Untuk MVP, **Jumlah Produk** adalah bilangan bulat positif pada setiap baris **Titipan**, **Penjualan**, atau **Retur**.
_Avoid_: Berat, karton, jumlah desimal

**Produk Aktif**:
**Produk** yang masih boleh dipilih untuk **Titipan** baru. **Produk** yang tidak aktif tetap muncul pada riwayat dan tetap dapat direkonsiliasi melalui **Penjualan** atau **Retur** jika masih memiliki **Stok Titipan**.
_Avoid_: Produk tersedia jika yang dimaksud adalah boleh dititipkan baru

**Harga Jual**:
Harga standar saat ini untuk satu **Produk**, berlaku untuk semua **Warung**, dan dinyatakan sebagai bilangan rupiah bulat. Untuk **Produk Aktif**, **Harga Jual** harus lebih dari nol.
_Avoid_: Tarif

**Harga Penjualan**:
Harga yang berlaku untuk **Penjualan** pada **Tanggal Kunjungan**, dinyatakan sebagai bilangan rupiah bulat. Harga ini menjadi dasar laporan historis meskipun **Harga Jual** produk berubah setelahnya, dan hanya **Owner** yang boleh mengoreksinya pada **Kunjungan** historis.
_Avoid_: Harga master jika membicarakan laporan historis

**Titipan**:
Produk roti yang diserahkan kepada **Warung** untuk dijual melalui sistem titip jual. Satu **Kunjungan** dapat memiliki Titipan untuk banyak **Produk**, dan **Titipan** dapat menjadi satu-satunya aktivitas transaksi pada **Kunjungan**.
_Avoid_: Pengiriman, shipment, batch, tanggal kedaluwarsa

**Stok Titipan**:
Jumlah produk yang masih berada di **Warung** dan belum tercatat sebagai **Penjualan** atau **Retur**. Untuk MVP, stok ini dihitung dari total **Titipan** dikurangi total **Penjualan** dan total **Retur**, per **Produk** dan per **Warung**, dan tidak boleh negatif pada titik mana pun dalam urutan **Tanggal Kunjungan**.
_Avoid_: Stok gudang, stok mobil

**Stok Awal Kunjungan**:
Jumlah **Stok Titipan** untuk **Produk** di **Warung** sebelum aktivitas pada **Kunjungan** berjalan. **Penjualan** dan **Retur** pada **Kunjungan** hanya mengurangi **Stok Awal Kunjungan**, sedangkan **Titipan** pada kunjungan itu menjadi stok setelah rekonsiliasi.
_Avoid_: Stok akhir jika membicarakan validasi penjualan atau retur kunjungan

**Penjualan**:
Jumlah produk titipan yang berhasil dijual oleh **Warung** kepada pembeli akhir. Satu **Kunjungan** dapat memiliki Penjualan untuk banyak **Produk**.
_Avoid_: Order, transaksi

**Nilai Penjualan**:
Nilai uang dari **Penjualan**, dihitung dari jumlah produk yang terjual dan **Harga Penjualan** produk tersebut.
_Avoid_: Omzet, pendapatan jika yang dimaksud adalah nilai produk terjual

**Retur**:
Produk titipan yang dikembalikan dari **Warung** karena tidak terjual atau tidak layak dijual. Satu **Kunjungan** dapat memiliki Retur untuk banyak **Produk**, dan **Retur** dapat dicatat tanpa **Penjualan** baru selama **Stok Awal Kunjungan** cukup.
_Avoid_: Pembatalan, cancellation

**Alasan Retur**:
Catatan opsional yang menjelaskan mengapa **Retur** terjadi, misalnya tidak laku atau rusak. Untuk MVP, alasan tidak dikategorikan secara baku.
_Avoid_: Kategori retur jika yang dimaksud hanya catatan bebas

**Pembayaran**:
Uang yang benar-benar diterima dari **Warung** atas produk titipan yang telah terjual. Untuk MVP, satu **Kunjungan** memiliki paling banyak satu total **Pembayaran**, nilainya dapat berbeda dari **Nilai Penjualan** pada **Kunjungan** yang sama, dapat bernilai nol ketika **Penjualan** ada tetapi uang belum diterima, dan dapat lebih dari nol meskipun tidak ada **Penjualan** baru.
_Avoid_: Invoice, billing, pendapatan jika yang dimaksud adalah uang diterima

**Metode Pembayaran**:
Cara **Pembayaran** diterima. Untuk MVP, metode yang digunakan adalah tunai, transfer, atau lainnya, dengan nominal dalam rupiah; metode lainnya wajib memiliki catatan.
_Avoid_: Payment gateway, multi-currency

**Total Pemasukan**:
Total **Pembayaran** yang benar-benar diterima pada konteks laporan atau dashboard.
_Avoid_: Pendapatan

**Selisih Pembayaran**:
Perbedaan sederhana antara total **Nilai Penjualan** dan total **Pembayaran** pada konteks laporan yang sama. Untuk MVP, **Selisih Pembayaran** tidak mengalokasikan pembayaran ke kunjungan atau tagihan tertentu.
_Avoid_: Piutang jika yang dimaksud adalah modul penagihan lengkap

## Example Dialogue

Dev: "Saat Sales datang ke Warung Sari Pagi, apa yang perlu dicatat?"

Domain expert: "Sales mulai dengan GPS Check-in. Itu membuat satu Kunjungan. Di dalamnya Sales mencatat Titipan baru, Penjualan dari titipan sebelumnya, Retur, dan Pembayaran yang diterima."

Dev: "Bolehkah Warung yang sama dikunjungi lebih dari sekali pada hari yang sama?"

Domain expert: "Boleh. Setiap kedatangan dicatat sebagai Kunjungan terpisah."

Dev: "Status Kunjungan apa saja yang dipakai?"

Domain expert: "Draft, selesai, dan dibatalkan. Hanya Kunjungan selesai yang memengaruhi stok dan laporan operasional. Draft yang tidak diselesaikan sampai akhir Hari Bisnis otomatis dibatalkan."

Dev: "Bolehkah Sales check-in ke warung lain saat masih punya draft?"

Domain expert: "Tidak. Satu Sales hanya boleh punya satu draft Kunjungan aktif."

Dev: "Kalau Sales datang tetapi warung tutup atau tidak ada aktivitas?"

Domain expert: "Tetap simpan sebagai Kunjungan selesai tanpa transaksi dengan catatan alasan. Itu masuk aktivitas Sales, tetapi tidak memengaruhi stok atau nilai operasional."

Dev: "Kalau lokasi GPS Sales jauh dari titik GPS Warung, apakah check-in ditolak?"

Domain expert: "Tidak untuk MVP. Sistem mencatat Jarak Check-in untuk audit Owner, tetapi tidak menolak Kunjungan berdasarkan radius."

Dev: "Apakah setiap check-in memperbarui Titik GPS Warung?"

Domain expert: "Tidak otomatis. Titik GPS Warung hanya berubah melalui edit eksplisit pada data Warung."

Dev: "Kalau Owner membuat Warung tanpa tahu lokasinya, apakah Sales masih bisa check-in?"

Domain expert: "Bisa. GPS Check-in pertama ke Warung itu menjadi Titik GPS Warung."

Dev: "Kalau Sales salah input jumlah saat masih hari yang sama?"

Domain expert: "Itu Koreksi Kunjungan dan Sales boleh memperbaikinya pada Tanggal Kunjungan yang sama. Setelah lewat hari itu, koreksi langsung dilakukan oleh Owner dan dicatat dalam audit log."

Dev: "Kalau harga pada Kunjungan historis salah, siapa yang boleh mengoreksi?"

Domain expert: "Owner boleh mengoreksi Harga Penjualan historis dan perubahan itu dicatat dalam audit log. Sales tidak boleh mengubah harga."

Dev: "Kalau Kunjungan dibatalkan, apakah stok dan laporan tetap menghitung datanya?"

Domain expert: "Tidak. Pembatalan Kunjungan membuat aktivitas di dalamnya tidak dihitung lagi, tetapi riwayat pembatalannya tetap tersimpan."

Dev: "Kalau Sales dinonaktifkan, apakah kunjungan lamanya hilang dari laporan?"

Domain expert: "Tidak. Sales tidak aktif tidak boleh login atau membuat Kunjungan baru, tetapi riwayatnya tetap muncul."

Dev: "Apakah Owner menentukan password permanen Sales?"

Domain expert: "Tidak. Owner memberi atau mereset Password Sementara, lalu Sales wajib menggantinya sebelum membuat Kunjungan."

Dev: "Berapa jenis akses pengguna di MVP?"

Domain expert: "Dua saja: Owner dan Sales. Owner mengelola sistem dan laporan, Sales mencatat Kunjungan."

Dev: "Apakah hanya Sales tertentu yang boleh mengunjungi Warung tertentu?"

Domain expert: "Tidak untuk MVP. Sales dapat memilih Warung aktif mana pun, dan laporan mencatat Sales yang melakukan Kunjungan."

Dev: "Bolehkah ada dua Warung bernama sama?"

Domain expert: "Boleh. Banyak warung punya nama generik, jadi Sales perlu melihat nama bersama alamat atau lokasi untuk memilih yang benar."

Dev: "Apakah Owner harus membuat Warung sebelum Sales bisa menitipkan roti?"

Domain expert: "Tidak. Owner boleh membuat Warung dari dashboard, tetapi Sales juga boleh melakukan Akuisisi Warung di lapangan. Nama warung dan GPS Check-in pertama wajib untuk akuisisi lapangan; data lain bisa dilengkapi belakangan."

Dev: "Kalau Warung dinonaktifkan, apakah riwayatnya hilang?"

Domain expert: "Tidak. Warung tidak aktif tetap muncul di laporan, tetapi tidak bisa dipilih untuk Kunjungan baru."

Dev: "Kalau Produk dinonaktifkan, apakah laporan lama hilang?"

Domain expert: "Tidak. Produk tidak aktif tetap muncul di riwayat. Sales tidak boleh membuat Titipan baru untuk produk itu, tetapi masih boleh mencatat Penjualan atau Retur untuk Stok Titipan yang tersisa."

Dev: "Apakah Owner harus menyiapkan SKU untuk semua produk?"

Domain expert: "Tidak. SKU Produk harus unik, tetapi sistem boleh membuatnya otomatis jika Owner tidak mengisi."

Dev: "Bolehkah ada dua Produk aktif dengan nama yang sama?"

Domain expert: "Tidak. Nama Produk aktif harus unik. Kalau ada varian, namanya dibuat eksplisit."

Dev: "Apakah jumlah roti bisa desimal atau dalam karton?"

Domain expert: "Tidak untuk MVP. Jumlah Produk dihitung pcs dan harus bilangan bulat positif pada setiap baris yang dicatat."

Dev: "Apakah setiap Retur harus punya alasan baku?"

Domain expert: "Tidak. Alasan Retur boleh dicatat bebas jika perlu, tetapi tidak wajib dan tidak memakai kategori baku di MVP."

Dev: "Bolehkah Kunjungan hanya berisi Retur?"

Domain expert: "Boleh, selama Stok Awal Kunjungan untuk produk tersebut cukup."

Dev: "Jadi Penjualan dan Pembayaran tidak berdiri sendiri?"

Domain expert: "Benar. Untuk MVP, keduanya dicatat sebagai bagian dari Kunjungan yang sama jika terjadi saat Sales berada di warung tersebut."

Dev: "Bagaimana sistem tahu stok yang masih ada di warung?"

Domain expert: "Stok Titipan dihitung dari Titipan yang pernah masuk ke warung, dikurangi Penjualan dan Retur untuk produk yang sama."

Dev: "Apakah Titipan perlu batch atau tanggal kedaluwarsa?"

Domain expert: "Tidak untuk MVP. Titipan cukup mencatat Produk dan Jumlah Produk."

Dev: "Bolehkah Kunjungan hanya berisi Titipan?"

Domain expert: "Boleh. Titipan saja valid, misalnya saat pengisian stok awal di Warung."

Dev: "Bolehkah Sales mencatat terjual dan retur lebih banyak daripada Stok Titipan?"

Domain expert: "Tidak. Penjualan dan Retur tidak boleh membuat Stok Titipan produk di warung menjadi negatif."

Dev: "Kalau Owner mengoreksi Kunjungan lama, bolehkah hasilnya membuat stok sempat negatif pada kunjungan berikutnya?"

Domain expert: "Tidak. Koreksi Kunjungan tetap harus menjaga Stok Titipan tidak negatif di sepanjang urutan Tanggal Kunjungan."

Dev: "Kalau pada kunjungan yang sama Sales membawa Titipan baru, apakah Titipan itu boleh dipakai untuk membenarkan Penjualan atau Retur?"

Domain expert: "Tidak. Penjualan dan Retur direkonsiliasi terhadap Stok Awal Kunjungan. Titipan baru menjadi stok setelah rekonsiliasi selesai."

Dev: "Kalau roti sebenarnya laku kemarin tetapi baru dicatat hari ini, masuk laporan tanggal mana?"

Domain expert: "Masuk Tanggal Kunjungan hari ini, karena itu tanggal Sales mendapatkan dan mencatat informasinya."

Dev: "Dashboard hari ini dihitung berdasarkan timezone apa?"

Domain expert: "Hari Bisnis mengikuti Asia/Jakarta untuk MVP."

Dev: "Kalau Nilai Penjualan saat kunjungan Rp100.000 tetapi Warung hanya membayar Rp80.000?"

Domain expert: "Catat Pembayaran Rp80.000 dan tampilkan Selisih Pembayaran Rp20.000. Itu belum menjadi modul Piutang lengkap di MVP."

Dev: "Kalau Warung melaporkan ada Penjualan tetapi belum membayar, apakah Pembayaran boleh nol?"

Domain expert: "Boleh. Pembayaran nol valid jika ada Penjualan, dan Selisih Pembayaran menunjukkan uang yang belum diterima pada konteks laporan."

Dev: "Kalau Warung membayar kekurangan dari kunjungan sebelumnya tanpa Penjualan baru, apakah boleh?"

Domain expert: "Boleh. Catat Pembayaran yang diterima, meskipun tidak ada Penjualan baru pada Kunjungan itu."

Dev: "Laporan Warung menampilkan pendapatan yang mana?"

Domain expert: "Hindari istilah pendapatan. Tampilkan Nilai Penjualan dan Pembayaran secara terpisah; Total Pemasukan berarti Pembayaran yang diterima."

Dev: "Apakah Pembayaran hari ini harus dialokasikan ke Penjualan tertentu?"

Domain expert: "Tidak untuk MVP. Selisih Pembayaran hanya dihitung dari total Nilai Penjualan dikurangi total Pembayaran pada konteks laporan."

Dev: "Kalau ada beberapa produk dalam satu Kunjungan, apakah Pembayaran dicatat per produk?"

Domain expert: "Tidak. Titipan, Penjualan, dan Retur dicatat per Produk, tetapi Pembayaran dicatat sebagai satu total untuk Kunjungan."

Dev: "Apakah Pembayaran punya pajak, diskon, biaya, atau mata uang lain?"

Domain expert: "Tidak untuk MVP. Pembayaran hanya nominal rupiah dengan metode tunai, transfer, atau lainnya."

Dev: "Kalau metode pembayaran dipilih lainnya, apakah perlu catatan?"

Domain expert: "Wajib, supaya Owner tahu cara pembayaran yang sebenarnya."

Dev: "Kalau Harga Jual Produk berubah setelah kunjungan dicatat, apakah laporan lama ikut berubah?"

Domain expert: "Tidak. Laporan memakai Harga Penjualan yang berlaku saat kunjungan dicatat."

Dev: "Apakah Harga Jual bisa desimal atau nol?"

Domain expert: "Tidak untuk Produk Aktif. Harga Jual memakai rupiah bulat dan harus lebih dari nol."

Dev: "Apakah tiap Warung bisa punya harga khusus?"

Domain expert: "Tidak untuk MVP. Harga Jual Produk berlaku untuk semua Warung."
