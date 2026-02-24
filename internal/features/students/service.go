package students

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
	"golang.org/x/sync/errgroup"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetGenders(ctx context.Context) ([]Gender, error) {
	genders, err := s.repo.GetGenders(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get genders: %w", err)
	}

	return genders, nil
}

func (s *Service) GetParentalStatusTypes(ctx context.Context) ([]ParentalStatusType, error) {
	statuses, err := s.repo.GetParentalStatusTypes(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get parental status types: %w", err)
	}

	return statuses, nil
}

func (s *Service) GetEnrollmentReasons(ctx context.Context) ([]EnrollmentReason, error) {
	reasons, err := s.repo.GetEnrollmentReasons(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get enrollment reasons: %w", err)
	}

	return reasons, nil
}

func (s *Service) GetIncomeRanges(ctx context.Context) ([]IncomeRange, error) {
	ranges, err := s.repo.GetIncomeRanges(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get income ranges: %w", err)
	}

	return ranges, nil
}

func (s *Service) GetStudentSupportTypes(ctx context.Context) ([]StudentSupportType, error) {
	supportTypes, err := s.repo.GetStudentSupportTypes(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get student support types: %w", err)
	}

	return supportTypes, nil
}

func (s *Service) GetSiblingSupportTypes(ctx context.Context) ([]SibilingSupportType, error) {
	supportTypes, err := s.repo.GetSiblingSupportTypes(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get sibling support types: %w", err)
	}

	return supportTypes, nil
}

func (s *Service) GetEducationalLevels(ctx context.Context) ([]EducationalLevel, error) {
	levels, err := s.repo.GetEducationalLevels(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get educational levels: %w", err)
	}

	return levels, nil
}

func (s *Service) GetCourses(ctx context.Context) ([]Course, error) {
	courses, err := s.repo.GetCourses(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get courses: %w", err)
	}

	return courses, nil
}

func (s *Service) GetCivilStatusTypes(ctx context.Context) ([]CivilStatusType, error) {
	statuses, err := s.repo.GetCivilStatusTypes(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get civil status types: %w", err)
	}

	return statuses, nil
}

func (s *Service) GetReligions(ctx context.Context) ([]Religion, error) {
	religions, err := s.repo.GetReligions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get religions: %w", err)
	}

	return religions, nil
}

func (s *Service) GetNatureOfResidenceTypes(ctx context.Context) ([]NatureOfResidenceType, error) {
	types, err := s.repo.GetNatureOfResidenceTypes(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get nature of residence types: %w", err)
	}

	return types, nil
}

