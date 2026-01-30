# Admin & User Profile Management API Documentation

## Overview

Sistem manajemen admin dan profil user untuk aplikasi POS Team Boolean. Fitur ini menyediakan kemampuan untuk:

- Edit profil user (nama, password)
- Melihat daftar admin (superadmin only)
- Mengedit akses admin (superadmin only)
- Logout
- Membuat admin baru dengan password otomatis via email

## Authentication

Semua endpoint (kecuali login) memerlukan JWT token di header:

```
Authorization: Bearer <token>
```

Token diperoleh dari endpoint Login dengan masa berlaku 24 jam.

---

## Endpoints

### 1. Login

**Endpoint:** `POST /api/v1/auth/login`

**Description:** User melakukan login dengan email dan password

**Request Body:**

```json
{
  "email": "admin@example.com",
  "password": "password123"
}
```

**Response Success (200):**

```json
{
  "code": 200,
  "message": "Login successful",
  "data": {
    "id": 1,
    "email": "admin@example.com",
    "name": "Admin User",
    "role": "superadmin",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_at": "2025-01-20T10:00:00Z"
  }
}
```

**Response Error (401):**

```json
{
  "code": 401,
  "message": "Invalid email or password"
}
```

---

### 2. Get User Profile

**Endpoint:** `GET /api/v1/profile`

**Description:** Mengambil profil user yang sedang login

**Headers:**

```
Authorization: Bearer <token>
```

**Response Success (200):**

```json
{
  "code": 200,
  "message": "Profil user berhasil diambil",
  "data": {
    "id": 1,
    "email": "admin@example.com",
    "name": "Admin User",
    "role": "superadmin",
    "status": "active",
    "created_at": "2025-01-15T10:00:00Z",
    "updated_at": "2025-01-15T10:00:00Z"
  }
}
```

---

### 3. Update User Profile

**Endpoint:** `PUT /api/v1/profile`

**Description:** Update profil user (nama dan/atau password)

**Headers:**

```
Authorization: Bearer <token>
```

**Request Body:**

```json
{
  "name": "Admin User Updated",
  "password": "newpassword123"
}
```

**Notes:**

- `name` dan `password` adalah optional, bisa salah satu atau keduanya
- `password` minimal 6 karakter

**Response Success (200):**

```json
{
  "code": 200,
  "message": "Profil user berhasil diubah",
  "data": {
    "id": 1,
    "email": "admin@example.com",
    "name": "Admin User Updated",
    "role": "superadmin",
    "status": "active",
    "created_at": "2025-01-15T10:00:00Z",
    "updated_at": "2025-01-20T10:05:00Z",
    "message": "Profil berhasil diubah"
  }
}
```

---

### 4. List Admins (Superadmin Only)

**Endpoint:** `GET /api/v1/admin/list`

**Description:** Mengambil daftar semua admin (hanya superadmin)

**Headers:**

```
Authorization: Bearer <token>
```

**Query Parameters:**

- `page` (optional, default: 1) - Nomor halaman
- `limit` (optional, default: 10, max: 100) - Jumlah data per halaman
- `role` (optional) - Filter by role (admin, superadmin, user)

**Examples:**

- `GET /api/v1/admin/list?page=1&limit=10`
- `GET /api/v1/admin/list?page=1&limit=10&role=admin`

**Response Success (200):**

```json
{
  "code": 200,
  "message": "Daftar admin berhasil diambil",
  "data": {
    "data": [
      {
        "id": 1,
        "email": "superadmin@example.com",
        "name": "Super Admin",
        "role": "superadmin",
        "status": "active",
        "created_at": "2025-01-15T10:00:00Z",
        "updated_at": "2025-01-15T10:00:00Z"
      },
      {
        "id": 2,
        "email": "admin@example.com",
        "name": "Admin User",
        "role": "admin",
        "status": "active",
        "created_at": "2025-01-16T10:00:00Z",
        "updated_at": "2025-01-16T10:00:00Z"
      }
    ],
    "total": 2,
    "page": 1,
    "limit": 10,
    "total_pages": 1
  }
}
```

**Response Error (403):**

```json
{
  "code": 403,
  "message": "Anda tidak memiliki akses untuk melihat daftar admin"
}
```

---

### 5. Edit Admin Access (Superadmin Only)

**Endpoint:** `PUT /api/v1/admin/:id/access`

**Description:** Edit role dan status admin (hanya superadmin)

**Headers:**

```
Authorization: Bearer <token>
```

**URL Parameters:**

- `id` - ID admin yang akan diedit

**Request Body:**

```json
{
  "role": "admin",
  "status": "active"
}
```

**Valid Values:**

- `role`: `admin`, `superadmin`, `user`
- `status`: `active`, `inactive`

**Response Success (200):**

```json
{
  "code": 200,
  "message": "Akses admin berhasil diubah",
  "data": {
    "id": 2,
    "email": "admin@example.com",
    "name": "Admin User",
    "role": "admin",
    "status": "inactive",
    "message": "Akses admin admin@example.com berhasil diubah"
  }
}
```

**Response Error (403):**

```json
{
  "code": 403,
  "message": "Hanya superadmin yang dapat mengedit akses admin"
}
```

