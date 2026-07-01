# 🏍️ MotoCare — Sistem Booking Layanan Servis Motor

## 🎯 Deskripsi Umum

MotoCare adalah sistem manajemen bengkel motor berbasis web yang memungkinkan pelanggan melakukan booking layanan servis motor secara online. Sistem terdiri dari **Backend API** menggunakan **Golang Fiber** dan **Frontend Dashboard** menggunakan **React Vite**. Sistem terhubung ke database **PostgreSQL berbasis Supabase**, memiliki fitur **autentikasi login dan register menggunakan JWT**, menyediakan **CRUD data layanan servis dan booking**, terdokumentasi menggunakan **Swagger**, serta dideploy secara online.

### Tema Aplikasi

**Sistem Booking Layanan Servis Motor (Motorcycle Service Booking & Management)**

---

## 👥 Anggota Kelompok

| Nama | NIM | Peran |
|------|-----|-------|
| Richard Firmansyah | 714240047 | Backend Developer |
| Wa Ode Nur Alia | 714240035 | Frontend Developer |

---

## 🔧 Teknologi

### Backend

- Bahasa: **Golang**
- Framework: **Fiber**
- ORM: **GORM**
- Database: **PostgreSQL menggunakan Supabase**
- Autentikasi: **JWT Bearer Token**
- Password Hashing: **bcrypt**
- API Documentation: **Swagger UI**
- Middleware: **CORS, Logger, Helmet, Recover, Rate Limiter, JWT Authentication, Role Authorization**
- Struktur project: **Modular Design (config, handlers, middlewares, models, repositories, routes, utils, seeders)**

### Frontend

- Framework: **React Vite**
- Routing: **React Router v7**
- Styling: **Tailwind CSS v4**
- Charts: **Recharts** (PieChart & BarChart)
- Alerts: **SweetAlert2**
- Icons: **Lucide React**
- Excel Export: **write-excel-file**
- Struktur project: **Atomic Design** (atoms, molecules, organisms, templates)
- Komunikasi API: **Fetch API** dengan client wrapper custom

### Deployment

- Backend: **Railway**
- Frontend: **Vercel** (tersedia `vercel.json`)
- Database: **Supabase Cloud**

---

## 🛠️ Spesifikasi

## 1. Backend API Golang Fiber

Backend dibuat menggunakan Golang Fiber dan terhubung ke PostgreSQL Supabase menggunakan GORM.

### A. Konfigurasi Database

Backend menggunakan database **PostgreSQL Supabase**.

Konfigurasi `.env`:

```env
SUPABASE_DSN=postgresql://username:password@host:5432/postgres
JWT_SECRET=your_jwt_secret_key
APP_ENV=development
RUN_MIGRATIONS=true
RUN_SEEDER=true
PORT=8080
FRONTEND_URL=http://localhost:5173
```

Catatan:
- `SUPABASE_DSN` digunakan untuk koneksi ke PostgreSQL Supabase.
- `JWT_SECRET` digunakan untuk signing token JWT.
- `RUN_MIGRATIONS=true` akan menjalankan AutoMigrate secara otomatis.
- `RUN_SEEDER=true` akan mengisi database dengan data contoh (seeder).
- Jangan push file `.env` ke GitHub.

---

### B. Endpoint Autentikasi

| Method | Endpoint | Keterangan |
|--------|----------|------------|
| `POST` | `/register` | Mendaftarkan user baru |
| `POST` | `/login` | Login user dan menghasilkan JWT |
| `GET` | `/me` | Mengambil data user yang sedang login |
| `PUT` | `/change-password` | Mengubah password user sendiri (user) atau user tertentu (admin) |

Data user:

- `id` — primary key
- `username` — unique, not null
- `email` — unique, not null
- `password` — disimpan dengan hash bcrypt
- `role` — `admin` atau `user` (default: `user`)

Ketentuan autentikasi:
- Password di-hash menggunakan **bcrypt**.
- Login menghasilkan token **JWT** dengan masa berlaku 24 jam.
- Token dikirim dari frontend melalui header: `Authorization: Bearer <token>`.
- Semua endpoint protected menggunakan middleware JWT.
- Rate limiter diterapkan pada endpoint `/register`, `/login`, dan `/change-password`.

---

### C. Endpoint CRUD Data Utama

