CREATE TABLE counselor_profiles(
    counselor_profile_id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT UNIQUE NOT NULL, 
    license_number VARCHAR(50), 
    specialization VARCHAR(100), -- e.g., 'Mental Health', 'Career'
    is_available BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);