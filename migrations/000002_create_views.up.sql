CREATE OR REPLACE VIEW simple_student_profiles AS
SELECT
    sr.student_record_id,
    usr.first_name,
    usr.middle_name,
    usr.last_name,       
    usr.email,
    sp.course,
    sp.year_level,
    sp.section
FROM student_records sr
JOIN users usr ON sr.user_id = usr.user_id
JOIN student_profiles sp ON sr.student_record_id = sp.student_record_id;

CREATE OR REPLACE VIEW guardians_info_view AS
SELECT
    g.guardian_id,
    g.educational_level_id,
    g.birth_date,
    g.last_name,
    g.first_name,
    g.middle_name,
    g.occupation,
    g.maiden_name,
    g.company_name,
    g.contact_number,
    sg.relationship_type_id,
    sg.is_primary_contact
FROM guardians g
JOIN student_guardians sg ON g.guardian_id = sg.guardian_id;