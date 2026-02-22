-- ======================================================
-- 1. CREATE GUIDANCE COUNSELOR
-- ======================================================
INSERT INTO users (
    role_id, first_name, middle_name, last_name,
    email, password_hash, is_active
) VALUES (
    2,
    'Liwanag',
    NULL,
    'Maliksi',
    'counselor@university.edu',
    '$2y$10$gxeDD.IKlEkqJmqmyVxy6eU9tFvC4ZK8KL3VZc2ex3BvNLo8DL5Dq',
    TRUE
);

SET @counselor_user_id = LAST_INSERT_ID();

INSERT INTO counselor_profiles (
    user_id, license_number, specialization, is_available
) VALUES (
    @counselor_user_id,
    'PRC-123456',
    'Mental Health and Career Guidance',
    TRUE
);

-- ======================================================
-- 2. STUDENT WITH COMPLETE PDS
-- ======================================================
INSERT INTO users (
    role_id, first_name, middle_name, last_name,
    email, password_hash, is_active
) VALUES (
    1,
    'Juan',
    'Santos',
    'Dela Cruz',
    'student1@university.edu',
    '$2y$10$gxeDD.IKlEkqJmqmyVxy6eU9tFvC4ZK8KL3VZc2ex3BvNLo8DL5Dq',
    TRUE
);

SET @complete_student_user_id = LAST_INSERT_ID();

INSERT INTO iir_records (
    user_id, is_submitted
) VALUES (
    @complete_student_user_id,
    TRUE
);

SET @complete_iir_id = LAST_INSERT_ID();

INSERT INTO student_personal_info (
    iir_id, gender_id, civil_status_id, religion_id,
    height_ft, weight_kg, complexion, student_number, course_id, section, year_level, high_school_gwa,
    place_of_birth, date_of_birth, is_employed, employer_name, employer_address, contact_number
) VALUES (
    @complete_iir_id,
    (SELECT id FROM genders WHERE gender_name = 'Male'),
    (SELECT id FROM civil_status_types WHERE status_name = 'Single'),
    (SELECT id FROM religions WHERE religion_name = 'Roman Catholic'),
    5.80,
    68.20,
    'Moreno',
    '2023-00123-TG-0',
    (SELECT id FROM courses WHERE code = 'BSIT' LIMIT 1),
    '4',
    4,
    92.50,
    'Manila',
    '2000-05-15',
    FALSE,
    NULL,
    NULL,
    '09171234567'
);

-- ======================================================
-- 3. ADDRESSES (STUDENT + RELATED PERSONS)
-- ======================================================
INSERT INTO addresses (region, city, barangay, street_detail)
VALUES
    ('Region IV-A', 'Calamba', 'Barangay 5', '123 Poblacion Road'),
    ('NCR', 'Quezon City', 'Batasan Hills', '456 Main Street, Unit 5B, Sunrise Towers'),
    ('NCR', 'Manila', 'Sampaloc', '12 Parent Street'),
    ('NCR', 'Manila', 'Sampaloc', '14 Parent Street'),
    ('Region IV-A', 'Calamba', 'Barangay 5', '125 Poblacion Road');

SET @provincial_address_id = (
    SELECT id FROM addresses WHERE street_detail = '123 Poblacion Road' ORDER BY id DESC LIMIT 1
);
SET @residential_address_id = (
    SELECT id FROM addresses WHERE street_detail = '456 Main Street, Unit 5B, Sunrise Towers' ORDER BY id DESC LIMIT 1
);
SET @father_address_id = (
    SELECT id FROM addresses WHERE street_detail = '12 Parent Street' ORDER BY id DESC LIMIT 1
);
SET @mother_address_id = (
    SELECT id FROM addresses WHERE street_detail = '14 Parent Street' ORDER BY id DESC LIMIT 1
);
SET @uncle_address_id = (
    SELECT id FROM addresses WHERE street_detail = '125 Poblacion Road' ORDER BY id DESC LIMIT 1
);

