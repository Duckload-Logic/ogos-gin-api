CREATE TABLE appointments (
    appointment_id INT AUTO_INCREMENT PRIMARY KEY,
    student_record_id INT NOT NULL,
    counselor_user_id INT, 
    appointment_type_id INT NOT NULL,
    scheduled_date DATE NOT NULL,
    scheduled_time TIME NOT NULL,
    concern_category VARCHAR(100), -- ex 'Academic', 'Personal', 'Career'
    status ENUM('Pending', 'Approved', 'Completed', 'Cancelled', 'Rescheduled') DEFAULT 'Pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (student_record_id) REFERENCES student_records(student_record_id) ON DELETE CASCADE,
    FOREIGN KEY (counselor_user_id) REFERENCES users(user_id) ON DELETE SET NULL,
    FOREIGN KEY (appointment_type_id) REFERENCES appointment_types(appointment_type_id)
);