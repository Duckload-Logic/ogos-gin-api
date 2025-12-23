CREATE TABLE student_addresses(
    student_address_id INT AUTO_INCREMENT PRIMARY KEY,
    student_record_id INT NOT NULL, 
    address_type_id INT NOT NULL,   
    region_name VARCHAR(100),
    province_name VARCHAR(100),
    city_name VARCHAR(100),
    barangay_name VARCHAR(100), 
    street_lot_blk VARCHAR(255),
    unit_no VARCHAR(50),         
    building_name VARCHAR(100),  
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (student_record_id) REFERENCES student_records(student_record_id) ON DELETE CASCADE,
    FOREIGN KEY (address_type_id) REFERENCES address_types(address_type_id)
);