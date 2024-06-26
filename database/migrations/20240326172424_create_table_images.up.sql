CREATE TABLE images (
    id INTEGER AUTO_INCREMENT,
    product_id INTEGER NOT NULL,
    file_name TEXT NOT NULL,
    type VARCHAR(10) NULL,
    position INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY (product_id) REFERENCES products (id)
) ENGINE = InnoDB;