CREATE OR REPLACE VIEW simple_student_profiles AS
SELECT
    sr.student_record_id,
    usr.first_name,
    usr.middle_name,
    usr.last_name,       
    usr.email,
    sp.course
FROM student_records sr
JOIN users usr ON sr.user_id = usr.user_id
JOIN student_profiles sp ON sr.student_record_id = sp.student_record_id;

CREATE OR REPLACE VIEW parents_info_view AS
SELECT
    p.parent_id,
    p.educational_level,
    p.birth_date,
    p.last_name,
    p.first_name,
    p.middle_name,
    p.occupation,
    p.company_name,
    sp.relationship
FROM parents p
JOIN student_parents sp ON p.parent_id = sp.parent_id;