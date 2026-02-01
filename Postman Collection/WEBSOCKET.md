# 🔌 WebSocket Documentation

## Dashboard Realtime WebSocket

WebSocket endpoint untuk mendapatkan data dashboard secara realtime (sales & revenue).

---

## 📌 Connection Info

| Property            | Value                                     |
| ------------------- | ----------------------------------------- |
| **URL**             | `ws://localhost:8080/api/v1/dashboard/ws` |
| **Protocol**        | WebSocket (WS)                            |
| **Update Interval** | Setiap 3 detik                            |
| **Authentication**  | Tidak diperlukan (public)                 |

---

## 📤 Response Format

Server akan mengirim data setiap 3 detik dalam format JSON:

```json
{
  "daily_sales": 1500000.0,
  "monthly_sales": 21000000.0,
  "daily_orders": 25,
  "monthly_orders": 350
}
```

### Response Fields

| Field            | Type  | Description             |
| ---------------- | ----- | ----------------------- |
| `daily_sales`    | float | Total revenue hari ini  |
| `monthly_sales`  | float | Total revenue bulan ini |
| `daily_orders`   | int   | Jumlah order hari ini   |
| `monthly_orders` | int   | Jumlah order bulan ini  |

---

## 🧪 Cara Testing WebSocket

### Metode 1: Postman WebSocket Request

1. Buka Postman
2. Klik **New** (tombol + di tab)
3. Pilih **WebSocket** (bukan HTTP Request)
4. Masukkan URL: `ws://localhost:8080/api/v1/dashboard/ws`
5. Klik **Connect**
6. Tunggu beberapa detik, data akan muncul otomatis

![Postman WebSocket Steps](https://learning.postman.com/docs/img/websocket-request.jpg)

> ⚠️ **Penting:** Jangan gunakan HTTP Request biasa (GET/POST) untuk WebSocket. Akan muncul error "Invalid protocol: ws:"

---

### Metode 2: wscat (Command Line)

```bash
# Install wscat (jika belum ada)
npm install -g wscat

# Connect ke WebSocket
wscat -c ws://localhost:8080/api/v1/dashboard/ws
```

**Output:**

```
Connected (press CTRL+C to quit)
< {"daily_sales":1500000,"monthly_sales":21000000,"daily_orders":25,"monthly_orders":350}
< {"daily_sales":1500000,"monthly_sales":21000000,"daily_orders":25,"monthly_orders":350}
< {"daily_sales":1550000,"monthly_sales":21050000,"daily_orders":26,"monthly_orders":351}
...
```

---

### Metode 3: Browser Console (Chrome/Firefox)

1. Buka browser (Chrome/Firefox/Edge)
2. Tekan **F12** untuk membuka Developer Tools
3. Masuk ke tab **Console**
4. Paste kode berikut:

```javascript
// Connect ke WebSocket
const ws = new WebSocket("ws://localhost:8080/api/v1/dashboard/ws");

// Event: Koneksi berhasil
ws.onopen = () => {
  console.log("✅ Connected to Dashboard WebSocket");
};

// Event: Menerima data
ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  console.log("📊 Realtime Data:", data);
  console.log(`   Daily Sales: Rp ${data.daily_sales.toLocaleString()}`);
  console.log(`   Monthly Sales: Rp ${data.monthly_sales.toLocaleString()}`);
  console.log(`   Daily Orders: ${data.daily_orders}`);
  console.log(`   Monthly Orders: ${data.monthly_orders}`);
  console.log("---");
};

// Event: Error
ws.onerror = (error) => {
  console.error("❌ WebSocket Error:", error);
};

// Event: Koneksi ditutup
ws.onclose = () => {
  console.log("🔌 Disconnected from Dashboard WebSocket");
};

// Untuk menutup koneksi manual, jalankan: ws.close()
```

5. Tekan **Enter**
6. Data akan muncul setiap 3 detik

---

### Metode 4: Online WebSocket Tester

1. Buka salah satu website berikut:
   - https://websocketking.com/
   - https://www.piesocket.com/websocket-tester
   - https://socketsbay.com/test-websockets

2. Masukkan URL: `ws://localhost:8080/api/v1/dashboard/ws`
3. Klik **Connect**
4. Data akan muncul secara otomatis

---

## ❓ Troubleshooting

### Error: "Invalid protocol: ws:"

**Penyebab:** Menggunakan HTTP request biasa untuk WebSocket
**Solusi:** Gunakan Postman WebSocket Request (New → WebSocket)

### Error: "WebSocket connection failed"

**Penyebab:** Server belum berjalan atau port salah
**Solusi:**

1. Pastikan server berjalan: `go run .`
2. Cek port: default 8080

### Error: "Connection refused"

**Penyebab:** Firewall atau server tidak aktif
**Solusi:**

1. Pastikan server aktif
2. Cek firewall tidak memblokir port 8080

### Tidak ada data yang muncul

**Penyebab:** Koneksi belum established
**Solusi:** Tunggu beberapa detik, data akan muncul setiap 3 detik

---

## 📝 Notes

- WebSocket update setiap **3 detik**
- Data diambil langsung dari database (realtime)
- Tidak memerlukan authentication
- Koneksi otomatis terputus jika client disconnect
- Server mendukung multiple concurrent connections

---

**Last Updated:** 1 February 2026
**Status:** ✅ READY FOR TESTING
