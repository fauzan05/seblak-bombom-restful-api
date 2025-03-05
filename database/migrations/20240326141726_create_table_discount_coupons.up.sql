CREATE TABLE discount_coupons (
    id INTEGER AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(255) NOT NULL unique,
    value FLOAT NOT NULL DEFAULT 0,
    type TINYINT(1) NOT NULL COMMENT '1 : nominal | 2 : percent',
    start DATETIME NULL DEFAULT NULL,
    end DATETIME NULL DEFAULT NULL,
    total_max_usage INT DEFAULT 1,
    max_usage_per_user INT DEFAULT 1,
    used_count INT DEFAULT 0,
    min_order_value INT DEFAULT 0,
    description TEXT NOT NULL,
    status BOOLEAN NOT NULL DEFAULT false,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL DEFAULT NULL,
    PRIMARY KEY (id)
) ENGINE = InnoDB;