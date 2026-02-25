-- ============================================================================
-- ROLLBACK: PERFORMANCE INDEXES
-- ============================================================================

-- Foreign Key Indexes
DROP INDEX idx_users_role_id ON users;
DROP INDEX idx_counselor_profiles_user_id ON counselor_profiles;
DROP INDEX idx_appointments_user_id ON appointments;

-- IIR Tables
DROP INDEX idx_iir_records_user_id ON iir_records;
DROP INDEX idx_iir_drafts_user_id ON iir_drafts;

-- Addresses
DROP INDEX idx_addresses_region_id ON addresses;
DROP INDEX idx_addresses_city_id ON addresses;
DROP INDEX idx_addresses_barangay_id ON addresses;
DROP INDEX idx_emergency_contacts_iir_id ON emergency_contacts;
DROP INDEX idx_emergency_contacts_relationship_id ON emergency_contacts;
DROP INDEX idx_emergency_contacts_address_id ON emergency_contacts;

-- Student Personal Info
DROP INDEX idx_student_personal_info_iir_id ON student_personal_info;
DROP INDEX idx_student_personal_info_gender_id ON student_personal_info;
DROP INDEX idx_student_personal_info_civil_status_id ON student_personal_info;
DROP INDEX idx_student_personal_info_religion_id ON student_personal_info;
DROP INDEX idx_student_personal_info_course_id ON student_personal_info;

-- Student Addresses
DROP INDEX idx_student_addresses_iir_id ON student_addresses;
DROP INDEX idx_student_addresses_address_id ON student_addresses;

-- Related Persons
DROP INDEX idx_student_related_persons_iir_id ON student_related_persons;
DROP INDEX idx_student_related_persons_related_person_id ON student_related_persons;
DROP INDEX idx_student_related_persons_relationship_id ON student_related_persons;

-- Family Background
DROP INDEX idx_family_backgrounds_iir_id ON family_backgrounds;
DROP INDEX idx_family_backgrounds_parental_status_id ON family_backgrounds;
DROP INDEX idx_family_backgrounds_nature_of_residence_id ON family_backgrounds;
DROP INDEX idx_student_sibling_supports_family_background_id ON student_sibling_supports;
DROP INDEX idx_student_sibling_supports_support_type_id ON student_sibling_supports;

-- Student Selected Reasons
DROP INDEX idx_student_selected_reasons_iir_id ON student_selected_reasons;
DROP INDEX idx_student_selected_reasons_reason_id ON student_selected_reasons;

-- Educational Background
DROP INDEX idx_educational_backgrounds_iir_id ON educational_backgrounds;
DROP INDEX idx_school_details_eb_id ON school_details;
DROP INDEX idx_school_details_educational_level_id ON school_details;

-- Health & Wellness
DROP INDEX idx_student_health_records_iir_id ON student_health_records;
DROP INDEX idx_student_consultations_iir_id ON student_consultations;
DROP INDEX idx_student_activities_iir_id ON student_activities;
DROP INDEX idx_student_activities_option_id ON student_activities;
DROP INDEX idx_student_subject_preferences_iir_id ON student_subject_preferences;
DROP INDEX idx_student_hobbies_iir_id ON student_hobbies;
DROP INDEX idx_test_results_iir_id ON test_results;
DROP INDEX idx_significant_notes_iir_id ON significant_notes;

-- Financial Support
DROP INDEX idx_student_finances_iir_id ON student_finances;
DROP INDEX idx_student_finances_income_range_id ON student_finances;
DROP INDEX idx_student_financial_supports_sf_id ON student_financial_supports;
DROP INDEX idx_student_financial_supports_support_type_id ON student_financial_supports;

-- Administrative
DROP INDEX idx_admission_slips_iir_id ON admission_slips;

-- Commonly Searched Columns
DROP INDEX idx_appointments_status ON appointments;
DROP INDEX idx_appointments_scheduled_date ON appointments;
DROP INDEX idx_cities_region_id ON cities;
DROP INDEX idx_barangays_city_id ON barangays;
