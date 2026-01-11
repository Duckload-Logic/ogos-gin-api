CREATE TABLE roles(
    role_id INT PRIMARY KEY,
    role_name VARCHAR(50) UNIQUE NOT NULL 
);

CREATE TABLE enrollment_reasons (
    reason_id INT AUTO_INCREMENT PRIMARY KEY,
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

CREATE TABLE student_records(
    student_record_id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT UNIQUE NOT NULL,
    is_submitted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

CREATE TABLE student_profiles(
    student_profile_id INT AUTO_INCREMENT PRIMARY KEY,
    student_record_id INT UNIQUE NOT NULL,   
    gender_id INT, 
    civil_status_type_id INT NOT NULL,
    religion VARCHAR(100) NOT NULL,
    height_ft DECIMAL(5,2) NOT NULL,               -- ex '5.5' (in ft)
    weight_kg DECIMAL(5,2) NOT NULL,               -- ex '65.50' (in kg)
    student_number VARCHAR(20) UNIQUE NOT NULL,
    high_school_gwa DECIMAL(4,2) NOT NULL,        -- ex '87.50'
    course VARCHAR(100) NOT NULL,      
    place_of_birth VARCHAR(255),
    birth_date DATE,
    contact_no VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (student_record_id) REFERENCES student_records(student_record_id) ON DELETE CASCADE,
    FOREIGN KEY (gender_id) REFERENCES genders(gender_id),
    FOREIGN KEY (civil_status_type_id) REFERENCES civil_status_types(civil_status_type_id) ON DELETE CASCADE
);

CREATE TABLE student_selected_reasons (
    student_record_id INT NOT NULL,
    reason_id INT NOT NULL,
    other_reason_text VARCHAR(255) DEFAULT NULL, -- To handle the "Others: Please specify" field
    PRIMARY KEY (student_record_id, reason_id),
    FOREIGN KEY (student_record_id) REFERENCES student_records(student_record_id) ON DELETE CASCADE,
    FOREIGN KEY (reason_id) REFERENCES enrollment_reasons(reason_id) ON DELETE CASCADE
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

CREATE TABLE parents(
    parent_id INT AUTO_INCREMENT PRIMARY KEY,
    educational_level VARCHAR(100) NOT NULL,
    birth_date DATE NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    middle_name VARCHAR(100) DEFAULT NULL,
    occupation VARCHAR(100) DEFAULT NULL,
    company_name VARCHAR(150) DEFAULT NULL
);

CREATE TABLE student_parents (
    student_record_id INT NOT NULL,
    parent_id INT NOT NULL,
    relationship ENUM('Father', 'Mother'),
    PRIMARY KEY (student_record_id, parent_id),
    FOREIGN KEY (student_record_id) REFERENCES student_records(student_record_id) ON DELETE CASCADE,
    FOREIGN KEY (parent_id) REFERENCES parents(parent_id) ON DELETE CASCADE
);

CREATE TABLE student_emergency_contacts (
    emergency_contact_id INT AUTO_INCREMENT PRIMARY KEY,
    student_record_id INT UNIQUE NOT NULL,
    parent_id INT, -- Can be NULL if emergency contact is not a guardian
    emergency_contact_first_name VARCHAR(255) NOT NULL,
    emergency_contact_middle_name VARCHAR(255),
    emergency_contact_last_name VARCHAR(255) NOT NULL,
    emergency_contact_phone VARCHAR(20) NOT NULL,
    emergency_contact_relationship VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (student_record_id) REFERENCES student_records(student_record_id) ON DELETE CASCADE,
    FOREIGN KEY (parent_id) REFERENCES parents(parent_id) ON DELETE SET NULL
);

CREATE TABLE family_backgrounds(
    family_background_id INT AUTO_INCREMENT PRIMARY KEY, 
    student_record_id INT UNIQUE NOT NULL,
    parental_status_id INT NOT NULL,
    parental_status_details VARCHAR(255) DEFAULT NULL,
    siblings_brothers INT NOT NULL,
    sibling_sisters INT NOT NULL,
    monthly_family_income VARCHAR(50) NOT NULL,  -- e.g. '45000-50000'
    guardian_first_name VARCHAR(255) NOT NULL,
    guardian_last_name VARCHAR(255) NOT NULL,
    guardian_middle_name VARCHAR(255),
    guardian_address VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (student_record_id) REFERENCES student_records(student_record_id) ON DELETE CASCADE,
    FOREIGN KEY (parental_status_id) REFERENCES parental_status_types(parental_status_id) ON DELETE CASCADE
);

CREATE TABLE student_addresses(
    student_address_id INT AUTO_INCREMENT PRIMARY KEY,
    student_record_id INT NOT NULL, 
    address_type ENUM('Residential', 'Provincial') NOT NULL,
    region_name VARCHAR(100),
    province_name VARCHAR(100),
    city_name VARCHAR(100),
    barangay_name VARCHAR(100), 
    street_lot_blk VARCHAR(255),
    unit_no VARCHAR(50),         
    building_name VARCHAR(100),  
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (student_record_id) REFERENCES student_records(student_record_id) ON DELETE CASCADE
);


CREATE TABLE educational_backgrounds(
    educational_background_id INT AUTO_INCREMENT PRIMARY KEY,
    student_record_id INT NOT NULL,
	educational_level ENUM(
        'Elementary', 
        'Junior High School', 
        'Senior High School'
    ) NOT NULL,  
    school_name VARCHAR(255) NOT NULL,
    location VARCHAR(255),
    school_type ENUM('Private', 'Public'),     -- ex Private/Public
    year_completed VARCHAR(10),
    awards TEXT,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (student_record_id) REFERENCES student_records(student_record_id) ON DELETE CASCADE
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
    user_id INT, 
    reason VARCHAR(255),
    scheduled_date DATE NOT NULL,
    scheduled_time TIME NOT NULL,
    concern_category VARCHAR(100), -- ex 'Academic', 'Personal', 'Career'
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
    finance_id INT NOT NULL AUTO_INCREMENT,
    student_record_id INT NOT NULL,
    employed_family_members_count INT,
    supports_studies_count INT,
    supports_family_count INT,
    financial_support VARCHAR(255) NOT NULL, 
    weekly_allowance DECIMAL(10,2),
    PRIMARY KEY (finance_id),
    FOREIGN KEY (student_record_id) REFERENCES student_records(student_record_id)
        ON DELETE CASCADE
);

CREATE TABLE student_health_records (
    health_id INT NOT NULL AUTO_INCREMENT,
    student_record_id INT UNIQUE NOT NULL,
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
    FOREIGN KEY (student_record_id) REFERENCES student_records(student_record_id)
        ON DELETE CASCADE
);

CREATE TABLE psychological_assessments (
    assessment_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    student_record_id INT NOT NULL,
    test_date DATE,
    test_name VARCHAR(255),
    raw_score VARCHAR(50), -- (RS)
    remarks TEXT,
    FOREIGN KEY (student_record_id) REFERENCES student_records(student_record_id)
        ON DELETE CASCADE
);