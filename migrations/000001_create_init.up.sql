-- ============================================================================
-- AUTH & ROLES
-- ============================================================================

CREATE TABLE user_roles(
    id INT PRIMARY KEY,
    `name` VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE users(
    id INT AUTO_INCREMENT PRIMARY KEY,
    role_id INT NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    middle_name VARCHAR(100) DEFAULT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (role_id) REFERENCES user_roles(id)
);

-- ============================================================================
-- REFERENCE DATA (Lookups & Enumerations)
-- ============================================================================

CREATE TABLE regions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL
);

CREATE TABLE cities (
    id INT AUTO_INCREMENT PRIMARY KEY,
    region_id INT NOT NULL,
    name VARCHAR(100),
    FOREIGN KEY (region_id) REFERENCES regions(id),
    CONSTRAINT unique_region_city UNIQUE KEY (region_id, name)
);

CREATE TABLE barangays (
    id INT AUTO_INCREMENT PRIMARY KEY,
    city_id INT NOT NULL,
    name VARCHAR(100),
    FOREIGN KEY (city_id) REFERENCES cities(id),
    CONSTRAINT unique_city_barangay UNIQUE KEY (city_id, name)
);

CREATE TABLE genders(
    id INT  PRIMARY KEY,
    gender_name VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE parental_status_types (
    id INT PRIMARY KEY,
    status_name VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE enrollment_reasons (
    id INT AUTO_INCREMENT PRIMARY KEY,
    reason_text VARCHAR(100) UNIQUE NOT NULL
);

CREATE TABLE income_ranges (
    id INT AUTO_INCREMENT PRIMARY KEY,
    range_text VARCHAR(100) NOT NULL
);


CREATE TABLE student_support_types (
    id INT AUTO_INCREMENT PRIMARY KEY,
    support_type_name VARCHAR(100) UNIQUE NOT NULL
);

CREATE TABLE educational_levels (
    id INT AUTO_INCREMENT PRIMARY KEY,
    level_name VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE courses (
    id INT AUTO_INCREMENT PRIMARY KEY,
    code VARCHAR(20) UNIQUE NOT NULL,
    course_name VARCHAR(100) UNIQUE NOT NULL
);

CREATE TABLE civil_status_types (
    id INT PRIMARY KEY,
    status_name VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE student_relationship_types (
    id INT AUTO_INCREMENT PRIMARY KEY,
    relationship_name VARCHAR(100) UNIQUE NOT NULL
);

CREATE TABLE nature_of_residence_types (
    id INT AUTO_INCREMENT PRIMARY KEY,
    residence_type_name VARCHAR(100) UNIQUE NOT NULL
);

CREATE TABLE religions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    religion_name VARCHAR(100) UNIQUE NOT NULL
);

CREATE TABLE sibling_support_types (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL
);

CREATE TABLE activity_options (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL, -- e.g., 'Math Club', 'Chess Club'
    category ENUM('academic', 'extra_curricular') NOT NULL,
    is_active BOOLEAN DEFAULT TRUE -- Allows you to "retire" old clubs without deleting data
);

-- ============================================================================
-- COUNSELING & APPOINTMENTS
-- ============================================================================
CREATE TABLE time_slots (
    id INT AUTO_INCREMENT PRIMARY KEY,
    time TIME NOT NULL UNIQUE
);

CREATE TABLE statuses (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    color_key ENUM('warning', 'danger', 'success', 'info', 'stale', 'notice') NOT NULL,
    status_type ENUM('appointment', 'slip', 'both') NOT NULL DEFAULT 'both'
);

CREATE TABLE appointment_categories (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE
);

CREATE TABLE appointments (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT,
    time_slot_id INT NOT NULL,
    when_date DATE NOT NULL,
    reason TEXT,
    admin_notes TEXT,
    appointment_category_id INT NOT NULL,
    status_id INT NOT NULL DEFAULT 1, -- Default to 'Pending'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (time_slot_id) REFERENCES time_slots(id),
    FOREIGN KEY (status_id) REFERENCES statuses(id),
    FOREIGN KEY (appointment_category_id) REFERENCES appointment_categories(id),

    -- PREVENTS DOUBLE BOOKING:
    UNIQUE KEY unique_appointment (when_date, time_slot_id)
);

CREATE TABLE counselor_profiles(
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT UNIQUE NOT NULL,
    license_number VARCHAR(50),
    specialization VARCHAR(100),
    is_available BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