Sistem memiliki **2 resource utama** yang saling berelasi: **Services** (Layanan Servis) dan **Bookings** (Booking).

#### Services (Layanan Servis)

| Method | Endpoint | Keterangan | Akses |
|--------|----------|------------|-------|
| `GET` | `/api/public/services` | List layanan aktif (public) | Public |
| `GET` | `/api/public/services/:id` | Detail layanan aktif (public) | Public |
| `GET` | `/api/services` | List semua layanan | Admin + User |
| `GET` | `/api/services/:id` | Detail layanan | Admin + User |
| `POST` | `/api/services` | Tambah layanan baru | Admin |
| `PUT` | `/api/services/:id` | Ubah layanan | Admin |
| `DELETE` | `/api/services/:id` | Hapus layanan | Admin |

#### Bookings (Booking Servis)

| Method | Endpoint | Keterangan | Akses |
|--------|----------|------------|-------|
| `GET` | `/api/bookings` | List booking (admin: semua, user: milik sendiri) | Admin + User |
| `GET` | `/api/bookings/:id` | Detail booking | Admin + User |
| `GET` | `/api/bookings/reserved-slots` | Slot yang sudah terisi | Admin + User |
| `POST` | `/api/bookings` | Tambah booking baru | Admin + User |
| `PUT` | `/api/bookings/:id` | Ubah booking (status / detail) | Admin + User |
| `DELETE` | `/api/bookings/:id` | Hapus booking | Admin |

#### Categories (Kategori Layanan)

| Method | Endpoint | Keterangan | Akses |
|--------|----------|------------|-------|
| `GET` | `/api/categories` | List kategori | Admin + User |
| `GET` | `/api/categories/:id` | Detail kategori | Admin + User |
| `POST` | `/api/categories` | Tambah kategori | Admin |
| `PUT` | `/api/categories/:id` | Ubah kategori | Admin |
| `DELETE` | `/api/categories/:id` | Hapus kategori | Admin |

#### Dashboard (Statistik)

| Method | Endpoint | Keterangan | Akses |
|--------|----------|------------|-------|
| `GET` | `/api/dashboard/stats` | Statistik dashboard | Admin |

---

### D. Relasi Database PostgreSQL

Database menggunakan relasi **foreign key** antar tabel.

Struktur tabel dan relasi:

```
users (1) ────< bookings (many)     // user_id FK -> users.id
service_categories (1) ────< services (many)  // category_id FK -> service_categories.id
services (1) ────< bookings (many)   // service_id FK -> services.id
```

**Tabel `users`**:
| Kolom | Tipe | Keterangan |
|-------|------|------------|
| id | SERIAL PK | Primary key |
| username | VARCHAR(100) UNIQUE | Nama pengguna |
| email | VARCHAR(150) UNIQUE | Email pengguna |
| password | VARCHAR(255) | Password (bcrypt hash) |
| role | VARCHAR(20) | admin / user |
| created_at | TIMESTAMP | Waktu dibuat |
| updated_at | TIMESTAMP | Waktu diubah |

**Tabel `service_categories`**:
| Kolom | Tipe | Keterangan |
|-------|------|------------|
| id | SERIAL PK | Primary key |
| name | VARCHAR(100) UNIQUE | Nama kategori |
| description | TEXT | Deskripsi kategori |
| created_at | TIMESTAMP | Waktu dibuat |
| updated_at | TIMESTAMP | Waktu diubah |

**Tabel `services`**:
| Kolom | Tipe | Keterangan |
|-------|------|------------|
| id | SERIAL PK | Primary key |
| category_id | INT FK | Foreign key ke service_categories.id |
| name | VARCHAR(150) | Nama layanan |
| description | TEXT | Deskripsi layanan |
| price | BIGINT | Harga layanan |
| duration_minutes | INT | Durasi pengerjaan (menit) |
| status | VARCHAR(20) | active / inactive |
| created_at | TIMESTAMP | Waktu dibuat |
| updated_at | TIMESTAMP | Waktu diubah |

