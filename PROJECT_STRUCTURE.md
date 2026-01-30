# PROJECT STRUCTURE - UPDATED

```
aplikasi-pos-team-boolean/
├── go.mod
├── go.sum
├── main.go
├── README.md
│
├── 📄 DOCUMENTATION FILES (NEW)
├── AUTHENTICATION_API.md ........................ API documentation
├── IMPLEMENTATION_STATUS_AUTH.md ............... Implementation status
├── SUMMARY_AUTH.md ............................. Quick summary
│
├── pkg/
│   ├── database/
│   │   ├── database.go ......................... Database connection
│   │   └── migration.go ........................ Migration + Seeding (UPDATED)
│   │
│   ├── middleware/
│   │   └── logging.go .......................... Logging middleware
│   │
│   └── utils/
│       ├── config.go ........................... Configuration reader
│       ├── logger.go ........................... Zap logger setup
│       ├── password_hash.go .................... Password hashing
│       ├── response.go ......................... Response helper
│       ├── token.go ............................ JWT token generation
│       ├── validator.go ........................ Input validation
│       ├── utils.go ............................ Utility functions
│       └── email_service.go (NEW) ............. SMTP email service
│
├── internal/
│   ├── adaptor/
│   │   ├── adaptor.go (UPDATED) ............... Main adaptor
│   │   ├── staff_adaptor.go ................... Staff HTTP handler
│   │   ├── inventories_adaptor.go ............ Inventories HTTP handler
│   │   ├── order_adaptor.go .................. Orders HTTP handler
│   │   └── auth_adaptor.go (NEW) ............. Auth HTTP handler
│   │
│   ├── data/
│   │   ├── entity/
│   │   │   ├── staff.go ....................... Staff model
│   │   │   ├── inventories.go ................ Inventories model
│   │   │   ├── order.go ...................... Order & OrderItem model
│   │   │   ├── payment_method.go ............. PaymentMethod model
│   │   │   ├── table.go ...................... Table model
│   │   │   └── user.go (NEW) ................. User & OTP model
│   │   │
│   │   └── repository/
│   │       ├── repository.go (UPDATED) ....... Main repository
│   │       ├── staff.go ....................... Staff repository
│   │       ├── inventories.go ................ Inventories repository
│   │       ├── order.go ...................... Order repository
│   │       └── auth.go (NEW) ................. Auth repository
│   │
│   ├── dto/
│   │   ├── staff.go ........................... Staff DTOs
│   │   ├── order.go ........................... Order DTOs
│   │   ├── inventories.go .................... Inventories DTOs
│   │   ├── pagination.go ..................... Pagination DTO
│   │   └── auth.go (NEW) ..................... Auth DTOs
│   │
│   ├── usecase/
│   │   ├── usecase.go (UPDATED) .............. Main usecase
│   │   ├── staff.go ........................... Staff usecase
│   │   ├── inventories.go .................... Inventories usecase
│   │   ├── order.go ........................... Order usecase
│   │   └── auth.go (NEW) ..................... Auth usecase
│   │
│   └── wire/
│       └── wire.go (UPDATED) ................. Dependency injection
│
├── .env.example (UPDATED) ..................... Configuration template
│
├── BACKUP SQL/
│   └── aplikasi-pos-team-boolean.sql ........ Database backup
│
├── SQL Queries/
│   ├── Create.sql ............................ Create table queries
│   └── Insert.sql ............................ Insert sample data
│
└── Postman Collection/
    ├── Authentication.postman_collection.json (NEW)
    └── aplikasi-pos-team-boolean Copy.postman_collection.json
```

## 📊 Architecture Overview

```
HTTP Request
    ↓
┌─────────────────────────────────────────┐
│         ADAPTOR (HTTP Handler)          │
│    - Parse request, return response     │
│    - Validasi input, error handling     │
│    - Logging request/response           │
└──────────────────┬──────────────────────┘
                   ↓
┌─────────────────────────────────────────┐
│        USECASE (Business Logic)         │
│    - Implement business rules           │
│    - Call repository methods            │
│    - Error handling & validation        │
│    - Service integration (email)        │
└──────────────────┬──────────────────────┘
                   ↓
┌─────────────────────────────────────────┐
│      REPOSITORY (Data Access)           │
│    - Query database dengan GORM         │
│    - Manage entities                    │
│    - Logging database operations        │
└──────────────────┬──────────────────────┘
                   ↓
┌─────────────────────────────────────────┐
│       DATABASE (PostgreSQL)             │
│    - Persistent data storage            │
│    - Tables: users, otps, staff, etc.   │
└─────────────────────────────────────────┘
```

## 🔄 Request Flow Example (Login)

```
1. POST /api/v1/auth/login
   ↓
2. AuthAdaptor.Login()
   - Parse JSON body
   - Validate input binding
   ↓
3. AuthUseCase.Login()
   - Get user dari database
   - Validate password
   - Check is_deleted flag
   - Generate JWT token
   ↓
4. AuthRepository.GetUserByEmail()
   - Query database
   - Return user data
   ↓
5. Response dengan JWT token
   {
     "status": true,
     "data": {
       "token": "eyJ...",
       "user": {...}
     }
   }
```

## 📦 Key Components

### **Authentication Flow**

```
Login → CheckEmail → SendOTP → ValidateOTP → ResetPassword
  ↓        ↓            ↓           ↓            ↓
User    Email      OTP Email    Verify       Password
Exists  Exists?    Sent         OTP Valid    Changed
```

### **Email Service**

```
SendOTP() → Generate 6-digit code → Save to DB → Send via SMTP
Validate() → Check code & expiry → Mark as used → Allow reset
```

### **Security**

```
Password → Hashing (bcrypt) → Store in DB
Token → Generate JWT → Include user info & role → Expire in 24h
OTP → Generate 6-digit → Expire in 10 min → Mark as used
User → is_deleted flag → Prevent login if deleted
```

## 🎯 Features Implemented

### ✅ Authentication (NEW)

- [x] Login API
- [x] Check Email API
- [x] Send OTP API
- [x] Validate OTP API
- [x] Reset Password API
- [x] Email Service
- [x] User Entity & Repository
- [x] OTP Entity & Repository
- [x] JWT Token Generation
- [x] Password Hashing

### ✅ Staff Management (EXISTING)

- [x] List Staff
- [x] Create Staff
- [x] Update Staff
- [x] Get Staff by ID
- [x] Delete Staff
- [x] Pagination & Sorting

### ✅ Inventory Management (EXISTING)

- [x] List Inventories
- [x] Create Inventory
- [x] Update Inventory
- [x] Delete Inventory
- [x] Filter by multiple criteria

### ✅ Orders Management (EXISTING)

- [x] List Orders
- [x] Create Order
- [x] Update Order
- [x] Delete Order
- [x] Get Available Tables
- [x] Get Payment Methods

---

**Total Files Created: 9**
**Total Files Modified: 6**
**Total API Endpoints: 5 (Authentication)**
**Status: ✅ COMPLETE & READY FOR TESTING**
