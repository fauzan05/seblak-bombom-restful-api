CREATE TABLE discount_coupons (
    id INTEGER AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(255) NOT NULL unique,
    value FLOAT NOT NULL DEFAULT 0,
    type ENUM("percent", "nominal") NOT NULL DEFAULT "nominal",
    start DATETIME DEFAULT CURRENT_TIMESTAMP,
    end DATETIME DEFAULT CURRENT_TIMESTAMP,
    total_max_usage INT DEFAULT 1,
    max_usage_per_user INT DEFAULT 1,
    used_count INT DEFAULT 0,
    min_order_value INT DEFAULT 0,
    description TEXT NOT NULL,
    status BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
) ENGINE = InnoDB;