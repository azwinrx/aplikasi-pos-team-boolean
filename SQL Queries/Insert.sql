-- 1. Insert Users (Password: 'password123' - ingat nanti di Go harus di-hash)
INSERT INTO users (name, email, password, role) VALUES
('Super Admin', 'super@pos.com', '$2a$10$wK/p.8f.0/..hashedpassword..', 'superadmin'),
('Manager Resto', 'admin@pos.com', '$2a$10$wK/p.8f.0/..hashedpassword..', 'admin'),
('Budi Staff', 'staff@pos.com', '$2a$10$wK/p.8f.0/..hashedpassword..', 'staff');

-- 2. Insert Categories
INSERT INTO categories (name) VALUES
('Makanan Berat'), ('Minuman'), ('Snack'), ('Dessert');

-- 3. Insert Payment Methods
INSERT INTO payment_methods (name) VALUES
('Cash'), ('QRIS'), ('Debit Card');

-- 4. Insert Tables
INSERT INTO tables (number, capacity, status) VALUES
('T01', 4, 'available'),
('T02', 2, 'occupied'),
('T03', 6, 'reserved'),
('T04', 4, 'available'),
('T05', 2, 'available');

-- 5. Insert Products (Campuran produk baru dan lama)
INSERT INTO products (category_id, name, description, price, created_at) VALUES
(1, 'Nasi Goreng Spesial', 'Nasi goreng dengan telur dan ayam', 25000, NOW() - INTERVAL '40 days'), -- Produk Lama
(1, 'Ayam Bakar Madu', 'Ayam bakar oles madu', 30000, NOW() - INTERVAL '5 days'), -- Produk Baru
(2, 'Es Teh Manis', 'Teh manis dingin segar', 5000, NOW() - INTERVAL '40 days'),
(2, 'Kopi Susu Gula Aren', 'Kopi kekinian', 18000, NOW() - INTERVAL '2 days'), -- Produk Baru
(3, 'Kentang Goreng', 'French fries original', 15000, NOW() - INTERVAL '40 days');

-- 6. Insert Inventory
INSERT INTO inventories (name, quantity, unit, min_stock) VALUES
('Beras', 50, 'kg', 10),
('Telur', 100, 'butir', 20),
('Minyak Goreng', 20, 'liter', 5),
('Gula Pasir', 15, 'kg', 5);

-- 7. Insert Notifications
INSERT INTO notifications (title, message, type) VALUES
('Stok Menipis', 'Stok Minyak Goreng tersisa 5 liter', 'alert'),
('Order Baru', 'Meja T02 melakukan pemesanan', 'order');

-- 8. Insert Orders (Dummy Data Penjualan untuk Laporan Revenue)
-- Order 1: Paid
INSERT INTO orders (user_id, table_id, payment_method_id, customer_name, total_amount, tax, status, created_at)
VALUES (3, 1, 1, 'Customer A', 55000, 5500, 'paid', NOW() - INTERVAL '2 days');

-- Order 2: Paid
INSERT INTO orders (user_id, table_id, payment_method_id, customer_name, total_amount, tax, status, created_at)
VALUES (3, 2, 2, 'Customer B', 30000, 3000, 'paid', NOW());

-- 9. Insert Order Items (Detail pesanan di atas)
-- Items untuk Order 1 (Total 55.000)
INSERT INTO order_items (order_id, product_id, quantity, price, subtotal) VALUES
(1, 1, 1, 25000, 25000), -- Nasi Goreng
(1, 2, 1, 30000, 30000); -- Ayam Bakar

-- Items untuk Order 2 (Total 30.000)
INSERT INTO order_items (order_id, product_id, quantity, price, subtotal) VALUES
(2, 2, 1, 30000, 30000); -- Ayam Bakar
