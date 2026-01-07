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

// ========================================
// |                                      |
// |      RETRIEVE SERVICE FUNCTIONS      |
// |                                      |
// ========================================

// ListStudents
func (s *Service) ListStudents(
	ctx context.Context, req ListStudentsRequest,
) (*ListStudentsResponse, error) {
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

// GetStudentRecordByStudentID
func (s *Service) GetStudentRecordByStudentID(
	ctx context.Context, userID int,
) (*StudentRecord, error) {
	studentRecord, err := s.repo.GetStudentRecordByStudentID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student record: %w", err)
	}

	if studentRecord == nil {
		return nil, nil
	}

	return studentRecord, nil
}

// GetStudentEnrollmentReasons
func (s *Service) GetStudentEnrollmentReasons(
	ctx context.Context, studentRecordID int,
) ([]StudentSelectedReason, error) {
	reasons, err := s.repo.GetStudentEnrollmentReasons(ctx, studentRecordID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student enrollment reasons: %w", err)
	}

	return reasons, nil
}

// GetBaseProfile - Combines StudentRecord and StudentProfile
func (s *Service) GetBaseProfile(
	ctx context.Context, studentRecordID int,
) (*StudentProfile, error) {
	studentProfile, err := s.repo.GetStudentProfileByStudentRecordID(
		ctx, studentRecordID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get student profile: %w", err)
	}

	return studentProfile, nil
}

// GetEmergencyContactInfo
func (s *Service) GetEmergencyContactInfo(
	ctx context.Context, studentRecordID int,
) (*StudentEmergencyContact, error) {
	emergencyContact, err := s.repo.GetEmergencyContact(ctx, studentRecordID)
	if err != nil {
		return nil, fmt.Errorf("failed to get emergency contact info: %w", err)
	}

	return emergencyContact, nil
}

// GetFamilyInfo
func (s *Service) GetFamilyInfo(
	ctx context.Context, studentRecordID int,
) (*FamilyBackground, error) {
	familyInfo, err := s.repo.GetFamily(ctx, studentRecordID)
	if err != nil {
		return nil, fmt.Errorf("failed to get family info: %w", err)
	}

	return familyInfo, nil
}

// GetParentsInfo
func (s *Service) GetParentsInfo(
	ctx context.Context, studentRecordID int,
) ([]ParentInfoView, error) {
	ParentInfo, err := s.repo.GetParents(ctx, studentRecordID)
	if err != nil {
		return nil, fmt.Errorf("failed to get Parent info: %w", err)
	}

	return ParentInfo, nil
}

// GetEducationInfo
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

// GetAddressInfo
func (s *Service) GetAddressInfo(
	ctx context.Context, studentRecordID int,
) ([]StudentAddress, error) {
	addressInfo, err := s.repo.GetAddresses(ctx, studentRecordID)
	if err != nil {
		return nil, fmt.Errorf("failed to get address info: %w", err)
	}

	return addressInfo, nil
}

// GetHealthInfo
func (s *Service) GetHealthInfo(
	ctx context.Context, studentRecordID int,
) (*StudentHealthRecord, error) {
	healthInfo, err := s.repo.GetHealthRecord(ctx, studentRecordID)
	if err != nil {
		return nil, fmt.Errorf("failed to get health info: %w", err)
	}

	return healthInfo, nil
}

// GetFinanceInfo
func (s *Service) GetFinanceInfo(
	ctx context.Context, studentRecordID int,
) (*StudentFinance, error) {
	finance, err := s.repo.GetFinance(ctx, studentRecordID)
	if err != nil {
		return nil, fmt.Errorf("failed to get finance info: %w", err)
	}

	return finance, nil
}

// ========================================
// |                                      |
// |       SAVE SERVICE FUNCTIONS         |
// |                                      |
// ========================================

func (s *Service) CreateStudentRecord(
	ctx context.Context, userID int,
) (int, error) {
	// Check if student record already exists
	existingRecord, err := s.repo.GetStudentRecordByStudentID(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to check existing student record: %w", err)
	}

	if existingRecord != nil {
		return existingRecord.ID, nil
	}

	// Create a new student record
	studentRecordID, err := s.repo.CreateStudentRecord(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to create student record: %w", err)
	}

	return studentRecordID, nil
}

// SaveEnrollmentReasons
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

func (s *Service) SaveBaseProfile(
	ctx context.Context, studentRecordID int, req CreateStudentRecordRequest,
) error {
	// Create student profile
	profile := &StudentProfile{
		StudentRecordID:   studentRecordID,
		GenderID:          req.GenderID,
		CivilStatusTypeID: req.CivilStatusTypeID,
		Religion:          req.Religion,
		HeightFt:          &req.HeightFt,
		WeightKg:          &req.WeightKg,
		StudentNumber:     req.StudentNumber,
		Course:            req.Course,
		HighSchoolGWA:     &req.HighSchoolGWA,
		PlaceOfBirth:      &req.PlaceOfBirth,
		BirthDate:         &req.BirthDate,
		ContactNo:         &req.ContactNo,
	}

	// Save the profile
	_, err := s.repo.SaveStudentProfile(ctx, profile)
	if err != nil {
		return fmt.Errorf("failed to save student profile: %w", err)
	}

	return nil
}