func (s *Service) GetStudentRelationshipTypes(ctx context.Context) ([]StudentRelationshipType, error) {
	types, err := s.repo.GetStudentRelationshipTypes(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get student relationship types: %w", err)
	}

	return types, nil
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
		req.Search,
		req.GetOffset(),
		req.PageSize,
		req.OrderBy,
		req.CourseID,
		req.GenderID,
		req.YearLevel,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list students: %w", err)
	}

	var wg sync.WaitGroup

	studentDTOs := make([]StudentProfileDTO, len(students))
	errChan := make(chan error, len(students))

	for i, st := range students {
		wg.Add(1)

		go func(i int, st StudentProfileView) {
			defer wg.Done()

			course, err := s.repo.GetCourseByID(ctx, st.CourseID)
			if err != nil {
				errChan <- fmt.Errorf("failed to get course for student %d: %w", st.UserID, err)
				return
			}

			gender, err := s.repo.GetGenderByID(ctx, st.GenderID)
			if err != nil {
				errChan <- fmt.Errorf("failed to get gender for student %d: %w", st.UserID, err)
				return
			}

			// Create DTO
			studentDTOs[i] = StudentProfileDTO{
				IIRID:         st.IIRID,
				UserID:        st.UserID,
				FirstName:     st.FirstName,
				MiddleName:    structs.FromSqlNull(st.MiddleName),
				LastName:      st.LastName,
				Gender:        *gender,
				Email:         st.Email,
				StudentNumber: st.StudentNumber,
				Course:        *course,
				Section:       st.Section,
				YearLevel:     st.YearLevel,
			}
		}(i, st)
	}

	wg.Wait()
	close(errChan)

	if len(errChan) > 0 {
		return nil, <-errChan // Return the first error encountered
	}

	// Get total count for pagination
	total, err := s.repo.GetTotalStudentsCount(
		ctx,
		req.Search,
		req.CourseID,
		req.GenderID,
		req.YearLevel,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get total students count: %w", err)
	}

	totalPages := (total + req.PageSize - 1) / req.PageSize
	if totalPages == 0 {
		totalPages = 1
	}

	return &ListStudentsResponse{
		Students:   studentDTOs,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *Service) GetStudentProfile(ctx context.Context, iirID int) (*ComprehensiveProfileDTO, error) {
	profile := &ComprehensiveProfileDTO{IIRID: iirID}
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		basicInfo, err := s.GetStudentBasicInfo(ctx, iirID)
		if err != nil {
			return err
		}

		profile.Student.BasicInfo = *basicInfo
		return nil
	})

	g.Go(func() error {
		personalInfo, err := s.GetStudentPersonalInfo(ctx, iirID)
		if err != nil {
			return err
		}

		profile.Student.StudentPersonalInfoDTO = *personalInfo
		return nil
	})

	g.Go(func() error {
		addresses, err := s.GetStudentAddresses(ctx, iirID)
		if err != nil {
			return err
		}

		profile.Student.Addresses = addresses
		return nil
	})

	g.Go(func() error {
		education, err := s.GetEducationalBackground(ctx, iirID)
		if err != nil {
			return err
		}

		profile.Education = *education
		return nil
	})

	g.Go(func() error {
		familyBackground, err := s.GetStudentFamilyBackground(ctx, iirID)
		if err != nil {
			return err
		}

		profile.Family.FamilyBackgroundDTO = *familyBackground
		return nil
	})

	g.Go(func() error {
		relatedPersons, err := s.GetStudentRelatedPersons(ctx, iirID)
		if err != nil {
			return err
		}

		profile.Family.RelatedPersons = relatedPersons
		return nil
	})

	g.Go(func() error {
		finance, err := s.GetStudentFinancialInfo(ctx, iirID)
		if err != nil {
			return err
		}

		profile.Family.Finance = *finance
		return nil
	})

	g.Go(func() error {
		healthRecord, err := s.GetStudentHealthRecord(ctx, iirID)
		if err != nil {
			return err
		}

		profile.Health.StudentHealthRecordDTO = *healthRecord
		return nil
	})

	g.Go(func() error {
		consultations, err := s.GetStudentConsultations(ctx, iirID)
		if err != nil {
			return err
		}

		profile.Health.Consultations = consultations
		return nil
	})

	g.Go(func() error {
		activities, err := s.GetStudentActivities(ctx, iirID)
		if err != nil {
			return err
		}

		profile.Interests.Activities = activities
		return nil
	})

	g.Go(func() error {
		subjectPreferences, err := s.GetStudentSubjectPreferences(ctx, iirID)
		if err != nil {
			return err
		}

		profile.Interests.SubjectPreferences = subjectPreferences
		return nil
	})

	g.Go(func() error {
		hobbies, err := s.GetStudentHobbies(ctx, iirID)
		if err != nil {
			return err
		}

		profile.Interests.Hobbies = hobbies
		return nil
	})

	g.Go(func() error {
		testResults, err := s.GetStudentTestResults(ctx, iirID)
		if err != nil {
			return err
		}

		profile.TestResults = testResults
		return nil
	})

	g.Go(func() error {
		significantNotes, err := s.GetStudentSignificantNotes(ctx, iirID)
		if err != nil {
			return err
		}

		profile.SignificantNotes = significantNotes
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return profile, nil
}

func (s *Service) GetStudentBasicInfo(ctx context.Context, iirID int) (*StudentBasicInfoViewDTO, error) {
	info, err := s.repo.GetStudentBasicInfo(ctx, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student basic info: %w", err)
	}

	basicInfo := &StudentBasicInfoViewDTO{
		ID:         info.ID,
		FirstName:  info.FirstName,
		MiddleName: structs.FromSqlNull(info.MiddleName),
		LastName:   info.LastName,
		Email:      info.Email,
	}

	return basicInfo, nil
}

func (s *Service) GetIIRDraft(ctx context.Context, userID int) (*ComprehensiveProfileDTO, error) {
	draft, err := s.repo.GetIIRDraftByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get IIR draft: %w", err)
	}

	if draft == nil || len(draft.Data) == 0 {
		return nil, nil
	}

	var draftData ComprehensiveProfileDTO
	if err := json.Unmarshal([]byte(draft.Data), &draftData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal IIR draft data: %w", err)
	}

	return &draftData, nil
}

func (s *Service) GetStudentIIRByUserID(ctx context.Context, userID int) (*IIRRecord, error) {
	iir, err := s.repo.GetStudentIIRByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student IIR by user ID: %w", err)
	}

	return iir, nil
}

func (s *Service) GetStudentIIR(ctx context.Context, iirID int) (*IIRRecord, error) {
	iir, err := s.repo.GetStudentIIR(ctx, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student IIR: %w", err)
	}

	return iir, nil
}

func (s *Service) GetStudentEnrollmentReasons(ctx context.Context, iirID int) ([]StudentSelectedReasonDTO, error) {
	selectedReasons, err := s.repo.GetStudentEnrollmentReasons(ctx, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student enrollment reasons: %w", err)
	}

	var reasons []StudentSelectedReasonDTO
	for _, r := range selectedReasons {
		reason, err := s.repo.GetEnrollmentReasonByID(ctx, r.ReasonID)
		if err != nil {
			return nil, fmt.Errorf("failed to get enrollment reason by ID: %w", err)
		}

		reasons = append(reasons, StudentSelectedReasonDTO{
			Reason:          *reason,
			OtherReasonText: r.OtherReasonText,
		})
	}

	return reasons, nil
}

