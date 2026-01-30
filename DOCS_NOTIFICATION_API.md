# Notification API Documentation

## Overview

Notification API adalah fitur untuk mengelola notifikasi pengguna dalam sistem POS. Sistem ini memungkinkan pengguna untuk melihat daftar notifikasi mereka, mengubah status notifikasi (dari "new" menjadi "readed"), dan menghapus notifikasi.

## Base URL

```
http://localhost:8080/api/v1/notifications
```

## Features

1. **List Notifications** - Mengambil daftar notifikasi dengan filtering dan pagination
2. **Update Notification Status** - Mengubah status notifikasi (new → readed)
3. **Delete Notification** - Menghapus notifikasi

## Authentication

Semua endpoint memerlukan user_id yang tersimpan di context dari JWT token. User hanya bisa mengakses notifikasi milik mereka sendiri.

---

## Endpoint Details

### 1. List Notifications

**Endpoint:** `GET /api/v1/notifications`

**Description:** Mengambil daftar notifikasi untuk user yang sedang login dengan opsi filter dan pagination.

**Query Parameters:**
| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| page | integer | No | 1 | Nomor halaman untuk pagination |
| limit | integer | No | 10 | Jumlah data per halaman (max: 100) |
| status | string | No | "" | Filter by status: "new", "readed", atau "all" untuk semua |
| type | string | No | "" | Filter by type: "order", "payment", "system", "alert" |
| sort_by | string | No | "created_at" | Field untuk sorting: "created_at", "status" |
| sort_order | string | No | "desc" | Urutan sort: "asc" atau "desc" |

**Example Request:**

```bash
curl -X GET "http://localhost:8080/api/v1/notifications?page=1&limit=10&status=new&type=order&sort_by=created_at&sort_order=desc" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response (200 OK):**

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
        "data": "{\"order_id\": 1, \"table_number\": \"T01\"}",
        "created_at": "2025-01-30T10:30:00Z",
        "updated_at": "2025-01-30T10:30:00Z"
      },
      {
        "id": 2,
        "user_id": 5,
        "title": "Pembayaran Dikonfirmasi",
        "message": "Pembayaran untuk pesanan #ORD001 sudah dikonfirmasi",
        "type": "payment",
        "status": "readed",
        "readed_at": "2025-01-30T10:35:00Z",
        "data": "{\"payment_method\": \"cash\", \"amount\": 150000}",
        "created_at": "2025-01-30T10:32:00Z",
        "updated_at": "2025-01-30T10:35:00Z"
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

**Error Response (401 Unauthorized):**

```json
{
  "code": 401,
  "message": "User tidak terautentikasi"
}
```

**Error Response (500 Internal Server Error):**

```json
{
  "code": 500,
  "message": "Gagal mengambil daftar notifikasi: [error details]"
}
```

---

### 2. Update Notification Status

**Endpoint:** `PUT /api/v1/notifications/:id/status`

**Description:** Mengubah status notifikasi dari "new" menjadi "readed". Hanya owner notifikasi yang dapat mengubah statusnya.

**URL Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| id | integer | Yes | ID dari notifikasi yang akan diupdate |

**Request Body:**

```json
{
  "status": "readed"
}
```

**Valid Status Values:**

- `"new"` - Notifikasi belum dibaca
- `"readed"` - Notifikasi sudah dibaca

**Example Request:**

```bash
curl -X PUT "http://localhost:8080/api/v1/notifications/1/status" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"status": "readed"}'
```

**Response (200 OK):**

```json
{
  "code": 200,
  "message": "Status notifikasi berhasil diubah",
  "data": {
    "id": 1,
    "status": "readed",
    "readed_at": "2025-01-30T10:35:00Z",
    "updated_at": "2025-01-30T10:35:00Z",
    "message": "Status notifikasi berhasil diubah menjadi readed"
  }
}
```

**Error Response (400 Bad Request):**

```json
{
  "code": 400,
  "message": "ID notifikasi tidak valid"
}
```

**Error Response (403 Forbidden):**

```json
{
  "code": 403,
  "message": "anda tidak memiliki akses ke notifikasi ini"
}
```

**Error Response (404 Not Found):**

```json
{
  "code": 404,
  "message": "notifikasi tidak ditemukan"
}
```

**Error Response (401 Unauthorized):**

```json
{
  "code": 401,
  "message": "User tidak terautentikasi"
}
```

---

### 3. Delete Notification

**Endpoint:** `DELETE /api/v1/notifications/:id`

**Description:** Menghapus notifikasi. Hanya owner notifikasi yang dapat menghapusnya. Penghapusan dilakukan dengan soft delete (DeletedAt akan diset).

**URL Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| id | integer | Yes | ID dari notifikasi yang akan dihapus |

**Example Request:**

```bash
curl -X DELETE "http://localhost:8080/api/v1/notifications/1" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response (200 OK):**

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

**Error Response (400 Bad Request):**

```json
{
  "code": 400,
  "message": "ID notifikasi tidak valid"
}
```

**Error Response (403 Forbidden):**

