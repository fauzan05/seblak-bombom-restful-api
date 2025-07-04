CREATE TABLE users (
    id INTEGER AUTO_INCREMENT,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    email_verified BOOLEAN DEFAULT FALSE,
    verification_token VARCHAR(64),
    token_expiry TIMESTAMP NULL DEFAULT NULL,
    phone VARCHAR(100) NOT NULL,
    password VARCHAR(100) NOT NULL,
    role ENUM ("admin", "customer") NOT NULL,
    user_profile TEXT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    PRIMARY KEY (id)
) ENGINE = InnoDB;