**Response Error (400):**

```json
{
  "code": 400,
  "message": "Tidak dapat menonaktifkan satu-satunya superadmin"
}
```

---

### 6. Create Admin (Superadmin Only)

**Endpoint:** `POST /api/v1/admin/create`

**Description:** Membuat admin baru dan mengirim password via email

**Headers:**

```
Authorization: Bearer <token>
```

**Request Body:**

```json
{
  "email": "newadmin@example.com",
  "name": "New Admin",
  "role": "admin"
}
```

**Validasi:**

- `email` - Email valid dan unik
- `name` - Minimal 3 karakter, maksimal 100 karakter
- `role` - `admin` atau `superadmin`

**Response Success (201):**

```json
{
  "code": 201,
  "message": "Admin berhasil dibuat. Password telah dikirim ke email",
  "data": {
    "id": 3,
    "email": "newadmin@example.com",
    "name": "New Admin",
    "role": "admin",
    "status": "active",
    "created_at": "2025-01-20T10:00:00Z",
    "updated_at": "2025-01-20T10:00:00Z"
  }
}
```

**Email Content:**

```
Halo New Admin,

Akun admin Anda telah berhasil dibuat di sistem POS.

Berikut adalah kredensial akun Anda:
Email: newadmin@example.com
Password: [random-generated-password]
Role: admin

Silakan login dan ubah password Anda di halaman profil.
Jangan bagikan password ini kepada orang lain.

Best regards,
POS System Administrator
```

**Response Error (403):**

```json
{
  "code": 403,
  "message": "Hanya superadmin yang dapat membuat admin baru"
}
```

**Response Error (400):**

```json
{
  "code": 400,
  "message": "Email sudah terdaftar"
}
```

---

### 7. Logout

**Endpoint:** `POST /api/v1/auth/logout`

**Description:** Logout user

**Headers:**

```
Authorization: Bearer <token>
```

**Request Body:** (empty)

**Response Success (200):**

```json
{
  "code": 200,
  "message": "Logout berhasil",
  "data": {
    "message": "Berhasil logout"
  }
}
```

---

## Error Responses

### 401 Unauthorized

```json
{
  "code": 401,
  "message": "Unauthorized" atau "Invalid or expired token"
}
```

### 403 Forbidden

```json
{
  "code": 403,
  "message": "Anda tidak memiliki akses untuk resource ini"
}
```

### 400 Bad Request

```json
{
  "code": 400,
  "message": "Deskripsi error spesifik"
}
```

### 500 Internal Server Error

```json
{
  "code": 500,
  "message": "Terjadi kesalahan pada server"
}
```

---

## Password Generation

Ketika membuat admin baru, password otomatis di-generate dengan kriteria:

- Panjang: 12 karakter
- Kombinasi: Huruf besar, huruf kecil, angka, dan simbol (!@#$%)
- Dikirim via email ke admin baru
- Admin harus mengubah password pada login pertama (rekomendasi)

---

## Role-Based Access Control

| Feature           | Admin | Superadmin |
| ----------------- | ----- | ---------- |
| Edit own profile  | ✓     | ✓          |
| View own profile  | ✓     | ✓          |
| List all admins   | ✗     | ✓          |
| Edit admin access | ✗     | ✓          |
| Create new admin  | ✗     | ✓          |
| Logout            | ✓     | ✓          |

---

## Testing with Postman

### Collection Variables

```json
{
  "base_url": "http://localhost:8080",
  "token": "your_jwt_token_here"
}
```

### Login Request

```
POST {{base_url}}/api/v1/auth/login
Content-Type: application/json

{
  "email": "superadmin@example.com",
  "password": "password123"
}
```

### Use Token

```
Authorization: Bearer {{token}}
```

---

## Implementation Notes

1. **JWT Token**
   - Signing Method: HMAC-SHA256
   - Expiration: 24 hours
   - Secret: Dari environment variable `JWT_SECRET`

2. **Password Security**
   - Hashing: bcrypt
   - Generated Password: Random 12 karakter
   - Stored: Hashed di database

3. **Email Service**
   - SMTP configuration dari environment
   - Email dikirim async (non-blocking)
   - Template: Plain text HTML

4. **Database**
   - User table dengan soft delete (is_deleted)
   - Role field untuk RBAC
   - Status field untuk enable/disable account

---

## Database Schema

### Users Table

```sql
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  email VARCHAR(255) UNIQUE NOT NULL,
  password VARCHAR(255) NOT NULL,
  name VARCHAR(100) NOT NULL,
  role VARCHAR(20) DEFAULT 'user',
  status VARCHAR(20) DEFAULT 'active',
  is_deleted BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_is_deleted ON users(is_deleted);
```

---

## File Structure

```
internal/
├── usecase/
│   └── admin.go (AdminUseCase interface dan implementation)
├── adaptor/
│   └── admin_adaptor.go (HTTP handlers)
├── dto/
│   └── admin.go (Request/Response DTOs)
└── data/repository/
    └── auth.go (dengan method tambahan)

pkg/
├── utils/
│   └── token.go (JWT generation & validation)
└── middleware/
    └── auth.go (Authentication middleware)
```
