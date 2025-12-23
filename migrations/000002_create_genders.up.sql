CREATE TABLE genders(
    gender_id INT AUTO_INCREMENT PRIMARY KEY,
    gender_name VARCHAR(50) UNIQUE NOT NULL -- ex 'Male', 'Female', 'Prefer not to say'
);
