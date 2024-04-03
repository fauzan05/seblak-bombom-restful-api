CREATE TABLE orders (
    id INTEGER AUTO_INCREMENT,
    invoice VARCHAR(255) NOT NULL,
    product_id INTEGER NOT NULL,
    product_name VARCHAR(255) NOT NULL,
    product_description TEXT NOT NULL,
    price INTEGER NOT NULL,
    quantity INTEGER NOT NULL,
    amount INTEGER NOT NULL,
    discount_value INTEGER NULL,
    discount_type ENUM("percent", "nominal") NOT NULL DEFAULT "nominal",
    user_id INTEGER NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(100) NOT NULL,
    payment_method ENUM("online", "onsite") NOT NULL DEFAULT "online",
    payment_status ENUM("pending", "paid", "failed") NOT NULL DEFAULT "pending",
    delivery_status ENUM("prepare", "on_the_way", "sent") NOT NULL DEFAULT "prepare",
    is_delivery BOOLEAN NOT NULL,
    deliver_cost INTEGER NOT NULL DEFAULT 0,
    category_name VARCHAR(100) NOT NULL,
    complete_address TEXT NOT NULL,
    google_map_link TEXT NOT NULL,
    distance INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
) ENGINE = InnoDB;