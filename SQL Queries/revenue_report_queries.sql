-- ============================================
-- REVENUE REPORT SQL QUERIES
-- Manual Testing & Verification Queries
-- Created: 2026-01-29
-- ============================================

-- ============================================
-- 1. REVENUE BY STATUS
-- ============================================

-- Get total revenue and breakdown by status
SELECT
    status,
    COALESCE(SUM(total_amount), 0) as total_revenue,
    COUNT(*) as order_count
FROM orders
WHERE deleted_at IS NULL
GROUP BY status
ORDER BY total_revenue DESC;

-- Get total revenue (all statuses)
SELECT
    COALESCE(SUM(total_amount), 0) as total_revenue,
    COUNT(*) as total_orders
FROM orders
WHERE deleted_at IS NULL;

-- ============================================
-- 2. REVENUE PER MONTH
-- ============================================

-- Get revenue per month for current year (2026)
SELECT
    EXTRACT(MONTH FROM created_at)::int as month,
    COALESCE(SUM(total_amount), 0) as total_revenue,
    COUNT(*) as order_count
FROM orders
WHERE deleted_at IS NULL
    AND EXTRACT(YEAR FROM created_at) = 2026
GROUP BY EXTRACT(MONTH FROM created_at)
ORDER BY month ASC;

-- Get revenue per month for specific year (parameterized)
SELECT
    EXTRACT(MONTH FROM created_at)::int as month,
    COALESCE(SUM(total_amount), 0) as total_revenue,
    COUNT(*) as order_count
FROM orders
WHERE deleted_at IS NULL
    AND EXTRACT(YEAR FROM created_at) = 2025  -- Change year here
GROUP BY EXTRACT(MONTH FROM created_at)
ORDER BY month ASC;

-- Get revenue summary per year
SELECT
    EXTRACT(YEAR FROM created_at)::int as year,
    COALESCE(SUM(total_amount), 0) as total_revenue,
    COUNT(*) as order_count
FROM orders
WHERE deleted_at IS NULL
GROUP BY EXTRACT(YEAR FROM created_at)
ORDER BY year DESC;

-- ============================================
-- 3. PRODUCT REVENUE LIST
-- ============================================

-- Get complete product revenue details
SELECT
    p.id as product_id,
    p.name as product_name,
    p.price as price,
    COALESCE(SUM(oi.subtotal), 0) as total_revenue,
    COALESCE(SUM(oi.quantity), 0) as total_sold,
    COUNT(DISTINCT oi.order_id) as order_count,
    MAX(o.created_at) as last_order_at
FROM products p
LEFT JOIN order_items oi ON p.id = oi.product_id AND oi.deleted_at IS NULL
LEFT JOIN orders o ON oi.order_id = o.id AND o.deleted_at IS NULL
WHERE p.deleted_at IS NULL
GROUP BY p.id, p.name, p.price
ORDER BY total_revenue DESC;

-- Get only products that have been sold
SELECT
    p.id as product_id,
    p.name as product_name,
    p.price as price,
    SUM(oi.subtotal) as total_revenue,
    SUM(oi.quantity) as total_sold,
    COUNT(DISTINCT oi.order_id) as order_count,
    MAX(o.created_at) as last_order_at
FROM products p
INNER JOIN order_items oi ON p.id = oi.product_id AND oi.deleted_at IS NULL
INNER JOIN orders o ON oi.order_id = o.id AND o.deleted_at IS NULL
WHERE p.deleted_at IS NULL
GROUP BY p.id, p.name, p.price
ORDER BY total_revenue DESC;

-- Get products that have NOT been sold
SELECT
    p.id as product_id,
    p.name as product_name,
    p.price as price
FROM products p
LEFT JOIN order_items oi ON p.id = oi.product_id AND oi.deleted_at IS NULL
WHERE p.deleted_at IS NULL
    AND oi.id IS NULL
ORDER BY p.name;

-- ============================================
-- 4. VERIFICATION QUERIES
-- ============================================

-- Verify order data
SELECT
    id,
    customer_name,
    total_amount,
    status,
    created_at
FROM orders
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT 10;

-- Verify order items
SELECT
    oi.id,
    oi.order_id,
    p.name as product_name,
    oi.quantity,
    oi.price,
    oi.subtotal,
    o.created_at
FROM order_items oi
INNER JOIN products p ON oi.product_id = p.id
INNER JOIN orders o ON oi.order_id = o.id
WHERE oi.deleted_at IS NULL
ORDER BY o.created_at DESC
LIMIT 10;

-- Check for deleted orders (should be excluded from revenue)
SELECT
    COUNT(*) as deleted_count
FROM orders
WHERE deleted_at IS NOT NULL;

-- ============================================
-- 5. ADVANCED ANALYTICS
-- ============================================

-- Revenue by status and month
SELECT
    status,
    EXTRACT(MONTH FROM created_at)::int as month,
    COALESCE(SUM(total_amount), 0) as total_revenue,
    COUNT(*) as order_count
FROM orders
WHERE deleted_at IS NULL
    AND EXTRACT(YEAR FROM created_at) = 2026
