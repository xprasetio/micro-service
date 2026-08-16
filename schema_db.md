-- =========================================================
-- SAYUR PROJECT DATABASE SCHEMA
-- PostgreSQL
-- =========================================================

-- =========================================================
-- 1. CATEGORIES
-- =========================================================

CREATE TABLE categories (
id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
parent_id BIGINT NULL,
name VARCHAR(255) NOT NULL,
icon VARCHAR(255),
status BOOLEAN DEFAULT TRUE,
slug VARCHAR(255) NOT NULL UNIQUE,
description TEXT,
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
deleted_at TIMESTAMP NULL,

    CONSTRAINT fk_categories_parent
        FOREIGN KEY (parent_id)
        REFERENCES categories(id)
        ON DELETE SET NULL

);

-- =========================================================
-- 2. ROLES
-- =========================================================

CREATE TABLE roles (
id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
name VARCHAR(255) NOT NULL UNIQUE
);

-- =========================================================
-- 3. USERS
-- =========================================================

CREATE TABLE users (
id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
name VARCHAR(255) NOT NULL,
email VARCHAR(255) NOT NULL UNIQUE,
password VARCHAR(255) NOT NULL,
address VARCHAR(255),
phone VARCHAR(50),
photo VARCHAR(255),
lat VARCHAR(100),
lng VARCHAR(100),
is_verified BOOLEAN DEFAULT FALSE,
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
deleted_at TIMESTAMP NULL
);

-- =========================================================
-- 4. USER ROLE
-- =========================================================

