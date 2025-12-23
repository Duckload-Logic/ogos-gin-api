CREATE TABLE student_records(
    student_record_id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT UNIQUE NOT NULL,       
    educational_level_id INT,
    student_number VARCHAR(20) UNIQUE NOT NULL,
    course VARCHAR(100) NOT NULL,      -- ex 'BS Information Technology or BSIT'
    year_level INT NOT NULL,           -- ex '3'
    section VARCHAR(50),               -- ex '3-1'
    good_moral_status BOOLEAN DEFAULT TRUE,   
    has_derogatory_record BOOLEAN DEFAULT FALSE, 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (educational_level_id) REFERENCES educational_levels(educational_level_id) ON DELETE CASCADE
);