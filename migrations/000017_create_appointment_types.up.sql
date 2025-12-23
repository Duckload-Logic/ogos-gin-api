CREATE TABLE appointment_types (
    appointment_type_id INT AUTO_INCREMENT PRIMARY KEY,
    appointment_type_name VARCHAR(50) UNIQUE NOT NULL -- ex 'Walk in', 'Online', 'Referral', or 'follow up'
);