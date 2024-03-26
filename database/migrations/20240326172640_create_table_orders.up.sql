CREATE TABLE orders (
    id INTEGER AUTO_INCREMENT,
    product_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    status ENUM("pending", "success", "failed") NOT NULL,
    amount INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY (product_id) REFERENCES products (id),
    FOREIGN KEY (user_id) REFERENCES users (id)
) ENGINE = InnoDB;