func (s *Service) GetStudentPersonalInfo(ctx context.Context, iirID int) (*StudentPersonalInfoDTO, error) {
	personalInfo, err := s.repo.GetStudentPersonalInfo(ctx, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student personal info: %w", err)
	}

	gender, err := s.repo.GetGenderByID(ctx, personalInfo.GenderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get gender by ID: %w", err)
	}

	civilStatus, err := s.repo.GetCivilStatusByID(ctx, personalInfo.CivilStatusID)
	if err != nil {
		return nil, fmt.Errorf("failed to get civil status by ID: %w", err)
	}

	religion, err := s.repo.GetReligionByID(ctx, personalInfo.ReligionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get religion by ID: %w", err)
	}

	course, err := s.repo.GetCourseByID(ctx, personalInfo.CourseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get course by ID: %w", err)
	}

	emergencyContact, err := s.repo.GetEmergencyContactByIIRID(ctx, personalInfo.IIRID)
	if err != nil {
		return nil, fmt.Errorf("failed to get emergency contact by IIR ID: %w", err)
	}

	emergencyContactRelationship, err := s.repo.GetStudentRelationshipByID(ctx, emergencyContact.RelationshipID)
	if err != nil {
		return nil, fmt.Errorf("failed to get emergency contact relationship by ID: %w", err)
	}

	emergencyContactAddress, err := s.repo.GetAddressByID(ctx, emergencyContact.AddressID)
	if err != nil {
		return nil, fmt.Errorf("failed to get emergency contact address by ID: %w", err)
	}

	emergencyContactDTO := EmergencyContactDTO{
		ID:            emergencyContact.ID,
		FirstName:     emergencyContact.FirstName,
		MiddleName:    structs.FromSqlNull(emergencyContact.MiddleName),
		LastName:      emergencyContact.LastName,
		ContactNumber: emergencyContact.ContactNumber,
		Relationship:  *emergencyContactRelationship,
		Address:       *emergencyContactAddress,
	}

	return &StudentPersonalInfoDTO{
		ID:               personalInfo.ID,
		StudentNumber:    personalInfo.StudentNumber,
		Gender:           *gender,
		CivilStatus:      *civilStatus,
		Religion:         *religion,
		HeightFt:         personalInfo.HeightFt,
		WeightKg:         personalInfo.WeightKg,
		Complexion:       personalInfo.Complexion,
		HighSchoolGWA:    personalInfo.HighSchoolGWA,
		Course:           *course,
		YearLevel:        personalInfo.YearLevel,
		Section:          personalInfo.Section,
		PlaceOfBirth:     personalInfo.PlaceOfBirth,
		DateOfBirth:      personalInfo.DateOfBirth,
		TelephoneNumber:  structs.FromSqlNull(personalInfo.TelephoneNumber),
		MobileNumber:     personalInfo.MobileNumber,
		IsEmployed:       personalInfo.IsEmployed,
		EmployerName:     structs.FromSqlNull(personalInfo.EmployerName),
		EmployerAddress:  structs.FromSqlNull(personalInfo.EmployerAddress),
		EmergencyContact: emergencyContactDTO,
	}, nil
}

func (s *Service) GetStudentAddresses(ctx context.Context, iirID int) ([]StudentAddressDTO, error) {
	studentAddresses, err := s.repo.GetStudentAddresses(ctx, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student addresses: %w", err)
	}

	var addresses []StudentAddressDTO
	for _, addr := range studentAddresses {
		a, err := s.repo.GetAddressByID(ctx, addr.AddressID)
		if err != nil {
			return nil, fmt.Errorf("failed to get address by ID: %w", err)
		}

		addresses = append(addresses, StudentAddressDTO{
			ID:          addr.ID,
			Address:     *a,
			AddressType: addr.AddressType,
			CreatedAt:   addr.CreatedAt,
			UpdatedAt:   addr.UpdatedAt,
		})
	}

	return addresses, nil
}

func (s *Service) GetStudentFamilyBackground(ctx context.Context, iirID int) (*FamilyBackgroundDTO, error) {
	studentFamily, err := s.repo.GetStudentFamilyBackground(ctx, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get family background: %w", err)
	}

	var family *FamilyBackgroundDTO

	parentalStatus, err := s.repo.GetParentalStatusByID(ctx, studentFamily.ParentalStatusID)
	if err != nil {
		return nil, fmt.Errorf("failed to get parental status by ID: %w", err)
	}

	natureOfResidence, err := s.repo.GetNatureOfResidenceByID(ctx, studentFamily.NatureOfResidenceId)
	if err != nil {
		return nil, fmt.Errorf("failed to get nature of residence by ID: %w", err)
	}

	siblingSupportTypes, err := s.repo.GetStudentSiblingSupport(ctx, studentFamily.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student sibling support: %w", err)
	}

	var supportTypes []SibilingSupportType
	for _, sst := range siblingSupportTypes {
		sibilingSupportType, err := s.repo.GetSiblingSupportTypeByID(ctx, sst.SupportTypeID)
		if err != nil {
			return nil, fmt.Errorf("failed to get sibling support type by ID: %w", err)
		}

		supportTypes = append(supportTypes, SibilingSupportType{
			ID:          sibilingSupportType.ID,
			SupportName: sibilingSupportType.SupportName,
		})
	}

	family = &FamilyBackgroundDTO{
		ID:                    studentFamily.ID,
		ParentalStatus:        *parentalStatus,
		ParentalStatusDetails: structs.FromSqlNull(studentFamily.ParentalStatusDetails),
		Brothers:              studentFamily.Brothers,
		Sisters:               studentFamily.Sisters,
		EmployedSiblings:      studentFamily.EmployedSiblings,
		OrdinalPosition:       studentFamily.OrdinalPosition,
		HaveQuietPlaceToStudy: studentFamily.HaveQuietPlaceToStudy,
		IsSharingRoom:         studentFamily.IsSharingRoom,
		SiblingSupportTypes:   supportTypes,
		RoomSharingDetails:    structs.FromSqlNull(studentFamily.RoomSharingDetails),
		NatureOfResidence:     *natureOfResidence,
	}

	return family, nil
}