**Tabel `bookings`**:
| Kolom | Tipe | Keterangan |
|-------|------|------------|
| id | SERIAL PK | Primary key |
| user_id | INT FK | Foreign key ke users.id |
| service_id | INT FK | Foreign key ke services.id |
| customer_name | VARCHAR(150) | Nama pelanggan |
| phone | VARCHAR(30) | Nomor telepon |
| vehicle_name | VARCHAR(150) | Nama kendaraan |
| vehicle_plate | VARCHAR(30) | Nomor plat |
| booking_date | TIMESTAMP | Jadwal booking |
| status | VARCHAR(30) | pending/confirmed/in_progress/completed/cancelled |
| notes | TEXT | Catatan |
| created_at | TIMESTAMP | Waktu dibuat |
| updated_at | TIMESTAMP | Waktu diubah |

Ketentuan data:
- Minimal **2 tabel utama** yang saling berelasi ✅ (4 tabel)
- Minimal **1 relasi foreign key** ✅ (3 relasi)
- Setiap tabel memiliki **minimal 10 data** ✅ (seeder menyediakan 2 user, 10 kategori, 10 layanan, 10 booking)

---

### E. Validasi Backend

Validasi dilakukan di backend menggunakan library `go-playground/validator/v10`.

Contoh validasi yang diterapkan:
- Field wajib diisi (`required`)
- Format email valid (`email`)
- Username dan email tidak boleh duplikat (pengecekan manual)
- Password minimal 6 karakter (`min=6`)
- Harga tidak boleh negatif (`gte=0`)
- Durasi harus lebih dari 0 (`gt=0`)
- Status harus valid (`oneof=active inactive` / `oneof=pending confirmed in_progress completed cancelled`)
- Foreign key harus valid (pengecekan `Exists()` di repository)
- Body tidak boleh kosong saat insert/update
- Booking tidak boleh double-booking di slot yang sama
- User hanya dapat mengakses booking miliknya sendiri
- User hanya dapat membuat booking dengan status `pending`
- User hanya dapat membatalkan booking (status `cancelled`)
- Booking final (`completed`/`cancelled`) tidak dapat diubah user

Response error yang jelas:

```json
{
  "message": "validasi gagal",
  "errors": {
    "name": "name wajib diisi",
    "price": "price tidak boleh negatif"
  }
}
```

---

### F. Middleware Backend

Backend menggunakan middleware berikut:

- **Logger** — mencatat setiap request HTTP
- **CORS** — mengizinkan akses dari frontend domain
- **Helmet** — mengamankan header HTTP
- **Recover** — mencegah panic crash
- **Rate Limiter** — membatasi request pada endpoint auth
- **JWT Authentication** — memvalidasi token pada endpoint protected
- **Role Authorization** — membedakan akses `admin` dan `user`

Ketentuan role:
- `admin` dapat melakukan tambah, ubah, hapus data layanan, kategori, dan booking
- `user` dapat melihat data layanan dan booking miliknya sendiri, membuat booking, dan mengedit/membatalkan booking miliknya

---

### G. Dokumentasi API Swagger

Dokumentasi Swagger dapat diakses melalui endpoint:

```
/docs
```

Swagger mendokumentasikan:
- Register (`POST /register`)
- Login (`POST /login`)
- Get Current User (`GET /me`)
- Change Password (`PUT /change-password`)
- List/Detail/Create/Update/Delete Categories
- List/Detail/Create/Update/Delete Services
- List/Detail/Create/Update/Delete Bookings
- Dashboard Stats
- Endpoint yang membutuhkan token diberi security `BearerAuth`

Format token di tombol Authorize Swagger:

