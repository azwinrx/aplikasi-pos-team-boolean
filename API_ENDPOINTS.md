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
| Menu - Categories    | 5              | ✅ NEW      |
| Menu - Products      | 6              | ✅ NEW      |
| Health Check         | 1              | ✅ Existing |
| **TOTAL**            | **35**         | ✅ READY    |

---

## 🍽️ MENU - CATEGORY ENDPOINTS (NEW ✨)

### 1. Get All Categories

```
GET /categories
```

**Query Parameters:**

- `page` (int): Page number (default: 1)
- `limit` (int): Items per page (default: 10)
- `search` (string): Search by category_name or description
- `sort_by` (string): category_name, created_at
- `sort_order` (string): asc, desc

**Response (200):**

```json
{
  "status": true,
  "message": "success get data",
  "data": [
    {
      "id": 1,
      "icon_category": "🍕",
      "category_name": "Pizza",
      "product_count": 20
    },
    {
      "id": 2,
      "icon_category": "🍔",
      "category_name": "Burger",
      "product_count": 15
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total_pages": 1,
    "total_items": 6
  }
}
```

---

### 2. Get Category by ID

```
GET /categories/:id
```

**Response (200):**

```json
{
  "status": true,
  "message": "success get category detail",
  "data": {
    "id": 1,
    "icon_category": "🍕",
    "category_name": "Pizza",
    "description": "Delicious pizza varieties",
    "product_count": 20,
    "created_at": "2026-01-31 10:00:00",
    "updated_at": "2026-01-31 10:00:00"
  }
}
```

---

### 3. Create Category

```
POST /categories
```

**Request:**

```json
{
  "icon_category": "🍕",
  "category_name": "Pizza",
  "description": "Delicious pizza varieties"
}
```

**Response (201):**

```json
{
  "status": true,
  "message": "success create category",
  "data": {
    "id": 1,
    "icon_category": "🍕",
    "category_name": "Pizza",
    "description": "Delicious pizza varieties",
    "product_count": 0,
    "created_at": "2026-01-31 10:00:00",
    "updated_at": "2026-01-31 10:00:00"
  }
}
```

---

### 4. Update Category

```
PUT /categories/:id
```

**Request:**

```json
{
  "icon_category": "🍕",
  "category_name": "Pizza Updated",
  "description": "Updated description"
}
```

**Response (200):**

```json
{
  "status": true,
  "message": "success update category",
  "data": {
    "id": 1,
    "icon_category": "🍕",
    "category_name": "Pizza Updated",
    "description": "Updated description",
    "product_count": 20,
    "created_at": "2026-01-31 10:00:00",
    "updated_at": "2026-01-31 11:00:00"
  }
}
```

---

### 5. Delete Category

```
DELETE /categories/:id
```

**Response (200):**

```json
{
  "status": true,
  "message": "success delete category",
  "data": null
}
```

**Note:** Category dengan products tidak bisa dihapus.

---

## 🍔 MENU - PRODUCT ENDPOINTS (NEW ✨)

### 1. Get All Products

```
GET /products
```

**Query Parameters:**

- `page` (int): Page number (default: 1)
- `limit` (int): Items per page (default: 10)
- `search` (string): Search by product_name or item_id
- `category_id` (int): Filter by category
- `is_available` (boolean): Filter by availability (true/false)
- `min_price` (float): Minimum price filter
- `max_price` (float): Maximum price filter
- `sort_order` (string): asc, desc

**Response (200):**

```json
{
  "status": true,
  "message": "success get data",
  "data": [
    {
      "id": 1,
      "product_image": "/images/chicken-parmesan.jpg",
      "product_name": "Chicken Parmesan",
      "item_id": "#22314644",
      "stock": 119,
      "category_name": "Chicken",
      "price": 55.0,
      "is_available": true,
      "availability": "in_stock"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total_pages": 1,
    "total_items": 7
  }
}
```

---

### 2. Get Products by Category

```
GET /products/category/:category_id
```

**Path Parameters:**

- `category_id` (int): Category ID to filter

**Query Parameters:**

- `page` (int): Page number (default: 1)
- `limit` (int): Items per page (default: 10)

**Response (200):**

```json
{
  "status": true,
  "message": "success get products by category",
  "data": [
    {
      "id": 1,
      "product_image": "/images/margherita-pizza.jpg",
      "product_name": "Margherita Pizza",
      "item_id": "#22314645",
      "stock": 85,
      "category_name": "Pizza",
      "price": 45.0,
      "is_available": true,
      "availability": "in_stock"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total_pages": 1,
    "total_items": 2
  }
}
```

---

### 3. Get Product by ID

```
GET /products/:id
```

**Response (200):**

```json
{
  "status": true,
  "message": "success get product detail",
  "data": {
    "id": 1,
    "product_image": "/images/chicken-parmesan.jpg",
    "product_name": "Chicken Parmesan",
    "item_id": "#22314644",
    "stock": 119,
    "category_id": 3,
    "category_name": "Chicken",
    "price": 55.0,
    "is_available": true,
    "availability": "in_stock",
    "created_at": "2026-01-31 10:00:00",
    "updated_at": "2026-01-31 10:00:00"
  }
}
```

---

### 4. Create Product

```
POST /products
```

**Request:**

```json
{
  "product_image": "/images/new-product.jpg",
  "product_name": "New Chicken Wings",
  "stock": 100,
  "category_id": 3,
  "price": 45.0
}
```

**Response (201):**