// SaveEmergencyContactInfo
func (s *Service) SaveEmergencyContactInfo(
	ctx context.Context, studentRecordID int, req UpdateEmergencyContactRequest,
) error {
	// Create emergency contact
	emergencyContact := &StudentEmergencyContact{
		StudentRecordID:              studentRecordID,
		ParentID:                     req.ParentID,
		EmergencyContactName:         req.EmergencyContactName,
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

// SaveFamilyInfo
func (s *Service) SaveFamilyInfo(
	ctx context.Context, studentRecordID int, req UpdateFamilyRequest,
) error {
	// Save family background
	family := &FamilyBackground{
		StudentRecordID:       studentRecordID,
		ParentalStatusID:      req.ParentalStatusID,
		ParentalStatusDetails: &req.ParentalStatusDetails,
		SiblingsBrothers:      *req.SiblingsBrothers,
		SiblingSisters:        *req.SiblingSisters,
		MonthlyFamilyIncome:   req.MonthlyFamilyIncome,
		GuardianName:          req.GuardianName,
		GuardianAddress:       req.GuardianAddress,
	}

	if err := s.repo.SaveFamilyInfo(ctx, family); err != nil {
		return fmt.Errorf("failed to save family info: %w", err)
	}

	var Parents []Parent
	var links []StudentParent

	for _, g := range req.Parents {
		ParentModel, linkModel := s.convertParentDTOToModel(
			g, g.Relationship,
		)

		Parents = append(Parents, ParentModel)
		links = append(links, linkModel)
	}
	// Save Parents
	return s.repo.SaveParentsInfo(ctx, studentRecordID, Parents, links)
}

func (s *Service) convertParentDTOToModel(
	dto ParentDTO, relationship int,
) (Parent, StudentParent) {
	Parent := Parent{
		EducationalLevel: dto.EducationalLevel,
		BirthDate:        &dto.BirthDate,
		LastName:         dto.LastName,
		FirstName:        dto.FirstName,
		MiddleName:       &dto.MiddleName,
		Occupation:       &dto.Occupation,
		CompanyName:      &dto.CompanyName,
	}

	link := StudentParent{
		StudentRecordID: 0, // Will be set by repository
		ParentID:        0, // Will be set by repository
		Relationship:    relationship,
	}

	return Parent, link
}

func (s *Service) SaveEducationInfo(
	ctx context.Context, studentRecordID int, req UpdateEducationRequest,
) error {
	var educations []EducationalBackground

	for _, e := range req.EducationalBGs {
		educations = append(educations, EducationalBackground{
			StudentRecordID:  studentRecordID,
			EducationalLevel: e.EducationalLevel,
			SchoolName:       e.SchoolName,
			Location:         &e.Location,
			SchoolType:       e.SchoolType,
			YearCompleted:    e.YearCompleted,
			Awards:           &e.Awards,
		})
	}

	return s.repo.SaveEducationInfo(ctx, studentRecordID, educations)
}

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
			StudentRecordID: studentRecordID,
			AddressType:     a.AddressType,
			RegionName:      &region,
			ProvinceName:    &province,
			CityName:        &city,
			BarangayName:    &brgy,
			StreetLotBlk:    &street,
			UnitNo:          &unit,
			BuildingName:    &building,
		})
	}

	return s.repo.SaveAddressInfo(ctx, studentRecordID, addresses)
}

func (s *Service) SaveHealthRecord(
	ctx context.Context, studentRecordID int, req UpdateHealthRecordRequest,
) error {
	fmt.Println(req)
	healthRecord := &StudentHealthRecord{
		StudentRecordID:       studentRecordID,
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

func (s *Service) SaveFinanceInfo(
	ctx context.Context, studentRecordID int, req UpdateFinanceRequest,
) error {
	// Create finance record
	finance := &StudentFinance{
		StudentRecordID:            studentRecordID,
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

func (s *Service) CompleteOnboarding(
	ctx context.Context, studentRecordID int,
) error {
	err := s.repo.MarkOnboardingComplete(ctx, studentRecordID)
	if err != nil {
		return fmt.Errorf("failed to complete onboarding: %w", err)
	}

	return nil
}

// =====================================
// |                                   |
// |      OWNERSHIP VERIFICATIONS      |
// |                                   |
// =====================================

func (s *Service) VerifyStudentRecordOwnership(
	ctx context.Context, userID int, resourceID int,
) (bool, error) {
	studentRecord, err := s.repo.GetStudentRecordByStudentID(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("Failed to get student record: %w", err)
	}

	if studentRecord == nil {
		return false, nil
	}

	return studentRecord.ID == resourceID, nil
}
