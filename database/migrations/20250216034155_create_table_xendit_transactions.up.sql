CREATE TABLE xendit_transactions (
    id VARCHAR(50) PRIMARY KEY,
    order_id INTEGER NOT NULL,
    reference_id VARCHAR(50) NOT NULL,
    amount BIGINT NOT NULL,
    currency VARCHAR(10) NOT NULL,
    payment_method VARCHAR(50) NOT NULL,
    payment_method_id VARCHAR(50),
    channel_code VARCHAR(50),
    qr_string TEXT,
    status VARCHAR(20) NOT NULL,
    description TEXT,
    expires_at TIMESTAMP NULL DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    failure_code VARCHAR(50),
    metadata JSON,
    FOREIGN KEY (order_id) REFERENCES orders (id)
) ENGINE = InnoDB;