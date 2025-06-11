CREATE TABLE wallet_transactions (
    id INTEGER AUTO_INCREMENT PRIMARY KEY,
    user_id INTEGER NOT NULL,
    order_id INTEGER NULL,
    amount DECIMAL(15, 2) NOT NULL,
    flow_type ENUM (
        'debit',
        -- uang keluar dari wallet
        'credit' -- uang masuk ke wallet
    ) NOT NULL,
    transaction_type ENUM (
        'top_up',
        -- user isi saldo
        'order_payment',
        -- bayar pesanan pakai wallet
        'order_refund',
        -- refund dari pembatalalan pesanan
        'withdraw',
        -- tarik dana ke rekening/cash
        'admin_adjustment',
        -- penyesuaian manual oleh admin
        'cashback',
        -- jika ada program cashback
        'transfer_in',
        -- terima transfer dari user lain
        'transfer_out' -- kirim transfer ke user lain
    ) NOT NULL,
    payment_method ENUM (
        'CASH',
        -- top up/withdraw tunai
        'QR_CODE',
        -- top up via payment gateway
        'WALLET',
        -- transaksi menggunakan wallet
        'EWALLET' -- transaksi sistem (refund, cashback, dll)
    ) NULL,
    status ENUM (
        'pending',
        -- menunggu konfirmasi (untuk withdraw/top up manual)
        'processing',
        -- sedang diproses
        'completed',
        -- selesai
        'failed',
        -- gagal
        'cancelled' -- dibatalkan
    ) NOT NULL DEFAULT 'completed',
    reference_number VARCHAR(100) NULL,
    -- nomor referensi dari xendit
    note TEXT,
    admin_note TEXT,
    -- catatan khusus admin
    processed_by INTEGER NULL,
    -- admin yang memproses (untuk manual)
    processed_at TIMESTAMP NULL,
    -- waktu diproses
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (order_id) REFERENCES orders (id),
    FOREIGN KEY (processed_by) REFERENCES users (id),
    INDEX idx_user_id (user_id),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at),
    INDEX idx_reference (reference_number)
) ENGINE = InnoDB;