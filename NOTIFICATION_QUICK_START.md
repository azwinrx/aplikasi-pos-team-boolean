# Notification System - Quick Start Guide

## Implementation Complete! ✅

Notification System dengan 3 fitur utama telah selesai diimplementasikan:

1. **List notifikasi** dengan filtering & pagination
2. **Update status notifikasi** (new ↔ readed)
3. **Hapus notifikasi** (soft delete)

---

## Files Created / Modified

### New Files ✨

```
internal/data/entity/notification.go
internal/data/repository/notification.go
internal/usecase/notification.go
internal/adaptor/notification_adaptor.go
internal/dto/notification.go (created earlier)
DOCS_NOTIFICATION_API.md
NOTIFICATION_IMPLEMENTATION_SUMMARY.md
Postman Collection/POS_Notification_API.postman_collection.json
```

### Modified Files 🔄

```
pkg/database/migration.go
internal/data/repository/repository.go
internal/usecase/usecase.go
internal/adaptor/adaptor.go
internal/wire/wire.go
```

---

## Database Schema

```sql
CREATE TABLE notifications (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    type VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'new',
    readed_at TIMESTAMP NULL,
    data JSONB NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,

    INDEX idx_user_id (user_id),
    INDEX idx_type (type),
    INDEX idx_status (status),
    INDEX idx_deleted_at (deleted_at)
);
```

---

## API Endpoints

### 1. List Notifications

```bash
GET /api/v1/notifications?page=1&limit=10&status=new&type=order&sort_by=created_at&sort_order=desc
Authorization: Bearer {JWT_TOKEN}
```

**Response:**

```json
{
  "code": 200,
  "message": "Daftar notifikasi berhasil diambil",
  "data": {
    "data": [
      {
        "id": 1,
        "user_id": 5,
        "title": "Pesanan Baru",
        "message": "Pesanan #ORD001 dari meja T01 sudah diterima",
        "type": "order",
        "status": "new",
        "readed_at": null,
        "data": "...",
        "created_at": "2025-01-30T10:30:00Z"
      }
    ],
    "total": 25,
    "page": 1,
    "limit": 10,
    "total_pages": 3,
    "unread_count": 8
  }
}
```

### 2. Update Status

```bash
PUT /api/v1/notifications/1/status
Authorization: Bearer {JWT_TOKEN}
Content-Type: application/json

{
  "status": "readed"
}
```

**Response:**

```json
{
  "code": 200,
  "message": "Status notifikasi berhasil diubah",
  "data": {
    "id": 1,
    "status": "readed",
    "readed_at": "2025-01-30T10:35:00Z",
    "message": "Status notifikasi berhasil diubah menjadi readed"
  }
}
```

### 3. Delete Notification

```bash
DELETE /api/v1/notifications/1
Authorization: Bearer {JWT_TOKEN}
```

**Response:**

```json
{
  "code": 200,
  "message": "Notifikasi berhasil dihapus",
  "data": {
    "id": 1,
    "message": "Notifikasi berhasil dihapus"
  }
}
```

---

## Testing dengan Postman

### Import Collection

1. Buka Postman
2. Click `Import` → `Upload File`
3. Select `Postman Collection/POS_Notification_API.postman_collection.json`
4. Collection akan di-import otomatis

### Setup Variables

1. Edit collection
2. Tab `Variables`
3. Set:
   - `base_url`: `http://localhost:8080`
   - `token`: Paste JWT token dari login endpoint

### Test Requests

Jalankan requests dalam urutan:

1. **List All Notifications** - Test list kosong atau ada data
2. **List Unread** - Filter status=new
3. **List by Type** - Filter type=order, payment, etc.
4. **Update Status** - Mark sebagai readed
5. **Delete** - Hapus notifikasi

---

## How to Use in Other Features

### Create Notification dari Order

```go
// Di order.go usecase/repository
notification := &entity.Notification{
    UserID:  staffUserID,
    Title:   "Pesanan Baru",
    Message: fmt.Sprintf("Pesanan #%d dari meja %s sudah diterima", orderID, tableNumber),
    Type:    "order",
    Status:  "new",
    Data:    `{"order_id": 1, "table_number": "T01"}`,
}

err := notificationUseCase.CreateNotification(ctx, notification)
```

### Create Notification dari Payment

```go
// Di payment validation
notification := &entity.Notification{
    UserID:  managerUserID,
    Title:   "Pembayaran Dikonfirmasi",
    Message: fmt.Sprintf("Pembayaran untuk pesanan #%d sebesar Rp%v sudah dikonfirmasi", orderID, amount),
    Type:    "payment",
    Status:  "new",
    Data:    `{"payment_method": "cash", "amount": 150000}`,
}

err := notificationUseCase.CreateNotification(ctx, notification)
```

---

## Notification Types & Examples

