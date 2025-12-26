INSERT INTO genders (gender_name) 
VALUES 
    ('Male'), 
    ('Female'), 
    ('Prefer not to say');

INSERT IGNORE INTO roles (role_name) 
VALUES 
    ('STUDENT'), 
    ('COUNSELOR'), 
    ('FRONTDESK');

INSERT INTO relationship_types (relationship_name) 
VALUES 
    ('Father'), ('Mother'), 
    ('Relative'), ('Legal Guardian');

INSERT INTO address_types (type_name) 
VALUES 
    ('Provincial'), ('Residential');

INSERT INTO educational_levels (level_name) 
VALUES 
    ('Elementary'), ('Junior High School'),
    ('Senior High School'), ('College');

INSERT INTO civil_status_types (status_name) 
VALUES 
    ('Single'), ('Married'), 
    ('Widowed'), ('Divorced');

INSERT INTO religion_types (religion_name) 
VALUES 
    ('Roman Catholicism'), 
    ('Islam'), 
    ('Iglesia ni Cristo'),
    ('Seventh-day Adventist'),
    ('Bible Baptist Church'),
    ('Philippine Independent Church'),
    ('Jehovahs Witnesses'),
    ('Buddhism'),
    ('Other');

INSERT INTO parental_status_types (status_name)
VALUES 
    ('Married and Living Together'), 
    ('Married but Living Separately'), 
    ('Father/Mother working Abroad'),
    ('Divorced or Annulled'), 
    ('Separated'), 
    ('Other');

INSERT INTO financial_support_types (support_type_name) 
VALUES 
    ('Scholarship'), 
    ('Self-funded'),
    ('Sponsored'), 
    ('Parental Support'),
    ('Others');

INSERT INTO health_remark_types (remark_name) 
VALUES 
    ('No problem'), 
    ('Issue');

INSERT INTO appointment_types (appointment_type_name) 
VALUES 
    ('Initial Interview'), 
    ('Mental Health Consultation'), 
    ('Career Guidance'), 
    ('Follow-up');
