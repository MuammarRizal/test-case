# Sales API - Sales, Tax & Discount Calculator

REST API untuk mengelola transaksi penjualan, perhitungan pajak, dan perhitungan diskon bertingkat.

## Struktur Project

```
sales-api/
├── cmd/
│   └── server/
│       └── main.go              # Entry point aplikasi
├── internal/
│   ├── config/
│   │   └── database.go          # Konfigurasi database
│   ├── handlers/
│   │   ├── penjualan.go         # Handler untuk penjualan
│   │   ├── pajak.go             # Handler untuk pajak
│   │   ├── diskon.go            # Handler untuk diskon
│   │   └── root.go              # Handler untuk root endpoint
│   ├── models/
│   │   ├── penjualan.go         # Model untuk penjualan
│   │   ├── pajak.go             # Model untuk pajak
│   │   ├── diskon.go            # Model untuk diskon
│   │   └── response.go          # Model untuk response API
│   ├── services/
│   │   ├── penjualan.go         # Business logic penjualan
│   │   ├── pajak.go             # Business logic pajak
│   │   └── diskon.go            # Business logic diskon
│   └── utils/
│       ├── response.go          # Utility untuk response
│       └── validation.go        # Utility untuk validasi
├── go.mod
└── README.md
```

## Prerequisites

- Go 1.21 atau lebih tinggi
- MySQL Database
- Database dengan nama `test_case`

## Setup Database

Pastikan Anda memiliki tabel berikut di database MySQL:

```sql
-- Tabel penjualan
CREATE TABLE penjualan (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nama_pelanggan VARCHAR(255) NOT NULL,
    tanggal DATE NOT NULL,
    jam TIME NOT NULL,
    total DECIMAL(15,2) NOT NULL,
    bayar_tunai DECIMAL(15,2) NOT NULL,
    kembali DECIMAL(15,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabel penjualan_item
CREATE TABLE penjualan_item (
    id INT AUTO_INCREMENT PRIMARY KEY,
    penjualan_id INT NOT NULL,
    item_id VARCHAR(255) NOT NULL,
    quantity DECIMAL(10,2) NOT NULL,
    harga DECIMAL(15,2) NOT NULL,
    sub_total DECIMAL(15,2) NOT NULL,
    FOREIGN KEY (penjualan_id) REFERENCES penjualan(id)
);
```

## Installation

1. Clone repository

```bash
git clone <repository-url>
cd sales-api
```

2. Install dependencies

```bash
go mod tidy
```

3. Konfigurasi database di `internal/config/database.go`

```go
dsn := "root:@tcp(127.0.0.1:3306)/test_case?parseTime=true"
```

4. Jalankan aplikasi

```bash
go run cmd/server/main.go
```

Server akan berjalan di `http://localhost:8080`

## API Endpoints

### 1. Root Endpoint

```bash
GET /
```

Menampilkan informasi API dan contoh penggunaan.

### 2. Penjualan (Sales Transaction)

```bash
POST /penjualan
```

**Request Body:**

```json
{
  "nama_pelanggan": "John Doe",
  "tanggal": "2024-01-15",
  "jam": "14:30",
  "total": 150000,
  "bayar_tunai": 200000,
  "kembali": 50000,
  "items": [
    {
      "item_id": "ITM001",
      "quantity": 2,
      "harga": 50000,
      "sub_total": 100000
    },
    {
      "item_id": "ITM002",
      "quantity": 1,
      "harga": 50000,
      "sub_total": 50000
    }
  ]
}
```

### 3. Perhitungan Pajak

#### GET Method

```bash
GET /hitung-pajak?total=22000&persen_pajak=10
```

#### POST Method

```bash
POST /hitung-pajak
```

**Request Body:**

```json
{
  "total": 22000,
  "persen_pajak": 10
}
```

**Response:**

```json
{
  "status": "success",
  "message": "Perhitungan pajak berhasil",
  "data": {
    "net_sales": 20000,
    "pajak_rp": 2000
  }
}
```

### 4. Perhitungan Diskon Bertingkat

```bash
POST /hitung-diskon
```

**Request Body:**

```json
{
  "discounts": [{ "diskon": "20" }, { "diskon": "10" }],
  "total_sebelum_diskon": 100000
}
```

#### Dengan Detail

```bash
POST /hitung-diskon?detail=true
```

**Response (tanpa detail):**

```json
{
  "status": "success",
  "message": "Perhitungan diskon bertingkat berhasil",
  "data": {
    "total_diskon": 28000,
    "total_harga_setelah_diskon": 72000
  }
}
```

**Response (dengan detail):**

```json
{
  "status": "success",
  "message": "Perhitungan diskon bertingkat berhasil",
  "data": {
    "total_diskon": 28000,
    "total_harga_setelah_diskon": 72000,
    "total_sebelum_diskon": 100000,
    "jumlah_tingkat_diskon": 2,
    "detail_perhitungan": [
      {
        "step": 1,
        "persen_diskon": 20,
        "harga_sebelum": 100000,
        "nominal_diskon": 20000,
        "harga_setelah": 80000
      },
      {
        "step": 2,
        "persen_diskon": 10,
        "harga_sebelum": 80000,
        "nominal_diskon": 8000,
        "harga_setelah": 72000
      }
    ]
  }
}
```

## Testing

### Contoh Curl Commands

```bash
# Test penjualan
curl -X POST http://localhost:8080/penjualan \
  -H "Content-Type: application/json" \
  -d '{
    "nama_pelanggan": "John Doe",
    "tanggal": "2024-01-15",
    "jam": "14:30",
    "total": 150000,
    "bayar_tunai": 200000,
    "kembali": 50000,
    "items": [
      {"item_id": "ITM001", "quantity": 2, "harga": 50000, "sub_total": 100000},
      {"item_id": "ITM002", "quantity": 1, "harga": 50000, "sub_total": 50000}
    ]
  }'

# Test pajak (GET)
curl "http://localhost:8080/hitung-pajak?total=22000&persen_pajak=10"

# Test pajak (POST)
curl -X POST http://localhost:8080/hitung-pajak \
  -H "Content-Type: application/json" \
  -d '{"total": 22000, "persen_pajak": 10}'

# Test diskon
curl -X POST http://localhost:8080/hitung-diskon \
  -H "Content-Type: application/json" \
  -d '{
    "discounts": [{"diskon": "20"}, {"diskon": "10"}],
    "total_sebelum_diskon": 100000
  }'

# Test diskon dengan detail
curl -X POST "http://localhost:8080/hitung-diskon?detail=true" \
  -H "Content-Type: application/json" \
  -d '{
    "discounts": [{"diskon": "20"}, {"diskon": "10"}],
    "total_sebelum_diskon": 100000
  }'
```

## Features

- ✅ Clean Architecture dengan separation of concerns
- ✅ Proper error handling dan validation
- ✅ Database transaction untuk penjualan
- ✅ Support GET dan POST untuk perhitungan pajak
- ✅ Perhitungan diskon bertingkat dengan detail
- ✅ JSON response format yang konsisten
- ✅ Input validation yang comprehensive
- ✅ Database connection pooling

## Development

Untuk menambah fitur baru:

1. Tambahkan model di `internal/models/`
2. Implementasikan business logic di `internal/services/`
3. Buat handler di `internal/handlers/`
4. Daftarkan route di `cmd/server/main.go`

## License

MIT License
