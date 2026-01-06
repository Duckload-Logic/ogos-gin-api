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

INSERT INTO relationship_types (relationship_type_id, relationship_name) 
VALUES 
    (1, 'Father'), (2, 'Mother'), 
    (3, 'Relative'), (4, 'Legal Guardian');

INSERT INTO address_types (address_type_id, type_name) 
VALUES 
    (1, 'Provincial'), (2, 'Residential');
INSERT INTO educational_levels (educational_level_id, level_name) 
VALUES 
    (1, 'Elementary'), (2, 'Junior High School'),
    (3, 'Senior High School'), (4, 'College');

INSERT INTO civil_status_types (civil_status_type_id, status_name) 
VALUES 
    (1, 'Single'), (2, 'Married'), 
    (3, 'Widowed'), (4, 'Divorced');

INSERT INTO religion_types (religion_type_id, religion_name) 
VALUES 
    (1, 'Roman Catholicism'), 
    (2, 'Islam'), 
    (3, 'Iglesia ni Cristo'),
    (4, 'Seventh-day Adventist'),
    (5, 'Bible Baptist Church'),
    (6, 'Philippine Independent Church'),
    (7, 'Jehovahs Witnesses'),
    (8, 'Buddhism'),
    (9, 'Other');

INSERT INTO parental_status_types (parental_status_id, status_name)
VALUES 
    (1, 'Married and Living Together'), 
    (2, 'Married but Living Separately'), 
    (3, 'Father/Mother working Abroad'),
    (4, 'Divorced or Annulled'), 
    (5, 'Separated'), 
    (6, 'Other');
INSERT INTO financial_support_types (financial_support_type_id, support_type_name) 
VALUES 
    (1, 'Scholarship'), 
    (2, 'Self-funded'),
    (3, 'Sponsored'), 
    (4, 'Parental Support'),
    (5, 'Others');

INSERT INTO health_remark_types (health_remark_type_id, remark_name) 
VALUES 
    (1, 'No problem'), 
    (2, 'Issue');

INSERT INTO appointment_types (appointment_type_id, appointment_type_name) 
VALUES 
    (1, 'Initial Interview'), 
    (2, 'Mental Health Consultation'), 
    (3, 'Career Guidance'), 
    (4, 'Follow-up');
