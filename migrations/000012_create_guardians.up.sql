CREATE TABLE guardians(
    guardian_id INT AUTO_INCREMENT PRIMARY KEY,
    guardian_type_id INT NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    middle_name VARCHAR(100) DEFAULT NULL,
    occupation VARCHAR(100) DEFAULT NULL,
    maiden_name VARCHAR(100) DEFAULT NULL, -- only if the parent type is mother
	FOREIGN KEY (guardian_type_id) REFERENCES guardian_types(guardian_type_id)
);