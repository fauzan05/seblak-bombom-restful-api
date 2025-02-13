CREATE TABLE midtrans_core_api_orders (
    id INTEGER AUTO_INCREMENT,
    order_id INTEGER NOT NULL,
    midtrans_order_id VARCHAR(100) NOT NULL,
    status_code VARCHAR(4) NOT NULL,
    status_message TEXT NOT NULL,
    transaction_id VARCHAR(255) NOT NULL,
    gross_amount FLOAT NOT NULL,
    currency VARCHAR(5) NOT NULL,
    payment_type VARCHAR(50),
    transaction_time TIMESTAMP NOT NULL,
    transaction_status VARCHAR(20),
    fraud_status VARCHAR(20),
    expiry_time TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY (order_id) REFERENCES orders (id)
) ENGINE = InnoDB;