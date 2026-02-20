CREATE TABLE roles(
    role_id INT PRIMARY KEY,
    role_name VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE enrollment_reasons (
    id INT AUTO_INCREMENT PRIMARY KEY,
    reason_text VARCHAR(100) UNIQUE NOT NULL
);

CREATE TABLE genders(
    gender_id INT  PRIMARY KEY,
    gender_name VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE civil_status_types (
    civil_status_type_id INT PRIMARY KEY,
    status_name VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE parental_status_types (
    parental_status_id INT PRIMARY KEY,
    status_name VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE users(
    user_id INT AUTO_INCREMENT PRIMARY KEY,
    role_id INT NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    middle_name VARCHAR(100),
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (role_id) REFERENCES roles(role_id)
);

CREATE TABLE iir_records(
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT UNIQUE NOT NULL,
    is_submitted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

CREATE TABLE student_profiles(
    id INT AUTO_INCREMENT PRIMARY KEY,
    student_record_id INT UNIQUE NOT NULL,
    gender_id INT,
    civil_status_type_id INT NOT NULL,
    religion VARCHAR(100) NOT NULL,
    height_ft DECIMAL(5,2) NOT NULL,
    weight_kg DECIMAL(5,2) NOT NULL,
    student_number VARCHAR(20) UNIQUE NOT NULL,
    high_school_gwa DECIMAL(4,2) NOT NULL,
    course VARCHAR(100) NOT NULL,
    year_level VARCHAR(50),
    section VARCHAR(50),
    place_of_birth VARCHAR(255),
    birth_date DATE,
    contact_no VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (student_record_id) REFERENCES iir_records(id) ON DELETE CASCADE,
    FOREIGN KEY (gender_id) REFERENCES genders(gender_id),
    FOREIGN KEY (civil_status_type_id) REFERENCES civil_status_types(civil_status_type_id) ON DELETE CASCADE
);

CREATE TABLE student_selected_reasons (
    iir_id INT NOT NULL,
    reason_id INT NOT NULL,
    other_reason_text VARCHAR(255) DEFAULT NULL,
    PRIMARY KEY (iir_id, reason_id),
    FOREIGN KEY (iir_id) REFERENCES iir_records(id) ON DELETE CASCADE,
    FOREIGN KEY (reason_id) REFERENCES enrollment_reasons(id) ON DELETE CASCADE
);

CREATE TABLE excuse_slips (
    excuse_slip_id INT AUTO_INCREMENT PRIMARY KEY,
    iir_id INT NOT NULL,
    reason TEXT NOT NULL,
    date_of_absence DATE NOT NULL,
    file_path VARCHAR(255) NOT NULL,
    excuse_slip_status ENUM('Pending', 'Approved', 'Rejected') DEFAULT 'Pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (iir_id) REFERENCES iir_records(id) ON DELETE CASCADE
);

CREATE TABLE addresses(
    id INT AUTO_INCREMENT PRIMARY KEY,
    region VARCHAR(100) NOT NULL,
    city VARCHAR(100) NOT NULL,
    barangay VARCHAR(100) NOT NULL,
    street_detail VARCHAR(255)
);

CREATE TABLE related_persons(
    id INT AUTO_INCREMENT PRIMARY KEY,
    address_id INT,
    educational_level VARCHAR(100) NOT NULL,
    birth_date DATE,
    last_name VARCHAR(100) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    middle_name VARCHAR(100) DEFAULT NULL,
    occupation VARCHAR(100) DEFAULT NULL,
    employer_name VARCHAR(150) DEFAULT NULL,
    employer_address VARCHAR(255) DEFAULT NULL,
    is_living BOOLEAN DEFAULT TRUE,
    FOREIGN KEY (address_id) REFERENCES addresses(id) ON DELETE SET NULL
);

CREATE TABLE student_related_persons (
    iir_id INT NOT NULL,
    related_person_id INT NOT NULL,
    relationship ENUM('Father', 'Mother', 'Guardian', 'Uncle', 'Aunt', 'Grandparent', 'Sibling', 'Other'),
    is_parent BOOLEAN DEFAULT FALSE,
    is_guardian BOOLEAN DEFAULT FALSE,
    is_emergency_contact BOOLEAN DEFAULT FALSE,
    PRIMARY KEY (iir_id, related_person_id),
    FOREIGN KEY (iir_id) REFERENCES iir_records(id) ON DELETE CASCADE,
    FOREIGN KEY (related_person_id) REFERENCES related_persons(id) ON DELETE CASCADE
);

CREATE TABLE student_emergency_contacts (
    emergency_contact_id INT AUTO_INCREMENT PRIMARY KEY,
    iir_id INT UNIQUE NOT NULL,
    related_person_id INT,
    emergency_contact_first_name VARCHAR(255) NOT NULL,
    emergency_contact_middle_name VARCHAR(255),
    emergency_contact_last_name VARCHAR(255) NOT NULL,
    emergency_contact_phone VARCHAR(20) NOT NULL,
    emergency_contact_relationship VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (iir_id) REFERENCES iir_records(id) ON DELETE CASCADE,
    FOREIGN KEY (related_person_id) REFERENCES related_persons(id) ON DELETE SET NULL
);

CREATE TABLE family_backgrounds(
    id INT AUTO_INCREMENT PRIMARY KEY,
    iir_id INT UNIQUE NOT NULL,
    parental_status_id INT NOT NULL,
    parental_status_details VARCHAR(255) DEFAULT NULL,
    siblings_brothers INT NOT NULL,
    sibling_sisters INT NOT NULL,
    monthly_family_income VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (iir_id) REFERENCES iir_records(id) ON DELETE CASCADE,
    FOREIGN KEY (parental_status_id) REFERENCES parental_status_types(parental_status_id) ON DELETE CASCADE
);

CREATE TABLE student_addresses(
    id INT AUTO_INCREMENT PRIMARY KEY,
    iir_id INT NOT NULL,
    address_id INT NOT NULL,
    address_type ENUM('Residential', 'Provincial') NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (iir_id) REFERENCES iir_records(id) ON DELETE CASCADE,
    FOREIGN KEY (address_id) REFERENCES addresses(id) ON DELETE CASCADE
);

CREATE TABLE educational_backgrounds(
    id INT AUTO_INCREMENT PRIMARY KEY,
    iir_id INT NOT NULL,
	educational_level ENUM(
        'Elementary',
        'Junior High School',
        'Senior High School'
    ) NOT NULL,
    school_name VARCHAR(255) NOT NULL,
    location VARCHAR(255),
    school_type ENUM('Private', 'Public'),
    year_completed VARCHAR(10),
    awards TEXT,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (iir_id) REFERENCES iir_records(id) ON DELETE CASCADE
);

CREATE TABLE counselor_profiles(
    counselor_profile_id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT UNIQUE NOT NULL,
    license_number VARCHAR(50),
    specialization VARCHAR(100),
    is_available BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

CREATE TABLE appointments (
    appointment_id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT,
    reason VARCHAR(255),
    scheduled_date DATE NOT NULL,
    scheduled_time TIME NOT NULL,
    concern_category VARCHAR(100),
    `status` ENUM('Pending', 'Approved', 'Completed', 'Cancelled', 'Rescheduled') DEFAULT 'Pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE SET NULL
);

CREATE TABLE session_notes (
    session_note_id INT AUTO_INCREMENT PRIMARY KEY,
    appointment_id INT UNIQUE NOT NULL,
    notes TEXT,
    recommendation TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (appointment_id) REFERENCES appointments(appointment_id) ON DELETE CASCADE
);

CREATE TABLE student_finances (
    iir_id INT PRIMARY KEY,
    employed_family_members_count INT,
    supports_studies_count INT,
    supports_family_count INT,
    financial_support VARCHAR(255) NOT NULL,
    weekly_allowance DECIMAL(10,2),
    FOREIGN KEY (iir_id) REFERENCES iir_records(id) ON DELETE CASCADE
);

CREATE TABLE student_health_records (
    health_id INT NOT NULL AUTO_INCREMENT,
    iir_id INT UNIQUE NOT NULL,
    vision_remark ENUM('No Problem', 'Issues'),
    hearing_remark ENUM('No Problem', 'Issues'),
    mobility_remark ENUM('No Problem', 'Issues'),
    speech_remark ENUM('No Problem', 'Issues'),
    general_health_remark ENUM('No Problem', 'Issues'),
    consulted_professional VARCHAR(255),
    consultation_reason TEXT,
    date_started DATE,
    num_sessions INT,
    date_concluded DATE,
    PRIMARY KEY (health_id),
    FOREIGN KEY (iir_id) REFERENCES iir_records(id) ON DELETE CASCADE
);

CREATE TABLE student_interests (
    interest_id INT AUTO_INCREMENT PRIMARY KEY,
    iir_id INT NOT NULL,
    interest_type VARCHAR(100) NOT NULL,
    interest_name VARCHAR(255) NOT NULL,
    is_favorite BOOLEAN DEFAULT FALSE,
    is_least_favorite BOOLEAN DEFAULT FALSE,
    `rank` INT,
    FOREIGN KEY (iir_id) REFERENCES iir_records(id) ON DELETE CASCADE
);

CREATE TABLE test_results (
    test_result_id INT AUTO_INCREMENT PRIMARY KEY,
    iir_id INT NOT NULL,
    test_date DATE,
    test_name VARCHAR(255),
    raw_score VARCHAR(50),
    percentile VARCHAR(50),
    description VARCHAR(255),
    FOREIGN KEY (iir_id) REFERENCES iir_records(id) ON DELETE CASCADE
);