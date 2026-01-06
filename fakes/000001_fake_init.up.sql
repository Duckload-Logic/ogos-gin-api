-- ======================================================
-- DUMMY DATA SEEDER FOR GUIDANCE SYSTEM
-- ======================================================
-- This seeder creates:
-- 1. 1 Guidance Counselor
-- 2. 1 Frontdesk Staff
-- 3. 1 Student with complete PDS filled up
-- 4. 1 Student with fresh user creation (no PDS filled up)
-- ======================================================

-- ======================================================
-- 1. CREATE GUIDANCE COUNSELOR
-- ======================================================
INSERT INTO users (
    role_id, first_name, middle_name, last_name, 
    email, password_hash, is_active
) VALUES (
    (SELECT role_id FROM roles WHERE role_name = 'COUNSELOR'),
    'Liwanag',
    NULL,
    'Maliksi',
    'liwanag.maliksi@university.edu',
    '$2y$10$gxeDD.IKlEkqJmqmyVxy6eU9tFvC4ZK8KL3VZc2ex3BvNLo8DL5Dq', -- password
    TRUE
);

SET @counselor_user_id = LAST_INSERT_ID();

-- Create counselor profile
INSERT INTO counselor_profiles (
    user_id, license_number, specialization, is_available
) VALUES (
    @counselor_user_id,
    'PRC-123456',
    'Mental Health and Career Guidance',
    TRUE
);

-- ======================================================
-- 2. CREATE FRONTDESK STAFF
-- ======================================================
INSERT INTO users (
    role_id, first_name, middle_name, last_name, 
    email, password_hash, is_active
) VALUES (
    (SELECT role_id FROM roles WHERE role_name = 'FRONTDESK'),
    'Anna',
    'Marie',
    'Cruz',
    'anna.cruz@university.edu',
    '$2y$10$gxeDD.IKlEkqJmqmyVxy6eU9tFvC4ZK8KL3VZc2ex3BvNLo8DL5Dq', -- password
    TRUE
);

SET @frontdesk_user_id = LAST_INSERT_ID();

-- ======================================================
-- 3. STUDENT WITH COMPLETE PDS FILLED UP
-- ======================================================
-- Create the student user
INSERT INTO users (
    role_id, first_name, middle_name, last_name, 
    email, password_hash, is_active
) VALUES (
    (SELECT role_id FROM roles WHERE role_name = 'STUDENT'),
    'Juan',
    'Santos',
    'Dela Cruz',
    'juan.delacruz@university.edu',
    '$2y$10$gxeDD.IKlEkqJmqmyVxy6eU9tFvC4ZK8KL3VZc2ex3BvNLo8DL5Dq', -- password
    TRUE
);

SET @complete_student_user_id = LAST_INSERT_ID();

-- Create student record
INSERT INTO student_records (
    user_id, civil_status_type_id, religion_type_id, 
    height_cm, weight_kg, student_number, 
    course, year_level, section, 
    good_moral_status, has_derogatory_record,
    gender_id, place_of_birth, birth_date, mobile_no
) VALUES (
    @complete_student_user_id,
    (SELECT civil_status_type_id FROM civil_status_types WHERE status_name = 'Single'),
    (SELECT religion_type_id FROM religion_types WHERE religion_name = 'Roman Catholicism'),
    175.5, -- height_cm
    68.2,  -- weight_kg
    '2023-00123', -- student_number
    'Bachelor of Science in Computer Science',
    3, -- year_level
    'CS-3A',
    TRUE,  -- good_moral_status
    FALSE,  -- has_derogatory_record
    (SELECT gender_id FROM genders WHERE gender_name = 'Male'),
    'Manila',
    '2000-05-15',
    '09171234567'
);

SET @complete_student_record_id = LAST_INSERT_ID();