GROUP BY status, EXTRACT(MONTH FROM created_at)
ORDER BY month ASC, status;

-- Top 5 products by revenue
SELECT
    p.id as product_id,
    p.name as product_name,
    COALESCE(SUM(oi.subtotal), 0) as total_revenue,
    COALESCE(SUM(oi.quantity), 0) as total_sold
FROM products p
LEFT JOIN order_items oi ON p.id = oi.product_id AND oi.deleted_at IS NULL
LEFT JOIN orders o ON oi.order_id = o.id AND o.deleted_at IS NULL
WHERE p.deleted_at IS NULL
GROUP BY p.id, p.name
ORDER BY total_revenue DESC
LIMIT 5;

-- Bottom 5 products by revenue
SELECT
    p.id as product_id,
    p.name as product_name,
    COALESCE(SUM(oi.subtotal), 0) as total_revenue,
    COALESCE(SUM(oi.quantity), 0) as total_sold
FROM products p
LEFT JOIN order_items oi ON p.id = oi.product_id AND oi.deleted_at IS NULL
LEFT JOIN orders o ON oi.order_id = o.id AND o.deleted_at IS NULL
WHERE p.deleted_at IS NULL
GROUP BY p.id, p.name
ORDER BY total_revenue ASC
LIMIT 5;

-- Average order value by status
SELECT
    status,
    COUNT(*) as order_count,
    COALESCE(AVG(total_amount), 0) as average_order_value,
    COALESCE(MIN(total_amount), 0) as min_order_value,
    COALESCE(MAX(total_amount), 0) as max_order_value,
    COALESCE(SUM(total_amount), 0) as total_revenue
FROM orders
WHERE deleted_at IS NULL
GROUP BY status
ORDER BY total_revenue DESC;

-- Daily revenue for current month
SELECT
    DATE(created_at) as order_date,
    COUNT(*) as order_count,
    COALESCE(SUM(total_amount), 0) as daily_revenue
FROM orders
WHERE deleted_at IS NULL
    AND EXTRACT(YEAR FROM created_at) = EXTRACT(YEAR FROM CURRENT_DATE)
    AND EXTRACT(MONTH FROM created_at) = EXTRACT(MONTH FROM CURRENT_DATE)
GROUP BY DATE(created_at)
ORDER BY order_date ASC;

-- ============================================
-- 6. DATA QUALITY CHECKS
-- ============================================

-- Check for orders without items
SELECT
    o.id,
    o.customer_name,
    o.total_amount,
    o.status
FROM orders o
LEFT JOIN order_items oi ON o.id = oi.order_id AND oi.deleted_at IS NULL
WHERE o.deleted_at IS NULL
    AND oi.id IS NULL;

-- Check for order items with NULL or zero values
SELECT
    id,
    order_id,
    product_id,
    quantity,
    price,
    subtotal
FROM order_items
WHERE deleted_at IS NULL
    AND (quantity IS NULL OR quantity <= 0
         OR price IS NULL OR price < 0
         OR subtotal IS NULL OR subtotal < 0);

-- Check for products with NULL price
SELECT
    id,
    name,
    price
FROM products
WHERE deleted_at IS NULL
    AND (price IS NULL OR price < 0);

-- ============================================
-- 7. PERFORMANCE TESTING
-- ============================================

-- Explain query plan for revenue by status
EXPLAIN ANALYZE
SELECT
    status,
    COALESCE(SUM(total_amount), 0) as total_revenue,
    COUNT(*) as order_count
FROM orders
WHERE deleted_at IS NULL
GROUP BY status
ORDER BY total_revenue DESC;

-- Explain query plan for product revenue
EXPLAIN ANALYZE
SELECT
    p.id as product_id,
    p.name as product_name,
    p.price as price,
    COALESCE(SUM(oi.subtotal), 0) as total_revenue,
    COALESCE(SUM(oi.quantity), 0) as total_sold,
    COUNT(DISTINCT oi.order_id) as order_count,
    MAX(o.created_at) as last_order_at
FROM products p
LEFT JOIN order_items oi ON p.id = oi.product_id AND oi.deleted_at IS NULL
LEFT JOIN orders o ON oi.order_id = o.id AND o.deleted_at IS NULL
WHERE p.deleted_at IS NULL
GROUP BY p.id, p.name, p.price
ORDER BY total_revenue DESC;

-- ============================================
-- 8. SAMPLE DATA INSERTION (FOR TESTING)
-- ============================================

-- Insert sample order (uncomment to use)
/*
INSERT INTO orders (user_id, table_id, payment_method_id, customer_name, total_amount, tax, status, created_at)
VALUES (1, 1, 1, 'Test Customer', 50000, 5000, 'paid', NOW());

-- Get last inserted order id
SELECT currval('orders_id_seq');

-- Insert sample order items (replace order_id with last inserted id)
INSERT INTO order_items (order_id, product_id, quantity, price, subtotal)
VALUES
    (LAST_INSERT_ID, 1, 2, 25000, 50000),
    (LAST_INSERT_ID, 2, 1, 30000, 30000);
*/

-- ============================================
-- END OF QUERIES
-- ============================================
