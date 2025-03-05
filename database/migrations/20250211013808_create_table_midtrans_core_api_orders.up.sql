CREATE TABLE midtrans_core_api_orders (
    id INTEGER AUTO_INCREMENT,
    order_id INTEGER NOT NULL,
    midtrans_order_id VARCHAR(100) NOT NULL UNIQUE,
    status_code VARCHAR(3) NOT NULL,
    status_message TEXT NOT NULL,
    transaction_id VARCHAR(255) NOT NULL UNIQUE,
    gross_amount DECIMAL(10,2) NOT NULL,
    currency VARCHAR(5) NOT NULL,
    payment_type VARCHAR(50),
    transaction_time DATETIME NOT NULL,
    transaction_status VARCHAR(20),
    fraud_status VARCHAR(20),
    expiry_time DATETIME NULL DEFAULT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE CASCADE
) ENGINE = InnoDB;
