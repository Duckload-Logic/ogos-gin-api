CREATE TABLE educational_backgrounds(
    educational_background_id INT AUTO_INCREMENT PRIMARY KEY,
    student_record_id INT NOT NULL,
	educational_level_id INT NOT NULL,  
    school_name VARCHAR(255) NOT NULL,
    location VARCHAR(255),
    school_type VARCHAR(100),     -- ex Private/Public
    year_completed VARCHAR(10),
    awards TEXT,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (student_record_id) REFERENCES student_records(student_record_id) ON DELETE CASCADE,
    FOREIGN KEY (educational_level_id) REFERENCES educational_levels(educational_level_id)
);