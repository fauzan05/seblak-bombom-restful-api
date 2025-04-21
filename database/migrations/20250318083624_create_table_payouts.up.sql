CREATE TABLE payouts (
    id INTEGER AUTO_INCREMENT PRIMARY KEY,
    user_id INTEGER NOT NULL,
    xendit_payout_id VARCHAR(50) NULL,
    amount DECIMAL(15, 2) NOT NULL,
    currency VARCHAR(10) NOT NULL DEFAULT 'IDR',
    method ENUM("offline", "online") NOT NULL COMMENT '0 : offline | 1 : online',
    status ENUM(
        "pending",
        "accepted",
        "cancelled",
        "failed",
        "succeeded",
        "refunded",
        "expired"
    ) NOT NULL,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (xendit_payout_id) REFERENCES xendit_payouts (id)
);