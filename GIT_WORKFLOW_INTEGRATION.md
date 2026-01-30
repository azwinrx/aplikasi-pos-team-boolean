# 🚀 Git Workflow - Integrasi Admin & User Management Features

Panduan lengkap untuk mengintegrasikan fitur Admin & User Management ke branch development tim Anda.

---

## 1. Pindah ke branch sendiri

```bash
git checkout dev-b
```

**Penjelasan:** Perintah ini memindahkan workspace Anda dari branch `main` ke branch `dev-b` (development branch Anda).

---

## 2. Ambil info terbaru dari server (Github/Gitlab)

Pastikan Dev A atau Dev B yang tahu sudah push ke repository terlebih dahulu.

```bash
git fetch origin
```

**Penjelasan:** Perintah ini mengambil semua informasi terbaru dari remote repository (Github/Gitlab) tanpa mengubah branch lokal Anda. Ini penting sebelum melakukan rebase atau merge.

---

## 3. Lakukan Rebase

Ini perintah ajaibnya.

```bash
git rebase origin/main
```

**Penjelasan:**

- Perintah ini mengambil semua commit dari `origin/main`
- Kemudian menerapkan semua commit dari branch `dev-b` Anda di atasnya
- Hasilnya adalah history yang linear dan clean
- Ini jauh lebih rapi daripada merge commit

**Jika ada conflict:**

```bash
# 1. Resolve conflicts di file yang bermasalah
# 2. Setelah fix, lanjutkan rebase dengan:
git rebase --continue

# Atau jika ingin membatalkan:
git rebase --abort
```

---

## 4. Push ke Repository

Setelah rebase selesai, push perubahan Anda.

```bash
git push origin dev-b --force-with-lease
```

**Penjelasan:**

- `--force-with-lease` lebih aman dari `--force`
- Ini akan reject jika ada orang lain yang push ke branch yang sama
- Lebih hati-hati dan menghindari overwrite perubahan orang lain

---

## 📋 Checklist Integrasi Admin Features

### Step 1: Siapkan Branch

```bash
# 1. Pastikan di branch lokal Anda
git status

# 2. Pindah ke dev-b jika belum
git checkout dev-b

# 3. Ambil update terbaru
git fetch origin

# 4. Lakukan rebase dengan main
git rebase origin/main
```

### Step 2: Verifikasi Files

Pastikan files berikut sudah di workspace Anda:

✅ `internal/usecase/admin.go`  
✅ `internal/adaptor/admin_adaptor.go`  
✅ `internal/dto/admin.go`  
✅ `pkg/middleware/auth.go`  
✅ `internal/wire/wire.go` (updated)  
✅ `internal/adaptor/adaptor.go` (updated)  
✅ `internal/usecase/usecase.go` (updated)  
✅ `internal/data/repository/auth.go` (updated)  
✅ `pkg/utils/token.go` (updated)

### Step 3: Install Dependencies

```bash
go mod download
go mod tidy
```

### Step 4: Test Compilation

```bash
go build -o main main.go
```

### Step 5: Push ke Dev Branch

```bash
git push origin dev-b --force-with-lease
```

### Step 6: Buat Pull Request (PR)

Buat PR dari `dev-b` ke `main` di Github/Gitlab untuk review sebelum merge.

---

## 🔄 Workflow Diagram

```
MAIN BRANCH (Production)
    ↑
    │  (PR - Code Review)
    │
    └─── DEV-B BRANCH (Development)
         ↑
         └─ Rebase dari main
         └─ Push ke origin
```

---

## ⚠️ Penting untuk Tim

### Jangan Lakukan Ini:

❌ Langsung push ke `main` tanpa PR  
❌ Gunakan `git push --force` (gunakan `--force-with-lease` saja)  
❌ Merge jika ada conflict tanpa resolve dulu  
❌ Push sebelum `go build` successful

### Harus Dilakukan:

✅ `git fetch origin` sebelum rebase  
✅ Resolve conflicts dengan teliti  
✅ Test compile dengan `go build`  
✅ Buat PR untuk code review  
✅ Semua tests passing baru merge

---

## 📝 Contoh Lengkap - Dari Awal sampai Akhir

