CREATE TABLE wallet_transactions (
    id INTEGER AUTO_INCREMENT PRIMARY KEY,
    user_id INTEGER NOT NULL,
    order_id INTEGER NULL,
    amount DECIMAL(15, 2) NOT NULL,
    type ENUM (
        'top_up',
        -- saat user top up (misal cash)
        'order_payment',
        -- saat user membayar pesanan
        'refund',
        -- saat pesanan dibatalkan/ditolak dan dana dikembalikan
        'withdraw' -- jika ada fitur penarikan dana
    ) NOT NULL,
    source ENUM (
        'xendit',
        -- dari pembayaran QRIS
        'cash',
        -- dari top up manual
        'refund_wallet',
        -- refund ke wallet
        'refund_xendit' -- refund dari pembayaran Xendit (tapi masuk ke wallet)
    ) NOT NULL,
    note TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (order_id) REFERENCES orders (id)
) ENGINE = InnoDB;