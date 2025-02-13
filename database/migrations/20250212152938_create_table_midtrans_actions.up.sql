CREATE TABLE midtrans_actions (
    id INTEGER AUTO_INCREMENT,
    midtrans_core_api_orders_id INTEGER NOT NULL,
    name VARCHAR(100) NOT NULL,
    method VARCHAR(4) NOT NULL,
    url TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY (midtrans_core_api_orders_id) REFERENCES midtrans_core_api_orders (id)
) ENGINE = InnoDB;