```bash
# 1. Lihat branch apa yang sedang aktif
git status

# Output: On branch main

# 2. Pindah ke development branch
git checkout dev-b

# Output: Switched to branch 'dev-b'

# 3. Ambil update terbaru dari server
git fetch origin

# Output: Dari ... ke ... (semua update terdownload)

# 4. Lihat apa saja yang berubah di main
git log --oneline main..origin/main

# 5. Lakukan rebase
git rebase origin/main

# Output:
# Successfully rebased and updated refs/heads/dev-b
# atau
# CONFLICT in file.go (jika ada conflict)

# 6. Jika ada conflict, resolve dulu:
# - Buka file yang conflict di VS Code
# - Pilih "Accept Current Change" atau "Accept Incoming Change"
# - Save file

# 7. Lanjutkan rebase
git rebase --continue

# 8. Verifikasi tidak ada error
go build -o main main.go

# 9. Push ke repository
git push origin dev-b --force-with-lease

# Output: Counting objects ... Everything up-to-date
```

---

## 🎯 Langkah Cepat (TL;DR)

Untuk yang ingin cepat tanpa penjelasan panjang:

```bash
git checkout dev-b
git fetch origin
git rebase origin/main
go build -o main main.go
git push origin dev-b --force-with-lease
```

---

## 🆘 Troubleshooting

### Problem: "fatal: Merging with a bare repository is not allowed"

**Solution:**

```bash
git fetch origin
git rebase origin/main
```

Pastikan format branch namanya benar: `origin/main` bukan hanya `main`

---

### Problem: "Conflict in auth.go" atau file lain

**Solution:**

```bash
# 1. Buka file dengan conflict di VS Code
# 2. Cari tanda <<<<<<< HEAD sampai >>>>>>>

# 3. Pilih conflict resolution:
#    - Accept Current Change (keep your changes)
#    - Accept Incoming Change (keep main changes)
#    - Accept Both Changes (merge keduanya)

# 4. Save file
# 5. Add file ke staging
git add internal/data/repository/auth.go

# 6. Lanjutkan rebase
git rebase --continue

# 7. Push
git push origin dev-b --force-with-lease
```

---

### Problem: "Rebase error, want to cancel"

**Solution:**

```bash
git rebase --abort
```

Ini akan membatalkan seluruh rebase dan kembali ke state sebelum rebase.

---

## 📊 Status Features

| Feature            | Status  | Implementer | Branch |
| ------------------ | ------- | ----------- | ------ |
| Edit Profil User   | ✅ Done | Copilot     | dev-b  |
| List Data Admin    | ✅ Done | Copilot     | dev-b  |
| Edit Akses Admin   | ✅ Done | Copilot     | dev-b  |
| Logout             | ✅ Done | Copilot     | dev-b  |
| Password via Email | ✅ Done | Copilot     | dev-b  |

---

## 📚 Dokumentasi Tambahan

Setelah push ke dev-b, silakan baca dokumentasi lengkap:

- **API Documentation:** `DOCUMENTATION/ADMIN_USER_MANAGEMENT.md`
- **Implementation Guide:** `DOCUMENTATION/IMPLEMENTATION_GUIDE.md`
- **Quick Start:** `QUICK_START_GUIDE.md`
- **Postman Collection:** `Postman Collection/Admin_User_Management.postman_collection.json`

---

## ✅ Final Checklist Sebelum Merge ke Main

- [ ] Semua code sudah di push ke `dev-b`
- [ ] `git fetch origin` sudah dilakukan
- [ ] `git rebase origin/main` berhasil tanpa conflict
- [ ] `go build -o main main.go` berhasil (zero errors)
- [ ] `go mod tidy` sudah dilakukan
- [ ] Semua files ada (7 files baru + 6 files updated)
- [ ] Postman collection dapat dijalankan
- [ ] PR sudah di-create untuk review
- [ ] Minimal 1 orang sudah approve PR
- [ ] Semua tests passing
- [ ] Ready to merge ke main!

---

## 🤝 Untuk Tim Member Lain

Jika member lain ingin update ke latest code dari main:

```bash
# Mereka jalankan ini di branch mereka
git fetch origin
git merge origin/dev-b

# Atau jika ingin rebase (cleaner):
git rebase origin/dev-b
```

---

**Last Updated:** January 30, 2026  
**Status:** Ready for Team Integration ✅
