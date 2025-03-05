CREATE TABLE addresses (
    id INTEGER AUTO_INCREMENT,
    user_id INTEGER NOT NULL,
    delivery_id INTEGER NOT NULL,
    complete_address TEXT DEFAULT NULL,
    google_maps_link TEXT DEFAULT NULL,
    is_main BOOLEAN NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL DEFAULT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (delivery_id) REFERENCES deliveries (id)
) ENGINE = InnoDB;