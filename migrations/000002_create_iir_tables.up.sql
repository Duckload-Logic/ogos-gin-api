-- ============================================================================
-- CORE IIR RECORDS
-- ============================================================================

CREATE TABLE iir_records(
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT UNIQUE NOT NULL,
    is_submitted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE iir_drafts(
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT UNIQUE NOT NULL,
    data JSON NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE addresses(
    id INT AUTO_INCREMENT PRIMARY KEY,
    region VARCHAR(100) NOT NULL,
    city VARCHAR(100) NOT NULL,
    barangay VARCHAR(100) NOT NULL,
    street_detail VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- ============================================================================
-- STUDENT PROFILES & ENROLLMENT
-- ============================================================================

CREATE TABLE emergency_contacts (
    id INT AUTO_INCREMENT PRIMARY KEY,
    iir_id INT NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    middle_name VARCHAR(100) DEFAULT NULL,
    last_name VARCHAR(100) NOT NULL,
    contact_number VARCHAR(20) NOT NULL,
    relationship_id INT NOT NULL,
    address_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (iir_id) REFERENCES iir_records(id) ON DELETE CASCADE,
    FOREIGN KEY (relationship_id) REFERENCES student_relationship_types(id),
    FOREIGN KEY (address_id) REFERENCES addresses(id)
);

CREATE TABLE student_personal_info(
    id INT AUTO_INCREMENT PRIMARY KEY,
    iir_id INT UNIQUE NOT NULL,
    student_number VARCHAR(20) UNIQUE NOT NULL,
    gender_id INT NOT NULL,
    civil_status_id INT NOT NULL,
    religion_id INT NOT NULL,
    height_ft DECIMAL(5,2) NOT NULL,
    weight_kg DECIMAL(5,2) NOT NULL,
    complexion VARCHAR(50) NOT NULL,
    high_school_gwa DECIMAL(4,2) NOT NULL,
    course_id INT NOT NULL,
    year_level INT NOT NULL,
    section INT NOT NULL,
    place_of_birth VARCHAR(255) NOT NULL,
    date_of_birth DATE NOT NULL,
    is_employed BOOLEAN DEFAULT FALSE,
    employer_name VARCHAR(255) DEFAULT NULL,
    employer_address VARCHAR(255) DEFAULT NULL,
    mobile_number VARCHAR(20) NOT NULL,
    telephone_number VARCHAR(20) DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (iir_id) REFERENCES iir_records(id) ON DELETE CASCADE,
    FOREIGN KEY (religion_id) REFERENCES religions(id),
    FOREIGN KEY (civil_status_id) REFERENCES civil_status_types(id),
    FOREIGN KEY (course_id) REFERENCES courses(id),
    FOREIGN KEY (gender_id) REFERENCES genders(id)
);

CREATE TABLE student_selected_reasons (
    iir_id INT NOT NULL,
    reason_id INT NOT NULL,
    other_reason_text VARCHAR(255) DEFAULT NULL,
    PRIMARY KEY (iir_id, reason_id),
    FOREIGN KEY (iir_id) REFERENCES iir_records(id) ON DELETE CASCADE,
    FOREIGN KEY (reason_id) REFERENCES enrollment_reasons(id) ON DELETE CASCADE
);

-- ============================================================================
-- LOCATIONS & ADDRESSES
-- ============================================================================

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

-- ============================================================================
-- FAMILY & RELATED PERSONS
-- ============================================================================

CREATE TABLE related_persons(
    id INT AUTO_INCREMENT PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,
    middle_name VARCHAR(100) DEFAULT NULL,
    last_name VARCHAR(100) NOT NULL,
    educational_level VARCHAR(100) NOT NULL,
    date_of_birth DATE NOT NULL,
    occupation VARCHAR(100) DEFAULT NULL,
    employer_name VARCHAR(150) DEFAULT NULL,
    employer_address VARCHAR(255) DEFAULT NULL,
    contact_number VARCHAR(20) DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE student_related_persons (
    iir_id INT NOT NULL,
    related_person_id INT NOT NULL,
    relationship_id INT,
    is_parent BOOLEAN DEFAULT FALSE,
    is_guardian BOOLEAN DEFAULT FALSE,
    is_living BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (iir_id, related_person_id),
    FOREIGN KEY (iir_id) REFERENCES iir_records(id) ON DELETE CASCADE,
    FOREIGN KEY (relationship_id) REFERENCES student_relationship_types(id),
    FOREIGN KEY (related_person_id) REFERENCES related_persons(id) ON DELETE CASCADE
);

CREATE TABLE family_backgrounds(
    id INT AUTO_INCREMENT PRIMARY KEY,
    iir_id INT UNIQUE NOT NULL,
    parental_status_id INT NOT NULL,
    parental_status_details VARCHAR(255) DEFAULT NULL,
    brothers INT NOT NULL,
    sisters INT NOT NULL,
    employed_siblings INT NOT NULL,
    ordinal_position INT NOT NULL,
    have_quiet_place_to_study BOOLEAN NOT NULL,
    is_sharing_room BOOLEAN NOT NULL,
    room_sharing_details VARCHAR(255) DEFAULT NULL,
    nature_of_residence_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (iir_id) REFERENCES iir_records(id) ON DELETE CASCADE,
    FOREIGN KEY (nature_of_residence_id) REFERENCES nature_of_residence_types(id),
    FOREIGN KEY (parental_status_id) REFERENCES parental_status_types(id)
);

CREATE TABLE student_sibling_supports (
    family_background_id INT NOT NULL,
    support_type_id INT NOT NULL,
    PRIMARY KEY (family_background_id, support_type_id),
    FOREIGN KEY (family_background_id) REFERENCES family_backgrounds(id) ON DELETE CASCADE,
    FOREIGN KEY (support_type_id) REFERENCES sibling_support_types(id) ON DELETE CASCADE
);

-- ============================================================================
-- EDUCATIONAL BACKGROUND
-- ============================================================================

CREATE TABLE educational_backgrounds(
    id INT AUTO_INCREMENT PRIMARY KEY,
    iir_id INT NOT NULL,
    nature_of_schooling ENUM('Continuous', 'Interrupted') NOT NULL,
    interrupted_details VARCHAR(255) DEFAULT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (iir_id) REFERENCES iir_records(id) ON DELETE CASCADE
);


CREATE TABLE school_details (
	id INT AUTO_INCREMENT PRIMARY KEY,
    eb_id INT NOT NULL,
    educational_level_id INT NOT NULL,
    school_name VARCHAR(255) NOT NULL,
    school_address VARCHAR(255) NOT NULL,
    school_type ENUM('Private', 'Public') NOT NULL,
    year_started SMALLINT NOT NULL,
    year_completed SMALLINT NOT NULL,
    awards TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (eb_id) REFERENCES educational_backgrounds(id) ON DELETE CASCADE,
    FOREIGN KEY (educational_level_id) REFERENCES educational_levels(id)
);

-- ============================================================================
-- HEALTH & WELLNESS
-- ============================================================================

CREATE TABLE student_health_records (
    id INT AUTO_INCREMENT PRIMARY KEY,
    iir_id INT UNIQUE NOT NULL,
    vision_has_problem BOOLEAN DEFAULT FALSE,
    vision_details VARCHAR(255) DEFAULT NULL,
    hearing_has_problem BOOLEAN DEFAULT FALSE,
    hearing_details VARCHAR(255) DEFAULT NULL,
    speech_has_problem BOOLEAN DEFAULT FALSE,
    speech_details VARCHAR(255) DEFAULT NULL,
    general_health_has_problem BOOLEAN DEFAULT FALSE,
    general_health_details VARCHAR(255) DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (iir_id) REFERENCES iir_records(id) ON DELETE CASCADE
);

CREATE TABLE student_consultations (
    id INT AUTO_INCREMENT PRIMARY KEY,
    iir_id INT NOT NULL,
    professional_type ENUM('Psychiatrist', 'Psychologist', 'Counselor') NOT NULL,
    has_consulted BOOLEAN DEFAULT FALSE,
    when_date VARCHAR(100) DEFAULT NULL,
    for_what TEXT DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (iir_id) REFERENCES iir_records(id) ON DELETE CASCADE
);

CREATE TABLE student_activities (
    id INT AUTO_INCREMENT PRIMARY KEY,
    iir_id INT NOT NULL,
    option_id INT,
    other_specification VARCHAR(255),
    role ENUM('Officer', 'Member', 'Other') DEFAULT 'Member',
    role_specification VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (iir_id) REFERENCES iir_records(id) ON DELETE CASCADE,
    FOREIGN KEY (option_id) REFERENCES activity_options(id)
);

CREATE TABLE student_subject_preferences (
    id INT AUTO_INCREMENT PRIMARY KEY,
    iir_id INT NOT NULL,
    subject_name VARCHAR(100) NOT NULL,
    is_favorite BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (iir_id) REFERENCES iir_records(id) ON DELETE CASCADE,
    UNIQUE KEY (iir_id, subject_name)
);

CREATE TABLE student_hobbies (
    id INT AUTO_INCREMENT PRIMARY KEY,
    iir_id INT NOT NULL,
    hobby_name VARCHAR(255) NOT NULL,
    priority_rank INT CHECK (priority_rank BETWEEN 1 AND 4),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (iir_id) REFERENCES iir_records(id) ON DELETE CASCADE
);

CREATE TABLE test_results (
    id INT AUTO_INCREMENT PRIMARY KEY,
    iir_id INT NOT NULL,
    test_date DATE,
    test_name VARCHAR(255),
    raw_score VARCHAR(50),
    percentile VARCHAR(50),
    description VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (iir_id) REFERENCES iir_records(id) ON DELETE CASCADE
);

CREATE TABLE significant_notes (
    id INT AUTO_INCREMENT PRIMARY KEY,
    iir_id INT NOT NULL,
    note_date DATE,
    incident_description TEXT,
    remarks TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (iir_id) REFERENCES iir_records(id) ON DELETE CASCADE
);

-- ============================================================================
-- FINANCIAL SUPPORT
-- ============================================================================

CREATE TABLE student_finances (
    id INT AUTO_INCREMENT PRIMARY KEY,
    iir_id INT UNIQUE NOT NULL,
    monthly_family_income_range_id INT,
    other_income_details VARCHAR(50) DEFAULT NULL,
    weekly_allowance DECIMAL(10,2) DEFAULT 0.00,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (iir_id) REFERENCES iir_records(id) ON DELETE CASCADE,
    FOREIGN KEY (monthly_family_income_range_id) REFERENCES income_ranges(id)
);

CREATE TABLE student_financial_supports (
    sf_id INT NOT NULL,
    support_type_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (support_type_id) REFERENCES student_support_types(id),
    FOREIGN KEY (sf_id) REFERENCES student_finances(id) ON DELETE CASCADE
);
-- ============================================================================
-- ADMINISTRATIVE & RECORDS
-- ============================================================================

CREATE TABLE admission_slips (
    id INT AUTO_INCREMENT PRIMARY KEY,
    iir_id INT NOT NULL,
    reason TEXT NOT NULL,
    date_of_absence DATE NOT NULL,
    file_path VARCHAR(255) NOT NULL,
    excuse_slip_status ENUM('Pending', 'Approved', 'Rejected') DEFAULT 'Pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (iir_id) REFERENCES iir_records(id) ON DELETE CASCADE
);