```text
Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

---

## 2. Frontend Dashboard

Frontend dibuat menggunakan **React Vite** dengan **Atomic Design Pattern** dan **Tailwind CSS**.

### A. Halaman

Frontend memiliki halaman berikut:

1. **Home** (`/home`) — Landing page publik dengan katalog layanan
2. **Login** (`/login`) — Form login dengan validasi
3. **Register** (`/register`) — Form pendaftaran akun baru
4. **Dashboard** (`/dashboard`) — Dashboard admin dengan statistik dan chart (Recharts)
5. **Services List** (`/services`) — Katalog layanan dengan pencarian dan filter
6. **Service Detail** (`/services/:id`) — Detail layanan dengan metrik
7. **Service Create** (`/services/create`) — Form tambah layanan (admin)
8. **Service Edit** (`/services/:id/edit`) — Form edit layanan (admin)
9. **Bookings List** (`/bookings`) — Daftar booking dengan filter dan status management
10. **Booking Create** (`/bookings/create`) — Form booking multi-layanan (user)
11. **Profile** (`/profile`) — Halaman profil user dan ubah password
12. **Logout** — Tombol logout dengan konfirmasi SweetAlert2

### B. Autentikasi Frontend

- Setelah login berhasil, token JWT dan data user disimpan di `localStorage`.
- Token dikirim otomatis melalui header `Authorization: Bearer <token>`.
- Token dicek masa berlaku (expiry) sebelum digunakan.
- Jika token expired, user diarahkan ke halaman login.
- Role-based routing: admin diarahkan ke `/dashboard`, user ke `/bookings`.
- Protected route mencegah akses halaman tanpa autentikasi.

### C. Fitur Data Table

Tabel data booking memiliki **9 kolom** (admin) dan **8 kolom** (user):

```
ID, User, Customer, Phone, Vehicle, Service, Jadwal Booking, Status, Actions
```

Fitur tabel:
- Menampilkan data dari API dengan server-side pagination
- Pencarian data (search by customer, phone, vehicle, plate)
- Filter data berdasarkan status
- Sorting berdasarkan tanggal, status
- Tombol detail (modal dialog)
- Tombol edit (user dengan batasan waktu)
- Tombol cancel (user untuk booking aktif)
- Status dropdown (admin untuk mengubah status)
- Export CSV & Excel (admin)
- Skeleton loading state
- Empty state dengan action suggestion

### D. Form dan Validasi Frontend

Form insert dan edit memiliki validasi di frontend:
- Input wajib diisi (nama, email, password, phone, dll.)
- Format email valid
- Password minimal 6 karakter
- Konfirmasi password harus sama
- Harga tidak boleh negatif
- Durasi harus > 0
- Minimal 1 layanan dipilih saat booking
- Slot booking yang sudah terisi tidak dapat dipilih

Feedback ditampilkan menggunakan:
- Pesan error/success inline pada form
- SweetAlert2 untuk konfirmasi dan notifikasi

---

## 3. Fitur Bonus

Fitur nilai tambahan yang telah diimplementasikan:

| Fitur | Teknologi | Keterangan |
|-------|-----------|------------|
| Visualisasi data | Recharts | PieChart (booking by status) dan BarChart (top services) |
| Export data ke Excel | write-excel-file | Export services dan bookings ke file .xlsx |
| Export data ke CSV | Custom implementation | Export services dan bookings ke file .csv |
| Pagination server-side | GORM + Fiber | Pagination dengan query params `page` dan `limit` |
| Sorting data | GORM + Fiber | Multi-field sorting (`sort_by`, `sort_order`) |
| Responsive layout | Tailwind CSS | Layout responsif untuk mobile dan desktop |
| Dark mode | CSS custom properties | Toggle light/dark theme |
| Dashboard statistik | Recharts | Total categories, services, bookings, revenue, booking status, top services |
| Middleware keamanan | Helmet + Rate Limiter | HTTP header security + rate limiting auth endpoints |
| Multi-service booking | Frontend | User dapat memilih beberapa layanan dalam satu booking |
| Slot management | Backend | Pengecekan slot booking yang sudah terisi |

---

## 📁 Struktur Project

### Backend

```text
motocare-backend/
├── config/
│   ├── database.go          # Koneksi database, connection pool, AutoMigrate
│   └── env.go               # Load environment variables
├── docs/
│   ├── docs.go              # Generated Swagger docs
│   ├── swagger.json         # Swagger JSON spec
│   └── swagger.yaml         # Swagger YAML spec
├── handlers/
│   ├── auth_handler.go      # Register, Login, Me, ChangePassword
│   ├── booking_handler.go   # CRUD Booking + ReservedSlots
│   ├── category_handler.go  # CRUD Category
│   ├── dashboard_handler.go # Dashboard stats
│   ├── service_handler.go   # CRUD Service + Public endpoints
│   ├── handler_helpers.go   # Helpers (parse ID, date, validasi)
│   ├── swagger_annotations.go
│   └── swagger_types.go
├── middlewares/
│   ├── jwt_middleware.go     # JWT authentication middleware
│   └── role_middleware.go    # Role-based authorization middleware
├── models/
│   ├── user.go              # User model
│   ├── service_category.go  # ServiceCategory model
│   ├── service.go           # Service model
│   └── booking.go           # Booking model
├── repositories/
│   ├── user_repository.go      # User DB operations
│   ├── category_repository.go  # Category DB operations
│   ├── service_repository.go   # Service DB operations (filter, sort, pagination)
│   ├── booking_repository.go   # Booking DB operations (filter, sort, pagination, slot check)
│   ├── dashboard_repository.go # Dashboard aggregation queries
│   └── query_helpers.go        # Shared query helpers
├── routes/
│   ├── auth_routes.go       # Auth + rate limiter routes
│   └── crud_routes.go       # CRUD + dashboard routes with JWT & role middleware
├── seeders/
│   └── database_seeder.go   # Seeder: 2 user, 10 kategori, 10 layanan, 10 booking
├── utils/
│   ├── jwt.go               # JWT generate & parse
│   ├── password.go           # bcrypt hash & compare
│   ├── validator.go          # Struct validation
│   ├── pagination.go         # Pagination parser & meta
│   └── response.go           # Standardized API responses
├── Dockerfile
├── go.mod
├── go.sum
├── main.go                  # Entry point: load env, connect DB, migrate, seed, routes, listen
└── .env                     # Environment variables (not committed)
```

### Frontend

```text
motocare-frontend/
├── src/
│   ├── components/
│   │   ├── atoms/
│   │   │   ├── StatusBadge.jsx    # Badge komponen untuk status
│   │   │   └── index.js
│   │   ├── molecules/
│   │   │   ├── EmptyState.jsx     # Empty state dengan ilustrasi
│   │   │   ├── ErrorBoundary.jsx  # React error boundary
│   │   │   └── index.js
│   │   ├── organisms/
│   │   │   ├── DashboardTopbar.jsx  # Topbar dengan user info & toggle
│   │   │   ├── ServiceForm.jsx      # Form reusable untuk create/edit service
│   │   │   ├── Sidebar.jsx          # Sidebar navigasi
│   │   │   └── index.js
│   │   └── templates/
│   │       ├── DashboardLayout.jsx  # Layout utama dashboard
│   │       └── index.js
│   ├── pages/
│   │   ├── Home.jsx            # Landing page publik
│   │   ├── Login.jsx           # Halaman login
│   │   ├── Register.jsx        # Halaman register
│   │   ├── Dashboard.jsx       # Admin dashboard dengan chart
│   │   ├── ServicesList.jsx    # Daftar layanan (card grid)
│   │   ├── ServiceDetail.jsx   # Detail layanan
│   │   ├── ServiceCreate.jsx   # Form tambah layanan
│   │   ├── ServiceEdit.jsx     # Form edit layanan
│   │   ├── BookingsList.jsx    # Daftar booking (table)
│   │   ├── BookingCreate.jsx   # Form booking multi-layanan
│   │   └── Profile.jsx         # Profil + ubah password
│   ├── routes/
│   │   ├── AppRoutes.jsx       # Definisi semua route
│   │   └── ProtectedRoute.jsx  # Route guard dengan role check
│   ├── services/
│   │   ├── api.js              # API client (fetch wrapper, auth header, error handling)
│   │   ├── services.js         # Service API calls
│   │   ├── bookings.js         # Booking API calls
│   │   └── dashboard.js        # Dashboard API calls
│   ├── utils/
│   │   ├── auth.js             # Token & session management
│   │   ├── theme.js            # Dark/light mode toggle
│   │   ├── validation.js       # Form validators
│   │   ├── bookingValidation.js
│   │   ├── serviceValidation.js
│   │   ├── csv.js              # CSV export utilities
│   │   ├── excel.js            # Excel export utilities
│   │   └── alerts.js           # SweetAlert2 wrappers
│   ├── styles/
│   │   └── typography.css
│   ├── App.jsx                 # Root component
│   ├── App.css                 # Global styles
│   ├── ui-fixes.css            # UI polish
│   ├── index.css               # Tailwind entry
│   └── main.jsx                # Entry point
├── index.html
├── package.json
├── vite.config.js
├── vercel.json                 # Vercel SPA routing config
└── .env                        # Environment variables
```

---

## 🚀 Cara Menjalankan

### Backend

```bash
cd motocare-backend

