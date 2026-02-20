CREATE OR REPLACE VIEW simple_student_profiles AS
SELECT
    sr.id AS student_record_id,
    usr.first_name,
    usr.middle_name,
    usr.last_name,
    usr.email,
    sp.course
FROM iir_records sr
JOIN users usr ON sr.user_id = usr.user_id
JOIN student_profiles sp ON sr.id = sp.student_record_id;

CREATE OR REPLACE VIEW parents_info_view AS
SELECT
    rp.id AS related_person_id,
    rp.educational_level,
    rp.birth_date,
    rp.last_name,
    rp.first_name,
    rp.middle_name,
    rp.occupation,
    rp.employer_name,
    srp.relationship
FROM related_persons rp
JOIN student_related_persons srp ON rp.id = srp.related_person_id;