INSERT INTO student_addresses (
    iir_id, address_id, address_type
) VALUES
    (@complete_iir_id, @provincial_address_id, 'Provincial'),
    (@complete_iir_id, @residential_address_id, 'Residential');

-- ======================================================
-- 4. RELATED PERSONS (FATHER/MOTHER/EMERGENCY CONTACT)
-- ======================================================
INSERT INTO related_persons (
    address_id, educational_level, date_of_birth, last_name, first_name,
    middle_name, occupation, employer_name, employer_address
) VALUES (
    @father_address_id,
    'College',
    '1975-03-15',
    'Dela Cruz',
    'Pedro',
    'Santos',
    'Software Engineer',
    'Tech Solutions Inc.',
    'BGC, Taguig City'
);

SET @father_related_person_id = LAST_INSERT_ID();

INSERT INTO related_persons (
    address_id, educational_level, date_of_birth, last_name, first_name,
    middle_name, occupation, employer_name, employer_address
) VALUES (
    @mother_address_id,
    'College',
    '1978-07-22',
    'Dela Cruz',
    'Maria',
    'Reyes',
    'Teacher',
    'St. Marys High School',
    'Manila City'
);

SET @mother_related_person_id = LAST_INSERT_ID();

INSERT INTO related_persons (
    address_id, educational_level, date_of_birth, last_name, first_name,
    middle_name, occupation, employer_name, employer_address
) VALUES (
    @uncle_address_id,
    'College',
    '1982-09-10',
    'Santos',
    'Luis',
    NULL,
    'Business Owner',
    'Santos Trading',
    'Calamba, Laguna'
);

SET @uncle_related_person_id = LAST_INSERT_ID();

INSERT INTO student_related_persons (
    iir_id, related_person_id, relationship_id,
    is_parent, is_guardian, is_living, is_emergency_contact
) VALUES
    (@complete_iir_id, @father_related_person_id, (SELECT id FROM student_relationship_types WHERE relationship_name = 'Father' LIMIT 1), TRUE, FALSE, TRUE, FALSE),
    (@complete_iir_id, @mother_related_person_id, (SELECT id FROM student_relationship_types WHERE relationship_name = 'Mother' LIMIT 1), TRUE, FALSE, TRUE, FALSE),
    (@complete_iir_id, @uncle_related_person_id, (SELECT id FROM student_relationship_types WHERE relationship_name = 'Relative' LIMIT 1), FALSE, TRUE, TRUE, TRUE);

-- ======================================================
-- 5. FAMILY BACKGROUND
-- ======================================================
INSERT INTO family_backgrounds (
    iir_id, parental_status_id, parental_status_details,
    brothers, sisters, employed_siblings, ordinal_position,
    have_quiet_place_to_study, is_sharing_room, room_sharing_details,
    nature_of_residence_id
) VALUES (
    @complete_iir_id,
    1,
    '', -- Only have when parental status is "Other"
    1,
    2,
    1,
    2,
    TRUE,
    TRUE,
    'Shares room with one sibling',
    (SELECT id FROM nature_of_residence_types WHERE residence_type_name = 'Family home' LIMIT 1)
);

-- ======================================================
-- 6. EDUCATIONAL BACKGROUND + SCHOOL DETAILS
-- ======================================================
INSERT INTO educational_backgrounds (
    iir_id, nature_of_schooling, interrupted_details
) VALUES (
    @complete_iir_id,
    'Continuous',
    NULL
);

SET @educational_background_id = LAST_INSERT_ID();

INSERT INTO school_details (
    eb_id, educational_level_id, school_name, school_address,
    school_type, year_started, year_completed, awards
) VALUES
    (@educational_background_id, (SELECT id FROM educational_levels WHERE level_name = 'Elementary' LIMIT 1), 'St. Marys Elementary School', 'Manila', 'Private', 2008, 2014, 'Valedictorian'),
    (@educational_background_id, (SELECT id FROM educational_levels WHERE level_name = 'High School' LIMIT 1), 'Manila Science High School', 'Manila', 'Public', 2014, 2020, 'With High Honors'),
    (@educational_background_id, (SELECT id FROM educational_levels WHERE level_name = 'College' LIMIT 1), 'PUP Taguig', 'Taguig City', 'Public', 2023, 2027, NULL);

