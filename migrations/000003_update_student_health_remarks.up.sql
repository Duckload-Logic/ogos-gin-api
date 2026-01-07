ALTER TABLE student_health_records
    MODIFY COLUMN vision_remark ENUM('No Problem', 'Issue'),
    MODIFY COLUMN hearing_remark ENUM('No Problem', 'Issue'),
    MODIFY COLUMN mobility_remark ENUM('No Problem', 'Issue'),
    MODIFY COLUMN speech_remark ENUM('No Problem', 'Issue'),
    MODIFY COLUMN general_health_remark ENUM('No Problem', 'Issue');