func (s *Service) GetStudentRelatedPersons(ctx context.Context, iirID int) ([]RelatedPersonDTO, error) {
	studentRelatedPersons, err := s.repo.GetStudentRelatedPersons(ctx, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student related persons: %w", err)
	}

	var related []RelatedPersonDTO
	for _, srp := range studentRelatedPersons {
		relatedPerson, err := s.repo.GetRelatedPersonByID(ctx, srp.RelatedPersonID)
		if err != nil {
			return nil, fmt.Errorf("failed to get related person by ID: %w", err)
		}

		relationship, err := s.repo.GetStudentRelationshipByID(ctx, srp.RelationshipID)
		if err != nil {
			return nil, fmt.Errorf("failed to get student relationship by ID: %w", err)
		}

		related = append(related, RelatedPersonDTO{
			ID:               relatedPerson.ID,
			FirstName:        relatedPerson.FirstName,
			LastName:         relatedPerson.LastName,
			MiddleName:       structs.FromSqlNull(relatedPerson.MiddleName),
			DateOfBirth:      relatedPerson.DateOfBirth,
			EducationalLevel: relatedPerson.EducationalLevel,
			Occupation:       structs.FromSqlNull(relatedPerson.Occupation),
			EmployerName:     structs.FromSqlNull(relatedPerson.EmployerName),
			EmployerAddress:  structs.FromSqlNull(relatedPerson.EmployerAddress),
			Relationship:     *relationship,
			IsParent:         srp.IsParent,
			IsGuardian:       srp.IsGuardian,
			IsLiving:         srp.IsLiving,
		})

	}

	return related, nil
}

func (s *Service) GetEducationalBackground(ctx context.Context, iirID int) (*EducationalBackgroundDTO, error) {
	educationalBackground, err := s.repo.GetStudentEducationalBackground(ctx, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get educational background: %w", err)
	}

	schools, err := s.repo.GetSchoolDetailsByEBID(ctx, educationalBackground.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get school details by EBID: %w", err)
	}

	var schoolDTOs []SchoolDetailsDTO

	for _, school := range schools {
		educationalLevel, err := s.repo.GetEducationalLevelByID(ctx, school.EducationalLevelID)
		if err != nil {
			return nil, fmt.Errorf("failed to get educational level by ID: %w", err)
		}

		schoolDTOs = append(schoolDTOs, SchoolDetailsDTO{
			ID:               school.ID,
			EducationalLevel: *educationalLevel,
			SchoolName:       school.SchoolName,
			SchoolAddress:    school.SchoolAddress,
			SchoolType:       school.SchoolType,
			YearStarted:      school.YearStarted,
			YearCompleted:    school.YearCompleted,
			Awards:           structs.FromSqlNull(school.Awards),
		})
	}

	return &EducationalBackgroundDTO{
		ID:                 educationalBackground.ID,
		NatureOfSchooling:  educationalBackground.NatureOfSchooling,
		InterruptedDetails: structs.FromSqlNull(educationalBackground.InterruptedDetails),
		School:             schoolDTOs,
		CreatedAt:          educationalBackground.CreatedAt,
		UpdatedAt:          educationalBackground.UpdatedAt,
	}, nil
}

func (s *Service) GetStudentFinancialInfo(ctx context.Context, iirID int) (*StudentFinanceDTO, error) {
	financialInfo, err := s.repo.GetStudentFinancialInfo(ctx, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student financial info: %w", err)
	}

	incomeRange, err := s.repo.GetIncomeRangeByID(ctx, financialInfo.MonthlyFamilyIncomeRangeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get income range by ID: %w", err)
	}

	financialSupportTypes, err := s.repo.GetFinancialSupportTypeByFinanceID(ctx, financialInfo.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get financial support types by finance ID: %w", err)
	}

	var supportTypes []StudentSupportType
	for _, fst := range financialSupportTypes {
		supportType, err := s.repo.GetStudentSupportByID(ctx, fst.SupportTypeID)
		if err != nil {
			return nil, fmt.Errorf("failed to get student support type by ID: %w", err)
		}
		supportTypes = append(supportTypes, *supportType)
	}

	return &StudentFinanceDTO{
		ID:                       financialInfo.ID,
		MonthlyFamilyIncomeRange: *incomeRange,
		OtherIncomeDetails:       structs.FromSqlNull(financialInfo.OtherIncomeDetails),
		FinancialSupportTypes:    supportTypes,
		WeeklyAllowance:          financialInfo.WeeklyAllowance,
	}, nil
}

