# API ENDPOINTS DOCUMENTATION

## 📌 Base URL

```
http://localhost:8080/api/v1
```

---

## 🔐 AUTHENTICATION ENDPOINTS (NEW ✨)

### 1. Login

```
POST /auth/login
```

**Request:**

```json
{
  "email": "admin@pos.com",
  "password": "admin123"
}
```

**Response (200):**

```json
{
  "status": true,
  "message": "Login successful",
  "data": {
    "id": 1,
    "email": "admin@pos.com",
    "name": "Admin User",
    "role": "admin",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_at": 1640000000
  }
}
```

---

### 2. Check Email

```
POST /auth/check-email
```

**Request:**

```json
{
  "email": "admin@pos.com"
}
```

**Response (200):**

```json
{
  "status": true,
  "message": "email already registered",
  "data": {
    "email": "admin@pos.com",
    "exists": true,
    "message": "email already registered"
  }
}
```

---

### 3. Send OTP

```
POST /auth/send-otp
```

**Request:**

```json
{
  "email": "admin@pos.com",
  "purpose": "password_reset"
}
```

**Response (200):**

```json
{
  "status": true,
  "message": "OTP has been sent to your email. Valid for 10 minutes.",
  "data": {
    "email": "admin@pos.com",
    "message": "OTP has been sent to your email. Valid for 10 minutes."
  }
}
```

---

### 4. Validate OTP

```
POST /auth/validate-otp
```

**Request:**

```json
{
  "email": "admin@pos.com",
  "otp_code": "123456",
  "purpose": "password_reset"
}
```

**Response (200):**

```json
{
  "status": true,
  "message": "OTP is valid",
  "data": {
    "valid": true,
    "message": "OTP is valid",
    "token": "123456"
  }
}
```

---

### 5. Reset Password

```
POST /auth/reset-password
```

**Request:**

```json
{
  "email": "admin@pos.com",
  "otp_code": "123456",
  "new_password": "newpassword123",
  "purpose": "password_reset"
}
```

**Response (200):**

```json
{
  "status": true,
  "message": "Password reset successfully. You can now login with your new password.",
  "data": {
    "email": "admin@pos.com",
    "message": "Password reset successfully. You can now login with your new password."
  }
}
```

---

## 👥 STAFF MANAGEMENT ENDPOINTS

### 6. Get All Staff

```
GET /staff?page=1&limit=10&search=&sort_by=name&sort_order=asc&role=admin
```

**Response (200):**

```json
{
  "status": true,
  "message": "success get data",
  "data": [...],
  "pagination": {
    "total_items": 5,
    "total_pages": 1,
    "current_page": 1,
    "limit": 10
  }
}
```

---

### 7. Create Staff

```
POST /staff
```

**Request:**

```json
{
  "full_name": "John Doe",
  "email": "john@example.com",
  "role": "staff",
  "phone_number": "081234567890",
  "salary": 3000000,
  "address": "Jakarta"
}
```

---

### 8. Get Staff by ID

```
GET /staff/:id
```

---

### 9. Get Staff by Email

```
GET /staff/email?email=john@example.com
```

---

### 10. Update Staff

```
PUT /staff/:id
```

---

### 11. Delete Staff

```
DELETE /staff/:id
```

---

## 📦 INVENTORY MANAGEMENT ENDPOINTS

### 12. Get All Inventories

```
GET /inventories
```

---

### 13. Filter Inventories

```
GET /inventories/filter?status=active&category=beverage&min_price=0&max_price=100
```

---

### 14. Create Inventory

```
POST /inventories
```

**Request:**

```json
{
  "name": "Coca Cola 1L",
  "category": "beverage",
  "quantity": 100,
  "unit": "litre",
  "min_stock": 50,
  "retail_price": 15.5,
  "status": "active"
}
```

---

### 15. Update Inventory

```
PUT /inventories/:id
```

---

### 16. Delete Inventory

```
DELETE /inventories/:id
```

---

## 🍽️ ORDERS MANAGEMENT ENDPOINTS

### 17. Get All Orders

```
GET /orders
```

---

### 18. Create Order

```
POST /orders
```

**Request:**

```json
{
  "customer_name": "Budi",
  "table_id": 1,
  "payment_method_id": 1,
  "items": [
    {
      "product_id": 1,
      "quantity": 2
    }
  ]
}
```

---

### 19. Update Order

```
PUT /orders/:id
```

---

### 20. Delete Order

```
DELETE /orders/:id
```

---

### 21. Get All Tables

```
GET /orders/tables
```

---

### 22. Get Payment Methods

```
GET /orders/payment-methods
```

---

### 23. Get Available Chairs

```
GET /orders/available-chairs
```

---

## ✅ HEALTH CHECK

### 24. Health Check

```
GET /health
```

**Response (200):**

```json
{
  "status": "healthy"
}
```

---

## 📊 SUMMARY

| Category             | Endpoint Count | Status      |
| -------------------- | -------------- | ----------- |
| Authentication       | 5              | ✅ NEW      |
| Staff Management     | 6              | ✅ Existing |
| Inventory Management | 5              | ✅ Existing |
| Orders Management    | 7              | ✅ Existing |
| Health Check         | 1              | ✅ Existing |
| **TOTAL**            | **24**         | ✅ READY    |

---

## 🔧 TESTING

### **Quick Test Flow:**

1. **Start Server:**

   ```bash
   go run main.go --migrate --seed
   ```

2. **Import Postman Collection:**
   - Open Postman
   - Import: `Postman Collection/Authentication.postman_collection.json`

3. **Test Authentication:**
   - Login → Get JWT token
   - Check Email → Validate email exists
   - Send OTP → Receive OTP code
   - Validate OTP → Confirm OTP valid
   - Reset Password → Change password

4. **Test Other Endpoints:**
   - Use existing Postman collection untuk staff, inventory, orders

---

## 📝 NOTES

- Semua endpoint return consistent JSON format
- Error response menggunakan HTTP status code yang sesuai
- JWT token berlaku 24 jam (configurable)
- OTP berlaku 10 menit
- Pagination default: page=1, limit=10
- Sorting: sort_by=name, sort_order=asc/desc

---

**Last Updated: 30 January 2026**
**Total API Endpoints: 24**
**Status: ✅ READY FOR TESTING**
