package students

import (
	"context"
	"database/sql"
	"fmt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// Retrieve - List
func (s *Service) ListStudents(
	ctx context.Context, req ListStudentsRequest,
) (*ListStudentsResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}

	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	students, err := s.repo.ListStudents(
		ctx,
		req.GetOffset(),
		req.PageSize,
		req.Course,
		req.GenderID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list students: %w", err)
	}

	// Get total count for pagination (you need to add this method to repository)
	total, err := s.repo.GetTotalStudentsCount(ctx, req.Course, req.GenderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get total students count: %w", err)
	}

	totalPages := (total + req.PageSize - 1) / req.PageSize
	if totalPages == 0 {
		totalPages = 1
	}

	return &ListStudentsResponse{
		Students:   students,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// Retrieve - Student Records
func (s *Service) GetInventoryRecordByStudentID(
	ctx context.Context, userID int,
) (*InventoryRecord, error) {
	studentRecord, err := s.repo.GetInventoryRecordByStudentID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student record: %w", err)
	}

	if studentRecord == nil {
		return nil, nil
	}

	return studentRecord, nil
}

// Retrieve - Enrollment Reasons
func (s *Service) GetStudentEnrollmentReasons(
	ctx context.Context, studentRecordID int,
) ([]StudentSelectedReason, error) {
	reasons, err := s.repo.GetStudentEnrollmentReasons(ctx, studentRecordID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student enrollment reasons: %w", err)
	}

	return reasons, nil
}

// Retrieve - Base Profile
func (s *Service) GetBaseProfile(
	ctx context.Context, studentRecordID int,
) (*StudentProfile, error) {
	studentProfile, err := s.repo.GetStudentProfileByInventoryRecordID(
		ctx, studentRecordID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get student profile: %w", err)
	}

	return studentProfile, nil
}

// Retrieve - Emergency Contact
func (s *Service) GetEmergencyContactInfo(
	ctx context.Context, studentRecordID int,
) (*StudentEmergencyContact, error) {
	emergencyContact, err := s.repo.GetEmergencyContact(ctx, studentRecordID)
	if err != nil {
		return nil, fmt.Errorf("failed to get emergency contact info: %w", err)
	}

	return emergencyContact, nil
}

// Retrieve - Family Info
func (s *Service) GetFamilyInfo(
	ctx context.Context, studentRecordID int,
) (*FamilyBackground, error) {
	familyInfo, err := s.repo.GetFamily(ctx, studentRecordID)
	if err != nil {
		return nil, fmt.Errorf("failed to get family info: %w", err)
	}

	return familyInfo, nil
}

// Retrieve - Related Persons Info
func (s *Service) GetRelatedPersonsInfo(
	ctx context.Context, studentRecordID int,
) ([]RelatedPersonInfoView, error) {
	relatedPersonsInfo, err := s.repo.GetRelatedPersons(ctx, studentRecordID)
	if err != nil {
		return nil, fmt.Errorf("failed to get related person info: %w", err)
	}

	return relatedPersonsInfo, nil
}

// Retrieve - Education Info
func (s *Service) GetEducationInfo(
	ctx context.Context, studentRecordID int,
) ([]EducationalBackground, error) {
	educationInfo, err := s.repo.GetEducationalBackgrounds(
		ctx, studentRecordID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get education info: %w", err)
	}

	return educationInfo, nil
}

// Retrieve - Address Info
func (s *Service) GetAddressInfo(
	ctx context.Context, studentRecordID int,
) ([]StudentAddress, error) {
	addressInfo, err := s.repo.GetAddresses(ctx, studentRecordID)
	if err != nil {
		return nil, fmt.Errorf("failed to get address info: %w", err)
	}

	return addressInfo, nil
}

// Retrieve - Health Info
func (s *Service) GetHealthInfo(
	ctx context.Context, studentRecordID int,
) (*StudentHealthRecord, error) {
	healthInfo, err := s.repo.GetHealthRecord(ctx, studentRecordID)
	if err != nil {
		return nil, fmt.Errorf("failed to get health info: %w", err)
	}

	return healthInfo, nil
}

// Retrieve - Finance Info
func (s *Service) GetFinanceInfo(
	ctx context.Context, studentRecordID int,
) (*StudentFinance, error) {
	finance, err := s.repo.GetFinance(ctx, studentRecordID)
	if err != nil {
		return nil, fmt.Errorf("failed to get finance info: %w", err)
	}

	return finance, nil
}

// Create - Student Record
func (s *Service) CreateInventoryRecord(
	ctx context.Context, userID int,
) (int, error) {
	// Check if student record already exists
	existingRecord, err := s.repo.GetInventoryRecordByStudentID(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to check existing student record: %w", err)
	}

	if existingRecord != nil {
		return existingRecord.ID, nil
	}

	// Create a new student record
	studentRecordID, err := s.repo.CreateInventoryRecord(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to create student record: %w", err)
	}

	return studentRecordID, nil
}

// Create/Update - Enrollment Reasons
func (s *Service) SaveEnrollmentReasons(
	ctx context.Context, studentRecordID int, req UpdateEnrollmentReasonsRequest,
) error {
	// Delete existing enrollment reasons for this student
	if err := s.repo.DeleteEnrollmentReasons(ctx, studentRecordID); err != nil {
		return fmt.Errorf("failed to delete existing enrollment reasons: %w", err)
	}

	// Insert new enrollment reasons
	if len(req.EnrollmentReasonIDs) > 0 {
		for _, reasonID := range req.EnrollmentReasonIDs {
			var otherReasonText sql.NullString
			if reasonID == 11 && req.OtherReasonText != "" {
				otherReasonText = sql.NullString{
					String: req.OtherReasonText, Valid: true,
				}
			}

			if err := s.repo.SaveEnrollmentReason(
				ctx, studentRecordID, reasonID, otherReasonText,
			); err != nil {
				return fmt.Errorf("failed to save enrollment reason: %w", err)
			}
		}
	}

	return nil
}

// Create/Update - Base Profile
func (s *Service) SaveBaseProfile(
	ctx context.Context, studentRecordID int, req CreateInventoryRecordRequest,
) error {
	// Create student profile
	profile := &StudentProfile{
		InventoryRecordID: studentRecordID,
		GenderID:          req.GenderID,
		CivilStatusTypeID: req.CivilStatusTypeID,
		Religion:          req.Religion,
		HeightFt:          &req.HeightFt,
		WeightKg:          &req.WeightKg,
		StudentNumber:     req.StudentNumber,
		Course:            req.Course,
		HighSchoolGWA:     &req.HighSchoolGWA,
		PlaceOfBirth:      &req.PlaceOfBirth,
		DateOfBirth:       &req.DateOfBirth,
		ContactNumber:     &req.ContactNumber,
	}

	// Save the profile
	_, err := s.repo.SaveStudentProfile(ctx, profile)
	if err != nil {
		return fmt.Errorf("failed to save student profile: %w", err)
	}

	return nil
}

// Create/Update - Emergency Contact
func (s *Service) SaveEmergencyContactInfo(
	ctx context.Context, studentRecordID int, req UpdateEmergencyContactRequest,
) error {
	// Create emergency contact
	emergencyContact := &StudentEmergencyContact{
		InventoryRecordID:            studentRecordID,
		ParentID:                     req.RelatedPersonID,
		EmergencyContactFirstName:    req.EmergencyContactFirstName,
		EmergencyContactMiddleName:   req.EmergencyContactMiddleName,
		EmergencyContactLastName:     req.EmergencyContactLastName,
		EmergencyContactPhone:        req.EmergencyContactPhone,
		EmergencyContactRelationship: req.EmergencyContactRelationship,
	}

	// Save the emergency contact
	err := s.repo.SaveEmergencyContact(ctx, emergencyContact)
	if err != nil {
		return fmt.Errorf("failed to save emergency contact: %w", err)
	}

	return nil
}

// Create/Update - Family Info
func (s *Service) SaveFamilyInfo(
	ctx context.Context, studentRecordID int, req UpdateFamilyRequest,
) error {
	// Save family background
	family := &FamilyBackground{
		InventoryRecordID:     studentRecordID,
		ParentalStatusID:      req.ParentalStatusID,
		ParentalStatusDetails: &req.ParentalStatusDetails,
		Brothers:              *req.Brothers,
		Sisters:               *req.Sisters,
		MonthlyFamilyIncome:   req.MonthlyFamilyIncome,
		GuardianFirstName:     req.GuardianFirstName,
		GuardianLastName:      req.GuardianLastName,
		GuardianMiddleName:    req.GuardianMiddleName,
		GuardianAddress:       req.GuardianAddress,
	}

	if err := s.repo.SaveFamilyInfo(ctx, family); err != nil {
		return fmt.Errorf("failed to save family info: %w", err)
	}

	var relatedPersons []RelatedPerson
	var links []StudentRelatedPerson

	for _, g := range req.RelatedPersons {
		relatedPersonModel, linkModel := s.convertRelatedPersonDTOToModel(
			g, g.Relationship,
		)

		relatedPersons = append(relatedPersons, relatedPersonModel)
		links = append(links, linkModel)
	}
	// Save Related Persons
	return s.repo.SaveRelatedPersonsInfo(ctx, studentRecordID, relatedPersons, links)
}

// Helper - Convert Related Person DTO to Related Person
func (s *Service) convertRelatedPersonDTOToModel(
	dto RelatedPersonDTO, relationship int,
) (RelatedPerson, StudentRelatedPerson) {
	relatedPerson := RelatedPerson{
		AddressID:        0,
		EducationalLevel: dto.EducationalLevel,
		DateOfBirth:      &dto.DateOfBirth,
		LastName:         dto.LastName,
		FirstName:        dto.FirstName,
		MiddleName:       &dto.MiddleName,
		Occupation:       &dto.Occupation,
		EmployerName:     &dto.CompanyName,
		EmployerAddress:  nil,
		IsLiving:         true,
	}

	link := StudentRelatedPerson{
		InventoryRecordID:  0, // Will be set by repository
		PersonID:           0, // Will be set by repository
		Relationship:       relationship,
		IsParent:           true,
		IsGuardian:         false,
		IsEmergencyContact: false,
	}

	return relatedPerson, link
}

// Create/Update - Education Info
func (s *Service) SaveEducationInfo(
	ctx context.Context, studentRecordID int, req UpdateEducationRequest,
) error {
	var educations []EducationalBackground

	for _, e := range req.EducationalBGs {
		educations = append(educations, EducationalBackground{
			InventoryRecordID: studentRecordID,
			EducationalLevel:  e.EducationalLevel,
			SchoolName:        e.SchoolName,
			Location:          &e.Location,
			SchoolType:        e.SchoolType,
			YearCompleted:     e.YearCompleted,
			Awards:            &e.Awards,
		})
	}

	return s.repo.SaveEducationInfo(ctx, studentRecordID, educations)
}

// Create/Update - Address Info
func (s *Service) SaveAddressInfo(
	ctx context.Context, studentRecordID int, input interface{},
) error {
	var dtoList []StudentAddressDTO

	// The Type Switch (Your "OR" logic)
	switch v := input.(type) {
	case UpdateAddressRequest:
		dtoList = v.Addresses
	case []StudentAddressDTO:
		dtoList = v
	default:
		return fmt.Errorf("invalid input type for SaveAddressInfo")
	}

	var addresses []StudentAddress

	for _, a := range dtoList {
		region := a.RegionName
		province := a.ProvinceName
		city := a.CityName
		brgy := a.BarangayName
		street := a.StreetLotBlk
		unit := a.UnitNo
		building := a.BuildingName

		addresses = append(addresses, StudentAddress{
			InventoryRecordID: studentRecordID,
			AddressType:       a.AddressType,
			RegionName:        &region,
			ProvinceName:      &province,
			CityName:          &city,
			BarangayName:      &brgy,
			StreetLotBlk:      &street,
			UnitNo:            &unit,
			BuildingName:      &building,
		})
	}

	return s.repo.SaveAddressInfo(ctx, studentRecordID, addresses)
}

// Create/Update - Health Record
func (s *Service) SaveHealthRecord(
	ctx context.Context, studentRecordID int, req UpdateHealthRecordRequest,
) error {
	fmt.Println(req)
	healthRecord := &StudentHealthRecord{
		InventoryRecordID:     studentRecordID,
		VisionRemark:          req.VisionRemark,
		HearingRemark:         req.HearingRemark,
		MobilityRemark:        req.MobilityRemark,
		SpeechRemark:          req.SpeechRemark,
		GeneralHealthRemark:   req.GeneralHealthRemark,
		ConsultedProfessional: req.ConsultedProfessional,
		ConsultationReason:    req.ConsultationReason,
		DateStarted:           req.DateStarted,
		NumberOfSessions:      req.NumberOfSessions,
		DateConcluded:         req.DateConcluded,
	}

	return s.repo.SaveHealthRecord(ctx, healthRecord)
}

// Create/Update - Finance Info
func (s *Service) SaveFinanceInfo(
	ctx context.Context, studentRecordID int, req UpdateFinanceRequest,
) error {
	// Create finance record
	finance := &StudentFinance{
		InventoryRecordID:          studentRecordID,
		EmployedFamilyMembersCount: &req.EmployedFamilyMembersCount,
		SupportsStudiesCount:       &req.SupportsStudiesCount,
		SupportsFamilyCount:        &req.SupportsFamilyCount,
		FinancialSupport:           req.FinancialSupport,
		WeeklyAllowance:            &req.WeeklyAllowance,
	}

	// Save the finance info
	err := s.repo.SaveFinanceInfo(ctx, finance)
	if err != nil {
		return fmt.Errorf("failed to save finance info: %w", err)
	}

	return nil
}

// Update - Complete Onboarding
func (s *Service) CompleteOnboarding(
	ctx context.Context, studentRecordID int,
) error {
	err := s.repo.MarkOnboardingComplete(ctx, studentRecordID)
	if err != nil {
		return fmt.Errorf("failed to complete onboarding: %w", err)
	}

	return nil
}

// Verify - Student Record Ownership
func (s *Service) VerifyInventoryRecordOwnership(
	ctx context.Context, userID int, resourceID int,
) (bool, error) {
	studentRecord, err := s.repo.GetInventoryRecordByStudentID(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("Failed to get student record: %w", err)
	}

	if studentRecord == nil {
		return false, nil
	}

	return studentRecord.ID == resourceID, nil
}

// Delete - Student Record
func (s *Service) DeleteInventoryRecord(
	ctx context.Context, studentRecordID int,
) error {
	err := s.repo.DeleteInventoryRecord(ctx, studentRecordID)
	if err != nil {
		return fmt.Errorf("failed to delete student record: %w", err)
	}

	return nil
}