```json
{
  "status": true,
  "message": "success create product",
  "data": {
    "id": 8,
    "product_image": "/images/new-product.jpg",
    "product_name": "New Chicken Wings",
    "item_id": "#12345678",
    "stock": 100,
    "category_id": 3,
    "category_name": "Chicken",
    "price": 45.0,
    "is_available": true,
    "availability": "in_stock",
    "created_at": "2026-01-31 10:00:00",
    "updated_at": "2026-01-31 10:00:00"
  }
}
```

**Note:** `item_id` di-generate secara otomatis.

---

### 5. Update Product

```
PUT /products/:id
```

**Request:**

```json
{
  "product_image": "/images/updated-product.jpg",
  "product_name": "Updated Chicken Wings",
  "stock": 50,
  "category_id": 3,
  "price": 48.0
}
```

**Response (200):**

```json
{
  "status": true,
  "message": "success update product",
  "data": {
    "id": 8,
    "product_image": "/images/updated-product.jpg",
    "product_name": "Updated Chicken Wings",
    "item_id": "#12345678",
    "stock": 50,
    "category_id": 3,
    "category_name": "Chicken",
    "price": 48.0,
    "is_available": true,
    "availability": "in_stock",
    "created_at": "2026-01-31 10:00:00",
    "updated_at": "2026-01-31 11:00:00"
  }
}
```

**Note:** `is_available` otomatis ter-update berdasarkan stock:

- `stock > 0` → `is_available: true`, `availability: "in_stock"`
- `stock = 0` → `is_available: false`, `availability: "out_of_stock"`

---

### 6. Delete Product

```
DELETE /products/:id
```

**Response (200):**

```json
{
  "status": true,
  "message": "success delete product",
  "data": null
}
```

---

## � DASHBOARD ENDPOINTS

### 1. Get Dashboard Summary

```
GET /dashboard/summary
```

**Description:** Get comprehensive dashboard summary including daily sales, monthly sales, and table status.

**Response (200):**

```json
{
  "status": true,
  "message": "Dashboard summary retrieved successfully",
  "data": {
    "daily_sales": {
      "total_orders": 25,
      "total_revenue": 1500000.0,
      "total_tax": 165000.0,
      "average_order": 60000.0,
      "paid_orders": 20,
      "pending_orders": 3,
      "cancelled_orders": 2
    },
    "monthly_sales": {
      "total_orders": 350,
      "total_revenue": 21000000.0,
      "total_tax": 2310000.0,
      "average_order": 60000.0,
      "paid_orders": 300,
      "pending_orders": 30,
      "cancelled_orders": 20
    },
    "table_summary": {
      "total_tables": 15,
      "available_tables": 8,
      "occupied_tables": 5,
      "reserved_tables": 2
    }
  }
}
```

---

### 2. Get Popular Products

```
GET /dashboard/popular-products?limit=10
```

**Description:** Get list of popular products based on total sales (most sold first).

**Query Parameters:**
| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| limit | int | No | 10 | Maximum number of products to return |

**Response (200):**

```json
{
  "status": true,
  "message": "Popular products retrieved successfully",
  "data": [
    {
      "product_id": 5,
      "product_image": "/images/nasi-goreng.jpg",
      "product_name": "Nasi Goreng Spesial",
      "item_id": "#12345678",
      "category_name": "Main Course",
      "price": 35000.0,
      "total_sold": 150,
      "total_revenue": 5250000.0
    },
    {
      "product_id": 3,
      "product_image": "/images/ayam-bakar.jpg",
      "product_name": "Ayam Bakar",
      "item_id": "#87654321",
      "category_name": "Chicken",
      "price": 45000.0,
      "total_sold": 120,
      "total_revenue": 5400000.0
    }
  ]
}
```

---

### 3. Get New Products

```
GET /dashboard/new-products?limit=10
```

**Description:** Get list of new products created in the last 30 days.

**Query Parameters:**
| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| limit | int | No | 10 | Maximum number of products to return |

**Response (200):**

```json
{
  "status": true,
  "message": "New products retrieved successfully",
  "data": [
    {
      "product_id": 10,
      "product_image": "/images/new-menu.jpg",
      "product_name": "Sate Kambing Premium",
      "item_id": "#11223344",
      "category_name": "Grilled",
      "price": 65000.0,
      "stock": 30,
      "availability": "in_stock",
      "created_at": "2026-01-28T10:30:00Z",
      "days_ago": 3
    },
    {
      "product_id": 9,
      "product_image": "/images/smoothie.jpg",
      "product_name": "Tropical Smoothie",
      "item_id": "#55667788",
      "category_name": "Beverages",
      "price": 28000.0,
      "stock": 0,
      "availability": "out_of_stock",
      "created_at": "2026-01-20T14:00:00Z",
      "days_ago": 11
    }
  ]
}
```

---

## �🔧 TESTING

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

4. **Test Menu Endpoints:**
   - GET /categories → List all categories
   - POST /categories → Create new category
   - GET /products?category_id=1 → Filter products by category
   - POST /products → Create new product

5. **Test Other Endpoints:**
   - Use existing Postman collection untuk staff, inventory, orders

---

## 📝 NOTES

- Semua endpoint return consistent JSON format
- Error response menggunakan HTTP status code yang sesuai
- JWT token berlaku 24 jam (configurable)
- OTP berlaku 10 menit
- Pagination default: page=1, limit=10
- Sorting: sort_by=name, sort_order=asc/desc
- Product status otomatis berdasarkan stock level
- Category dengan products tidak bisa dihapus

---

**Last Updated: 31 January 2026**
**Total API Endpoints: 37**
**Status: ✅ READY FOR TESTING**
