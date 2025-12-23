CREATE TABLE educational_levels(
    educational_level_id INT AUTO_INCREMENT PRIMARY KEY,
    level_name VARCHAR(100) UNIQUE NOT NULL -- ex 'College', 'Vocational'
);