# Setup environment
cp .env.example .env
# Edit .env dengan SUPABASE_DSN dan JWT_SECRET

# Install dependencies
go mod download

# Generate Swagger docs (jika ada perubahan)
go install github.com/swaggo/swag/cmd/swag@latest
swag init

# Run
go run main.go
```

Server berjalan di `http://localhost:8080`. Swagger UI di `http://localhost:8080/docs`.

### Frontend

```bash
cd motocare-frontend

# Setup environment
echo "VITE_API_BASE_URL=http://localhost:8080" > .env

# Install dependencies
npm install

# Development
npm run dev

# Build
npm run build
npm run preview
```

Frontend berjalan di `http://localhost:5173`.

### Seeder

Set `RUN_SEEDER=true` di `.env` backend untuk mengisi database dengan data contoh:
- **2 user**: `admin@motocare.test` / `user@motocare.test` (password: `password123`)
- **10 kategori servis**
- **10 layanan servis**
- **10 data booking**

---

## ✅ Checklist Fitur Wajib

| Fitur | Status |
|-------|--------|
| Backend Golang Fiber | ✅ |
| Database PostgreSQL Supabase | ✅ |
| Koneksi database GORM | ✅ |
| Autentikasi JWT | ✅ |
| Password hashing bcrypt | ✅ |
| Register dan Login | ✅ |
| Middleware JWT | ✅ |
| Role authorization (admin & user) | ✅ |
| CRUD minimal 1 resource utama | ✅ (Services + Bookings) |
| Minimal 2 tabel berelasi (FK) | ✅ (4 tabel, 3 relasi) |
| Minimal 10 data per tabel | ✅ (via seeder) |
| Validasi backend | ✅ |
| Validasi frontend | ✅ |
| Dokumentasi Swagger | ✅ |
| Frontend dashboard | ✅ |
| Pencarian dan filter data | ✅ |
| Deploy backend | ✅ |
| Deploy frontend | ✅ |
| Dokumentasi PDF | ✅ |

