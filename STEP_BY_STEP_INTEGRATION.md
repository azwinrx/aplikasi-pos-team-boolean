# 📖 LANGKAH-LANGKAH INTEGRASI ADMIN FEATURES

## UNTUK DILIHAT TEMAN-TEMAN TEAM ANDA

---

## 1️⃣ Pindah ke branch sendiri

```bash
git checkout dev-b
```

**Apa artinya?**  
Perintah ini mengubah branch lokal Anda dari `main` ke `dev-b`. Semua perubahan yang ada di `dev-b` akan tampil di file explorer Anda.

---

## 2️⃣ Ambil info terbaru dari server (Github/Gitlab)

Supaya laptop Dev B tidak ketinggalan data, dan semuanya sudah up-to-date.

```bash
git fetch origin
```

**Apa artinya?**  
Perintah ini download semua update terbaru dari remote server. Ini seperti mengecek "ada update apa dari server?" tanpa mengubah file Anda sendiri.

**Hasil yang akan Anda lihat:**

```
From https://github.com/azwinrx/aplikasi-pos-team-boolean
 * [new branch]      dev-b      -> origin/dev-b
   4a7c8f9..5b9d2e0  main       -> origin/main
```

---

## 3️⃣ Lakukan Rebase

Ini perintah ajaibnya.

```bash
git rebase origin/main
```

**Apa artinya?**

- Rebase mengambil semua commit dari `origin/main`
- Kemudian "memutar ulang" semua commit Anda di atasnya
- Hasilnya adalah git history yang **linear** dan **clean**
- Ini lebih bagus daripada merge yang menghasilkan commit tambahan

**Hasil yang akan Anda lihat:**

```
Successfully rebased and updated refs/heads/dev-b.
```

**Jika ada CONFLICT (perubahan yang sama):**

```
CONFLICT (content): Merge conflict in internal/usecase/auth.go

Resolve all conflicts manually, then run "git rebase --continue".

As a decision point for the rebase, you can also run "git rebase --abort".
```

**Kalau ada conflict, lakukan ini:**

```bash
# 1. Buka file yang conflict di VS Code
# 2. Cari simbol: <<<<<<< HEAD sampai >>>>>>>
# 3. Pilih mana yang mau: "Accept Current Change" atau "Accept Incoming Change"
# 4. Save file
# 5. Lanjutkan rebase:

git add .
git rebase --continue
```

---

## 4️⃣ Push ke Repository

Kirim semua perubahan Anda ke remote server agar semua orang bisa lihat.

```bash
git push origin dev-b --force-with-lease
```

**Apa artinya?**

- `git push` = kirim ke server
- `origin` = server/repository Anda
- `dev-b` = branch tujuan
- `--force-with-lease` = "paksa push, tapi jangan overwrite punya orang lain"

**Hasil yang akan Anda lihat:**

```
Counting objects: 87, done.
Delta compression using up to 8 threads.
Compressing objects: 100% (45/45), done.
Writing objects: 100% (87/87), 18.5 KB | 2.3 MB/s, done.
Total 87 (delta 42), reused 0 (delta 0)
remote: Resolving deltas: 100% (42/42), done.
To https://github.com/azwinrx/aplikasi-pos-team-boolean
 5b9d2e0..7c3f1a2  dev-b -> dev-b
```

---

# ⚡ QUICK REFERENCE UNTUK TEAM

Copy-paste ini untuk eksekusi cepat:

```bash
git checkout dev-b && git fetch origin && git rebase origin/main && git push origin dev-b --force-with-lease
```

**Atau jalankan satu-satu jika ada masalah:**

```bash
# 1. Pindah branch
git checkout dev-b

# 2. Ambil update
git fetch origin

# 3. Rebase dengan main
git rebase origin/main

# 4. Jika ada conflict, resolve dulu
# (edit files di VS Code)
git add .
git rebase --continue

# 5. Push ke server
git push origin dev-b --force-with-lease
```

---

# ✅ CHECKLIST SEBELUM & SESUDAH

## SEBELUM MULAI:

- [ ] Pastikan tidak ada uncommitted changes: `git status`
- [ ] Jika ada, commit dulu atau stash: `git stash`
- [ ] Pastikan di branch yang benar: `git branch` (tandanya _, misalnya `_ dev-b`)

## SAAT EKSEKUSI:

- [ ] `git checkout dev-b` - Selesai? Lihat pada `-b` di depan branch list
- [ ] `git fetch origin` - Ada update? Output tidak error?
- [ ] `git rebase origin/main` - Berhasil? Lihat "Successfully rebased"
- [ ] Jika conflict: resolve dulu
- [ ] `git push origin dev-b --force-with-lease` - Selesai push?