-- ======================================================
-- CREATE GUARDIANS FOR COMPLETE STUDENT
-- ======================================================
-- Father
INSERT INTO guardians (
    educational_level_id, birth_date, last_name, first_name,
    middle_name, occupation, maiden_name, company_name, contact_number
) VALUES (
    (SELECT educational_level_id FROM educational_levels WHERE level_name = 'College'),
    '1975-03-15',
    'Dela Cruz',
    'Pedro',
    'Santos',
    'Software Engineer',
    NULL,
    'Tech Solutions Inc.',
    '09175554444'
);

SET @father_guardian_id = LAST_INSERT_ID();

-- Mother
INSERT INTO guardians (
    educational_level_id, birth_date, last_name, first_name,
    middle_name, occupation, maiden_name, company_name, contact_number
) VALUES (
    (SELECT educational_level_id FROM educational_levels WHERE level_name = 'College'),
    '1978-07-22',
    'Dela Cruz',
    'Maria',
    'Reyes',
    'Teacher',
    'Garcia',
    'St. Marys High School',
    '09176665555'
);

SET @mother_guardian_id = LAST_INSERT_ID();

-- Link guardians to student
INSERT INTO student_guardians (
    student_record_id, guardian_id, relationship_type_id, is_primary_contact
) VALUES 
    (@complete_student_record_id, @father_guardian_id, 
     (SELECT relationship_type_id FROM relationship_types WHERE relationship_name = 'Father'), 
     TRUE), -- Father is primary contact
    (@complete_student_record_id, @mother_guardian_id, 
     (SELECT relationship_type_id FROM relationship_types WHERE relationship_name = 'Mother'), 
     FALSE);

-- ======================================================
-- CREATE FAMILY BACKGROUND FOR COMPLETE STUDENT
-- ======================================================
INSERT INTO family_backgrounds (
    student_record_id, parental_status_id, parental_status_details,
    siblings_brothers, sibling_sisters, monthly_family_income
) VALUES (
    @complete_student_record_id,
    (SELECT parental_status_id FROM parental_status_types WHERE status_name = 'Married and Living Together'),
    'Both parents are working professionals',
    1, -- siblings_brothers
    2, -- sibling_sisters
    65000.50
);

-- ======================================================
-- CREATE EDUCATIONAL BACKGROUNDS FOR COMPLETE STUDENT
-- ======================================================
-- Elementary
INSERT INTO educational_backgrounds (
    student_record_id, educational_level_id, school_name,
    location, school_type, year_completed, awards
) VALUES (
    @complete_student_record_id,
    (SELECT educational_level_id FROM educational_levels WHERE level_name = 'Elementary'),
    'St. Marys Elementary School',
    'Manila',
    'Private',
    '2014',
    'Valedictorian'
);

-- Junior High School
INSERT INTO educational_backgrounds (
    student_record_id, educational_level_id, school_name,
    location, school_type, year_completed, awards
) VALUES (
    @complete_student_record_id,
    (SELECT educational_level_id FROM educational_levels WHERE level_name = 'Junior High School'),
    'Manila Science High School',
    'Manila',
    'Public',
    '2018',
    'With Honors'
);

-- Senior High School
INSERT INTO educational_backgrounds (
    student_record_id, educational_level_id, school_name,
    location, school_type, year_completed, awards
) VALUES (
    @complete_student_record_id,
    (SELECT educational_level_id FROM educational_levels WHERE level_name = 'Senior High School'),
    'Manila Science High School',
    'Manila',
    'Public',
    '2020',
    'With High Honors'
);

-- ======================================================
-- CREATE ADDRESSES FOR COMPLETE STUDENT
-- ======================================================
-- Provincial Address
INSERT INTO student_addresses (
    student_record_id, address_type_id, region_name,
    province_name, city_name, barangay_name,
    street_lot_blk, unit_no, building_name
) VALUES (
    @complete_student_record_id,
    (SELECT address_type_id FROM address_types WHERE type_name = 'Provincial'),
    'Region IV-A',
    'Laguna',
    'Calamba',
    'Barangay 5',
    '123 Poblacion Road',
    '',
    ''
);

