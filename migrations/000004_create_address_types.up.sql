CREATE TABLE address_types(
    address_type_id INT AUTO_INCREMENT PRIMARY KEY,
    type_name VARCHAR(50) UNIQUE NOT NULL -- ex 'Residential', 'Permanent'
);