func (s *Service) GetStudentHealthRecord(ctx context.Context, iirID int) (*StudentHealthRecordDTO, error) {
	healthRecord, err := s.repo.GetStudentHealthRecord(ctx, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student health record: %w", err)
	}

	return &StudentHealthRecordDTO{
		ID:                      healthRecord.ID,
		VisionHasProblem:        healthRecord.VisionHasProblem,
		VisionDetails:           structs.FromSqlNull(healthRecord.VisionDetails),
		HearingHasProblem:       healthRecord.HearingHasProblem,
		HearingDetails:          structs.FromSqlNull(healthRecord.HearingDetails),
		SpeechHasProblem:        healthRecord.SpeechHasProblem,
		SpeechDetails:           structs.FromSqlNull(healthRecord.SpeechDetails),
		GeneralHealthHasProblem: healthRecord.GeneralHealthHasProblem,
		GeneralHealthDetails:    structs.FromSqlNull(healthRecord.GeneralHealthDetails),
	}, nil
}

func (s *Service) GetStudentConsultations(ctx context.Context, iirID int) ([]StudentConsultationDTO, error) {
	consultations, err := s.repo.GetStudentConsultations(ctx, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student consultations: %w", err)
	}

	var consultationDTOs []StudentConsultationDTO
	for _, c := range consultations {
		consultationDTOs = append(consultationDTOs, StudentConsultationDTO{
			ID:               c.ID,
			ProfessionalType: c.ProfessionalType,
			HasConsulted:     c.HasConsulted,
			WhenDate:         structs.FromSqlNull(c.WhenDate),
			ForWhat:          structs.FromSqlNull(c.ForWhat),
		})
	}

	return consultationDTOs, nil
}

func (s *Service) GetStudentActivities(ctx context.Context, iirID int) ([]StudentActivityDTO, error) {
	activities, err := s.repo.GetStudentActivities(ctx, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student activities: %w", err)
	}

	var wg sync.WaitGroup
	activityDTOs := make([]StudentActivityDTO, len(activities))
	errChan := make(chan error, len(activities))

	for i, a := range activities {
		wg.Add(1)
		go func(i int, a StudentActivity) {
			defer wg.Done()
			option, err := s.repo.GetActivityOptionByID(ctx, a.OptionID)
			if err != nil {
				errChan <- fmt.Errorf("failed to get activity option by ID: %w", err)
				return
			}
			activityDTOs[i] = StudentActivityDTO{
				ID:                 a.ID,
				ActivityOption:     *option,
				OtherSpecification: structs.FromSqlNull(a.OtherSpecification),
				Role:               a.Role,
				RoleSpecification:  structs.FromSqlNull(a.RoleSpecification),
			}
		}(i, a)
	}

	wg.Wait()
	close(errChan)

	if len(errChan) > 0 {
		return nil, <-errChan // Return the first error encountered
	}

	return activityDTOs, nil
}

func (s *Service) GetStudentSubjectPreferences(ctx context.Context, iirID int) ([]StudentSubjectPreferenceDTO, error) {
	preferences, err := s.repo.GetStudentSubjectPreferences(ctx, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student subject preferences: %w", err)
	}

	var preferenceDTOs []StudentSubjectPreferenceDTO
	for _, p := range preferences {
		preferenceDTOs = append(preferenceDTOs, StudentSubjectPreferenceDTO{
			ID:          p.ID,
			SubjectName: p.SubjectName,
			IsFavorite:  p.IsFavorite,
		})
	}

	return preferenceDTOs, nil
}

func (s *Service) GetStudentHobbies(ctx context.Context, iirID int) ([]StudentHobbyDTO, error) {
	hobbies, err := s.repo.GetStudentHobbies(ctx, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student hobbies: %w", err)
	}

	var hobbyDTOs []StudentHobbyDTO
	for _, h := range hobbies {
		hobbyDTOs = append(hobbyDTOs, StudentHobbyDTO{
			ID:           h.ID,
			HobbyName:    h.HobbyName,
			PriorityRank: h.PriorityRank,
		})
	}

	return hobbyDTOs, nil
}

func (s *Service) GetStudentTestResults(ctx context.Context, iirID int) ([]TestResultDTO, error) {
	testResults, err := s.repo.GetStudentTestResults(ctx, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student test results: %w", err)
	}

	var testResultDTOs []TestResultDTO
	for _, tr := range testResults {
		testResultDTOs = append(testResultDTOs, TestResultDTO{
			ID:          tr.ID,
			TestDate:    tr.TestDate,
			TestName:    tr.TestName,
			RawScore:    tr.RawScore,
			Percentile:  tr.Percentile,
			Description: tr.Description,
		})
	}

	return testResultDTOs, nil
}

func (s *Service) GetStudentSignificantNotes(ctx context.Context, iirID int) ([]SignificantNoteDTO, error) {
	notes, err := s.repo.GetStudentSignificantNotes(ctx, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student significant notes: %w", err)
	}

	var noteDTOs []SignificantNoteDTO
	for _, n := range notes {
		noteDTOs = append(noteDTOs, SignificantNoteDTO{
			ID:                  n.ID,
			NoteDate:            n.NoteDate,
			IncidentDescription: n.IncidentDescription,
			Remarks:             n.Remarks,
			CreatedAt:           n.CreatedAt,
			UpdatedAt:           n.UpdatedAt,
		})
	}

	return noteDTOs, nil
}

