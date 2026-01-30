# 🚀 QUICK START GUIDE - AUTHENTICATION API

## Prerequisites

- Go 1.18+
- PostgreSQL database
- Postman (untuk testing)
- SMTP credentials (optional, untuk email)

---

## 1️⃣ Setup Database

### Step 1: Create Database

```sql
CREATE DATABASE "database-aplikasi-pos";
```

### Step 2: Configure .env

Copy `.env.example` ke `.env`:

```bash
cp .env.example .env
```

Edit `.env` dengan database credentials Anda:

```env
DATABASE_USERNAME=postgres
DATABASE_PASSWORD=your_password
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_NAME=database-aplikasi-pos

# Optional: SMTP untuk email
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_EMAIL=your_email@gmail.com
SMTP_PASSWORD=your_app_password
```

### Step 3: Run Migration & Seeding

```bash
go run main.go --migrate --seed
```

Output akan menunjukkan:

```
Starting database auto migration...
   Table migrated: users
   Table migrated: otps
   Table migrated: staff
   ...
Database auto migration completed successfully!

Starting database seeding...
   Seeding users data...
   Seeded 3 users
   Seeding payment methods data...
   Seeded 5 payment methods
   ...
Database seeding completed!
```

---

## 2️⃣ Start Server

```bash
go run main.go
```

Output:

```
========================================
🚀 Server running on http://localhost:8080
========================================
```

---

## 3️⃣ Test API

### Option 1: Using Postman (Recommended)

1. Open Postman
2. Click "Import"
3. Select file: `Postman Collection/Authentication.postman_collection.json`
4. Collection akan loaded dengan 5 endpoints

### Option 2: Using cURL

```bash
# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@pos.com","password":"admin123"}'

# Check Email
curl -X POST http://localhost:8080/api/v1/auth/check-email \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@pos.com"}'

# Send OTP
curl -X POST http://localhost:8080/api/v1/auth/send-otp \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@pos.com","purpose":"password_reset"}'
```

---

## 4️⃣ Test Flow

### Complete Authentication Flow:

#### 1. Login (Dapatkan Token)

```bash
POST /api/v1/auth/login
Body: {
  "email": "admin@pos.com",
  "password": "admin123"
}
```

✅ Response: JWT token di `data.token`

#### 2. Check Email (Validasi Email)

```bash
POST /api/v1/auth/check-email
Body: {
  "email": "admin@pos.com"
}
```

✅ Response: `data.exists = true`

#### 3. Send OTP (Minta OTP)

```bash
POST /api/v1/auth/send-otp
Body: {
  "email": "admin@pos.com",
  "purpose": "password_reset"
}
```

✅ Response: OTP dikirim ke email (check logs untuk development)

#### 4. Get OTP Code

**Development:** Check di logs atau database:

```sql
SELECT otp_code FROM otps
WHERE email = 'admin@pos.com'
ORDER BY created_at DESC
LIMIT 1;
```

#### 5. Validate OTP (Validasi OTP)

```bash
POST /api/v1/auth/validate-otp
Body: {
  "email": "admin@pos.com",
  "otp_code": "123456",  # Replace dengan OTP dari step 4
  "purpose": "password_reset"
}
```

✅ Response: `data.valid = true`

#### 6. Reset Password (Ubah Password)

```bash
POST /api/v1/auth/reset-password
Body: {
  "email": "admin@pos.com",
  "otp_code": "123456",  # Same OTP as step 5
  "new_password": "newpassword123",
  "purpose": "password_reset"
}
```

✅ Response: `"Password reset successfully"`

#### 7. Login dengan Password Baru

```bash
POST /api/v1/auth/login
Body: {
  "email": "admin@pos.com",
  "password": "newpassword123"
}
```

✅ Response: New JWT token

---

## 🔑 Default Credentials

| User    | Email           | Password   |
| ------- | --------------- | ---------- |
| Admin   | admin@pos.com   | admin123   |
| Manager | manager@pos.com | manager123 |
| Staff   | staff@pos.com   | staff123   |

---

## 📝 Important Notes

### OTP Handling

- OTP berlaku **10 menit** dari saat dikirim
- OTP adalah **6 digit** random number
- OTP otomatis di-mark sebagai "used" setelah validasi
- Tidak bisa reuse OTP yang sudah digunakan

### Password Requirements

- Minimum **6 characters**
- No special format requirement (bisa alphanumeric only)

### JWT Token

- Berlaku **24 jam** (configurable di `.env`)
- Token included dalam setiap request ke protected endpoints (future implementation)
- Format: `Bearer <token>` (di Authorization header)

### Email Service

- Jika SMTP tidak configured, email tidak akan dikirim (tapi tidak error)
- Untuk development, check logs untuk melihat OTP code
- Untuk production, setup SMTP credentials di `.env`

---

## 🐛 Troubleshooting

### Error: "Failed to connect to database"

```bash
# Check database connection:
# 1. Pastikan PostgreSQL running
# 2. Check .env credentials
# 3. Check database name
```

### Error: "invalid or expired OTP"

```bash
# 1. Check OTP belum expired (10 menit)
# 2. Check OTP code benar (6 digit)
# 3. Check OTP belum digunakan
# 4. Check purpose sama (password_reset/email_verification)
```

### Error: "email not registered"

```bash
# 1. Check email terdaftar di database
# 2. Check user tidak di-delete (is_deleted = false)
# 3. Gunakan salah satu default emails: admin@pos.com, manager@pos.com, staff@pos.com
```

### No logs showing up

```bash
# Check DEBUG flag di .env
# DEBUG=true  # Untuk show debug logs
```

---

## 📚 Documentation Files

Lihat file-file berikut untuk informasi lebih detail:

- **AUTHENTICATION_API.md** - Full API documentation
- **API_ENDPOINTS.md** - Semua endpoints yang tersedia
- **PROJECT_STRUCTURE.md** - Project structure & architecture
- **IMPLEMENTATION_STATUS_AUTH.md** - Implementation details
- **IMPLEMENTATION_CHECKLIST_AUTH.md** - Complete checklist

---

## ✅ Verification Checklist

Sebelum go to production, pastikan:

- [ ] Database migration berhasil
- [ ] Seeding data berhasil
- [ ] Server starting tanpa error
- [ ] Health check endpoint working: `GET /health`
- [ ] Login endpoint working dengan default credentials
- [ ] OTP functionality working (at least validation)
- [ ] Password reset flow complete
- [ ] Error handling proper
- [ ] Logs showing up correctly
- [ ] SMTP configured (if using email features)

---

## 🎯 Next Steps

Setelah authentication API working:

1. **Implement Authentication Middleware** - Protect routes dengan JWT
2. **Add User Profile Management** - Edit user info
3. **Add Admin Access Management** - Control user permissions
4. **Implement Logout** - Invalidate tokens
5. **Add Token Refresh** - Extend session
6. **Dashboard Features** - Sales analytics
7. **Menu Management** - Product categories & items
8. **Unit Tests** - Test coverage untuk auth functions

---

## 📞 Support

Jika ada masalah atau pertanyaan:

1. Check logs di `./logs/` directory
2. Check database dengan SQL client
3. Review AUTHENTICATION_API.md documentation
4. Review source code di `internal/` folder

---

**Happy Testing! 🚀**

**Last Updated: 30 January 2026**
