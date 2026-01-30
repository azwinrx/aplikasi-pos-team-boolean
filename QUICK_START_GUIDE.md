# Quick Start Guide - Admin & User Management

## Installation Steps

### 1. Install JWT Dependency

```bash
cd d:\Lumoshive\aplikasi-pos-team-boolean
go get github.com/golang-jwt/jwt/v5
```

### 2. Set Environment Variables

```bash
# Create .env file or set in your environment
JWT_SECRET=your-secret-key-change-in-production
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_FROM=noreply@pos-system.com
```

### 3. Run Application

```bash
go run main.go
```

---

## Test Flow with Postman

### Step 1: Import Collection

1. Open Postman
2. Click "Import"
3. Select `Admin_User_Management.postman_collection.json`
4. Set `base_url` variable to `http://localhost:8080`

### Step 2: Login

```
POST /api/v1/auth/login

Body:
{
  "email": "superadmin@example.com",
  "password": "password123"
}

Save token from response to {{token}} variable
```

### Step 3: Get Profile

```
GET /api/v1/profile
Header: Authorization: Bearer {{token}}

Response: User profile data
```

### Step 4: Update Profile

```
PUT /api/v1/profile
Header: Authorization: Bearer {{token}}

Body:
{
  "name": "Updated Name",
  "password": "newpassword123"
}

Response: Updated profile
```

### Step 5: List Admins (Superadmin Only)

```
GET /api/v1/admin/list?page=1&limit=10
Header: Authorization: Bearer {{token}}

Response: List of all admins with pagination
```

### Step 6: Create New Admin

```
POST /api/v1/admin/create
Header: Authorization: Bearer {{token}}

Body:
{
  "email": "newadmin@example.com",
  "name": "New Admin",
  "role": "admin"
}

Response: Admin created + password sent to email
```

### Step 7: Edit Admin Access

```
PUT /api/v1/admin/2/access
Header: Authorization: Bearer {{token}}

Body:
{
  "role": "admin",
  "status": "active"
}

Response: Admin access updated
```

### Step 8: Logout

```
POST /api/v1/auth/logout
Header: Authorization: Bearer {{token}}

Response: Logout successful
```

---

## API Reference

### User Profile Endpoints

| Method | Endpoint          | Description         | Auth   |
| ------ | ----------------- | ------------------- | ------ |
| GET    | `/api/v1/profile` | Get user profile    | ✅ JWT |
| PUT    | `/api/v1/profile` | Update user profile | ✅ JWT |

### Admin Management Endpoints

| Method | Endpoint                   | Description       | Auth   | Role       |
| ------ | -------------------------- | ----------------- | ------ | ---------- |
| GET    | `/api/v1/admin/list`       | List all admins   | ✅ JWT | Superadmin |
| POST   | `/api/v1/admin/create`     | Create new admin  | ✅ JWT | Superadmin |
| PUT    | `/api/v1/admin/:id/access` | Edit admin access | ✅ JWT | Superadmin |

### Authentication Endpoints

| Method | Endpoint                      | Description    | Auth   |
| ------ | ----------------------------- | -------------- | ------ |
| POST   | `/api/v1/auth/login`          | Login          | ❌ No  |
| POST   | `/api/v1/auth/logout`         | Logout         | ✅ JWT |
| POST   | `/api/v1/auth/check-email`    | Check email    | ❌ No  |
| POST   | `/api/v1/auth/send-otp`       | Send OTP       | ❌ No  |
| POST   | `/api/v1/auth/validate-otp`   | Validate OTP   | ❌ No  |
| POST   | `/api/v1/auth/reset-password` | Reset password | ❌ No  |

---

## Troubleshooting

### Error: "Missing authorization header"

**Solution:** Add header `Authorization: Bearer <token>` to request

### Error: "Invalid or expired token"

**Solution:** Login again to get new token

### Error: "Anda tidak memiliki akses"

**Solution:** User must be superadmin for admin management endpoints

### Error: "Email sudah terdaftar"

**Solution:** Use different email for new admin

### Error: "SMTP connection failed"

**Solution:** Check SMTP credentials in environment variables

---

## File Structure

```
d:\Lumoshive\aplikasi-pos-team-boolean\
├── internal/
│   ├── usecase/
│   │   └── admin.go                    ← Business logic
│   ├── adaptor/
│   │   └── admin_adaptor.go            ← HTTP handlers
│   ├── dto/
│   │   └── admin.go                    ← Request/Response models
│   ├── data/repository/
│   │   └── auth.go                     ← Database operations
│   └── wire/
│       └── wire.go                     ← Routes setup
├── pkg/
│   ├── utils/
│   │   └── token.go                    ← JWT functions
│   └── middleware/
│       └── auth.go                     ← Auth middleware
└── DOCUMENTATION/
    ├── ADMIN_USER_MANAGEMENT.md        ← Full API docs
    └── IMPLEMENTATION_GUIDE.md         ← Implementation details
```

---

## Key Features

### ✅ User Profile Management

- Get own profile
- Update name
- Update password
- Auto-hash password with bcrypt

### ✅ Admin Management (Superadmin Only)

- List all admins with pagination
- Create new admin with auto-generated password
- Send password via email
- Change admin role and status
- Prevent disabling only superadmin

### ✅ Security

- JWT authentication (24h expiry)
- Role-based access control
- Input validation
- Password hashing
- Email protection

### ✅ Email Integration

- Auto-send password on admin creation
- SMTP configuration support
- Professional email templates

---

## Common Operations

### Create New Admin User

```bash
curl -X POST http://localhost:8080/api/v1/admin/create \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "name": "Admin User",
    "role": "admin"
  }'
```

### List All Admins

```bash
curl -X GET "http://localhost:8080/api/v1/admin/list?page=1&limit=10" \
  -H "Authorization: Bearer <token>"
```

### Update User Profile

```bash
curl -X PUT http://localhost:8080/api/v1/profile \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "New Name",
    "password": "newpassword123"
  }'
```

### Change Admin Status

```bash
curl -X PUT http://localhost:8080/api/v1/admin/2/access \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "role": "admin",
    "status": "inactive"
  }'
```

---

## Database Seeding (Initial Setup)

### Create Initial Superadmin

```sql
-- Hash password: "password123" with bcrypt
INSERT INTO users (email, password, name, role, status, is_deleted)
VALUES (
  'superadmin@example.com',
  '$2a$10$...hashed_password...',
  'Super Admin',
  'superadmin',
  'active',
  false
);
```

---

## Important Notes

⚠️ **Security**

- Change `JWT_SECRET` in production!
- Use strong SMTP password
- Don't commit .env file
- Update password hash algorithm if needed

⚠️ **Email**

- Gmail requires App Passwords (not regular password)
- Check spam folder for test emails
- May need to enable "Less secure app access"

⚠️ **JWT Token**

- Token expires in 24 hours
- Must login again after expiration
- Token sent in response body (save to use)

---

## Support

📚 **Documentation**

- Full API Docs: `DOCUMENTATION/ADMIN_USER_MANAGEMENT.md`
- Implementation Guide: `DOCUMENTATION/IMPLEMENTATION_GUIDE.md`

🧪 **Testing**

- Postman Collection: `Postman Collection/Admin_User_Management.postman_collection.json`
- 9 pre-configured test requests

💬 **Code Comments**

- Inline documentation in all files
- Clear error messages
- Structured logging

---

**Status:** Ready for Production ✅  
**Last Updated:** January 20, 2025
