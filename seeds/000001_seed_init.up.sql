INSERT INTO genders (gender_id, gender_name)
VALUES
    (1, 'Male'),
    (2, 'Female'),
    (3, 'Prefer not to say');

INSERT IGNORE INTO roles (role_id, role_name)
VALUES
    (1, 'STUDENT'),
    (2, 'COUNSELOR');

INSERT INTO enrollment_reasons (er_id, reason_text) VALUES
    (1, 'Lower tuition fee'),
    (2, 'Safety of the place'),
    (3, 'Spacious Campus'),
    (4, 'Nearness of home to school'),
    (5, 'Accessible to transportation'),
    (6, 'Better quality of education'),
    (7, 'Adequate School Facilities'),
    (8, 'Son / Daughter of PUP Employee'),
    (9, 'Closer Student-Faculty Relations'),
    (10, 'Expecting Scholarship Offer');

INSERT INTO student_support_types (support_type_name) VALUES
    ('Parents'),
    ('Brother/Sister'),
    ('Spouse'),
    ('Scholarship'),
    ('Relatives'),
    ('Self-supporting/working student');

INSERT INTO income_ranges (ir_id, range_text) VALUES
    (1, 'Below Php 5,000'),
    (2, 'Php 5,001 - Php 10,000'),
    (3, 'Php 10,001 - Php 15,000'),
    (4, 'Php 15,001 - Php 20,000'),
    (5, 'Php 20,001 - Php 30,000'),
    (6, 'Php 30,001 - Php 35,000'),
    (7, 'Php 35,001 - Php 40,000'),
    (8, 'Php 40,001 - Php 45,000'),
    (9, 'Php 45,001 - Php 50,000'),
    (10, 'Above Php 50,001');

-- INSERT INTO appointment_types (appointment_type_id, appointment_type_name)
-- VALUES
--     (1, 'Initial Interview'),
--     (2, 'Mental Health Consultation'),
--     (3, 'Career Guidance'),
--     (4, 'Follow-up');
