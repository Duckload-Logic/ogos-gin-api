CREATE TABLE roles(
    role_id INT AUTO_INCREMENT PRIMARY KEY,
    role_name VARCHAR(50) UNIQUE NOT NULL -- ex 'Student', 'Counselor', 'Admin'
);