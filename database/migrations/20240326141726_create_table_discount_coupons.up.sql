CREATE TABLE discount_coupons (
    id INTEGER AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(255) NOT NULL unique,
    value DECIMAL(15, 2) NOT NULL,
    type ENUM("nominal", "percent") NOT NULL,
    start TIMESTAMP NULL DEFAULT NULL,
    end TIMESTAMP NULL DEFAULT NULL,
    max_usage_per_user INT DEFAULT 1,
    used_count INT DEFAULT 0,
    min_order_value DECIMAL(15, 2) DEFAULT 0,
    description TEXT NOT NULL,
    status BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    PRIMARY KEY (id)
) ENGINE = InnoDB;