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
		req.YearLevel,
		req.GenderID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list students: %w", err)
	}

	// Get total count for pagination (you need to add this method to repository)
	total, err := s.repo.GetTotalStudentsCount(ctx, req.Course, req.YearLevel, req.GenderID)
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
) (int, error) {
	studentRecord, err := s.repo.GetStudentRecordByStudentID(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get student record: %w", err)
	}

	if studentRecord == nil {
		return 0, nil
	}

	return studentRecord.ID, nil
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

// GetGuardiansInfo
func (s *Service) GetGuardiansInfo(
	ctx context.Context, studentRecordID int,
) ([]GuardianInfoView, error) {
	guardianInfo, err := s.repo.GetGuardians(ctx, studentRecordID)
	if err != nil {
		return nil, fmt.Errorf("failed to get guardian info: %w", err)
	}

	return guardianInfo, nil
}

// GetPrimaryGuardianInfo
func (s *Service) GetPrimaryGuardianInfo(
	ctx context.Context, studentRecordID int,
) (*Guardian, error) {
	primaryGuardian, err := s.repo.GetPrimaryGuardian(ctx, studentRecordID)
	if err != nil {
		return nil, fmt.Errorf("failed to get primary guardian: %w", err)
	}

	return primaryGuardian, nil
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
		StudentRecordID:     studentRecordID,
		GenderID:            req.GenderID,
		CivilStatusTypeID:   req.CivilStatusTypeID,
		ReligionTypeID:      req.ReligionTypeID,
		HeightCm:            &req.HeightCm,
		WeightKg:            &req.WeightKg,
		StudentNumber:       req.StudentNumber,
		Course:              req.Course,
		YearLevel:           req.YearLevel,
		Section:             &req.Section,
		GoodMoralStatus:     req.GoodMoralStatus,
		HasDerogatoryRecord: req.HasDerogatoryRecord,
		PlaceOfBirth:        &req.PlaceOfBirth,
		BirthDate:           &req.BirthDate,
		MobileNo:            &req.MobileNo,
	}

	// Save the profile
	_, err := s.repo.SaveStudentProfile(ctx, profile)
	if err != nil {
		return fmt.Errorf("failed to save student profile: %w", err)
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
		MonthlyFamilyIncome:   &req.MonthlyFamilyIncome,
	}

	if err := s.repo.SaveFamilyInfo(ctx, family); err != nil {
		return fmt.Errorf("failed to save family info: %w", err)
	}

	var guardians []Guardian
	var links []StudentGuardian

	for _, g := range req.Guardians {
		guardianModel, linkModel := s.convertGuardianDTOToModel(
			g, g.RelationshipTypeID,
			g.IsPrimary,
		)

		guardians = append(guardians, guardianModel)
		links = append(links, linkModel)
	}
	// Save guardians
	return s.repo.SaveGuardiansInfo(ctx, studentRecordID, guardians, links)
}

func (s *Service) convertGuardianDTOToModel(
	dto GuardianDTO, relationshipTypeID int, isPrimaryContact bool,
) (Guardian, StudentGuardian) {
	guardian := Guardian{
		EducationalLevelID: dto.EducationalLevelID,
		BirthDate:          &dto.BirthDate,
		LastName:           dto.LastName,
		FirstName:          dto.FirstName,
		MiddleName:         &dto.MiddleName,
		Occupation:         &dto.Occupation,
		MaidenName:         &dto.MaidenName,
		CompanyName:        &dto.CompanyName,
		ContactNumber:      &dto.ContactNumber,
	}

	link := StudentGuardian{
		StudentRecordID:    0, // Will be set by repository
		GuardianID:         0, // Will be set by repository
		RelationshipTypeID: relationshipTypeID,
		IsPrimaryContact:   isPrimaryContact,
	}

	return guardian, link
}

func (s *Service) SaveEducationInfo(
	ctx context.Context, studentRecordID int, req UpdateEducationRequest,
) error {
	var educations []EducationalBackground

	for _, e := range req.EducationalBGs {
		educations = append(educations, EducationalBackground{
			StudentRecordID:    studentRecordID,
			EducationalLevelID: e.EducationalLevelID,
			SchoolName:         e.SchoolName,
			Location:           &e.Location,
			SchoolType:         e.SchoolType,
			YearCompleted:      e.YearCompleted,
			Awards:             &e.Awards,
		})
	}

	return s.repo.SaveEducationInfo(ctx, studentRecordID, educations)
}

func (s *Service) SaveAddressInfo(
	ctx context.Context, studentRecordID int, req UpdateAddressRequest,
) error {
	var addresses []StudentAddress

	for _, a := range req.Addresses {
		addresses = append(addresses, StudentAddress{
			StudentRecordID: studentRecordID,
			AddressType:     getAddressTypeFromID(a.AddressTypeID), // Convert ID to enum value
			RegionName:      &a.RegionName,
			ProvinceName:    &a.ProvinceName,
			CityName:        &a.CityName,
			BarangayName:    &a.BarangayName,
			StreetLotBlk:    &a.StreetLotBlk,
			UnitNo:          &a.UnitNo,
			BuildingName:    &a.BuildingName,
		})
	}

	return s.repo.SaveAddressInfo(ctx, studentRecordID, addresses)
}

// Helper function to convert address type ID to enum value
func getAddressTypeFromID(id int) string {
	switch id {
	case 1:
		return "Residential"
	case 2:
		return "Provincial"
	default:
		return "Residential"
	}
}

func (s *Service) SaveHealthRecord(
	ctx context.Context, studentRecordID int, req UpdateHealthRecordRequest,
) error {
	// Convert IDs to enum values
	visionRemark := getHealthRemarkFromID(req.VisionRemarkID)
	hearingRemark := getHealthRemarkFromID(req.HearingRemarkID)
	mobilityRemark := getHealthRemarkFromID(req.MobilityRemarkID)
	speechRemark := getHealthRemarkFromID(req.SpeechRemarkID)
	generalHealthRemark := getHealthRemarkFromID(req.GeneralHealthRemarkID)

	healthRecord := &StudentHealthRecord{
		StudentRecordID:       studentRecordID,
		VisionRemark:          visionRemark,
		HearingRemark:         hearingRemark,
		MobilityRemark:        mobilityRemark,
		SpeechRemark:          speechRemark,
		GeneralHealthRemark:   generalHealthRemark,
		ConsultedProfessional: req.ConsultedProfessional,
		ConsultationReason:    req.ConsultationReason,
		DateStarted:           req.DateStarted,
		NumberOfSessions:      req.NumberOfSessions,
		DateConcluded:         req.DateConcluded,
	}

	return s.repo.SaveHealthRecord(ctx, healthRecord)
}

// Helper function to convert health remark ID to enum value
func getHealthRemarkFromID(id int) string {
	switch id {
	case 1:
		return "No Problem"
	case 2:
		return "Issues"
	default:
		return "No Problem"
	}
}

func (s *Service) SaveFinanceInfo(
	ctx context.Context, studentRecordID int, req UpdateFinanceRequest,
) error {
	// Create finance record
	finance := &StudentFinance{
		StudentRecordID:        studentRecordID,
		IsEmployed:             &req.IsEmployed,
		SupportsStudies:        &req.SupportsStudies,
		SupportsFamily:         &req.SupportsFamily,
		FinancialSupportTypeID: req.FinancialSupportTypeID,
		WeeklyAllowance:        &req.WeeklyAllowance,
	}

	// Save the finance info
	err := s.repo.SaveFinanceInfo(ctx, finance)
	if err != nil {
		return fmt.Errorf("failed to save finance info: %w", err)
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
