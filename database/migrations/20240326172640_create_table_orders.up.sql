CREATE TABLE orders (
    id INTEGER AUTO_INCREMENT,
    invoice VARCHAR(255) NOT NULL,
    amount FLOAT NOT NULL,
    discount_value FLOAT NULL,
    discount_type TINYINT(1) NOT NULL COMMENT '1 : nominal | 2 : percent',
    total_discount FLOAT NULL,
    user_id INTEGER NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(100) NOT NULL,
    payment_gateway VARCHAR(20) NOT NULL,
    payment_method VARCHAR(20) NOT NULL,
    channel_code VARCHAR(20) NOT NULL,
    payment_status TINYINT(1) NOT NULL COMMENT '-1 : cancelled | 0 : unpaid | 1 : pending | 2 : paid',
    order_status TINYINT(1) NOT NULL COMMENT '1 : pending order | 2 : order received | 3 : order being delivered | 4 : order delivered | 5 : ready for pickup | 0 : order rejected',
    is_delivery BOOLEAN NOT NULL,
    delivery_cost FLOAT NOT NULL DEFAULT 0,
    complete_address TEXT NOT NULL,
    note TEXT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
) ENGINE = InnoDB;