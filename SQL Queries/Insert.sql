-- 1. Insert Users (Password: 'password123' - ingat nanti di Go harus di-hash)
INSERT INTO users (name, email, password, role) VALUES
('Super Admin', 'super@pos.com', '$2a$10$wK/p.8f.0/..hashedpassword..', 'superadmin'),
('Manager Resto', 'admin@pos.com', '$2a$10$wK/p.8f.0/..hashedpassword..', 'admin'),
('Budi Staff', 'staff@pos.com', '$2a$10$wK/p.8f.0/..hashedpassword..', 'staff');

-- 2. Insert Categories
INSERT INTO categories (icon_category, category_name, description) VALUES
('🍕', 'Pizza', 'Delicious pizza varieties'),
('🍔', 'Burger', 'Juicy burgers and sandwiches'),
('🍗', 'Chicken', 'Crispy fried chicken'),
('🥐', 'Bakery', 'Fresh baked goods'),
('🥤', 'Beverage', 'Refreshing drinks'),
('🦐', 'Seafood', 'Fresh seafood dishes');

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

-- 5. Insert Products (Menu Items)
INSERT INTO products (product_image, product_name, item_id, stock, category_id, price, is_available) VALUES
('/images/chicken-parmesan.jpg', 'Chicken Parmesan', '#22314644', 119, 3, 55.00, true),
('/images/margherita-pizza.jpg', 'Margherita Pizza', '#22314645', 85, 1, 45.00, true),
('/images/pepperoni-pizza.jpg', 'Pepperoni Pizza', '#22314646', 72, 1, 50.00, true),
('/images/classic-burger.jpg', 'Classic Burger', '#22314647', 95, 2, 35.00, true),
('/images/cheese-burger.jpg', 'Cheese Burger', '#22314648', 8, 2, 40.00, true),
('/images/cola.jpg', 'Cola', '#22314649', 200, 5, 5.00, true),
('/images/orange-juice.jpg', 'Orange Juice', '#22314650', 0, 5, 8.00, false);

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
