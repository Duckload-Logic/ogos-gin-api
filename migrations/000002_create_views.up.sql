CREATE OR REPLACE VIEW simple_student_profiles AS
SELECT
    iir.iir_id AS student_record_id,
    usr.user_id,
    usr.first_name,
    usr.middle_name,
    usr.last_name,
    usr.email,
    sp.student_number,
    sp.course,
    sp.section,
    sp.year_level
FROM iir_records iir
JOIN users usr ON iir.user_id = usr.user_id
JOIN student_profiles sp ON iir.iir_id = sp.iir_id;

CREATE OR REPLACE VIEW parents_info_view AS
SELECT
    rp.rp_id AS related_person_id,
    rp.educational_level,
    rp.date_of_birth,
    rp.last_name,
    rp.first_name,
    rp.middle_name,
    rp.occupation,
    rp.employer_name,
    srp.relationship
FROM related_persons rp
JOIN student_related_persons srp ON rp.rp_id = srp.related_person_id;