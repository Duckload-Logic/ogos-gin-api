SET FOREIGN_KEY_CHECKS = 0;

TRUNCATE TABLE users;
TRUNCATE TABLE counselor_profiles;
TRUNCATE TABLE student_records;
TRUNCATE TABLE family_backgrounds;
TRUNCATE TABLE student_guardians;
TRUNCATE TABLE guardians;
TRUNCATE TABLE educational_backgrounds;
TRUNCATE TABLE student_addresses;
TRUNCATE TABLE student_health_records;
TRUNCATE TABLE appointments;
TRUNCATE TABLE psychological_assessments;
TRUNCATE TABLE session_notes;

SET FOREIGN_KEY_CHECKS = 1;

DROP TABLE fake_migrations