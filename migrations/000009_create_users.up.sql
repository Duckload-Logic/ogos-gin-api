CREATE TABLE users(
    user_id INT AUTO_INCREMENT PRIMARY KEY,
    role_id INT NOT NULL,
    gender_id INT, 
    first_name VARCHAR(100) NOT NULL,
    middle_name VARCHAR(100),
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    place_of_birth VARCHAR(255),
    birth_date DATE,
    mobile_no VARCHAR(20),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (role_id) REFERENCES roles(role_id),
    FOREIGN KEY (gender_id) REFERENCES genders(gender_id)
);