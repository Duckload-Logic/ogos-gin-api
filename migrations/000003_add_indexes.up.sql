-- ============================================================================
-- PERFORMANCE INDEXES
-- ============================================================================

-- Foreign Key Indexes (Critical for JOIN performance)
CREATE INDEX idx_users_role_id ON users(role_id);
CREATE INDEX idx_counselor_profiles_user_id ON counselor_profiles(user_id);
CREATE INDEX idx_appointments_user_id ON appointments(user_id);

-- IIR Tables
CREATE INDEX idx_iir_records_user_id ON iir_records(user_id);
CREATE INDEX idx_iir_drafts_user_id ON iir_drafts(user_id);

-- Addresses
CREATE INDEX idx_addresses_region_id ON addresses(region_id);
CREATE INDEX idx_addresses_city_id ON addresses(city_id);
CREATE INDEX idx_addresses_barangay_id ON addresses(barangay_id);
CREATE INDEX idx_emergency_contacts_iir_id ON emergency_contacts(iir_id);
CREATE INDEX idx_emergency_contacts_relationship_id ON emergency_contacts(relationship_id);
CREATE INDEX idx_emergency_contacts_address_id ON emergency_contacts(address_id);

-- Student Personal Info
CREATE INDEX idx_student_personal_info_iir_id ON student_personal_info(iir_id);
CREATE INDEX idx_student_personal_info_gender_id ON student_personal_info(gender_id);
CREATE INDEX idx_student_personal_info_civil_status_id ON student_personal_info(civil_status_id);
CREATE INDEX idx_student_personal_info_religion_id ON student_personal_info(religion_id);
CREATE INDEX idx_student_personal_info_course_id ON student_personal_info(course_id);

-- Student Addresses
CREATE INDEX idx_student_addresses_iir_id ON student_addresses(iir_id);
CREATE INDEX idx_student_addresses_address_id ON student_addresses(address_id);

-- Related Persons
CREATE INDEX idx_student_related_persons_iir_id ON student_related_persons(iir_id);
CREATE INDEX idx_student_related_persons_related_person_id ON student_related_persons(related_person_id);
CREATE INDEX idx_student_related_persons_relationship_id ON student_related_persons(relationship_id);

-- Family Background
CREATE INDEX idx_family_backgrounds_iir_id ON family_backgrounds(iir_id);
CREATE INDEX idx_family_backgrounds_parental_status_id ON family_backgrounds(parental_status_id);
CREATE INDEX idx_family_backgrounds_nature_of_residence_id ON family_backgrounds(nature_of_residence_id);
CREATE INDEX idx_student_sibling_supports_family_background_id ON student_sibling_supports(family_background_id);
CREATE INDEX idx_student_sibling_supports_support_type_id ON student_sibling_supports(support_type_id);

-- Student Selected Reasons
CREATE INDEX idx_student_selected_reasons_iir_id ON student_selected_reasons(iir_id);
CREATE INDEX idx_student_selected_reasons_reason_id ON student_selected_reasons(reason_id);

-- Educational Background
CREATE INDEX idx_educational_backgrounds_iir_id ON educational_backgrounds(iir_id);
CREATE INDEX idx_school_details_eb_id ON school_details(eb_id);
CREATE INDEX idx_school_details_educational_level_id ON school_details(educational_level_id);

-- Health & Wellness
CREATE INDEX idx_student_health_records_iir_id ON student_health_records(iir_id);
CREATE INDEX idx_student_consultations_iir_id ON student_consultations(iir_id);
CREATE INDEX idx_student_activities_iir_id ON student_activities(iir_id);
CREATE INDEX idx_student_activities_option_id ON student_activities(option_id);
CREATE INDEX idx_student_subject_preferences_iir_id ON student_subject_preferences(iir_id);
CREATE INDEX idx_student_hobbies_iir_id ON student_hobbies(iir_id);
CREATE INDEX idx_test_results_iir_id ON test_results(iir_id);
CREATE INDEX idx_significant_notes_iir_id ON significant_notes(iir_id);

-- Financial Support
CREATE INDEX idx_student_finances_iir_id ON student_finances(iir_id);
CREATE INDEX idx_student_finances_income_range_id ON student_finances(monthly_family_income_range_id);
CREATE INDEX idx_student_financial_supports_sf_id ON student_financial_supports(sf_id);
CREATE INDEX idx_student_financial_supports_support_type_id ON student_financial_supports(support_type_id);

-- Administrative
CREATE INDEX idx_admission_slips_iir_id ON admission_slips(iir_id);

-- Commonly Searched Columns
CREATE INDEX idx_appointments_status_id ON appointments(status_id);
CREATE INDEX idx_appointments_when_date ON appointments(when_date);
CREATE INDEX idx_cities_region_id ON cities(region_id);
CREATE INDEX idx_barangays_city_id ON barangays(city_id);