```json
{
  "code": 403,
  "message": "anda tidak memiliki akses untuk menghapus notifikasi ini"
}
```

**Error Response (404 Not Found):**

```json
{
  "code": 404,
  "message": "notifikasi tidak ditemukan"
}
```

**Error Response (401 Unauthorized):**

```json
{
  "code": 401,
  "message": "User tidak terautentikasi"
}
```

---

## Notification Types

Notification dapat memiliki tipe sebagai berikut:

| Type      | Description                   | Example                                             |
| --------- | ----------------------------- | --------------------------------------------------- |
| `order`   | Notifikasi terkait pesanan    | "Pesanan baru diterima", "Pesanan sedang disiapkan" |
| `payment` | Notifikasi terkait pembayaran | "Pembayaran dikonfirmasi", "Pembayaran tertunda"    |
| `system`  | Notifikasi sistem             | "Database backup selesai", "Update sistem tersedia" |
| `alert`   | Notifikasi penting/alert      | "Stok menipis", "Error sistem"                      |

---

## Notification Statuses

| Status   | Description                       |
| -------- | --------------------------------- |
| `new`    | Notifikasi baru yang belum dibaca |
| `readed` | Notifikasi sudah dibaca           |

---

## Data Field Format

Field `data` di notification dapat berisi informasi tambahan dalam format JSON string. Contoh:

```json
{
  "order_id": 1,
  "table_number": "T01",
  "customer_name": "John Doe",
  "total_price": 150000
}
```

---

## Common Use Cases

### 1. Get Unread Notifications

```bash
curl -X GET "http://localhost:8080/api/v1/notifications?status=new&limit=5" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 2. Mark Multiple Notifications as Read

```bash
# Notifikasi perlu diupdate satu per satu
curl -X PUT "http://localhost:8080/api/v1/notifications/1/status" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"status": "readed"}'

curl -X PUT "http://localhost:8080/api/v1/notifications/2/status" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"status": "readed"}'
```

### 3. Get Order Notifications Only

```bash
curl -X GET "http://localhost:8080/api/v1/notifications?type=order&status=new" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 4. Get All Notifications Sorted by Newest First

```bash
curl -X GET "http://localhost:8080/api/v1/notifications?sort_by=created_at&sort_order=desc&limit=20" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## Implementation Details

### Database Schema

```sql
CREATE TABLE notifications (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    type VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'new',
    readed_at TIMESTAMP,
    data JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    -- Indexes
    INDEX idx_user_id (user_id),
    INDEX idx_type (type),
    INDEX idx_status (status),
    INDEX idx_deleted_at (deleted_at)
);
```

### Architecture Layers

1. **Entity Layer** (`internal/data/entity/notification.go`)
   - Defines notification data model
   - Handles database schema

2. **Repository Layer** (`internal/data/repository/notification.go`)
   - Handles database operations
   - Methods: GetNotificationsByUserID, UpdateNotificationStatus, DeleteNotification, etc.

3. **UseCase Layer** (`internal/usecase/notification.go`)
   - Business logic
   - Validation and authorization
   - Methods: ListNotifications, UpdateNotificationStatus, DeleteNotification

4. **Adaptor Layer** (`internal/adaptor/notification_adaptor.go`)
   - HTTP request handlers
   - Request/response binding
   - Error handling

---

## Security Considerations

1. **User Isolation**: User hanya dapat mengakses notifikasi milik mereka sendiri
2. **Authorization**: Setiap operasi memverifikasi bahwa notifikasi milik user yang sedang login
3. **Soft Delete**: Notifikasi dihapus dengan soft delete, tidak benar-benar dihapus dari database
4. **JWT Authentication**: Semua endpoint memerlukan JWT token yang valid

---

## Error Handling

### Common Error Codes

| Code | Status                | Meaning                   |
| ---- | --------------------- | ------------------------- |
| 200  | OK                    | Request berhasil          |
| 400  | Bad Request           | Input tidak valid         |
| 401  | Unauthorized          | User tidak terautentikasi |
| 403  | Forbidden             | User tidak memiliki akses |
| 404  | Not Found             | Resource tidak ditemukan  |
| 500  | Internal Server Error | Error pada server         |

---

## Best Practices

1. **Pagination**: Selalu gunakan pagination untuk list notifications, jangan load semua data sekaligus
2. **Filtering**: Gunakan status dan type filters untuk mengurangi data yang dikirimkan
3. **Sorting**: Default sorting adalah created_at DESC (notifikasi terbaru di atas)
4. **Real-time Updates**: Untuk notifikasi real-time, pertimbangkan menggunakan WebSocket atau polling
5. **Cleanup**: Implement cleanup job untuk menghapus notifikasi lama (> 30 hari)

---

## Related Features

- **Authentication API** - Untuk login dan validasi user
- **Order API** - Untuk membuat notifikasi pesanan
- **Staff API** - Untuk manajemen staff yang menerima notifikasi

---

## Changelog

### Version 1.0.0 (2025-01-30)

- Initial release
- Implemented List, Update Status, and Delete notifications
- Added filtering and pagination support
- Soft delete implementation