| Type      | Use Case         | Example                          |
| --------- | ---------------- | -------------------------------- |
| `order`   | Order related    | Pesanan baru, order siap         |
| `payment` | Payment related  | Pembayaran dikonfirmasi, pending |
| `system`  | System events    | Database backup, update tersedia |
| `alert`   | Important alerts | Stok menipis, error sistem       |

---

## Query Parameter Examples

```bash
# Get only unread notifications
GET /api/v1/notifications?status=new

# Get only order type
GET /api/v1/notifications?type=order

# Get unread orders
GET /api/v1/notifications?status=new&type=order

# Pagination
GET /api/v1/notifications?page=2&limit=20

# Custom sorting
GET /api/v1/notifications?sort_by=created_at&sort_order=asc

# Combination
GET /api/v1/notifications?status=new&type=payment&page=1&limit=10&sort_by=created_at&sort_order=desc
```

---

## Error Responses

### 400 Bad Request

```json
{
  "code": 400,
  "message": "ID notifikasi tidak valid"
}
```

### 401 Unauthorized

```json
{
  "code": 401,
  "message": "User tidak terautentikasi"
}
```

### 403 Forbidden

```json
{
  "code": 403,
  "message": "anda tidak memiliki akses ke notifikasi ini"
}
```

### 404 Not Found

```json
{
  "code": 404,
  "message": "notifikasi tidak ditemukan"
}
```

### 500 Internal Server Error

```json
{
  "code": 500,
  "message": "Gagal mengambil daftar notifikasi: [error details]"
}
```

---

## Security Checklist

- ✅ User hanya bisa akses notifikasi mereka sendiri
- ✅ Ownership verification di setiap operasi
- ✅ JWT authentication required
- ✅ Input validation
- ✅ Soft delete (tidak permanent delete)
- ✅ Proper HTTP status codes
- ✅ Zap logging untuk audit trail

---

## Performance Considerations

1. **Pagination**: Always use pagination, don't load all notifications at once
2. **Indexes**: Database has indexes on user_id, type, status, deleted_at
3. **Filtering**: Use status/type filters to reduce data transfer
4. **Cleanup**: Consider implementing job to delete old notifications (> 30 days)

---

## Integration Checklist

- [ ] Database migrated (notification table created)
- [ ] Postman collection imported
- [ ] Test List endpoint
- [ ] Test Update endpoint
- [ ] Test Delete endpoint
- [ ] Test filtering and pagination
- [ ] Test error handling
- [ ] Create notifications from Order API
- [ ] Create notifications from Payment API
- [ ] Implement frontend notification widget

---

## Code Structure

```
notification_adaptor.go (HTTP handlers)
         ↓
notification_usecase.go (Business logic)
         ↓
notification_repository.go (Database)
         ↓
notification_entity.go (Model)
```

---

## Troubleshooting

### Token expired or invalid

- Get new token from Auth API login endpoint
- Update {{token}} variable di Postman

### Notification not found in list

- Verify user_id is correct
- Check if notification was deleted
- Verify deleted_at is NULL in database

### Permission denied

- Verify you own the notification
- Check user_id matches in request

### Database error

- Run migration: `db.AutoMigrate(&entity.Notification{})`
- Verify PostgreSQL is running
- Check database connection string

---

## Next: Creating Notifications from Other Features

After integration test, implement notification creation in:

1. Order API - Create notification saat order created
2. Payment API - Create notification saat payment confirmed
3. Staff Management - Notification saat staff added/updated
4. Inventory - Alert saat stok menipis

---

## Documentation Files

1. **DOCS_NOTIFICATION_API.md** - Full API documentation
2. **NOTIFICATION_IMPLEMENTATION_SUMMARY.md** - Implementation details
3. **This file** - Quick start guide

---

## Git Commit Message Template

```
feat(notification): Implement notification system with CRUD operations

- Add notification entity with user isolation
- Implement notification repository with filtering & pagination
- Add notification usecase with business logic
- Create notification adaptor for HTTP handlers
- Register notification routes in wire.go
- Update database migration
- Add comprehensive API documentation
- Add Postman collection for testing

Implements:
- GET /api/v1/notifications (List with filters)
- PUT /api/v1/notifications/:id/status (Update status)
- DELETE /api/v1/notifications/:id (Delete notification)

Features:
- Pagination support (page, limit)
- Filtering by status (new, readed)
- Filtering by type (order, payment, system, alert)
- Custom sorting
- Unread count tracking
- User isolation & authorization
- Soft delete implementation
- Proper error handling
```

---

## Success Criteria ✅

- [x] 3 API endpoints implemented (List, Update, Delete)
- [x] Database schema created with proper indexes
- [x] Clean architecture followed
- [x] User isolation & authorization working
- [x] Error handling implemented
- [x] Comprehensive documentation
- [x] Postman collection ready
- [x] No compilation errors
- [x] Ready for integration with other features

---

**Status**: Ready for Testing & Integration ✨