-- Residential Address
INSERT INTO student_addresses (
    student_record_id, address_type_id, region_name,
    province_name, city_name, barangay_name,
    street_lot_blk, unit_no, building_name
) VALUES (
    @complete_student_record_id,
    (SELECT address_type_id FROM address_types WHERE type_name = 'Residential'),
    'NCR',
    'Metro Manila',
    'Quezon City',
    'Batasan Hills',
    '456 Main Street',
    'Unit 5B',
    'Sunrise Towers'
);

-- ======================================================
-- CREATE FINANCE INFO FOR COMPLETE STUDENT
-- ======================================================
INSERT INTO student_finances (
    student_record_id, is_employed, supports_studies,
    supports_family, financial_support_type_id, weekly_allowance
) VALUES (
    @complete_student_record_id,
    0, -- is_employed (false)
    1, -- supports_studies (true, parents support studies)
    0, -- supports_family (false)
    (SELECT financial_support_type_id FROM financial_support_types WHERE support_type_name = 'Parental Support'),
    2000.00 -- weekly_allowance
);

-- ======================================================
-- CREATE HEALTH RECORD FOR COMPLETE STUDENT
-- ======================================================
INSERT INTO student_health_records (
    student_record_id, vision_remark_id, hearing_remark_id,
    mobility_remark_id, speech_remark_id, general_health_remark_id,
    consulted_professional, consultation_reason, date_started,
    num_sessions, date_concluded
) VALUES (
    @complete_student_record_id,
    (SELECT health_remark_type_id FROM health_remark_types WHERE remark_name = 'No problem'),
    (SELECT health_remark_type_id FROM health_remark_types WHERE remark_name = 'No problem'),
    (SELECT health_remark_type_id FROM health_remark_types WHERE remark_name = 'No problem'),
    (SELECT health_remark_type_id FROM health_remark_types WHERE remark_name = 'No problem'),
    (SELECT health_remark_type_id FROM health_remark_types WHERE remark_name = 'No problem'),
    'Dr. Maria Santos',
    'Annual check-up',
    '2024-01-15',
    1,
    '2024-01-15'
);

-- ======================================================
-- CREATE PSYCHOLOGICAL ASSESSMENT FOR COMPLETE STUDENT
-- ======================================================
INSERT INTO psychological_assessments (
    student_record_id, test_date, test_name, raw_score, remarks
) VALUES (
    @complete_student_record_id,
    '2024-02-01',
    'Personality Test',
    '75',
    'Well-adjusted student with good social skills'
);

-- ======================================================
-- CREATE APPOINTMENT FOR COMPLETE STUDENT (Optional)
-- ======================================================
INSERT INTO appointments (
    student_record_id, counselor_user_id, appointment_type_id,
    scheduled_date, scheduled_time, concern_category, status
) VALUES (
    @complete_student_record_id,
    @counselor_user_id,
    (SELECT appointment_type_id FROM appointment_types WHERE appointment_type_name = 'Career Guidance'),
    '2024-03-20',
    '14:00:00',
    'Career',
    'Approved'
);

SET @appointment_id = LAST_INSERT_ID();

-- Create session note for the appointment
INSERT INTO session_notes (
    appointment_id, notes, recommendation
) VALUES (
    @appointment_id,
    'Student is interested in software development career. Discussed various tech career paths.',
    'Recommended to attend coding workshops and consider internship opportunities.'
);

-- ======================================================
-- 4. STUDENT WITH FRESH USER CREATION (NO PDS FILLED UP)
-- ======================================================
INSERT INTO users (
    role_id, first_name, middle_name, last_name, 
    email, password_hash, is_active
) VALUES (
    (SELECT role_id FROM roles WHERE role_name = 'STUDENT'),
    'Maria',
    'Clara',
    'Santos',
    'maria.santos@university.edu',
    '$2y$10$gxeDD.IKlEkqJmqmyVxy6eU9tFvC4ZK8KL3VZc2ex3BvNLo8DL5Dq', -- password
    TRUE
);

-- Note: This student only has a user account, no student_record or PDS data