-- ======================================================
-- 7. FINANCIAL SUPPORT
-- ======================================================
INSERT INTO student_finances (
    iir_id, monthly_family_income_range_id, other_income_details,
    weekly_allowance
) VALUES (
    @complete_iir_id,
    (SELECT id FROM income_ranges WHERE range_text = 'Above Php 50,001' LIMIT 1),
    NULL,
    600.00
);

INSERT INTO student_financial_supports (sf_id, support_type_id) VALUES
    ((SELECT id FROM student_finances WHERE iir_id = @complete_iir_id LIMIT 1), (SELECT id FROM student_support_types WHERE support_type_name = 'Parents' LIMIT 1)),
    ((SELECT id FROM student_finances WHERE iir_id = @complete_iir_id LIMIT 1), (SELECT id FROM student_support_types WHERE support_type_name = 'Scholarship' LIMIT 1));

-- ======================================================
-- 8. HEALTH + CONSULTATION + TEST RESULTS
-- ======================================================
INSERT INTO student_health_records (
    iir_id, vision_has_problem, vision_details,
    hearing_has_problem, hearing_details,
    speech_has_problem, speech_details,
    general_health_has_problem, general_health_details
) VALUES (
    @complete_iir_id,
    FALSE, NULL,
    FALSE, NULL,
    FALSE, NULL,
    FALSE, 'No active medical concern'
);

INSERT INTO student_consultations (
    iir_id, professional_type, has_consulted, when_date, for_what
) VALUES (
    @complete_iir_id,
    'Counselor',
    TRUE,
    '2024-01-15',
    'Annual wellness check'
);

INSERT INTO test_results (
    iir_id, test_date, test_name, raw_score, percentile, description
) VALUES (
    @complete_iir_id,
    '2024-02-01',
    'Personality Test',
    '75',
    '80',
    'Well-adjusted student with good social skills'
);

-- ======================================================
-- 9. APPOINTMENT + SESSION NOTE
-- ======================================================
INSERT INTO appointments (
    user_id, reason, scheduled_date, scheduled_time, concern_category, status
) VALUES (
    @complete_student_user_id,
    'Career Guidance',
    '2026-01-20',
    '14:00:00',
    'Career',
    'Approved'
);

SET @appointment_id = LAST_INSERT_ID();

INSERT INTO significant_notes (
    iir_id, note_date, incident_description, remarks
) VALUES (
    @complete_iir_id,
    '2026-01-20',
    'Career guidance session completed.',
    'Recommended to attend coding workshops and consider internship opportunities.'
);

-- ======================================================
-- 10. STUDENT WITH FRESH USER CREATION (NO PDS)
-- ======================================================
INSERT INTO users (
    role_id, first_name, middle_name, last_name,
    email, password_hash, is_active
) VALUES (
    1,
    'Maria',
    'Clara',
    'Santos',
    'student2@university.edu',
    '$2y$10$gxeDD.IKlEkqJmqmyVxy6eU9tFvC4ZK8KL3VZc2ex3BvNLo8DL5Dq',
    TRUE
);

SET @fresh_student_user_id = LAST_INSERT_ID();

INSERT INTO iir_records (
    user_id
) VALUES (
    @fresh_student_user_id
);

-- Link reasons to complete student
INSERT INTO student_selected_reasons (iir_id, reason_id)
SELECT @complete_iir_id, id
FROM enrollment_reasons
WHERE reason_text IN ('Lower tuition fee', 'Nearness of home to school');

-- ======================================================
-- 11. ADMISSION SLIP FOR TESTING
-- ======================================================
INSERT INTO admission_slips (
    iir_id, reason, date_of_absence, file_path, excuse_slip_status
) VALUES (
    @complete_iir_id,
    'Medical appointment with doctor',
    '2024-03-15',
    '/uploads/excuse_slips/slip_2024001.pdf',
    'Approved'
);