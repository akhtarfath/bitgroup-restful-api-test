# simple-restful-api

CRUD API sederhana untuk **Products** dan **Categories** menggunakan Go + Fiber v3, dengan penyimpanan data ke file JSON lokal (`data/products.json`, `data/categories.json`).

Endpoint CRUD dilindungi **Bearer token authentication**. Token didapat dari endpoint login yang dilindungi **rate limiter**.

## Fitur

- CRUD lengkap untuk produk & kategori
- Penyimpanan data persisten ke file JSON (via Viper untuk pembacaan)
- Login dengan Bearer token (kredensial di `.env`)
- Rate limiter khusus endpoint login
- Konfigurasi memakai Viper (`.env`)
- Arsitektur berlapis: `controller` → `service` → `store`
- Unit test + integration test HTTP
- Koleksi Postman/Insomnia untuk testing

## Teknologi

- [Go](https://go.dev/) 1.25
- [Fiber v3](https://docs.gofiber.io/) — web framework
- [Viper](https://github.com/spf13/viper) — konfigurasi & pembacaan file JSON
- [Google UUID](https://github.com/google/uuid) — generate ID

## Struktur Proyek

```
.
├── main.go                    # entrypoint
├── config/                    # konfigurasi viper (.env)
├── routes/                    # wiring route + dependency
├── domain/                    # struct & helper response
├── features/
│   ├── auth/                  # login, token service, middleware, rate limiter
│   ├── products/              # CRUD produk
│   └── categories/            # CRUD kategori
├── data/                      # file JSON hasil penyimpanan (auto-generate)
├── postman_collection.json    # koleksi Postman/Insomnia
└── .env.example               # contoh konfigurasi
```

## Cara Menjalankan

```bash
go run .
```

Secara default server berjalan di `http://localhost:3000`.

### Kredensial Login

Salin `.env.example` menjadi `.env`, lalu sesuaikan (opsional):

```
APP_PORT=3000
AUTH_USERNAME=admin
AUTH_PASSWORD=admin
TOKEN_TTL_HOURS=24
RATE_LIMIT_MAX=5
RATE_LIMIT_WINDOW_MINUTES=1
```

File `.env` tidak ikut di-commit (dikecualikan di `.gitignore`). Jika `.env` tidak ada, nilai default `admin`/`admin` tetap dipakai.

## Endpoint

### Autentikasi

| Method | Path     | Deskripsi                                   |
|--------|----------|---------------------------------------------|
| POST   | `/login` | Login, mengembalikan Bearer token. **Rate-limited** (5x/menit per IP). |

Body:

```json
{
  "username": "admin",
  "password": "admin"
}
```

Respons:

```json
{
  "code": 200,
  "status": "success",
  "message": "Login successful",
  "data": {
    "token": "<token>",
    "token_type": "Bearer",
    "expires_at": "2026-08-29T19:57:07+07:00"
  }
}
```

### Categories (butuh token)

| Method | Path                 | Deskripsi            |
|--------|----------------------|----------------------|
| GET    | `/categories`        | Daftar semua kategori|
| POST   | `/categories`        | Buat kategori        |
| GET    | `/categories/:id`    | Detail kategori      |
| PUT    | `/categories/:id`    | Ubah kategori        |
| DELETE | `/categories/:id`    | Hapus kategori       |

### Products (butuh token)

| Method | Path               | Deskripsi            |
|--------|--------------------|----------------------|
| GET    | `/products`        | Daftar semua produk  |
| POST   | `/products`        | Buat produk          |
| GET    | `/products/:id`    | Detail produk        |
| PUT    | `/products/:id`    | Ubah produk          |
| DELETE | `/products/:id`    | Hapus produk         |

Semua endpoint di atas wajib menyertakan header:

```
Authorization: Bearer <token>
```

Contoh body membuat produk:

```json
{
  "name": "Clean Code",
  "price": 35,
  "category_id": "<id-kategori>"
}
```

## Format Respons

Sukses dengan data:

```json
{
  "code": 200,
  "status": "success",
  "message": "Products retrieved successfully",
  "data": []
}
```

Error:

```json
{
  "code": 404,
  "status": "error",
  "message": "Product not found"
}
```

## Contoh Curl

```bash
# 1. Login untuk mendapatkan token
TOKEN=$(curl -s -X POST localhost:3000/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"admin"}' | jq -r '.data.token')

# 2. Ambil daftar produk (butuh token)
curl -s localhost:3000/products -H "Authorization: Bearer $TOKEN"

# 3. Buat kategori
curl -s -X POST localhost:3000/categories \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"name":"Books"}'
```

## Testing

```bash
go test ./...          # unit + integration test
go test -race ./...    # dengan race detector
```

## Testing dengan Postman / Insomnia

Import file `postman_collection.json`:

1. **Postman**: tombol *Import* → pilih file.
2. **Insomnia**: *Import* → pilih file (format Postman Collection terdeteksi otomatis).

Panggil **Auth → Login** terlebih dahulu; token akan otomatis tersimpan ke variabel `token` dan dipakai oleh seluruh request CRUD.