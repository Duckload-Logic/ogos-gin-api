-- ============================================================================
-- AUTH & ROLES
-- ============================================================================

CREATE TABLE roles(
    role_id INT PRIMARY KEY,
    role_name VARCHAR(50) UNIQUE NOT NULL
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

-- ============================================================================
-- REFERENCE DATA (Lookups & Enumerations)
-- ============================================================================

CREATE TABLE genders(
    gender_id INT  PRIMARY KEY,
    gender_name VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE parental_status_types (
    ps_id INT PRIMARY KEY,
    status_name VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE enrollment_reasons (
    er_id INT AUTO_INCREMENT PRIMARY KEY,
    reason_text VARCHAR(100) UNIQUE NOT NULL
);

CREATE TABLE income_ranges (
    ir_id INT AUTO_INCREMENT PRIMARY KEY,
    range_text VARCHAR(100) NOT NULL
);

-- ============================================================================
-- CORE IIR RECORDS
-- ============================================================================

CREATE TABLE iir_records(
    iir_id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT UNIQUE NOT NULL,
    is_submitted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

-- ============================================================================
-- STUDENT PROFILES & ENROLLMENT
-- ============================================================================

CREATE TABLE student_profiles(
    sp_id INT AUTO_INCREMENT PRIMARY KEY,
    iir_id INT UNIQUE NOT NULL,
    student_number VARCHAR(20) UNIQUE NOT NULL,
    gender_id INT,
    civil_status ENUM('Single', 'Married', 'Widowed', 'Separated', 'Divorced') DEFAULT 'Single',
    religion VARCHAR(100) NOT NULL,
    height_ft DECIMAL(5,2) NOT NULL,
    weight_kg DECIMAL(5,2) NOT NULL,
    high_school_gwa DECIMAL(4,2) NOT NULL,
    course VARCHAR(100) NOT NULL,
    year_level VARCHAR(50),
    section VARCHAR(50),
    place_of_birth VARCHAR(255),
    date_of_birth DATE,
    is_employed BOOLEAN DEFAULT FALSE,
    employer_name VARCHAR(255) DEFAULT NULL,
    employer_address VARCHAR(255) DEFAULT NULL,
    contact_no VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (iir_id) REFERENCES iir_records(iir_id) ON DELETE CASCADE,
    FOREIGN KEY (gender_id) REFERENCES genders(gender_id)
);

CREATE TABLE student_selected_reasons (
    iir_id INT NOT NULL,
    reason_id INT NOT NULL,
    other_reason_text VARCHAR(255) DEFAULT NULL,
    PRIMARY KEY (iir_id, reason_id),
    FOREIGN KEY (iir_id) REFERENCES iir_records(iir_id) ON DELETE CASCADE,
    FOREIGN KEY (reason_id) REFERENCES enrollment_reasons(er_id) ON DELETE CASCADE
);

-- ============================================================================
-- LOCATIONS & ADDRESSES
-- ============================================================================

CREATE TABLE addresses(
    address_id INT AUTO_INCREMENT PRIMARY KEY,
    region VARCHAR(100) NOT NULL,
    city VARCHAR(100) NOT NULL,
    barangay VARCHAR(100) NOT NULL,
    street_detail VARCHAR(255)
);

CREATE TABLE student_addresses(
    sa_id INT AUTO_INCREMENT PRIMARY KEY,
    iir_id INT NOT NULL,
    address_id INT NOT NULL,
    address_type ENUM('Residential', 'Provincial') NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (iir_id) REFERENCES iir_records(iir_id) ON DELETE CASCADE,
    FOREIGN KEY (address_id) REFERENCES addresses(address_id) ON DELETE CASCADE
);

-- ============================================================================
-- FAMILY & RELATED PERSONS
-- ============================================================================

CREATE TABLE related_persons(
    rp_id INT AUTO_INCREMENT PRIMARY KEY,
    address_id INT,
    educational_level VARCHAR(100) NOT NULL,
    date_of_birth DATE,
    last_name VARCHAR(100) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    middle_name VARCHAR(100) DEFAULT NULL,
    occupation VARCHAR(100) DEFAULT NULL,
    employer_name VARCHAR(150) DEFAULT NULL,
    employer_address VARCHAR(255) DEFAULT NULL,
    FOREIGN KEY (address_id) REFERENCES addresses(address_id) ON DELETE SET NULL
);

CREATE TABLE student_related_persons (
    iir_id INT NOT NULL,
    related_person_id INT NOT NULL,
    relationship ENUM(
        'Father',
        'Mother',
        'Guardian',
        'Uncle',
        'Aunt',
        'Grandparent',
        'Sibling',
        'Other'
    ),
    is_parent BOOLEAN DEFAULT FALSE,
    is_guardian BOOLEAN DEFAULT FALSE,
    is_living BOOLEAN DEFAULT TRUE,
    is_emergency_contact BOOLEAN DEFAULT FALSE,
    PRIMARY KEY (iir_id, related_person_id),
    FOREIGN KEY (iir_id) REFERENCES iir_records(iir_id) ON DELETE CASCADE,
    FOREIGN KEY (related_person_id) REFERENCES related_persons(rp_id) ON DELETE CASCADE
);

CREATE TABLE family_backgrounds(
    fb_id INT AUTO_INCREMENT PRIMARY KEY,
    iir_id INT UNIQUE NOT NULL,
    parental_status ENUM(
        'Married and staying together',
        'Not Married but living together',
        'Single Parent',
        'Married but Separated',
        'Other'
    ) NOT NULL,
    parental_status_details VARCHAR(255) DEFAULT NULL,
    brothers INT NOT NULL,
    sisters INT NOT NULL,
    employed_siblings INT NOT NULL,
    ordinal_position INT NOT NULL,
    have_quiet_place_to_study BOOLEAN NOT NULL,
    is_sharing_room BOOLEAN NOT NULL,
    room_sharing_details VARCHAR(255) DEFAULT NULL,
    nature_of_residence ENUM(
        'Family home',
        "Relative's house",
        'Bed spacer',
        'House of married brother/sister',
        'Rented apartment/house',
        'Dormitory',
        'Shares apartment with friends/relatives'
    ) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (iir_id) REFERENCES iir_records(iir_id) ON DELETE CASCADE
);

-- ============================================================================
-- EDUCATIONAL BACKGROUND
-- ============================================================================

CREATE TABLE educational_backgrounds(
    eb_id INT AUTO_INCREMENT PRIMARY KEY,
    iir_id INT NOT NULL,
    nature_of_schooling ENUM('Continuous', 'Interrupted') NOT NULL,
    interrupted_details VARCHAR(255) DEFAULT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (iir_id) REFERENCES iir_records(iir_id) ON DELETE CASCADE
);


CREATE TABLE school_details (
	sd_id INT AUTO_INCREMENT PRIMARY KEY,
    eb_id INT NOT NULL,
    educational_level ENUM(
        'Pre-Elementary',
        'Elementary',
        'High School',
        'Vocational',
        'College'
    ) NOT NULL,
    school_name VARCHAR(255) NOT NULL,
    school_address VARCHAR(255) NOT NULL,
    school_type ENUM('Private', 'Public') NOT NULL,
    year_started SMALLINT NOT NULL,
    year_completed SMALLINT NOT NULL,
    awards TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (eb_id) REFERENCES educational_backgrounds(eb_id) ON DELETE CASCADE
);

-- ============================================================================
-- HEALTH & WELLNESS
-- ============================================================================

CREATE TABLE student_health_records (
    health_id INT NOT NULL AUTO_INCREMENT,
    student_record_id INT UNIQUE NOT NULL,
    vision_has_problem BOOLEAN DEFAULT FALSE,
    vision_details VARCHAR(255) DEFAULT NULL,
    hearing_has_problem BOOLEAN DEFAULT FALSE,
    hearing_details VARCHAR(255) DEFAULT NULL,
    speech_has_problem BOOLEAN DEFAULT FALSE,
    speech_details VARCHAR(255) DEFAULT NULL,
    general_health_has_problem BOOLEAN DEFAULT FALSE,
    general_health_details VARCHAR(255) DEFAULT NULL,
    PRIMARY KEY (health_id),
    FOREIGN KEY (student_record_id) REFERENCES iir_records(iir_id) ON DELETE CASCADE
);

CREATE TABLE psychological_consultations (
    consultation_id INT NOT NULL AUTO_INCREMENT,
    student_record_id INT NOT NULL,
    professional_type ENUM('Psychiatrist', 'Psychologist', 'Counselor') NOT NULL,
    has_consulted BOOLEAN DEFAULT FALSE,
    when_date VARCHAR(100) DEFAULT NULL,
    for_what TEXT DEFAULT NULL,
    PRIMARY KEY (consultation_id),
    FOREIGN KEY (student_record_id) REFERENCES iir_records(iir_id) ON DELETE CASCADE
);

CREATE TABLE student_interests (
    interest_id INT AUTO_INCREMENT PRIMARY KEY,
    iir_id INT NOT NULL,
    interest_type VARCHAR(100) NOT NULL,
    interest_name VARCHAR(255) NOT NULL,
    is_favorite BOOLEAN DEFAULT FALSE,
    is_least_favorite BOOLEAN DEFAULT FALSE,
    `rank` INT,
    FOREIGN KEY (iir_id) REFERENCES iir_records(iir_id) ON DELETE CASCADE
);

CREATE TABLE test_results (
    tr_id INT AUTO_INCREMENT PRIMARY KEY,
    iir_id INT NOT NULL,
    test_date DATE,
    test_name VARCHAR(255),
    raw_score VARCHAR(50),
    percentile VARCHAR(50),
    description VARCHAR(255),
    FOREIGN KEY (iir_id) REFERENCES iir_records(iir_id) ON DELETE CASCADE
);

-- ============================================================================
-- FINANCIAL SUPPORT
-- ============================================================================

CREATE TABLE student_support_types (
    sst_id INT AUTO_INCREMENT PRIMARY KEY,
    support_type_name VARCHAR(100) UNIQUE NOT NULL
);

CREATE TABLE student_finances (
    sf_id INT AUTO_INCREMENT PRIMARY KEY,
    iir_id INT UNIQUE NOT NULL,
    monthly_family_income_range_id INT,
    other_income_details VARCHAR(50) DEFAULT NULL,
    financial_support_type_id INT DEFAULT NULL,
    weekly_allowance DECIMAL(10,2),
    FOREIGN KEY (iir_id) REFERENCES iir_records(iir_id) ON DELETE CASCADE,
    FOREIGN KEY (monthly_family_income_range_id) REFERENCES income_ranges(ir_id) ON DELETE SET NULL,
    FOREIGN KEY (financial_support_type_id) REFERENCES student_support_types(sst_id) ON DELETE SET NULL
);


-- ============================================================================
-- COUNSELING & APPOINTMENTS
-- ============================================================================

CREATE TABLE counselor_profiles(
    cp_id INT AUTO_INCREMENT PRIMARY KEY,
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
    `status` ENUM(
        'Pending',
        'Approved',
        'Completed',
        'Cancelled',
        'Rescheduled'
    ) DEFAULT 'Pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE SET NULL
);

CREATE TABLE session_notes (
    sn_id INT AUTO_INCREMENT PRIMARY KEY,
    appointment_id INT UNIQUE NOT NULL,
    notes TEXT,
    recommendation TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (appointment_id) REFERENCES appointments(appointment_id) ON DELETE CASCADE
);

-- ============================================================================
-- ADMINISTRATIVE & RECORDS
-- ============================================================================

CREATE TABLE admission_slips (
    as_id INT AUTO_INCREMENT PRIMARY KEY,
    iir_id INT NOT NULL,
    reason TEXT NOT NULL,
    date_of_absence DATE NOT NULL,
    file_path VARCHAR(255) NOT NULL,
    excuse_slip_status ENUM('Pending', 'Approved', 'Rejected') DEFAULT 'Pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (iir_id) REFERENCES iir_records(iir_id) ON DELETE CASCADE
);