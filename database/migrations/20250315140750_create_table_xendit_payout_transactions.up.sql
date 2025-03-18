CREATE TABLE xendit_payouts (
    id VARCHAR(50) PRIMARY KEY,
    user_id INTEGER NOT NULL,
    business_id VARCHAR(50) NOT NULL,
    reference_id VARCHAR(255) UNIQUE NOT NULL,
    amount DECIMAL(15, 2) NOT NULL,
    currency VARCHAR(10) NOT NULL,
    description TEXT,
    channel_code VARCHAR(50) NOT NULL,
    account_number VARCHAR(50) NOT NULL,
    account_holder_name VARCHAR(100) NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    estimated_arrival TIMESTAMP NULL DEFAULT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);