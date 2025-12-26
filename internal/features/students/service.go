package students

import (
	"context"
	"database/sql"
	"fmt"
	"time"
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

// GetBaseProfile
func (s *Service) GetBaseProfile(
	ctx context.Context, userID int,
) (*StudentRecord, error) {
	studentRecordInfo, err := s.repo.GetStudentRecordByStudentID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student record: %w", err)
	}

	return studentRecordInfo, nil
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
) ([]Guardian, error) {
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

// ========================================
// |                                      |
// |       SAVE SERVICE FUNCTIONS         |
// |                                      |
// ========================================

// SaveBaseProfile
func (s *Service) SaveBaseProfile(
	ctx context.Context, req CreateStudentRecordRequest,
) (int, error) {
	record := &StudentRecord{
		UserID:              req.UserID,
		CivilStatusTypeID:   req.CivilStatusTypeID,
		ReligionTypeID:      req.ReligionTypeID,
		HeightCm:            sql.NullFloat64{Float64: req.HeightCm, Valid: true},
		WeightKg:            sql.NullFloat64{Float64: req.WeightKg, Valid: true},
		StudentNumber:       req.StudentNumber,
		Course:              req.Course,
		YearLevel:           req.YearLevel,
		Section:             sql.NullString{String: req.Section, Valid: req.Section != ""},
		GoodMoralStatus:     req.GoodMoralStatus,
		HasDerogatoryRecord: req.HasDerogatoryRecord,
	}

	return s.repo.SaveBaseProfileInfo(ctx, record)
}

// SaveFamilyInfo
func (s *Service) SaveFamilyInfo(
	ctx context.Context, studentRecordID int, req UpdateFamilyRequest,
) error {
	// Save family background
	family := &FamilyBackground{
		StudentRecordID:       studentRecordID,
		ParentalStatusID:      req.ParentalStatusID,
		ParentalStatusDetails: sql.NullString{String: req.ParentalStatusDetails, Valid: req.ParentalStatusDetails != ""},
		SiblingsBrothers:      req.SiblingsBrothers,
		SiblingSisters:        req.SiblingSisters,
		MonthlyFamilyIncome:   sql.NullFloat64{Float64: req.MonthlyFamilyIncome, Valid: true},
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
	// Parse birth date
	var birthDate sql.NullTime
	if dto.BirthDate != "" {
		// "2006-01-02" is the reference layout for parsing dates
		t, err := time.Parse("2006-01-02", dto.BirthDate)
		if err == nil {
			birthDate = sql.NullTime{Time: t, Valid: true}
		}
	}

	// Create contact number
	var contactNumber sql.NullString
	if dto.ContactNumber != "" {
		contactNumber = sql.NullString{String: dto.ContactNumber, Valid: true}
	}

	guardian := Guardian{
		EducationalLevelID: dto.EducationalLevelID,
		BirthDate:          birthDate,
		LastName:           dto.LastName,
		FirstName:          dto.FirstName,
		MiddleName:         sql.NullString{String: dto.MiddleName, Valid: dto.MiddleName != ""},
		Occupation:         sql.NullString{String: dto.Occupation, Valid: dto.Occupation != ""},
		MaidenName:         sql.NullString{String: dto.MaidenName, Valid: dto.MaidenName != ""},
		CompanyName:        sql.NullString{String: dto.CompanyName, Valid: dto.CompanyName != ""},
		ContactNumber:      contactNumber,
		RelationshipTypeID: relationshipTypeID,
		IsPrimaryContact:   isPrimaryContact,
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
			Location:           sql.NullString{String: e.Location, Valid: e.Location != ""},
			SchoolType:         e.SchoolType,
			YearCompleted:      e.YearCompleted,
			Awards:             sql.NullString{String: e.Awards, Valid: e.Awards != ""},
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
			AddressTypeID:   a.AddressTypeID,
			RegionName:      sql.NullString{String: a.RegionName, Valid: a.RegionName != ""},
			ProvinceName:    sql.NullString{String: a.ProvinceName, Valid: a.ProvinceName != ""},
			CityName:        sql.NullString{String: a.CityName, Valid: a.CityName != ""},
			BarangayName:    sql.NullString{String: a.BarangayName, Valid: a.BarangayName != ""},
			StreetLotBlk:    sql.NullString{String: a.StreetLotBlk, Valid: a.StreetLotBlk != ""},
			UnitNo:          sql.NullString{String: a.UnitNo, Valid: a.UnitNo != ""},
			BuildingName:    sql.NullString{String: a.BuildingName, Valid: a.BuildingName != ""},
		})
	}

	return s.repo.SaveAddressInfo(ctx, studentRecordID, addresses)
}

func (s *Service) SaveHealthRecord(
	ctx context.Context, studentRecordID int, req UpdateHealthRecordRequest,
) error {
	healthRecord := &StudentHealthRecord{
		StudentRecordID:       studentRecordID,
		VisionRemarkID:        req.VisionRemarkID,
		HearingRemarkID:       req.HearingRemarkID,
		MobilityRemarkID:      req.MobilityRemarkID,
		SpeechRemarkID:        req.SpeechRemarkID,
		GeneralHealthRemarkID: req.GeneralHealthRemarkID,
		ConsultedProfessional: sql.NullString{String: "", Valid: false},
		ConsultationReason:    sql.NullString{String: "", Valid: false},
		DateStarted:           sql.NullString{String: "", Valid: false},
		NumberOfSessions:      sql.NullInt64{Int64: 0, Valid: false},
		DateConcluded:         sql.NullString{String: "", Valid: false},
	}

	if req.ConsultedProfessional != nil {
		healthRecord.ConsultedProfessional = sql.NullString{
			String: *req.ConsultedProfessional,
			Valid:  true,
		}
	}

	if req.ConsultationReason != nil {
		healthRecord.ConsultationReason = sql.NullString{
			String: *req.ConsultationReason,
			Valid:  true,
		}
	}

	if req.DateStarted != nil {
		healthRecord.DateStarted = sql.NullString{
			String: *req.DateStarted,
			Valid:  true,
		}
	}

	if req.NumberOfSessions != nil {
		healthRecord.NumberOfSessions = sql.NullInt64{
			Int64: int64(*req.NumberOfSessions),
			Valid: true,
		}
	}

	if req.DateConcluded != nil {
		healthRecord.DateConcluded = sql.NullString{
			String: *req.DateConcluded,
			Valid:  true,
		}
	}

	return s.repo.SaveHealthRecord(ctx, healthRecord)
}
