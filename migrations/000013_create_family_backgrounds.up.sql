CREATE TABLE family_backgrounds(
    family_background_id INT AUTO_INCREMENT PRIMARY KEY, 
    student_record_id INT UNIQUE NOT NULL,
    father_id INT NOT NULL,
    mother_id INT NOT NULL,
    guardian_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (student_record_id) REFERENCES student_records(student_record_id) ON DELETE CASCADE,
    FOREIGN KEY (father_id) REFERENCES guardians(guardian_id) ON DELETE CASCADE,
    FOREIGN KEY (mother_id) REFERENCES guardians(guardian_id) ON DELETE CASCADE,
    FOREIGN KEY (guardian_id) REFERENCES guardians(guardian_id) ON DELETE CASCADE
);