CREATE TABLE user_role (
id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
role_id BIGINT NOT NULL,
user_id BIGINT NOT NULL,

    CONSTRAINT fk_user_role_role
        FOREIGN KEY (role_id)
        REFERENCES roles(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_user_role_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_user_role
        UNIQUE (role_id, user_id)

);

-- =========================================================
-- 5. PRODUCTS
-- =========================================================

CREATE TABLE products (
id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    category_slug VARCHAR(255) NOT NULL,
    parent_id BIGINT NULL,

    name VARCHAR(255) NOT NULL,
    image TEXT,
    description VARCHAR(255),

    reguler_price BIGINT NOT NULL DEFAULT 0,
    sale_price BIGINT NOT NULL DEFAULT 0,

    unit VARCHAR(100),
    weight INTEGER DEFAULT 0,
    stock INTEGER DEFAULT 0,
    varian INTEGER DEFAULT 0,

    -- Di ERD tertulis STATUS_PRODUCT ENUM.
    -- Untuk sementara menggunakan VARCHAR.
    status VARCHAR(50),

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,

    CONSTRAINT fk_products_category
        FOREIGN KEY (category_slug)
        REFERENCES categories(slug)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,

    CONSTRAINT fk_products_parent
        FOREIGN KEY (parent_id)
        REFERENCES products(id)
        ON DELETE SET NULL

);

-- =========================================================
-- 6. ORDERS
-- =========================================================

CREATE TABLE orders (
id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    order_code VARCHAR(60) NOT NULL UNIQUE,
    buyer_id BIGINT NOT NULL,

    order_date DATE NOT NULL DEFAULT CURRENT_DATE,

    -- STATUS_ORDER pada ERD
    status VARCHAR(50),

    total_amount BIGINT NOT NULL DEFAULT 0,

    -- SHIPPING_TYPE pada ERD
    shipping_type VARCHAR(50),

    shipping_fee INTEGER NOT NULL DEFAULT 0,

    order_time TIME,
    remarks VARCHAR(255),

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,

    CONSTRAINT fk_orders_buyer
        FOREIGN KEY (buyer_id)
        REFERENCES users(id)
        ON DELETE RESTRICT

);

-- =========================================================
-- 7. ORDER ITEMS
-- =========================================================

CREATE TABLE order_items (
id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    product_id BIGINT NOT NULL,
    order_id BIGINT NOT NULL,

    quantity INTEGER NOT NULL DEFAULT 1,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,

    CONSTRAINT fk_order_items_product
        FOREIGN KEY (product_id)
        REFERENCES products(id)
        ON DELETE RESTRICT,

    CONSTRAINT fk_order_items_order
        FOREIGN KEY (order_id)
        REFERENCES orders(id)
        ON DELETE CASCADE

);

-- =========================================================
-- 8. PAYMENTS
-- =========================================================

CREATE TABLE payments (
id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    order_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,

    payment_method VARCHAR(50),
    payment_status VARCHAR(50),
    payment_gateway_id VARCHAR(50),

    gross_amount NUMERIC(15,2) DEFAULT 0,
    payment_url TEXT,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,

    CONSTRAINT fk_payments_order
        FOREIGN KEY (order_id)
        REFERENCES orders(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_payments_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE RESTRICT

);

-- =========================================================
-- 9. PAYMENT LOGS
-- =========================================================

CREATE TABLE payment_logs (
id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    payment_id BIGINT NOT NULL,
    status VARCHAR(50),

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_payment_logs_payment
        FOREIGN KEY (payment_id)
        REFERENCES payments(id)
        ON DELETE CASCADE

);

-- =========================================================
-- 10. VERIFICATION TOKENS
-- =========================================================

CREATE TABLE verification_tokens (
id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    user_id BIGINT NOT NULL,
    token TEXT NOT NULL,

    -- TOKEN_TYPE pada ERD
    token_type VARCHAR(50),

    expires_at TIMESTAMP,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,

    CONSTRAINT fk_verification_tokens_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE

);

-- =========================================================
-- 11. NOTIFICATIONS
-- =========================================================

CREATE TABLE notifications (
id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    -- NOTIF_TYPE pada ERD
    notification_type VARCHAR(50),

    receiver_id BIGINT NOT NULL,

    subject VARCHAR(255),
    message TEXT,

    sent_at TIMESTAMP,
    read_at TIMESTAMP,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,

    CONSTRAINT fk_notifications_receiver
        FOREIGN KEY (receiver_id)
        REFERENCES users(id)
        ON DELETE CASCADE

);

-- =========================================================
-- INDEXES
-- =========================================================

CREATE INDEX idx_categories_parent_id
ON categories(parent_id);

CREATE INDEX idx_products_category_slug
ON products(category_slug);

CREATE INDEX idx_products_parent_id
ON products(parent_id);

CREATE INDEX idx_products_status
ON products(status);

CREATE INDEX idx_orders_buyer_id
ON orders(buyer_id);

CREATE INDEX idx_orders_status
ON orders(status);

CREATE INDEX idx_orders_order_date
ON orders(order_date);

CREATE INDEX idx_order_items_product_id
ON order_items(product_id);

CREATE INDEX idx_order_items_order_id
ON order_items(order_id);

CREATE INDEX idx_payments_order_id
ON payments(order_id);

CREATE INDEX idx_payments_user_id
ON payments(user_id);

CREATE INDEX idx_payments_status
ON payments(payment_status);

CREATE INDEX idx_payment_logs_payment_id
ON payment_logs(payment_id);

CREATE INDEX idx_verification_tokens_user_id
ON verification_tokens(user_id);

CREATE INDEX idx_verification_tokens_token
ON verification_tokens(token);

CREATE INDEX idx_notifications_receiver_id
ON notifications(receiver_id);

CREATE INDEX idx_notifications_read_at
ON notifications(read_at);

-- =========================================================
-- OPTIONAL: updated_at trigger
-- =========================================================

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
NEW.updated_at = CURRENT_TIMESTAMP;
RETURN NEW;
END;

$$
LANGUAGE plpgsql;


CREATE TRIGGER trg_categories_updated_at
BEFORE UPDATE ON categories
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();


CREATE TRIGGER trg_products_updated_at
BEFORE UPDATE ON products
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();


CREATE TRIGGER trg_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();


CREATE TRIGGER trg_orders_updated_at
BEFORE UPDATE ON orders
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();


CREATE TRIGGER trg_order_items_updated_at
BEFORE UPDATE ON order_items
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();


CREATE TRIGGER trg_payments_updated_at
BEFORE UPDATE ON payments
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();


CREATE TRIGGER trg_payment_logs_updated_at
BEFORE UPDATE ON payment_logs
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();


CREATE TRIGGER trg_verification_tokens_updated_at
BEFORE UPDATE ON verification_tokens
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();


CREATE TRIGGER trg_notifications_updated_at
BEFORE UPDATE ON notifications
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
$$
