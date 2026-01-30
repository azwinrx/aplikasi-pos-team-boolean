# Notification System - Integration & Testing Guide

## 📋 Table of Contents

1. [Pre-Testing Setup](#pre-testing-setup)
2. [Manual Testing Guide](#manual-testing-guide)
3. [Integration Tests](#integration-tests)
4. [Troubleshooting](#troubleshooting)
5. [Performance Testing](#performance-testing)

---

## Pre-Testing Setup

### 1. Environment Configuration

Ensure `.env` file contains:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=pos_database

JWT_SECRET=your_secret_key
JWT_EXPIRY=24h

SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your_email@gmail.com
SMTP_PASSWORD=your_app_password
```

### 2. Database Verification

```bash
# Connect to PostgreSQL
psql -U postgres -d pos_database

# Check if notifications table exists
\dt notifications

# Expected output shows table with columns: id, user_id, title, message, type, status, readed_at, data, created_at, updated_at, deleted_at
```

### 3. Start Application

```bash
# Build and run
go build -o pos-app main.go
./pos-app

# Expected output:
# INFO    database       Auto migration completed successfully!
# Routes registered successfully
# Server is running on :8080
```

### 4. Postman Setup

```bash
1. Open Postman
2. Click "Import"
3. Select: Postman Collection/POS_Notification_API.postman_collection.json
4. Set variables:
   - base_url: http://localhost:8080
   - token: [from login endpoint]
```

---

## Manual Testing Guide

### Test 1: Authentication Setup

**Objective**: Get JWT token for testing

```bash
# Step 1: Login
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "admin@pos.com",
  "password": "admin123"
}

# Expected Response (200 OK):
{
  "code": 200,
  "message": "Login berhasil",
  "data": {
    "user": {
      "id": 1,
      "email": "admin@pos.com",
      "name": "Admin User",
      "role": "admin"
    },
    "token": "eyJhbGciOiJIUzI1NiIs..."
  }
}

# Step 2: Copy token value
# Set in Postman: {{token}} = eyJhbGciOiJIUzI1NiIs...
```

### Test 2: List Notifications (Empty State)

**Objective**: Verify list endpoint works, database is connected

```bash
# Request
GET /api/v1/notifications?page=1&limit=10
Authorization: Bearer {{token}}

# Expected Response (200 OK):
{
  "code": 200,
  "message": "Daftar notifikasi berhasil diambil",
  "data": {
    "data": [],
    "total": 0,
    "page": 1,
    "limit": 10,
    "total_pages": 1,
    "unread_count": 0
  }
}

# Verification:
# ✅ Status 200
# ✅ Empty data array
# ✅ Total 0
# ✅ Unread count 0
```

### Test 3: Create Notification (via SQL)

**Objective**: Prepare test data

```sql
-- Insert test notification
INSERT INTO notifications (user_id, title, message, type, status, created_at, updated_at)
VALUES (
  1,
  'Test Pesanan Baru',
  'Pesanan #ORD001 dari meja T01 sudah diterima',
  'order',
  'new',
  CURRENT_TIMESTAMP,
  CURRENT_TIMESTAMP
);

-- Insert more test data
INSERT INTO notifications (user_id, title, message, type, status, created_at, updated_at)
VALUES
  (1, 'Pembayaran Dikonfirmasi', 'Pembayaran sudah diterima', 'payment', 'new', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  (1, 'Stok Menipis', 'Coca Cola stok tinggal 5', 'alert', 'new', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
  (1, 'Pesanan Siap', 'Pesanan #ORD001 sudah siap', 'order', 'readed', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Verify
SELECT COUNT(*) FROM notifications WHERE user_id = 1;
-- Expected: 4
```

### Test 4: List Notifications (With Data)

**Objective**: Verify list endpoint returns correct data

```bash
# Request
GET /api/v1/notifications?page=1&limit=10
Authorization: Bearer {{token}}

# Expected Response (200 OK):
{
  "code": 200,
  "message": "Daftar notifikasi berhasil diambil",
  "data": {
    "data": [
      {
        "id": 1,
        "user_id": 1,
        "title": "Test Pesanan Baru",
        "message": "Pesanan #ORD001 dari meja T01 sudah diterima",
        "type": "order",
        "status": "new",
        "readed_at": null,
        "created_at": "2025-01-30T...",
        "updated_at": "2025-01-30T..."
      },
      ...
    ],
    "total": 4,
    "page": 1,
    "limit": 10,
    "total_pages": 1,
    "unread_count": 3
  }
}

# Verification:
# ✅ Status 200
# ✅ Data array has 4 items
# ✅ Unread count is 3 (only "new" status)
# ✅ Pagination info correct
# ✅ Each notification has all required fields
```

### Test 5: Filter by Status (New)

**Objective**: Verify status filtering

```bash
# Request
GET /api/v1/notifications?status=new&limit=10
Authorization: Bearer {{token}}

# Expected Response:
# ✅ Only notifications with status="new" returned
# ✅ unread_count = total = 3
# ✅ Notifications with status="readed" not included
```

### Test 6: Filter by Type

**Objective**: Verify type filtering

```bash
# Request - Order only
GET /api/v1/notifications?type=order&limit=10
Authorization: Bearer {{token}}

# Expected Response:
# ✅ Only type="order" notifications returned
# ✅ Payment and alert notifications filtered out

# Request - Payment only
GET /api/v1/notifications?type=payment&limit=10
Authorization: Bearer {{token}}

# Expected Response:
# ✅ Only type="payment" notifications returned
```

### Test 7: Pagination

**Objective**: Verify pagination works correctly

```bash
# Request - Page 1 with limit 2
GET /api/v1/notifications?page=1&limit=2
Authorization: Bearer {{token}}

# Expected Response:
# ✅ Data array has 2 items
# ✅ page = 1
# ✅ limit = 2
# ✅ total = 4
# ✅ total_pages = 2

# Request - Page 2 with limit 2
GET /api/v1/notifications?page=2&limit=2
Authorization: Bearer {{token}}

# Expected Response:
# ✅ Data array has 2 items
# ✅ page = 2
# ✅ First item is different from page 1
# ✅ total_pages = 2
```

### Test 8: Sorting

**Objective**: Verify sorting works

```bash
# Request - Ascending sort
GET /api/v1/notifications?sort_by=created_at&sort_order=asc&limit=10
Authorization: Bearer {{token}}

# Expected Response:
# ✅ Notifications sorted by created_at ascending (oldest first)

# Request - Descending sort
GET /api/v1/notifications?sort_by=created_at&sort_order=desc&limit=10
Authorization: Bearer {{token}}

# Expected Response:
# ✅ Notifications sorted by created_at descending (newest first)
```

### Test 9: Update Status to Readed

**Objective**: Mark notification as read

```bash
# Before Update
# Notification 1: status = "new", readed_at = null

# Request
PUT /api/v1/notifications/1/status
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "status": "readed"
}

# Expected Response (200 OK):
{
  "code": 200,
  "message": "Status notifikasi berhasil diubah",
  "data": {
    "id": 1,
    "status": "readed",
    "readed_at": "2025-01-30T15:30:00Z",
    "message": "Status notifikasi berhasil diubah menjadi readed"
  }
}

# Verification:
# ✅ Status 200
# ✅ Status changed to "readed"
# ✅ ReadedAt timestamp set to current time

# Verify in database:
SELECT status, readed_at FROM notifications WHERE id = 1;
-- Expected: "readed" | 2025-01-30 15:30:00
```

### Test 10: Update Status Back to New

**Objective**: Verify status can be changed back

```bash
# Request
PUT /api/v1/notifications/1/status
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "status": "new"
}

# Expected Response (200 OK):
# ✅ Status changed back to "new"
# ✅ ReadedAt should still be set (not cleared)

# Verification:
SELECT status, readed_at FROM notifications WHERE id = 1;
-- Expected: "new" | 2025-01-30 15:30:00 (unchanged)
```

### Test 11: Delete Notification

**Objective**: Delete notification (soft delete)

```bash
# Before Delete
SELECT deleted_at FROM notifications WHERE id = 1;
-- Expected: NULL

# Request
DELETE /api/v1/notifications/1
Authorization: Bearer {{token}}

# Expected Response (200 OK):
{
  "code": 200,
  "message": "Notifikasi berhasil dihapus",
  "data": {
    "id": 1,
    "message": "Notifikasi berhasil dihapus"
  }
}

# Verification:
# ✅ Status 200
# ✅ Response message confirms deletion

# Check database (soft delete):
SELECT * FROM notifications WHERE id = 1;
-- Expected: deleted_at = 2025-01-30 15:35:00 (set to current time)

# Check list endpoint:
GET /api/v1/notifications
-- Expected: Deleted notification NOT in list (filtered by deleted_at IS NULL)
```

### Test 12: Authorization Test

**Objective**: Verify user isolation

```bash
# Create notification for user 2
INSERT INTO notifications (user_id, title, message, type, status, created_at, updated_at)
VALUES (2, 'User 2 Notification', 'This belongs to user 2', 'order', 'new', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

# Login as user 1
# Get notifications as user 1
GET /api/v1/notifications
Authorization: Bearer {{token_user1}}

# Verification:
# ✅ User 1 can only see notifications with user_id = 1
# ✅ User 2's notification NOT in list

# Try to update user 2's notification as user 1
PUT /api/v1/notifications/5/status (where id=5 belongs to user 2)
Authorization: Bearer {{token_user1}}
Content-Type: application/json

{
  "status": "readed"
}

# Expected Response (403 Forbidden):
{
  "code": 403,
  "message": "anda tidak memiliki akses ke notifikasi ini"
}

# Verification:
# ✅ Status 403
# ✅ Cannot access other user's notification
# ✅ Authorization working correctly
```

### Test 13: Error Handling

**Objective**: Verify error responses

```bash
# Test 1: Invalid notification ID
DELETE /api/v1/notifications/99999
Authorization: Bearer {{token}}

# Expected Response (404 Not Found):
{
  "code": 404,
  "message": "notifikasi tidak ditemukan"
}

# Test 2: Missing authorization header
GET /api/v1/notifications

# Expected Response (401 Unauthorized):
{
  "code": 401,
  "message": "User tidak terautentikasi"
}

# Test 3: Invalid status value
PUT /api/v1/notifications/1/status
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "status": "invalid_status"
}

# Expected Response (400 Bad Request):
{
  "code": 400,
  "message": "status 'invalid_status' tidak valid. gunakan 'new' atau 'readed'"
}

# Test 4: Invalid limit value
GET /api/v1/notifications?limit=999
Authorization: Bearer {{token}}

# Expected Response:
# ✅ Limit capped to 100
# ✅ Returns 100 notifications max
```

---

## Integration Tests

### Integration Test 1: Order → Notification

**Objective**: Create notification when order is created

```bash
# In order.go repository or usecase, add:
notification := &entity.Notification{
    UserID:  staffID,
    Title:   "Pesanan Baru",
    Message: "Pesanan sudah diterima",
    Type:    "order",
    Status:  "new",
}
notificationRepo.CreateNotification(ctx, notification)

# Test:
# 1. Create order
# 2. Check if notification appears in list
GET /api/v1/notifications
Authorization: Bearer {{token}}
# ✅ Notification should appear in list
```

### Integration Test 2: Payment → Notification

**Objective**: Create notification when payment is confirmed

```bash
# Similar to order integration
# Trigger payment confirmation
# Verify notification created with type="payment"
```

---

## Troubleshooting

### Issue 1: 401 Unauthorized

**Problem**: Getting 401 when token is present
**Solution**:

```bash
1. Verify token format: "Bearer <token>"
2. Check if token expired: Re-login
3. Verify Authorization header exists
4. Check user exists in database
```

### Issue 2: 500 Internal Server Error

**Problem**: Getting 500 error
**Solution**:

```bash
1. Check application logs
2. Verify database connection
3. Verify table exists: SELECT * FROM notifications LIMIT 1;
4. Check user_id exists in users table
```

### Issue 3: Notifications Not Appearing

**Problem**: List endpoint returns empty even after creating notification
**Solution**:

```bash
1. Check database directly:
   SELECT COUNT(*) FROM notifications WHERE user_id = 1;
2. Verify user_id matches logged-in user
3. Check deleted_at is NULL
4. Verify database migration ran: \dt notifications
```

### Issue 4: Cannot Update Status

**Problem**: Update status returns 403
**Solution**:

```bash
1. Verify you own the notification (user_id matches)
2. Verify notification exists: SELECT * FROM notifications WHERE id = X;
3. Check token belongs to correct user
```

---

## Performance Testing

### Load Test 1: List Large Number of Notifications

```bash
# Create 1000 test notifications
INSERT INTO notifications (user_id, title, message, type, status, created_at, updated_at)
SELECT 1, 'Notification ' || generate_series(1, 1000), 'Message', 'order', 'new', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP;

# Test with pagination
GET /api/v1/notifications?page=1&limit=100
Authorization: Bearer {{token}}

# Measure response time
# Expected: < 500ms

# Test with filters
GET /api/v1/notifications?status=new&type=order&page=1&limit=100
Authorization: Bearer {{token}}

# Expected: < 300ms (indexes should speed up filtering)
```

### Load Test 2: Concurrent Updates

```bash
# Use load testing tool (ab, hey, or JMeter)
# Simulate 10 concurrent users updating notifications

hey -n 100 -c 10 -m PUT \
  -H "Authorization: Bearer {{token}}" \
  -d '{"status":"readed"}' \
  http://localhost:8080/api/v1/notifications/1/status

# Expected:
# ✅ All requests succeed
# ✅ Average latency < 100ms
# ✅ No data corruption
```

---

## Checklist: Testing Complete

- [ ] Pre-testing setup verified
- [ ] Database connection working
- [ ] Application running on port 8080
- [ ] Postman collection imported
- [ ] JWT token obtained
- [ ] Test 1-4: Basic operations (List, Get, Create, Read)
- [ ] Test 5-7: Filtering and pagination
- [ ] Test 8: Sorting
- [ ] Test 9-10: Update operations
- [ ] Test 11: Delete operation
- [ ] Test 12: Authorization
- [ ] Test 13: Error handling
- [ ] Integration tests passed
- [ ] Performance acceptable
- [ ] Documentation reviewed
- [ ] Ready for production

---

## Success Criteria

✅ **All Tests Passing**

- List endpoint returns correct data
- Filtering works correctly
- Pagination working
- Update operations successful
- Delete operations successful
- Authorization enforced
- Error handling proper
- Performance acceptable

✅ **Ready for Production**

- No critical errors
- Data integrity verified
- User isolation confirmed
- Logging working
- Documentation complete

---

**Status**: Ready for Comprehensive Testing ✨