func (s *Service) SaveIIRDraft(ctx context.Context, userID int, req ComprehensiveProfileDTO) (int, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return 0, fmt.Errorf("failed to encode draft data: %w", err)
	}

	draft := IIRDraft{
		UserID: userID,
		Data:   string(jsonData),
	}

	draftID, err := s.repo.UpsertIIRDraft(ctx, draft)
	if err != nil {
		return 0, fmt.Errorf("failed to save IIR draft: %w", err)
	}

	return draftID, nil
}

func (s *Service) SubmitStudentIIR(ctx context.Context, userID int, req ComprehensiveProfileDTO) (int, error) {
	tx, err := s.repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	iirRecord := &IIRRecord{
		UserID:      userID,
		IsSubmitted: false,
	}
	iirID, err := s.repo.UpsertIIRRecord(ctx, tx, iirRecord)
	if err != nil {
		return 0, fmt.Errorf("failed to create IIR record: %w", err)
	}

	// 1. Save Student Personal Info
	if err := s.repo.UpsertStudentPersonalInfo(ctx, tx, &StudentPersonalInfo{
		IIRID:           iirID,
		StudentNumber:   req.Student.StudentPersonalInfoDTO.StudentNumber,
		GenderID:        req.Student.StudentPersonalInfoDTO.Gender.ID,
		CivilStatusID:   req.Student.StudentPersonalInfoDTO.CivilStatus.ID,
		ReligionID:      req.Student.StudentPersonalInfoDTO.Religion.ID,
		HeightFt:        req.Student.StudentPersonalInfoDTO.HeightFt,
		WeightKg:        req.Student.StudentPersonalInfoDTO.WeightKg,
		Complexion:      req.Student.StudentPersonalInfoDTO.Complexion,
		HighSchoolGWA:   req.Student.StudentPersonalInfoDTO.HighSchoolGWA,
		CourseID:        req.Student.StudentPersonalInfoDTO.Course.ID,
		YearLevel:       req.Student.StudentPersonalInfoDTO.YearLevel,
		Section:         req.Student.StudentPersonalInfoDTO.Section,
		PlaceOfBirth:    req.Student.StudentPersonalInfoDTO.PlaceOfBirth,
		DateOfBirth:     req.Student.StudentPersonalInfoDTO.DateOfBirth,
		IsEmployed:      req.Student.StudentPersonalInfoDTO.IsEmployed,
		EmployerName:    structs.ToSqlNull(req.Student.StudentPersonalInfoDTO.EmployerName),
		EmployerAddress: structs.ToSqlNull(req.Student.StudentPersonalInfoDTO.EmployerAddress),
		MobileNumber:    req.Student.StudentPersonalInfoDTO.MobileNumber,
		TelephoneNumber: structs.ToSqlNull(req.Student.StudentPersonalInfoDTO.TelephoneNumber),
	}); err != nil {
		return 0, fmt.Errorf("failed to upsert student personal info: %w", err)
	}

	// 2. Save Emergency Contact
	ec := req.Student.StudentPersonalInfoDTO.EmergencyContact
	addressID, err := s.repo.UpsertAddress(ctx, tx, &Address{
		Region:       ec.Address.Region,
		City:         ec.Address.City,
		Barangay:     ec.Address.Barangay,
		StreetDetail: ec.Address.StreetDetail,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to upsert emergency contact address: %w", err)
	}

	if _, err := s.repo.UpsertEmergencyContact(ctx, tx, &EmergencyContact{
		IIRID:          iirID,
		FirstName:      ec.FirstName,
		MiddleName:     structs.ToSqlNull(ec.MiddleName),
		LastName:       ec.LastName,
		ContactNumber:  ec.ContactNumber,
		RelationshipID: ec.Relationship.ID,
		AddressID:      addressID,
	}); err != nil {
		return 0, fmt.Errorf("failed to upsert emergency contact: %w", err)
	}

	// 3. Save Student Addresses
	for _, addrDTO := range req.Student.Addresses {
		addressID, err := s.repo.UpsertAddress(ctx, tx, &Address{
			Region:       addrDTO.Address.Region,
			City:         addrDTO.Address.City,
			Barangay:     addrDTO.Address.Barangay,
			StreetDetail: addrDTO.Address.StreetDetail,
		})
		if err != nil {
			return 0, fmt.Errorf("failed to upsert student address: %w", err)
		}

		if _, err := s.repo.UpsertStudentAddress(ctx, tx, &StudentAddress{
			IIRID:       iirID,
			AddressID:   addressID,
			AddressType: addrDTO.AddressType,
		}); err != nil {
			return 0, fmt.Errorf("failed to save student address relation: %w", err)
		}
	}

	// 4. Save Educational Background
	ebID, err := s.repo.UpsertEducationalBackground(ctx, tx, &EducationalBackground{
		IIRID:              iirID,
		NatureOfSchooling:  req.Education.NatureOfSchooling,
		InterruptedDetails: structs.ToSqlNull(req.Education.InterruptedDetails),
	})
	if err != nil {
		return 0, fmt.Errorf("failed to upsert educational background: %w", err)
	}

	if err := s.repo.DeleteSchoolDetailsByEBID(ctx, tx, ebID); err != nil {
		return 0, fmt.Errorf("failed to delete existing school details: %w", err)
	}

	for _, schoolDTO := range req.Education.School {
		if _, err := s.repo.UpsertSchoolDetails(ctx, tx, &SchoolDetails{
			EBID:               ebID,
			EducationalLevelID: schoolDTO.EducationalLevel.ID,
			SchoolName:         schoolDTO.SchoolName,
			SchoolAddress:      schoolDTO.SchoolAddress,
			SchoolType:         schoolDTO.SchoolType,
			YearStarted:        schoolDTO.YearStarted,
			YearCompleted:      schoolDTO.YearCompleted,
			Awards:             structs.ToSqlNull(schoolDTO.Awards),
		}); err != nil {
			return 0, fmt.Errorf("failed to save school details: %w", err)
		}
	}

	// 5. Save Family Background
	fbID, err := s.repo.UpsertFamilyBackground(ctx, tx, &FamilyBackground{
		IIRID:                 iirID,
		ParentalStatusID:      req.Family.FamilyBackgroundDTO.ParentalStatus.ID,
		ParentalStatusDetails: structs.ToSqlNull(req.Family.FamilyBackgroundDTO.ParentalStatusDetails),
		Brothers:              req.Family.FamilyBackgroundDTO.Brothers,
		Sisters:               req.Family.FamilyBackgroundDTO.Sisters,
		EmployedSiblings:      req.Family.FamilyBackgroundDTO.EmployedSiblings,
		OrdinalPosition:       req.Family.FamilyBackgroundDTO.OrdinalPosition,
		HaveQuietPlaceToStudy: req.Family.FamilyBackgroundDTO.HaveQuietPlaceToStudy,
		IsSharingRoom:         req.Family.FamilyBackgroundDTO.IsSharingRoom,
		RoomSharingDetails:    structs.ToSqlNull(req.Family.FamilyBackgroundDTO.RoomSharingDetails),
		NatureOfResidenceId:   req.Family.FamilyBackgroundDTO.NatureOfResidence.ID,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to upsert family background: %w", err)
	}

	// Save Sibling Supports
	if err := s.repo.DeleteStudentSiblingSupportsByFamilyID(ctx, tx, fbID); err != nil {
		return 0, fmt.Errorf("failed to delete existing sibling supports: %w", err)
	}

	for _, supportType := range req.Family.FamilyBackgroundDTO.SiblingSupportTypes {
		if err := s.repo.CreateStudentSiblingSupport(ctx, tx, &StudentSiblingSupport{
			FamilyBackgroundID: fbID,
			SupportTypeID:      supportType.ID,
		}); err != nil {
			return 0, fmt.Errorf("failed to save sibling support: %w", err)
		}
	}

	// Save Related Persons
	if err := s.repo.DeleteStudentRelatedPersons(ctx, tx, iirID); err != nil {
		return 0, fmt.Errorf("failed to delete existing related persons: %w", err)
	}

	for _, relPersonDTO := range req.Family.RelatedPersons {
		relPersonID, err := s.repo.UpsertRelatedPerson(ctx, tx, &RelatedPerson{
			FirstName:        relPersonDTO.FirstName,
			LastName:         relPersonDTO.LastName,
			MiddleName:       structs.ToSqlNull(relPersonDTO.MiddleName),
			DateOfBirth:      relPersonDTO.DateOfBirth,
			EducationalLevel: relPersonDTO.EducationalLevel,
			Occupation:       structs.ToSqlNull(relPersonDTO.Occupation),
			EmployerName:     structs.ToSqlNull(relPersonDTO.EmployerName),
			EmployerAddress:  structs.ToSqlNull(relPersonDTO.EmployerAddress),
		})
		if err != nil {
			return 0, fmt.Errorf("failed to save related person: %w", err)
		}

		if err := s.repo.UpsertStudentRelatedPerson(ctx, tx, &StudentRelatedPerson{
			IIRID:           iirID,
			RelatedPersonID: relPersonID,
			RelationshipID:  relPersonDTO.Relationship.ID,
			IsParent:        relPersonDTO.IsParent,
			IsGuardian:      relPersonDTO.IsGuardian,
			IsLiving:        relPersonDTO.IsLiving,
		}); err != nil {
			return 0, fmt.Errorf("failed to save student related person relation: %w", err)
		}
	}

	// 6. Save Health Record
	if _, err := s.repo.UpsertStudentHealthRecord(ctx, tx, &StudentHealthRecord{
		IIRID:                   iirID,
		VisionHasProblem:        req.Health.StudentHealthRecordDTO.VisionHasProblem,
		VisionDetails:           structs.ToSqlNull(req.Health.StudentHealthRecordDTO.VisionDetails),
		HearingHasProblem:       req.Health.StudentHealthRecordDTO.HearingHasProblem,
		HearingDetails:          structs.ToSqlNull(req.Health.StudentHealthRecordDTO.HearingDetails),
		SpeechHasProblem:        req.Health.StudentHealthRecordDTO.SpeechHasProblem,
		SpeechDetails:           structs.ToSqlNull(req.Health.StudentHealthRecordDTO.SpeechDetails),
		GeneralHealthHasProblem: req.Health.StudentHealthRecordDTO.GeneralHealthHasProblem,
		GeneralHealthDetails:    structs.ToSqlNull(req.Health.StudentHealthRecordDTO.GeneralHealthDetails),
	}); err != nil {
		return 0, fmt.Errorf("failed to upsert student health record: %w", err)
	}

	// 7. Save Consultations
	for _, consultationDTO := range req.Health.Consultations {
		if _, err := s.repo.UpsertStudentConsultation(ctx, tx, &StudentConsultation{
			IIRID:            iirID,
			ProfessionalType: consultationDTO.ProfessionalType,
			HasConsulted:     consultationDTO.HasConsulted,
			WhenDate:         structs.ToSqlNull(consultationDTO.WhenDate),
			ForWhat:          structs.ToSqlNull(consultationDTO.ForWhat),
		}); err != nil {
			return 0, fmt.Errorf("failed to save student consultation: %w", err)
		}
	}

	// 8. Save Financial Info
	sfID, err := s.repo.UpsertStudentFinance(ctx, tx, &StudentFinance{
		IIRID:                      iirID,
		MonthlyFamilyIncomeRangeID: req.Family.Finance.MonthlyFamilyIncomeRange.ID,
		OtherIncomeDetails:         structs.ToSqlNull(req.Family.Finance.OtherIncomeDetails),
		WeeklyAllowance:            req.Family.Finance.WeeklyAllowance,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to upsert student finance: %w", err)
	}

	if err := s.repo.DeleteStudentFinancialSupportsByFinanceID(ctx, tx, sfID); err != nil {
		return 0, fmt.Errorf("failed to delete existing financial supports: %w", err)
	}

	for _, supportType := range req.Family.Finance.FinancialSupportTypes {
		if err := s.repo.CreateStudentFinancialSupport(ctx, tx, &StudentFinancialSupport{
			StudentFinanceID: sfID,
			SupportTypeID:    supportType.ID,
		}); err != nil {
			return 0, fmt.Errorf("failed to save student financial support: %w", err)
		}
	}

	// 9. Save Activities
	if err := s.repo.DeleteStudentActivitiesByIIRID(ctx, tx, iirID); err != nil {
		return 0, fmt.Errorf("failed to delete existing activities: %w", err)
	}

	for _, activityDTO := range req.Interests.Activities {
		if _, err := s.repo.CreateStudentActivity(ctx, tx, &StudentActivity{
			IIRID:              iirID,
			OptionID:           activityDTO.ActivityOption.ID,
			OtherSpecification: structs.ToSqlNull(activityDTO.OtherSpecification),
			Role:               activityDTO.Role,
			RoleSpecification:  structs.ToSqlNull(activityDTO.RoleSpecification),
		}); err != nil {
			return 0, fmt.Errorf("failed to save student activity: %w", err)
		}
	}

	// 10. Save Subject Preferences
	if err := s.repo.DeleteStudentSubjectPreferencesByIIRID(ctx, tx, iirID); err != nil {
		return 0, fmt.Errorf("failed to delete existing subject preferences: %w", err)
	}

	for _, prefDTO := range req.Interests.SubjectPreferences {
		if _, err := s.repo.CreateStudentSubjectPreference(ctx, tx, &StudentSubjectPreference{
			IIRID:       iirID,
			SubjectName: prefDTO.SubjectName,
			IsFavorite:  prefDTO.IsFavorite,
		}); err != nil {
			return 0, fmt.Errorf("failed to save student subject preference: %w", err)
		}
	}

	// 11. Save Hobbies
	if err := s.repo.DeleteStudentHobbiesByIIRID(ctx, tx, iirID); err != nil {
		return 0, fmt.Errorf("failed to delete existing hobbies: %w", err)
	}

	for _, hobbyDTO := range req.Interests.Hobbies {
		if _, err := s.repo.CreateStudentHobby(ctx, tx, &StudentHobby{
			IIRID:        iirID,
			HobbyName:    hobbyDTO.HobbyName,
			PriorityRank: hobbyDTO.PriorityRank,
		}); err != nil {
			return 0, fmt.Errorf("failed to save student hobby: %w", err)
		}
	}

	// 12. Save Test Results
	if err := s.repo.DeleteTestResultsByIIRID(ctx, tx, iirID); err != nil {
		return 0, fmt.Errorf("failed to delete existing test results: %w", err)
	}

	for _, testResultDTO := range req.TestResults {
		if _, err := s.repo.CreateTestResult(ctx, tx, &TestResult{
			IIRID:       iirID,
			TestDate:    testResultDTO.TestDate,
			TestName:    testResultDTO.TestName,
			RawScore:    testResultDTO.RawScore,
			Percentile:  testResultDTO.Percentile,
			Description: testResultDTO.Description,
		}); err != nil {
			return 0, fmt.Errorf("failed to save test result: %w", err)
		}
	}

	// 13. Save Significant Notes
	if err := s.repo.DeleteSignificantNotesByIIRID(ctx, tx, iirID); err != nil {
		return 0, fmt.Errorf("failed to delete existing significant notes: %w", err)
	}

	IIRRecordUpdate := &IIRRecord{
		ID:          iirID,
		UserID:      userID,
		IsSubmitted: true,
	}

	if _, err := s.repo.UpsertIIRRecord(ctx, tx, IIRRecordUpdate); err != nil {
		return 0, fmt.Errorf("failed to update IIR record as submitted: %w", err)
	}

	// All operations succeeded - commit transaction
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return iirID, nil
}