---

## 🌟 Checklist Fitur Bonus

| Fitur | Status |
|-------|--------|
| Visualisasi data (Recharts) | ✅ |
| Export data ke Excel | ✅ |
| Export data ke CSV | ✅ |
| Pagination server-side | ✅ |
| Sorting data | ✅ |
| Responsive layout | ✅ |
| Dark mode | ✅ |
| Dashboard statistik | ✅ |
| Middleware keamanan tambahan | ✅ (Helmet, Rate Limiter) |

---

## 📎 Link

| Item | Link |
|------|------|
| Repository GitHub Backend | https://github.com/copetpasarsenin/motocare-backend |
| Repository GitHub Frontend | https://github.com/copetpasarsenin/motocare-frontend |
| Deploy Backend | https://motocare-backend-production.up.railway.app/ |
| Deploy Frontend | https://motocare-frontend.vercel.app/ |
| Swagger UI | https://motocare-backend-production.up.railway.app/docs/index.html |
| Supabase Dashboard | https://supabase.com/dashboard/project/updmlcnmucfozgiuovan/editor |

---

## 📸 Screenshot

Dokumentasi screenshot lengkap tersedia di file **DOKUMEN_MOTOCARE_TUGAS_BESAR.docx** (11 screenshot asli dari deploy + placeholder Supabase):

| # | Halaman | URL |
|---|---------|-----|
| 1 | Home — Landing Page | https://motocare-frontend.vercel.app/ |
| 2 | Login — Form Masuk | https://motocare-frontend.vercel.app/login |
| 3 | Register — Form Daftar | https://motocare-frontend.vercel.app/register |
| 4 | Services List — Katalog | https://motocare-frontend.vercel.app/services |
| 5 | Service Detail | https://motocare-frontend.vercel.app/services/1 |
| 6 | Dashboard Admin — Statistik + Chart | https://motocare-frontend.vercel.app/dashboard |
| 7 | Bookings List — Tabel + Filter | https://motocare-frontend.vercel.app/bookings |
| 8 | Profile — Akun + Password | https://motocare-frontend.vercel.app/profile |
| 9 | Swagger UI — Dokumentasi API | https://motocare-backend-production.up.railway.app/docs/index.html |
| 10 | JWT Token di localStorage | https://motocare-frontend.vercel.app/ (setelah login) |
| 11 | Dashboard — Auth Terverifikasi | https://motocare-frontend.vercel.app/dashboard (setelah login) |
| 12 | ⏳ Supabase Table Editor | https://supabase.com/dashboard/project/updmlcnmucfozgiuovan/editor |

---

Selamat mengerjakan tugas besar. MotoCare siap dipresentasikan!