## SESUDAH PUSH:

- [ ] Cek di Github/Gitlab apakah sudah muncul di `dev-b`
- [ ] Lihat branch history dengan: `git log --oneline -10`
- [ ] Compile test: `go build -o main main.go`
- [ ] Tidak ada error? Siap untuk di-merge ke main!

---

# 🎓 PENJELASAN UNTUK YANG BARU GITHUB

### Branch adalah apa?

Branch adalah "versi copy" dari project Anda. Seperti:

- `main` = Versi Production (yang live di server)
- `dev-b` = Versi Development (tempat testing fitur baru)

### Rebase vs Merge - Apa bedanya?

**MERGE:**

```
    main:    A -- B -- C
                       \
dev-b:               D -- E -- M (merge commit)
```

Hasilnya ada extra commit M

**REBASE:**

```
    main:    A -- B -- C
                       |
dev-b:                 D' -- E' (commit D dan E diulang di atas C)
```

Hasilnya lebih clean, no extra commit

### Kenapa `--force-with-lease`?

```bash
# Bahaya:
git push --force

# Lebih aman:
git push --force-with-lease
```

`--force-with-lease` lebih hati-hati. Dia tidak akan push jika ada perubahan dari orang lain di branch yang sama di server.

---

# 🆘 COMMON PROBLEMS & SOLUTIONS

### ❌ Error: "You are not currently on a branch"

**Solusi:**

```bash
git checkout dev-b
```

### ❌ Error: "fatal: 'origin/main' does not appear to be a 'git' repository"

**Solusi:**

```bash
git fetch origin
# atau
git fetch origin main
```

### ❌ Error: "Your local changes will be overwritten"

**Solusi:**

```bash
git stash
git checkout dev-b
git stash pop
```

### ❌ Rebase stuck (ada conflict yang tidak bisa resolve)

**Solusi:**

```bash
# Cancel rebase
git rebase --abort

# Start dari awal
git fetch origin
git rebase origin/main
```

### ❌ Sudah push tapi mau undo

**Solusi:**

```bash
# Reset ke commit sebelumnya
git reset --soft HEAD~1

# atau
git revert <commit-hash>
```

---

# 📊 FILES YANG BERUBAH

## FILES BARU (7):

```
✅ internal/usecase/admin.go
✅ internal/adaptor/admin_adaptor.go
✅ internal/dto/admin.go
✅ pkg/middleware/auth.go
✅ DOCUMENTATION/ADMIN_USER_MANAGEMENT.md
✅ DOCUMENTATION/IMPLEMENTATION_GUIDE.md
✅ Postman Collection/Admin_User_Management.postman_collection.json
```

## FILES DIUPDATE (6):

```
🔄 internal/usecase/usecase.go
🔄 internal/adaptor/adaptor.go
🔄 internal/data/repository/auth.go
🔄 internal/wire/wire.go
🔄 pkg/utils/token.go
🔄 internal/usecase/auth.go
```

---

# 🚀 ENDPOINTS BARU

Setelah rebase selesai, ada 7 endpoint baru:

```
✅ GET    /api/v1/profile              - Get user profile
✅ PUT    /api/v1/profile              - Update user profile
✅ POST   /api/v1/auth/logout          - Logout user
✅ GET    /api/v1/admin/list           - List admins (superadmin)
✅ PUT    /api/v1/admin/:id/access     - Edit admin access (superadmin)
✅ POST   /api/v1/admin/create         - Create admin (superadmin)
```

Semua endpoint ini sudah ada di Postman collection yang baru.

---

# 📚 DOKUMENTASI TAMBAHAN

Setelah berhasil, baca ini untuk lebih detail:

1. **Quick Start:** `QUICK_START_GUIDE.md`
2. **API Docs:** `DOCUMENTATION/ADMIN_USER_MANAGEMENT.md`
3. **Implementation:** `DOCUMENTATION/IMPLEMENTATION_GUIDE.md`
4. **Postman:** `Postman Collection/Admin_User_Management.postman_collection.json`

---

# ✨ TIPS BONUS UNTUK EXPERT

### View rebase history:

```bash
git log --oneline --graph --all
```

### Undo last commit (belum push):

```bash
git reset --soft HEAD~1
```

### Edit last commit message:

```bash
git commit --amend
```

### Cherry-pick specific commit:

```bash
git cherry-pick <commit-hash>
```

### Stash changes sementara:

```bash
git stash
git stash list
git stash pop
```

---

**Status:** ✅ SIAP UNTUK TEAM LIHAT  
**Last Updated:** January 30, 2026  
**Untuk Dipindahkan ke:** Team Communication Channel (Slack, Discord, atau Docs Internal)
