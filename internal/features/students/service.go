package students

import (
	"context"
	"fmt"
	"sync"
	"time"

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
				MiddleName:    st.MiddleName,
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

func (s *Service) GetStudentBasicInfo(ctx context.Context, iirID int) (*StudentBasicInfoView, error) {
	info, err := s.repo.GetStudentBasicInfo(ctx, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student basic info: %w", err)
	}

	return info, nil
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
		MiddleName:    emergencyContact.MiddleName,
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
		TelephoneNumber:  personalInfo.TelephoneNumber,
		MobileNumber:     personalInfo.MobileNumber,
		IsEmployed:       personalInfo.IsEmployed,
		EmployerName:     personalInfo.EmployerName,
		EmployerAddress:  personalInfo.EmployerAddress,
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
		ParentalStatusDetails: studentFamily.ParentalStatusDetails,
		Brothers:              studentFamily.Brothers,
		Sisters:               studentFamily.Sisters,
		EmployedSiblings:      studentFamily.EmployedSiblings,
		OrdinalPosition:       studentFamily.OrdinalPosition,
		HaveQuietPlaceToStudy: studentFamily.HaveQuietPlaceToStudy,
		IsSharingRoom:         studentFamily.IsSharingRoom,
		SiblingSupportTypes:   supportTypes,
		RoomSharingDetails:    studentFamily.RoomSharingDetails,
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
			MiddleName:       relatedPerson.MiddleName,
			DateOfBirth:      relatedPerson.DateOfBirth,
			EducationalLevel: relatedPerson.EducationalLevel,
			Occupation:       relatedPerson.Occupation,
			EmployerName:     relatedPerson.EmployerName,
			EmployerAddress:  relatedPerson.EmployerAddress,
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
			Awards:           school.Awards,
		})
	}

	return &EducationalBackgroundDTO{
		ID:                 educationalBackground.ID,
		NatureOfSchooling:  educationalBackground.NatureOfSchooling,
		InterruptedDetails: educationalBackground.InterruptedDetails,
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
		OtherIncomeDetails:       financialInfo.OtherIncomeDetails,
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
		VisionDetails:           healthRecord.VisionDetails,
		HearingHasProblem:       healthRecord.HearingHasProblem,
		HearingDetails:          healthRecord.HearingDetails,
		SpeechHasProblem:        healthRecord.SpeechHasProblem,
		SpeechDetails:           healthRecord.SpeechDetails,
		GeneralHealthHasProblem: healthRecord.GeneralHealthHasProblem,
		GeneralHealthDetails:    healthRecord.GeneralHealthDetails,
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
			WhenDate:         c.WhenDate,
			ForWhat:          c.ForWhat,
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
				OtherSpecification: a.OtherSpecification,
				Role:               a.Role,
				RoleSpecification:  a.RoleSpecification,
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

func (s *Service) SubmitStudentIIR(ctx context.Context, userID int, req ComprehensiveProfileDTO) (int, error) {
	now := time.Now()
	iirRecord := &IIRRecord{
		UserID:      userID,
		IsSubmitted: false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	iirID, err := s.repo.CreateIIRRecord(ctx, iirRecord)
	if err != nil {
		return 0, fmt.Errorf("failed to create IIR record: %w", err)
	}

	// Use errgroup for concurrent operations
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		if req.Student.StudentPersonalInfoDTO.ID == 0 {
			req.Student.StudentPersonalInfoDTO.IIRID = iirID
		}
		return s.repo.UpsertStudentPersonalInfo(ctx, &StudentPersonalInfo{
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
			EmployerName:    req.Student.StudentPersonalInfoDTO.EmployerName,
			EmployerAddress: req.Student.StudentPersonalInfoDTO.EmployerAddress,
			MobileNumber:    req.Student.StudentPersonalInfoDTO.MobileNumber,
			TelephoneNumber: req.Student.StudentPersonalInfoDTO.TelephoneNumber,
			CreatedAt:       now,
			UpdatedAt:       now,
		})
	})

	g.Go(func() error {
		ec := req.Student.StudentPersonalInfoDTO.EmergencyContact
		// First upsert the address
		addressID, err := s.repo.UpsertAddress(ctx, &Address{
			Region:       ec.Address.Region,
			City:         ec.Address.City,
			Barangay:     ec.Address.Barangay,
			StreetDetail: ec.Address.StreetDetail,
			CreatedAt:    now,
			UpdatedAt:    now,
		})
		if err != nil {
			return fmt.Errorf("failed to upsert emergency contact address: %w", err)
		}

		// Then upsert the emergency contact
		_, err = s.repo.UpsertEmergencyContact(ctx, &EmergencyContact{
			IIRID:          iirID,
			FirstName:      ec.FirstName,
			MiddleName:     ec.MiddleName,
			LastName:       ec.LastName,
			ContactNumber:  ec.ContactNumber,
			RelationshipID: ec.Relationship.ID,
			AddressID:      addressID,
			CreatedAt:      now,
			UpdatedAt:      now,
		})
		return err
	})

	g.Go(func() error {
		for _, addrDTO := range req.Student.Addresses {
			addressID, err := s.repo.UpsertAddress(ctx, &Address{
				Region:       addrDTO.Address.Region,
				City:         addrDTO.Address.City,
				Barangay:     addrDTO.Address.Barangay,
				StreetDetail: addrDTO.Address.StreetDetail,
				CreatedAt:    now,
				UpdatedAt:    now,
			})
			if err != nil {
				return fmt.Errorf("failed to upsert student address: %w", err)
			}

			_, err = s.repo.UpsertStudentAddress(ctx, &StudentAddress{
				IIRID:       iirID,
				AddressID:   addressID,
				AddressType: addrDTO.AddressType,
				CreatedAt:   now,
				UpdatedAt:   now,
			})
			if err != nil {
				return fmt.Errorf("failed to save student address relation: %w", err)
			}
		}
		return nil
	})

	g.Go(func() error {
		ebID, err := s.repo.UpsertEducationalBackground(ctx, &EducationalBackground{
			IIRID:              iirID,
			NatureOfSchooling:  req.Education.NatureOfSchooling,
			InterruptedDetails: req.Education.InterruptedDetails,
			CreatedAt:          now,
			UpdatedAt:          now,
		})
		if err != nil {
			return fmt.Errorf("failed to upsert educational background: %w", err)
		}

		// Delete existing school details
		if err := s.repo.DeleteSchoolDetailsByEBID(ctx, ebID); err != nil {
			return fmt.Errorf("failed to delete existing school details: %w", err)
		}

		// Save new school details
		for _, schoolDTO := range req.Education.School {
			_, err := s.repo.UpsertSchoolDetails(ctx, &SchoolDetails{
				EBID:               ebID,
				EducationalLevelID: schoolDTO.EducationalLevel.ID,
				SchoolName:         schoolDTO.SchoolName,
				SchoolAddress:      schoolDTO.SchoolAddress,
				SchoolType:         schoolDTO.SchoolType,
				YearStarted:        schoolDTO.YearStarted,
				YearCompleted:      schoolDTO.YearCompleted,
				Awards:             schoolDTO.Awards,
				CreatedAt:          now,
				UpdatedAt:          now,
			})
			if err != nil {
				return fmt.Errorf("failed to save school details: %w", err)
			}
		}
		return nil
	})

	g.Go(func() error {
		fbID, err := s.repo.UpsertFamilyBackground(ctx, &FamilyBackground{
			IIRID:                 iirID,
			ParentalStatusID:      req.Family.FamilyBackgroundDTO.ParentalStatus.ID,
			ParentalStatusDetails: req.Family.FamilyBackgroundDTO.ParentalStatusDetails,
			Brothers:              req.Family.FamilyBackgroundDTO.Brothers,
			Sisters:               req.Family.FamilyBackgroundDTO.Sisters,
			EmployedSiblings:      req.Family.FamilyBackgroundDTO.EmployedSiblings,
			OrdinalPosition:       req.Family.FamilyBackgroundDTO.OrdinalPosition,
			HaveQuietPlaceToStudy: req.Family.FamilyBackgroundDTO.HaveQuietPlaceToStudy,
			IsSharingRoom:         req.Family.FamilyBackgroundDTO.IsSharingRoom,
			RoomSharingDetails:    req.Family.FamilyBackgroundDTO.RoomSharingDetails,
			NatureOfResidenceId:   req.Family.FamilyBackgroundDTO.NatureOfResidence.ID,
			CreatedAt:             now,
			UpdatedAt:             now,
		})
		if err != nil {
			return fmt.Errorf("failed to upsert family background: %w", err)
		}

		// Delete existing sibling supports
		if err := s.repo.DeleteStudentSiblingSupportsByFamilyID(ctx, fbID); err != nil {
			return fmt.Errorf("failed to delete existing sibling supports: %w", err)
		}

		// Save new sibling supports
		for _, supportType := range req.Family.FamilyBackgroundDTO.SiblingSupportTypes {
			if err := s.repo.CreateStudentSiblingSupport(ctx, &StudentSiblingSupport{
				FamilyBackgroundID: fbID,
				SupportTypeID:      supportType.ID,
			}); err != nil {
				return fmt.Errorf("failed to save sibling support: %w", err)
			}
		}

		// Delete existing related persons
		if err := s.repo.DeleteStudentRelatedPersons(ctx, iirID); err != nil {
			return fmt.Errorf("failed to delete existing related persons: %w", err)
		}

		// Save new related persons
		for _, relPersonDTO := range req.Family.RelatedPersons {
			relPersonID, err := s.repo.UpsertRelatedPerson(ctx, &RelatedPerson{
				FirstName:        relPersonDTO.FirstName,
				LastName:         relPersonDTO.LastName,
				MiddleName:       relPersonDTO.MiddleName,
				DateOfBirth:      relPersonDTO.DateOfBirth,
				EducationalLevel: relPersonDTO.EducationalLevel,
				Occupation:       relPersonDTO.Occupation,
				EmployerName:     relPersonDTO.EmployerName,
				EmployerAddress:  relPersonDTO.EmployerAddress,
				CreatedAt:        now,
				UpdatedAt:        now,
			})
			if err != nil {
				return fmt.Errorf("failed to save related person: %w", err)
			}

			if err := s.repo.UpsertStudentRelatedPerson(ctx, &StudentRelatedPerson{
				IIRID:           iirID,
				RelatedPersonID: relPersonID,
				RelationshipID:  relPersonDTO.Relationship.ID,
				IsParent:        relPersonDTO.IsParent,
				IsGuardian:      relPersonDTO.IsGuardian,
				IsLiving:        relPersonDTO.IsLiving,
				CreatedAt:       now,
				UpdatedAt:       now,
			}); err != nil {
				return fmt.Errorf("failed to save student related person relation: %w", err)
			}
		}
		return nil
	})

	g.Go(func() error {
		_, err := s.repo.UpsertStudentHealthRecord(ctx, &StudentHealthRecord{
			IIRID:                   iirID,
			VisionHasProblem:        req.Health.StudentHealthRecordDTO.VisionHasProblem,
			VisionDetails:           req.Health.StudentHealthRecordDTO.VisionDetails,
			HearingHasProblem:       req.Health.StudentHealthRecordDTO.HearingHasProblem,
			HearingDetails:          req.Health.StudentHealthRecordDTO.HearingDetails,
			SpeechHasProblem:        req.Health.StudentHealthRecordDTO.SpeechHasProblem,
			SpeechDetails:           req.Health.StudentHealthRecordDTO.SpeechDetails,
			GeneralHealthHasProblem: req.Health.StudentHealthRecordDTO.GeneralHealthHasProblem,
			GeneralHealthDetails:    req.Health.StudentHealthRecordDTO.GeneralHealthDetails,
			CreatedAt:               now,
			UpdatedAt:               now,
		})
		return err
	})

	g.Go(func() error {
		for _, consultationDTO := range req.Health.Consultations {
			_, err := s.repo.UpsertStudentConsultation(ctx, &StudentConsultation{
				IIRID:            iirID,
				ProfessionalType: consultationDTO.ProfessionalType,
				HasConsulted:     consultationDTO.HasConsulted,
				WhenDate:         consultationDTO.WhenDate,
				ForWhat:          consultationDTO.ForWhat,
				CreatedAt:        now,
				UpdatedAt:        now,
			})
			if err != nil {
				return fmt.Errorf("failed to save student consultation: %w", err)
			}
		}
		return nil
	})

	g.Go(func() error {
		sfID, err := s.repo.UpsertStudentFinance(ctx, &StudentFinance{
			IIRID:                      iirID,
			MonthlyFamilyIncomeRangeID: req.Family.Finance.MonthlyFamilyIncomeRange.ID,
			OtherIncomeDetails:         req.Family.Finance.OtherIncomeDetails,
			WeeklyAllowance:            req.Family.Finance.WeeklyAllowance,
			CreatedAt:                  now,
			UpdatedAt:                  now,
		})
		if err != nil {
			return fmt.Errorf("failed to upsert student finance: %w", err)
		}

		// Delete existing financial supports
		if err := s.repo.DeleteStudentFinancialSupportsByFinanceID(ctx, sfID); err != nil {
			return fmt.Errorf("failed to delete existing financial supports: %w", err)
		}

		// Save new financial supports
		for _, supportType := range req.Family.Finance.FinancialSupportTypes {
			if err := s.repo.CreateStudentFinancialSupport(ctx, &StudentFinancialSupport{
				StudentFinanceID: sfID,
				SupportTypeID:    supportType.ID,
				CreatedAt:        now,
				UpdatedAt:        now,
			}); err != nil {
				return fmt.Errorf("failed to save student financial support: %w", err)
			}
		}
		return nil
	})

	g.Go(func() error {
		// Delete existing activities
		if err := s.repo.DeleteStudentActivitiesByIIRID(ctx, iirID); err != nil {
			return fmt.Errorf("failed to delete existing activities: %w", err)
		}

		// Save new activities
		for _, activityDTO := range req.Interests.Activities {
			_, err := s.repo.CreateStudentActivity(ctx, &StudentActivity{
				IIRID:              iirID,
				OptionID:           activityDTO.ActivityOption.ID,
				OtherSpecification: activityDTO.OtherSpecification,
				Role:               activityDTO.Role,
				RoleSpecification:  activityDTO.RoleSpecification,
				CreatedAt:          now,
				UpdatedAt:          now,
			})
			if err != nil {
				return fmt.Errorf("failed to save student activity: %w", err)
			}
		}
		return nil
	})

	g.Go(func() error {
		// Delete existing subject preferences
		if err := s.repo.DeleteStudentSubjectPreferencesByIIRID(ctx, iirID); err != nil {
			return fmt.Errorf("failed to delete existing subject preferences: %w", err)
		}

		// Save new subject preferences
		for _, prefDTO := range req.Interests.SubjectPreferences {
			_, err := s.repo.CreateStudentSubjectPreference(ctx, &StudentSubjectPreference{
				IIRID:       iirID,
				SubjectName: prefDTO.SubjectName,
				IsFavorite:  prefDTO.IsFavorite,
				CreatedAt:   now,
				UpdatedAt:   now,
			})
			if err != nil {
				return fmt.Errorf("failed to save student subject preference: %w", err)
			}
		}
		return nil
	})

	g.Go(func() error {
		// Delete existing hobbies
		if err := s.repo.DeleteStudentHobbiesByIIRID(ctx, iirID); err != nil {
			return fmt.Errorf("failed to delete existing hobbies: %w", err)
		}

		// Save new hobbies
		for _, hobbyDTO := range req.Interests.Hobbies {
			_, err := s.repo.CreateStudentHobby(ctx, &StudentHobby{
				IIRID:        iirID,
				HobbyName:    hobbyDTO.HobbyName,
				PriorityRank: hobbyDTO.PriorityRank,
				CreatedAt:    now,
				UpdatedAt:    now,
			})
			if err != nil {
				return fmt.Errorf("failed to save student hobby: %w", err)
			}
		}
		return nil
	})

	g.Go(func() error {
		// Delete existing test results
		if err := s.repo.DeleteTestResultsByIIRID(ctx, iirID); err != nil {
			return fmt.Errorf("failed to delete existing test results: %w", err)
		}

		// Save new test results
		for _, testResultDTO := range req.TestResults {
			_, err := s.repo.CreateTestResult(ctx, &TestResult{
				IIRID:       iirID,
				TestDate:    testResultDTO.TestDate,
				TestName:    testResultDTO.TestName,
				RawScore:    testResultDTO.RawScore,
				Percentile:  testResultDTO.Percentile,
				Description: testResultDTO.Description,
				CreatedAt:   now,
				UpdatedAt:   now,
			})
			if err != nil {
				return fmt.Errorf("failed to save test result: %w", err)
			}
		}
		return nil
	})

	g.Go(func() error {
		// Delete existing significant notes
		if err := s.repo.DeleteSignificantNotesByIIRID(ctx, iirID); err != nil {
			return fmt.Errorf("failed to delete existing significant notes: %w", err)
		}

		// Save new significant notes
		for _, noteDTO := range req.SignificantNotes {
			_, err := s.repo.CreateSignificantNote(ctx, &SignificantNote{
				IIRID:               iirID,
				NoteDate:            noteDTO.NoteDate,
				IncidentDescription: noteDTO.IncidentDescription,
				Remarks:             noteDTO.Remarks,
				CreatedAt:           now,
				UpdatedAt:           now,
			})
			if err != nil {
				return fmt.Errorf("failed to save significant note: %w", err)
			}
		}
		return nil
	})

	return iirID, g.Wait()
}
