CREATE TABLE roles(
    role_id INT AUTO_INCREMENT PRIMARY KEY,
    role_name VARCHAR(50) UNIQUE NOT NULL 
);

CREATE TABLE genders(
    gender_id INT AUTO_INCREMENT PRIMARY KEY,
    gender_name VARCHAR(50) UNIQUE NOT NULL 
);

CREATE TABLE educational_levels(
    educational_level_id INT AUTO_INCREMENT PRIMARY KEY,
    level_name VARCHAR(100) UNIQUE NOT NULL -- ex 'College', 'Vocational'
);

CREATE TABLE address_types(
    address_type_id INT AUTO_INCREMENT PRIMARY KEY,
    type_name VARCHAR(50) UNIQUE NOT NULL -- ex 'Residential', 'Permanent'
);

CREATE TABLE civil_status_types (
    civil_status_type_id INT AUTO_INCREMENT PRIMARY KEY,
    status_name VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE religion_types (
    religion_type_id INT AUTO_INCREMENT PRIMARY KEY,
    religion_name VARCHAR(100) UNIQUE NOT NULL
);

CREATE TABLE relationship_types (
    relationship_type_id INT AUTO_INCREMENT PRIMARY KEY,
    relationship_name VARCHAR(50) UNIQUE NOT NULL -- ex: 'Father', 'Mother', 'Step-father', 'Aunt'
);

CREATE TABLE parental_status_types (
    parental_status_id INT AUTO_INCREMENT PRIMARY KEY,
    status_name VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE financial_support_types (
    financial_support_type_id INT AUTO_INCREMENT PRIMARY KEY,
    support_type_name VARCHAR(100) UNIQUE NOT NULL
);

CREATE TABLE health_remark_types (
    health_remark_type_id INT AUTO_INCREMENT PRIMARY KEY,
    remark_name VARCHAR(100) UNIQUE NOT NULL
);

CREATE TABLE appointment_types (
    appointment_type_id INT AUTO_INCREMENT PRIMARY KEY,
    appointment_type_name VARCHAR(50) UNIQUE NOT NULL -- ex 'Walk in', 'Online', 'Referral', or 'follow up'
);

CREATE TABLE users(
    user_id INT AUTO_INCREMENT PRIMARY KEY,
    role_id INT NOT NULL,
    gender_id INT, 
    first_name VARCHAR(100) NOT NULL,
    middle_name VARCHAR(100),
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    place_of_birth VARCHAR(255),
    birth_date DATE,
    mobile_no VARCHAR(20),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (role_id) REFERENCES roles(role_id),
    FOREIGN KEY (gender_id) REFERENCES genders(gender_id)
);

CREATE TABLE student_records(
    student_record_id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT UNIQUE NOT NULL,       
    civil_status_type_id INT NOT NULL,
    religion_type_id INT NOT NULL,
    height_cm DECIMAL(5,2) NOT NULL,               -- ex '165.50' (in cm)
    weight_kg DECIMAL(5,2) NOT NULL,               -- ex '65.50' (in kg)
    student_number VARCHAR(20) UNIQUE NOT NULL,
    course VARCHAR(100) NOT NULL,      -- ex 'BS Information Technology or BSIT'
    year_level INT NOT NULL,           -- ex '3'
    section VARCHAR(50),               -- ex '3-1'
    good_moral_status BOOLEAN DEFAULT TRUE,   
    has_derogatory_record BOOLEAN DEFAULT FALSE, 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (civil_status_type_id) REFERENCES civil_status_types(civil_status_type_id) ON DELETE CASCADE,
    FOREIGN KEY (religion_type_id) REFERENCES religion_types(religion_type_id) ON DELETE CASCADE
);

CREATE TABLE excuse_slips (
    excuse_slip_id INT AUTO_INCREMENT PRIMARY KEY,
    student_record_id INT NOT NULL,
    reason TEXT NOT NULL,
    date_of_absence DATE NOT NULL,
    file_path VARCHAR(255) NOT NULL,
    excuse_slip_status ENUM('Pending', 'Approved', 'Rejected') DEFAULT 'Pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (student_record_id) REFERENCES student_records(student_record_id) ON DELETE CASCADE
);

CREATE TABLE guardians(
    guardian_id INT AUTO_INCREMENT PRIMARY KEY,
    educational_level_id INT NOT NULL,
    birth_date DATE NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    middle_name VARCHAR(100) DEFAULT NULL,
    occupation VARCHAR(100) DEFAULT NULL,
    maiden_name VARCHAR(100) DEFAULT NULL, -- only if the parent type is mother
    company_name VARCHAR(150) DEFAULT NULL,
    contact_number VARCHAR(20) DEFAULT NULL,
	FOREIGN KEY (educational_level_id) REFERENCES educational_levels(educational_level_id)
);

CREATE TABLE student_guardians (
    student_record_id INT NOT NULL,
    guardian_id INT NOT NULL,
    relationship_type_id INT NOT NULL,
    is_primary_contact BOOLEAN DEFAULT FALSE,
    PRIMARY KEY (student_record_id, guardian_id),
    FOREIGN KEY (student_record_id) REFERENCES student_records(student_record_id) ON DELETE CASCADE,
    FOREIGN KEY (guardian_id) REFERENCES guardians(guardian_id) ON DELETE CASCADE,
    FOREIGN KEY (relationship_type_id) REFERENCES relationship_types(relationship_type_id)
);

CREATE TABLE family_backgrounds(
    family_background_id INT AUTO_INCREMENT PRIMARY KEY, 
    student_record_id INT UNIQUE NOT NULL,
    parental_status_id INT NOT NULL,
    parental_status_details VARCHAR(255) DEFAULT NULL,
    siblings_brothers INT NOT NULL,
    sibling_sisters INT NOT NULL,
    monthly_family_income DECIMAL(10,2) NOT NULL,  -- e.g. '50000.00'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (student_record_id) REFERENCES student_records(student_record_id) ON DELETE CASCADE,
    FOREIGN KEY (parental_status_id) REFERENCES parental_status_types(parental_status_id) ON DELETE CASCADE
);

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

CREATE TABLE educational_backgrounds(
    educational_background_id INT AUTO_INCREMENT PRIMARY KEY,
    student_record_id INT NOT NULL,
	educational_level_id INT NOT NULL,  
    school_name VARCHAR(255) NOT NULL,
    location VARCHAR(255),
    school_type VARCHAR(100),     -- ex Private/Public
    year_completed VARCHAR(10),
    awards TEXT,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (student_record_id) REFERENCES student_records(student_record_id) ON DELETE CASCADE,
    FOREIGN KEY (educational_level_id) REFERENCES educational_levels(educational_level_id)
);

CREATE TABLE counselor_profiles(
    counselor_profile_id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT UNIQUE NOT NULL, 
    license_number VARCHAR(50), 
    specialization VARCHAR(100), -- e.g., 'Mental Health', 'Career'
    is_available BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

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

CREATE TABLE session_notes (
    session_note_id INT AUTO_INCREMENT PRIMARY KEY,
    appointment_id INT UNIQUE NOT NULL,
    notes TEXT,
    recommendation TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (appointment_id) REFERENCES appointments(appointment_id) ON DELETE CASCADE
);

CREATE TABLE student_finances (
    finance_id INT NOT NULL AUTO_INCREMENT,
    student_record_id INT NOT NULL,
    is_employed TINYINT(1),
    supports_studies TINYINT(1),
    supports_family TINYINT(1),
    financial_support_type_id INT NOT NULL, 
    weekly_allowance DECIMAL(10,2),
    PRIMARY KEY (finance_id),
    FOREIGN KEY (student_record_id) REFERENCES student_records(student_record_id)
        ON DELETE CASCADE,
    FOREIGN KEY (financial_support_type_id) REFERENCES financial_support_types(financial_support_type_id)
        ON DELETE CASCADE   
);

CREATE TABLE student_health_records (
    health_id INT NOT NULL AUTO_INCREMENT,
    student_record_id INT NOT NULL,
    vision_remark_id INT NOT NULL,
    hearing_remark_id INT NOT NULL,
    mobility_remark_id INT NOT NULL,
    speech_remark_id INT NOT NULL,
    general_health_remark_id INT NOT NULL,
    consulted_professional VARCHAR(255),
    consultation_reason TEXT,
    date_started DATE,
    num_sessions INT,
    date_concluded DATE,
    PRIMARY KEY (health_id),
    FOREIGN KEY (student_record_id) REFERENCES student_records(student_record_id)
        ON DELETE CASCADE,
    FOREIGN KEY (vision_remark_id) REFERENCES health_remark_types(health_remark_type_id)
        ON DELETE CASCADE,
    FOREIGN KEY (hearing_remark_id) REFERENCES health_remark_types(health_remark_type_id)
        ON DELETE CASCADE,
    FOREIGN KEY (mobility_remark_id) REFERENCES health_remark_types(health_remark_type_id)
        ON DELETE CASCADE,
    FOREIGN KEY (speech_remark_id) REFERENCES health_remark_types(health_remark_type_id)
        ON DELETE CASCADE,
    FOREIGN KEY (general_health_remark_id) REFERENCES health_remark_types(health_remark_type_id)
        ON DELETE CASCADE
);

CREATE TABLE psychological_assessments (
    assessment_id INT NOT NULL AUTO_INCREMENT,
    student_record_id INT NOT NULL,
    test_date DATE,
    test_name VARCHAR(255),
    raw_score VARCHAR(50), -- (RS)
    remarks TEXT,
    PRIMARY KEY (assessment_id),
    FOREIGN KEY (student_record_id) REFERENCES student_records(student_record_id)
        ON DELETE CASCADE
);