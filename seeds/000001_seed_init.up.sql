INSERT INTO genders (gender_id, gender_name)
VALUES
    (1, 'Male'),
    (2, 'Female'),
    (3, 'Prefer not to say');

INSERT IGNORE INTO roles (role_id, role_name)
VALUES
    (1, 'STUDENT'),
    (2, 'COUNSELOR'),
    (3, 'FRONTDESK');

INSERT INTO enrollment_reasons (id, reason_text) VALUES
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

INSERT INTO civil_status_types (civil_status_type_id, status_name)
VALUES
    (1, 'Single'), (2, 'Married'),
    (3, 'Widowed'), (4, 'Divorced');

INSERT INTO parental_status_types (parental_status_id, status_name)
VALUES
    (1, 'Married and Living Together'),
    (2, 'Married but Living Separately'),
    (3, 'Father/Mother working Abroad'),
    (4, 'Divorced or Annulled'),
    (5, 'Separated'),
    (6, 'Other');

-- INSERT INTO appointment_types (appointment_type_id, appointment_type_name)
-- VALUES
--     (1, 'Initial Interview'),
--     (2, 'Mental Health Consultation'),
--     (3, 'Career Guidance'),
--     (4, 'Follow-up');
