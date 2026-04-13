package students

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
	"github.com/olazo-johnalbert/duckload-api/internal/core/config"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/core/pdf"
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
	"github.com/olazo-johnalbert/duckload-api/internal/features/locations"
	"github.com/olazo-johnalbert/duckload-api/internal/features/users"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
	"golang.org/x/sync/errgroup"
)

// Service provides student-related business logic and data access.
type Service struct {
	repo         RepositoryInterface
	locationsSvc locations.ServiceInterface
	userService  users.ServiceInterface
	logService   audit.Logger
	notifService audit.Notifier
	cfg          *config.Config
	pdfService   *pdf.Service
}

// NewService creates a new student service instance.
func NewService(
	repo RepositoryInterface,
	locationsSvc locations.ServiceInterface,
	userService users.ServiceInterface,
	logService audit.Logger,
	notifService audit.Notifier,
	cfg *config.Config,
	pdfService *pdf.Service,
) *Service {
	return &Service{
		repo:         repo,
		locationsSvc: locationsSvc,
		userService:  userService,
		logService:   logService,
		notifService: notifService,
		cfg:          cfg,
		pdfService:   pdfService,
	}
}

// GetGenders retrieves all available gender types.
func (s *Service) GetGenders(ctx context.Context) ([]Gender, error) {
	genders, err := s.repo.GetGenders(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get genders: %w", err)
	}

	return genders, nil
}

// GetParentalStatusTypes retrieves all available parental status types.
func (s *Service) GetParentalStatusTypes(
	ctx context.Context,
) ([]ParentalStatusType, error) {
	statuses, err := s.repo.GetParentalStatusTypes(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get parental status types: %w", err)
	}

	return statuses, nil
}

// GetEnrollmentReasons retrieves all available enrollment reason types.
func (s *Service) GetEnrollmentReasons(
	ctx context.Context,
) ([]EnrollmentReason, error) {
	reasons, err := s.repo.GetEnrollmentReasons(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get enrollment reasons: %w", err)
	}

	return reasons, nil
}

// GetIncomeRanges retrieves all available family income range types.
func (s *Service) GetIncomeRanges(ctx context.Context) ([]IncomeRange, error) {
	ranges, err := s.repo.GetIncomeRanges(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get income ranges: %w", err)
	}

	return ranges, nil
}

// GetStudentSupportTypes retrieves all available student support types.
func (s *Service) GetStudentSupportTypes(
	ctx context.Context,
) ([]StudentSupportType, error) {
	supportTypes, err := s.repo.GetStudentSupportTypes(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get student support types: %w", err)
	}

	return supportTypes, nil
}

// GetSiblingSupportTypes retrieves all available sibling support types.
func (s *Service) GetSiblingSupportTypes(
	ctx context.Context,
) ([]SibilingSupportType, error) {
	supportTypes, err := s.repo.GetSiblingSupportTypes(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get sibling support types: %w", err)
	}

	return supportTypes, nil
}

// GetEducationalLevels retrieves all available educational levels.
func (s *Service) GetEducationalLevels(
	ctx context.Context,
) ([]EducationalLevel, error) {
	levels, err := s.repo.GetEducationalLevels(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get educational levels: %w", err)
	}

	return levels, nil
}

// GetCourses retrieves all available academic courses.
func (s *Service) GetCourses(
	ctx context.Context,
) ([]Course, error) {
	courses, err := s.repo.GetCourses(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get courses: %w", err)
	}

	return courses, nil
}

// GetCivilStatusTypes retrieves all available civil status types.
func (s *Service) GetCivilStatusTypes(
	ctx context.Context,
) ([]CivilStatusType, error) {
	statuses, err := s.repo.GetCivilStatusTypes(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get civil status types: %w", err)
	}

	return statuses, nil
}

// GetReligions retrieves all available religions.
func (s *Service) GetReligions(
	ctx context.Context,
) ([]Religion, error) {
	religions, err := s.repo.GetReligions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get religions: %w", err)
	}

	return religions, nil
}

// GetNatureOfResidenceTypes retrieves all available nature of residence types.
func (s *Service) GetNatureOfResidenceTypes(
	ctx context.Context,
) ([]NatureOfResidenceType, error) {
	types, err := s.repo.GetNatureOfResidenceTypes(ctx)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get nature of residence types: %w",
			err,
		)
	}

	return types, nil
}

// GetActivityOptions retrieves all available student activity options.
func (s *Service) GetActivityOptions(
	ctx context.Context,
) ([]ActivityOption, error) {
	options, err := s.repo.GetActivityOptions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get activity options: %w", err)
	}

	return options, nil
}

// GetStudentRelationshipTypes retrieves all available student relationship
// types.
func (s *Service) GetStudentRelationshipTypes(
	ctx context.Context,
) ([]StudentRelationshipType, error) {
	types, err := s.repo.GetStudentRelationshipTypes(ctx)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get student relationship types: %w",
			err,
		)
	}

	return types, nil
}

// Retrieve - List
// ListStudents retrieves a paginated list of student profiles.
func (s *Service) ListStudents(
	ctx context.Context, req ListStudentsRequest,
) (*ListStudentsResponse, error) {
	req.SetDefaults("last_name")

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
				errChan <- fmt.Errorf(
					"failed to get course for student %s: %w",
					st.UserID,
					err,
				)
				return
			}

			gender, err := s.repo.GetGenderByID(ctx, st.GenderID)
			if err != nil {
				errChan <- fmt.Errorf(
					"failed to get gender for student %s: %w",
					st.UserID,
					err,
				)
				return
			}

			// Create DTO
			studentDTOs[i] = StudentProfileDTO{
				IIRID:         st.IIRID,
				UserID:        st.UserID,
				FirstName:     st.FirstName,
				MiddleName:    structs.FromSqlNull(st.MiddleName),
				LastName:      st.LastName,
				SuffixName:    structs.FromSqlNull(st.SuffixName),
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

	return &ListStudentsResponse{
		Students: studentDTOs,
		Meta:     structs.CalculateMetadata(total, req.Page, req.PageSize),
	}, nil
}

// GetStudentProfile retrieves the full comprehensive profile of a student.
func (s *Service) GetStudentProfile(
	ctx context.Context,
	iirID string,
) (*ComprehensiveProfileDTO, error) {
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

	// g.Go(func() error {
	// 	testResults, err := s.GetStudentTestResults(ctx, iirID)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	profile.TestResults = testResults
	// 	return nil
	// })

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return profile, nil
}

// GetStudentBasicInfo retrieves basic identification info for a student.
func (s *Service) GetStudentBasicInfo(
	ctx context.Context,
	iirID string,
) (*StudentBasicInfoViewDTO, error) {
	info, err := s.repo.GetStudentBasicInfo(ctx, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student basic info: %w", err)
	}

	basicInfo := &StudentBasicInfoViewDTO{
		Email:      info.Email,
		FirstName:  info.FirstName,
		MiddleName: structs.FromSqlNull(info.MiddleName),
		LastName:   info.LastName,
	}

	return basicInfo, nil
}

// GetIIRDraft retrieves the latest IIR draft for a user.
func (s *Service) GetIIRDraft(
	ctx context.Context,
	userID string,
) (*ComprehensiveProfileDTO, error) {
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

// GetStudentIIRByUserID retrieves the IIR record for a specific user ID.
func (s *Service) GetStudentIIRByUserID(
	ctx context.Context,
	userID string,
) (*IIRRecord, error) {
	iir, err := s.repo.GetStudentIIRByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student IIR by user ID: %w", err)
	}

	return iir, nil
}

// GetStudentIIR retrieves a specific IIR record by its ID.
func (s *Service) GetStudentIIR(
	ctx context.Context,
	iirID string,
) (*IIRRecord, error) {
	iir, err := s.repo.GetStudentIIR(ctx, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student IIR: %w", err)
	}

	return iir, nil
}

// GetStudentEnrollmentReasons retrieves the reasons why a student enrolled.
func (s *Service) GetStudentEnrollmentReasons(
	ctx context.Context,
	iirID string,
) ([]StudentSelectedReasonDTO, error) {
	selectedReasons, err := s.repo.GetStudentEnrollmentReasons(ctx, iirID)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get student enrollment reasons: %w",
			err,
		)
	}

	var reasons []StudentSelectedReasonDTO
	for _, r := range selectedReasons {
		reason, err := s.repo.GetEnrollmentReasonByID(ctx, r.ReasonID)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to get enrollment reason by ID: %w",
				err,
			)
		}

		reasons = append(reasons, StudentSelectedReasonDTO{
			Reason:          *reason,
			OtherReasonText: r.OtherReasonText,
		})
	}

	return reasons, nil
}

// GetStudentPersonalInfo retrieves detailed personal information for a student.
func (s *Service) GetStudentPersonalInfo(
	ctx context.Context,
	iirID string,
) (*StudentPersonalInfoDTO, error) {
	personalInfo, err := s.repo.GetStudentPersonalInfo(ctx, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student personal info: %w", err)
	}

	gender, err := s.repo.GetGenderByID(ctx, personalInfo.GenderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get gender by ID: %w", err)
	}

	civilStatus, err := s.repo.GetCivilStatusByID(
		ctx,
		personalInfo.CivilStatusID,
	)
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

	emergencyContact, err := s.repo.GetEmergencyContactByIIRID(
		ctx,
		personalInfo.IIRID,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get emergency contact by IIR ID: %w",
			err,
		)
	}

	emergencyContactRelationship, err := s.repo.GetStudentRelationshipByID(
		ctx,
		emergencyContact.RelationshipID,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get emergency contact relationship by ID: %w",
			err,
		)
	}

	emergencyAddressDTO, err := s.locationsSvc.GetAddressByID(
		ctx,
		emergencyContact.AddressID,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get emergency contact address by ID: %w",
			err,
		)
	}

	emergencyContactDTO := EmergencyContactDTO{
		ID:            emergencyContact.ID,
		FirstName:     emergencyContact.FirstName,
		MiddleName:    structs.FromSqlNull(emergencyContact.MiddleName),
		LastName:      emergencyContact.LastName,
		ContactNumber: emergencyContact.ContactNumber,
		Relationship:  *emergencyContactRelationship,
		Address:       emergencyAddressDTO,
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

// GetStudentAddresses retrieves all addresses associated with a student.
func (s *Service) GetStudentAddresses(
	ctx context.Context,
	iirID string,
) ([]StudentAddressDTO, error) {
	studentAddresses, err := s.repo.GetStudentAddresses(ctx, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student addresses: %w", err)
	}

	var addresses []StudentAddressDTO
	for _, addr := range studentAddresses {
		// Use locations service to get detailed address information
		addrDTO, err := s.locationsSvc.GetAddressByID(ctx, addr.AddressID)
		if err != nil {
			return nil, fmt.Errorf("failed to get address by ID: %w", err)
		}

		addresses = append(addresses, StudentAddressDTO{
			ID:          addr.ID,
			Address:     addrDTO,
			AddressType: addr.AddressType,
			CreatedAt:   addr.CreatedAt,
			UpdatedAt:   addr.UpdatedAt,
		})
	}

	return addresses, nil
}

// GetStudentFamilyBackground retrieves the family background information for a
// student.
func (s *Service) GetStudentFamilyBackground(
	ctx context.Context,
	iirID string,
) (*FamilyBackgroundDTO, error) {
	studentFamily, err := s.repo.GetStudentFamilyBackground(ctx, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get family background: %w", err)
	}

	var family *FamilyBackgroundDTO

	parentalStatus, err := s.repo.GetParentalStatusByID(
		ctx,
		studentFamily.ParentalStatusID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get parental status by ID: %w", err)
	}

	natureOfResidence, err := s.repo.GetNatureOfResidenceByID(
		ctx,
		studentFamily.NatureOfResidenceId,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get nature of residence by ID: %w",
			err,
		)
	}

	siblingSupportTypes, err := s.repo.GetStudentSiblingSupport(
		ctx,
		studentFamily.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get student sibling support: %w", err)
	}

	var supportTypes []SibilingSupportType
	for _, sst := range siblingSupportTypes {
		sibilingSupportType, err := s.repo.GetSiblingSupportTypeByID(
			ctx,
			sst.SupportTypeID,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to get sibling support type by ID: %w",
				err,
			)
		}

		supportTypes = append(supportTypes, SibilingSupportType{
			ID:          sibilingSupportType.ID,
			SupportName: sibilingSupportType.SupportName,
		})
	}

	family = &FamilyBackgroundDTO{
		ID:             studentFamily.ID,
		ParentalStatus: *parentalStatus,
		ParentalStatusDetails: structs.FromSqlNull(
			studentFamily.ParentalStatusDetails,
		),
		Brothers:              &studentFamily.Brothers,
		Sisters:               &studentFamily.Sisters,
		EmployedSiblings:      &studentFamily.EmployedSiblings,
		OrdinalPosition:       studentFamily.OrdinalPosition,
		HaveQuietPlaceToStudy: studentFamily.HaveQuietPlaceToStudy,
		IsSharingRoom:         studentFamily.IsSharingRoom,
		SiblingSupportTypes:   supportTypes,
		RoomSharingDetails: structs.FromSqlNull(
			studentFamily.RoomSharingDetails,
		),
		NatureOfResidence: *natureOfResidence,
	}

	return family, nil
}

// GetStudentRelatedPersons retrieves all related persons (parents, guardians,
// etc.) for a student.
func (s *Service) GetStudentRelatedPersons(
	ctx context.Context,
	iirID string,
) ([]RelatedPersonDTO, error) {
	studentRelatedPersons, err := s.repo.GetStudentRelatedPersons(ctx, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student related persons: %w", err)
	}

	var related []RelatedPersonDTO
	for _, srp := range studentRelatedPersons {
		relatedPerson, err := s.repo.GetRelatedPersonByID(
			ctx,
			srp.RelatedPersonID,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to get related person by ID: %w",
				err,
			)
		}

		relationship, err := s.repo.GetStudentRelationshipByID(
			ctx,
			srp.RelationshipID,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to get student relationship by ID: %w",
				err,
			)
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
			EmployerAddress: structs.FromSqlNull(
				relatedPerson.EmployerAddress,
			),
			Relationship: *relationship,
			IsParent:     srp.IsParent,
			IsGuardian:   srp.IsGuardian,
			IsLiving:     srp.IsLiving,
		})

	}

	return related, nil
}

// GetEducationalBackground retrieves the educational history for a student.
func (s *Service) GetEducationalBackground(
	ctx context.Context,
	iirID string,
) (*EducationalBackgroundDTO, error) {
	educationalBackground, err := s.repo.GetStudentEducationalBackground(
		ctx,
		iirID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get educational background: %w", err)
	}

	schools, err := s.repo.GetSchoolDetailsByEBID(ctx, educationalBackground.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get school details by EBID: %w", err)
	}

	var schoolDTOs []SchoolDetailsDTO

	for _, school := range schools {
		educationalLevel, err := s.repo.GetEducationalLevelByID(
			ctx,
			school.EducationalLevelID,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to get educational level by ID: %w",
				err,
			)
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
		ID:                educationalBackground.ID,
		NatureOfSchooling: educationalBackground.NatureOfSchooling,
		InterruptedDetails: structs.FromSqlNull(
			educationalBackground.InterruptedDetails,
		),
		School:    schoolDTOs,
		CreatedAt: educationalBackground.CreatedAt,
		UpdatedAt: educationalBackground.UpdatedAt,
	}, nil
}

// GetStudentFinancialInfo retrieves religious and financial info for a student.
func (s *Service) GetStudentFinancialInfo(
	ctx context.Context,
	iirID string,
) (*StudentFinanceDTO, error) {
	financialInfo, err := s.repo.GetStudentFinancialInfo(ctx, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student financial info: %w", err)
	}

	incomeRange, err := s.repo.GetIncomeRangeByID(
		ctx,
		financialInfo.MonthlyFamilyIncomeRangeID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get income range by ID: %w", err)
	}

	financialSupportTypes, err := s.repo.GetFinancialSupportTypeByFinanceID(
		ctx,
		financialInfo.ID,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get financial support types by finance ID: %w",
			err,
		)
	}

	var supportTypes []StudentSupportType
	for _, fst := range financialSupportTypes {
		supportType, err := s.repo.GetStudentSupportByID(ctx, fst.SupportTypeID)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to get student support type by ID: %w",
				err,
			)
		}
		supportTypes = append(supportTypes, *supportType)
	}

	return &StudentFinanceDTO{
		ID:                       financialInfo.ID,
		MonthlyFamilyIncomeRange: *incomeRange,
		OtherIncomeDetails: structs.FromSqlNull(
			financialInfo.OtherIncomeDetails,
		),
		FinancialSupportTypes: supportTypes,
		WeeklyAllowance:       financialInfo.WeeklyAllowance,
	}, nil
}

// GetStudentHealthRecord retrieves the health record for a student.
func (s *Service) GetStudentHealthRecord(
	ctx context.Context,
	iirID string,
) (*StudentHealthRecordDTO, error) {
	healthRecord, err := s.repo.GetStudentHealthRecord(ctx, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student health record: %w", err)
	}

	return &StudentHealthRecordDTO{
		ID:               healthRecord.ID,
		VisionHasProblem: healthRecord.VisionHasProblem,
		VisionDetails: structs.FromSqlNull(
			healthRecord.VisionDetails,
		),
		HearingHasProblem: healthRecord.HearingHasProblem,
		HearingDetails: structs.FromSqlNull(
			healthRecord.HearingDetails,
		),
		SpeechHasProblem: healthRecord.SpeechHasProblem,
		SpeechDetails: structs.FromSqlNull(
			healthRecord.SpeechDetails,
		),
		GeneralHealthHasProblem: healthRecord.GeneralHealthHasProblem,
		GeneralHealthDetails: structs.FromSqlNull(
			healthRecord.GeneralHealthDetails,
		),
	}, nil
}

// GetStudentConsultations retrieves past consultations for a student.
func (s *Service) GetStudentConsultations(
	ctx context.Context,
	iirID string,
) ([]StudentConsultationDTO, error) {
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

// GetStudentActivities retrieves all extracurricular activities for a student.
func (s *Service) GetStudentActivities(
	ctx context.Context,
	iirID string,
) ([]StudentActivityDTO, error) {
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

// GetStudentSubjectPreferences retrieves academic subject preferences for a
// student.
func (s *Service) GetStudentSubjectPreferences(
	ctx context.Context,
	iirID string,
) ([]StudentSubjectPreferenceDTO, error) {
	preferences, err := s.repo.GetStudentSubjectPreferences(ctx, iirID)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get student subject preferences: %w",
			err,
		)
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

// GetStudentHobbies retrieves hobbies and interests for a student.
func (s *Service) GetStudentHobbies(
	ctx context.Context,
	iirID string,
) ([]StudentHobbyDTO, error) {
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

// GetStudentTestResults retrieves standardized test results for a student.
func (s *Service) GetStudentTestResults(
	ctx context.Context,
	iirID string,
) ([]TestResultDTO, error) {
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

// SaveIIRDraft persists a student's IIR draft data to the datastore.
func (s *Service) SaveIIRDraft(
	ctx context.Context,
	userID string,
	req ComprehensiveProfileDTO,
) (int, error) {
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
		audit.Dispatch(ctx, s.logService, s.notifService, audit.DispatchParams{
			Log: &audit.LogParams{
				Level:    audit.LevelError,
				Category: audit.CategoryAudit,
				Action:   audit.ActionIIRUpdateFailed,
				Message: fmt.Sprintf(
					"Failed to save IIR draft for User #%s",
					userID,
				),
				Metadata: &audit.LogMetadata{
					EntityType: constants.IIREntityType,
					NewValues:  req,
					Error:      err.Error(),
				},
			},
		})
		return 0, fmt.Errorf("failed to save IIR draft: %w", err)
	}

	audit.Dispatch(ctx, s.logService, s.notifService, audit.DispatchParams{
		Log: &audit.LogParams{
			Level:    audit.LevelInfo,
			Category: audit.CategoryAudit,
			Action:   audit.ActionIIRDraftSaved,
			Message: fmt.Sprintf(
				"IIR draft saved for User #%s",
				userID,
			),
			Metadata: &audit.LogMetadata{
				EntityType: constants.IIREntityType,
				NewValues:  req,
			},
		},
	})

	return draftID, nil
}

// SubmitStudentIIR processes and persists a full student IIR submission.
func (s *Service) SubmitStudentIIR(
	ctx context.Context,
	userID string,
	req ComprehensiveProfileDTO,
) (string, error) {
	var iirID string
	err := datastore.RunInTransaction(
		ctx,
		s.repo.GetDB(),
		func(tx datastore.DB) error {
			iirRecord := &IIRRecord{
				ID:          uuid.New().String(),
				UserID:      userID,
				IsSubmitted: false,
			}
			var err error
			iirID, err = s.repo.UpsertIIRRecord(ctx, tx, iirRecord)
			if err != nil {
				return fmt.Errorf("failed to create IIR record: %w", err)
			}

			// Save Student Personal Info
			if err := s.saveStudentPersonalInfo(ctx, tx, iirID, req.Student.StudentPersonalInfoDTO); err != nil {
				return err
			}

			// Save Student Addresses
			if err := s.saveStudentAddresses(ctx, tx, iirID, req.Student.Addresses); err != nil {
				return err
			}

			// Save Educational Background
			if err := s.saveEducationalBackground(ctx, tx, iirID, req.Education); err != nil {
				return err
			}

			// Save Family Background (Background, Sibling Support, and
			// Related Persons)
			if err := s.saveFamilyBackground(ctx, tx, iirID, req); err != nil {
				return err
			}

			// Save Health Record and Consultations
			if err := s.saveStudentHealthRecord(ctx, tx, iirID, req); err != nil {
				return err
			}

			// Save Financial Info
			if err := s.saveStudentFinance(ctx, tx, iirID, req.Family.Finance); err != nil {
				return err
			}

			// Save Interests (Activities, Subject Preferences, Hobbies)
			if err := s.saveStudentInterests(ctx, tx, iirID, req); err != nil {
				return err
			}

			IIRRecordUpdate := &IIRRecord{
				ID:          iirID,
				UserID:      userID,
				IsSubmitted: true,
			}

			if _, err := s.repo.UpsertIIRRecord(ctx, tx, IIRRecordUpdate); err != nil {
				return fmt.Errorf(
					"failed to update IIR record as submitted: %w",
					err,
				)
			}

			return nil
		},
	)
	if err != nil {
		audit.Dispatch(ctx, s.logService, s.notifService, audit.DispatchParams{
			Log: &audit.LogParams{
				Level:    audit.LevelError,
				Category: audit.CategoryAudit,
				Action:   audit.ActionIIRCreateFailed,
				Message: fmt.Sprintf(
					"Failed to submit IIR for User #%s",
					userID,
				),
				Metadata: &audit.LogMetadata{
					EntityType: constants.IIREntityType,
					NewValues:  req,
					Error:      err.Error(),
				},
			},
		})
		return "", err
	}

	// Fetch personalized notification targets
	student, _ := s.userService.GetUserByID(ctx, userID)
	studentName := "A student"
	if student != nil {
		studentName = fmt.Sprintf("%s %s", student.FirstName, student.LastName)
	}

	counselorIDs, _ := s.userService.GetUserIDsByRole(
		ctx,
		int(constants.AdminRoleID),
	)

	notifications := []audit.NotificationParams{
		{
			ReceiverID: structs.StringToNullableString(userID),
			TargetID:   structs.StringToNullableString(iirID),
			TargetType: structs.StringToNullableString(constants.IIREntityType),
			Title:      "IIR Submitted Successfully",
			Message: fmt.Sprintf(
				"Your Individual Inventory Record (#%s) has been submitted.",
				iirID,
			),
			Type: constants.IIREntityType,
		},
	}

	for _, cid := range counselorIDs {
		notifications = append(notifications, audit.NotificationParams{
			ReceiverID: structs.StringToNullableString(cid),
			TargetID:   structs.StringToNullableString(iirID),
			TargetType: structs.StringToNullableString(constants.IIREntityType),
			Title:      "New IIR Submission",
			Message: fmt.Sprintf(
				"Student %s has submitted their IIR (#%s).",
				studentName,
				iirID,
			),
			Type: constants.IIREntityType,
		})
	}

	audit.Dispatch(ctx, s.logService, s.notifService, audit.DispatchParams{
		Log: &audit.LogParams{
			Level:    audit.LevelInfo,
			Category: audit.CategoryAudit,
			Action:   audit.ActionIIRSubmitted,
			Message:  fmt.Sprintf("IIR #%s submitted", iirID),
			TargetID: structs.StringToNullableString(iirID),
			TargetType: structs.StringToNullableString(
				constants.IIREntityType,
			),
			Metadata: &audit.LogMetadata{
				EntityType: constants.IIREntityType,
				EntityID:   iirID,
				NewValues:  req,
			},
		},
		Notifications: notifications,
	})

	return iirID, nil
}

func (s *Service) saveStudentPersonalInfo(
	ctx context.Context,
	tx datastore.DB,
	iirID string,
	dto StudentPersonalInfoDTO,
) error {
	if err := s.repo.UpsertStudentPersonalInfo(ctx, tx, &StudentPersonalInfo{
		IIRID: iirID,
		StudentNumber: dto.StudentNumber,
		GenderID:      dto.Gender.ID,
		CivilStatusID: dto.CivilStatus.ID,
		ReligionID:    dto.Religion.ID,
		HeightFt:      dto.HeightFt,
		WeightKg:      dto.WeightKg,
		Complexion:    dto.Complexion,
		HighSchoolGWA: dto.HighSchoolGWA,
		CourseID:      dto.Course.ID,
		YearLevel:     dto.YearLevel,
		Section:       dto.Section,
		PlaceOfBirth:  dto.PlaceOfBirth,
		DateOfBirth:   dto.DateOfBirth,
		IsEmployed:    dto.IsEmployed,
		EmployerName: structs.ToSqlNull(
			dto.EmployerName,
		),
		EmployerAddress: structs.ToSqlNull(
			dto.EmployerAddress,
		),
		MobileNumber: dto.MobileNumber,
		TelephoneNumber: structs.ToSqlNull(
			dto.TelephoneNumber,
		),
	}); err != nil {
		return fmt.Errorf("failed to upsert student personal info: %w", err)
	}

	ec := dto.EmergencyContact
	var ecProvinceCode *string
	if ec.Address.Province != nil {
		ecProvinceCode = &ec.Address.Province.Code
	}
	addressID, err := s.locationsSvc.SaveAddress(
		ctx,
		tx,
		&locations.Address{
			RegionCode:   ec.Address.Region.Code,
			ProvinceCode: ecProvinceCode,
			CityCode:     ec.Address.City.Code,
			BarangayCode: ec.Address.Barangay.Code,
			StreetDetail: &ec.Address.StreetDetail,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to upsert emergency contact address: %w", err)
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
		return fmt.Errorf("failed to upsert emergency contact: %w", err)
	}

	return nil
}

func (s *Service) saveStudentAddresses(
	ctx context.Context,
	tx datastore.DB,
	iirID string,
	dtos []StudentAddressDTO,
) error {
	for _, addrDTO := range dtos {
		var addrProvinceCode *string
		if addrDTO.Address.Province != nil {
			addrProvinceCode = &addrDTO.Address.Province.Code
		}
		addressID, err := s.locationsSvc.SaveAddress(
			ctx,
			tx,
			&locations.Address{
				RegionCode:   addrDTO.Address.Region.Code,
				ProvinceCode: addrProvinceCode,
				CityCode:     addrDTO.Address.City.Code,
				BarangayCode: addrDTO.Address.Barangay.Code,
				StreetDetail: &addrDTO.Address.StreetDetail,
			},
		)
		if err != nil {
			return fmt.Errorf("failed to upsert student address: %w", err)
		}

		if _, err := s.repo.UpsertStudentAddress(ctx, tx, &StudentAddress{
			IIRID:       iirID,
			AddressID:   addressID,
			AddressType: addrDTO.AddressType,
		}); err != nil {
			return fmt.Errorf(
				"failed to save student address relation: %w",
				err,
			)
		}
	}
	return nil
}

func (s *Service) saveEducationalBackground(
	ctx context.Context,
	tx datastore.DB,
	iirID string,
	dto EducationalBackgroundDTO,
) error {
	ebID, err := s.repo.UpsertEducationalBackground(
		ctx,
		tx,
		&EducationalBackground{
			IIRID:             iirID,
			NatureOfSchooling: dto.NatureOfSchooling,
			InterruptedDetails: structs.ToSqlNull(
				dto.InterruptedDetails,
			),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to upsert educational background: %w", err)
	}

	if err := s.repo.DeleteSchoolDetailsByEBID(ctx, tx, ebID); err != nil {
		return fmt.Errorf("failed to delete existing school details: %w", err)
	}

	for _, schoolDTO := range dto.School {
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
			return fmt.Errorf("failed to save school details: %w", err)
		}
	}
	return nil
}

func (s *Service) saveFamilyBackground(
	ctx context.Context,
	tx datastore.DB,
	iirID string,
	dto ComprehensiveProfileDTO,
) error {
	fbID, err := s.repo.UpsertFamilyBackground(
		ctx,
		tx,
		&FamilyBackground{
			IIRID:            iirID,
			ParentalStatusID: dto.Family.FamilyBackgroundDTO.ParentalStatus.ID,
			ParentalStatusDetails: structs.ToSqlNull(
				dto.Family.FamilyBackgroundDTO.ParentalStatusDetails,
			),
			Brothers:              *dto.Family.FamilyBackgroundDTO.Brothers,
			Sisters:               *dto.Family.FamilyBackgroundDTO.Sisters,
			EmployedSiblings:      *dto.Family.FamilyBackgroundDTO.EmployedSiblings,
			OrdinalPosition:       dto.Family.FamilyBackgroundDTO.OrdinalPosition,
			HaveQuietPlaceToStudy: dto.Family.FamilyBackgroundDTO.HaveQuietPlaceToStudy,
			IsSharingRoom:         dto.Family.FamilyBackgroundDTO.IsSharingRoom,
			RoomSharingDetails: structs.ToSqlNull(
				dto.Family.FamilyBackgroundDTO.RoomSharingDetails,
			),
			NatureOfResidenceId: dto.Family.FamilyBackgroundDTO.NatureOfResidence.ID,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to upsert family background: %w", err)
	}

	// Save Sibling Supports
	if err := s.repo.DeleteStudentSiblingSupportsByFamilyID(
		ctx,
		tx,
		fbID,
	); err != nil {
		return fmt.Errorf("failed to delete existing sibling supports: %w", err)
	}

	for _, supportType := range dto.Family.FamilyBackgroundDTO.SiblingSupportTypes {
		if err := s.repo.CreateStudentSiblingSupport(ctx, tx, &StudentSiblingSupport{
			FamilyBackgroundID: fbID,
			SupportTypeID:      supportType.ID,
		}); err != nil {
			return fmt.Errorf("failed to save sibling support: %w", err)
		}
	}

	// Save Related Persons
	if err := s.repo.DeleteStudentRelatedPersons(ctx, tx, iirID); err != nil {
		return fmt.Errorf("failed to delete existing related persons: %w", err)
	}

	for _, relPersonDTO := range dto.Family.RelatedPersons {
		relPersonID, err := s.repo.UpsertRelatedPerson(
			ctx,
			tx,
			&RelatedPerson{
				FirstName: relPersonDTO.FirstName,
				LastName:  relPersonDTO.LastName,
				MiddleName: structs.ToSqlNull(
					relPersonDTO.MiddleName,
				),
				DateOfBirth:      relPersonDTO.DateOfBirth,
				EducationalLevel: relPersonDTO.EducationalLevel,
				Occupation: structs.ToSqlNull(
					relPersonDTO.Occupation,
				),
				EmployerName: structs.ToSqlNull(
					relPersonDTO.EmployerName,
				),
				EmployerAddress: structs.ToSqlNull(
					relPersonDTO.EmployerAddress,
				),
			},
		)
		if err != nil {
			return fmt.Errorf("failed to save related person: %w", err)
		}

		if err := s.repo.UpsertStudentRelatedPerson(
			ctx,
			tx,
			&StudentRelatedPerson{
				IIRID:           iirID,
				RelatedPersonID: relPersonID,
				RelationshipID:  relPersonDTO.Relationship.ID,
				IsParent:        relPersonDTO.IsParent,
				IsGuardian:      relPersonDTO.IsGuardian,
				IsLiving:        relPersonDTO.IsLiving,
			},
		); err != nil {
			return fmt.Errorf(
				"failed to save student related person relation: %w",
				err,
			)
		}
	}

	return nil
}

func (s *Service) saveStudentHealthRecord(
	ctx context.Context,
	tx datastore.DB,
	iirID string,
	dto ComprehensiveProfileDTO,
) error {
	if _, err := s.repo.UpsertStudentHealthRecord(
		ctx,
		tx,
		&StudentHealthRecord{
			IIRID:                   iirID,
			VisionHasProblem:        dto.Health.StudentHealthRecordDTO.VisionHasProblem,
			VisionDetails:           structs.ToSqlNull(dto.Health.StudentHealthRecordDTO.VisionDetails),
			HearingHasProblem:       dto.Health.StudentHealthRecordDTO.HearingHasProblem,
			HearingDetails:          structs.ToSqlNull(dto.Health.StudentHealthRecordDTO.HearingDetails),
			SpeechHasProblem:        dto.Health.StudentHealthRecordDTO.SpeechHasProblem,
			SpeechDetails:           structs.ToSqlNull(dto.Health.StudentHealthRecordDTO.SpeechDetails),
			GeneralHealthHasProblem: dto.Health.StudentHealthRecordDTO.GeneralHealthHasProblem,
			GeneralHealthDetails:    structs.ToSqlNull(dto.Health.StudentHealthRecordDTO.GeneralHealthDetails),
		},
	); err != nil {
		return fmt.Errorf("failed to upsert student health record: %w", err)
	}

	// Save Consultations
	for _, consultationDTO := range dto.Health.Consultations {
		if _, err := s.repo.UpsertStudentConsultation(ctx, tx, &StudentConsultation{
			IIRID:            iirID,
			ProfessionalType: consultationDTO.ProfessionalType,
			HasConsulted:     consultationDTO.HasConsulted,
			WhenDate:         structs.ToSqlNull(consultationDTO.WhenDate),
			ForWhat:          structs.ToSqlNull(consultationDTO.ForWhat),
		}); err != nil {
			return fmt.Errorf("failed to save student consultation: %w", err)
		}
	}

	return nil
}

func (s *Service) saveStudentFinance(
	ctx context.Context,
	tx datastore.DB,
	iirID string,
	dto StudentFinanceDTO,
) error {
	sfID, err := s.repo.UpsertStudentFinance(
		ctx,
		tx,
		&StudentFinance{
			IIRID:                      iirID,
			MonthlyFamilyIncomeRangeID: dto.MonthlyFamilyIncomeRange.ID,
			OtherIncomeDetails: structs.ToSqlNull(
				dto.OtherIncomeDetails,
			),
			WeeklyAllowance: dto.WeeklyAllowance,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to upsert student finance: %w", err)
	}

	if err := s.repo.DeleteStudentFinancialSupportsByFinanceID(
		ctx,
		tx,
		sfID,
	); err != nil {
		return fmt.Errorf(
			"failed to delete existing financial supports: %w",
			err,
		)
	}

	for _, supportType := range dto.FinancialSupportTypes {
		if err := s.repo.CreateStudentFinancialSupport(
			ctx,
			tx,
			&StudentFinancialSupport{
				StudentFinanceID: sfID,
				SupportTypeID:    supportType.ID,
			},
		); err != nil {
			return fmt.Errorf(
				"failed to save student financial support: %w",
				err,
			)
		}
	}

	return nil
}

func (s *Service) saveStudentInterests(
	ctx context.Context,
	tx datastore.DB,
	iirID string,
	dto ComprehensiveProfileDTO,
) error {
	// Save Activities
	if err := s.repo.DeleteStudentActivitiesByIIRID(
		ctx,
		tx,
		iirID,
	); err != nil {
		return fmt.Errorf("failed to delete existing activities: %w", err)
	}

	for _, activityDTO := range dto.Interests.Activities {
		if _, err := s.repo.CreateStudentActivity(ctx, tx, &StudentActivity{
			IIRID:              iirID,
			OptionID:           activityDTO.ActivityOption.ID,
			OtherSpecification: structs.ToSqlNull(activityDTO.OtherSpecification),
			Role:               activityDTO.Role,
			RoleSpecification:  structs.ToSqlNull(activityDTO.RoleSpecification),
		}); err != nil {
			return fmt.Errorf("failed to save student activity: %w", err)
		}
	}

	// Save Subject Preferences
	if err := s.repo.DeleteStudentSubjectPreferencesByIIRID(
		ctx,
		tx,
		iirID,
	); err != nil {
		return fmt.Errorf(
			"failed to delete existing subject preferences: %w",
			err,
		)
	}

	for _, prefDTO := range dto.Interests.SubjectPreferences {
		if _, err := s.repo.CreateStudentSubjectPreference(
			ctx,
			tx,
			&StudentSubjectPreference{
				IIRID:       iirID,
				SubjectName: prefDTO.SubjectName,
				IsFavorite:  prefDTO.IsFavorite,
			},
		); err != nil {
			return fmt.Errorf(
				"failed to save student subject preference: %w",
				err,
			)
		}
	}

	// Save Hobbies
	if err := s.repo.DeleteStudentHobbiesByIIRID(ctx, tx, iirID); err != nil {
		return fmt.Errorf("failed to delete existing hobbies: %w", err)
	}

	for _, hobbyDTO := range dto.Interests.Hobbies {
		if _, err := s.repo.CreateStudentHobby(ctx, tx, &StudentHobby{
			IIRID:        iirID,
			HobbyName:    hobbyDTO.HobbyName,
			PriorityRank: hobbyDTO.PriorityRank,
		}); err != nil {
			return fmt.Errorf("failed to save student hobby: %w", err)
		}
	}

	return nil
}

// GenerateIIR generates a student's IIR as a PDF using an HTML template and
// Gotenberg.
func (s *Service) GenerateIIR(
	ctx context.Context,
	iirID string,
) ([]byte, string, error) {
	profile, err := s.GetStudentProfile(ctx, iirID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get student profile: %w", err)
	}

	templatePath := filepath.Join(
		"internal",
		"features",
		"students",
		"assets",
		"iirf.html",
	)

	data := struct {
		Profile   *ComprehensiveProfileDTO
		DateToday string
	}{
		Profile:   profile,
		DateToday: s.GetFormattedDate(time.Now().Format("2006-01-02")),
	}

	pdfBytes, err := s.pdfService.GenerateFromTemplate(ctx, templatePath, data)
	if err != nil {
		return nil, "", fmt.Errorf(
			"failed to generate pdf from template: %w",
			err,
		)
	}

	fileName := fmt.Sprintf(
		"IIR_%s_%s.pdf",
		profile.Student.StudentNumber,
		time.Now().Format("20060102"),
	)

	return pdfBytes, fileName, nil
}

func (s *Service) GetFormattedDate(date string) string {
	const inputLayout = "2006-01-02"
	const outputLayout = "January 02, 2006"
	const emptyString = ""

	t, err := time.Parse(inputLayout, date)
	if err != nil {
		// [HandlerName] {Specific Step}: error message
		fmt.Printf("[GetFormattedDOB] {Date Parsing}: %v\n", err)
		return emptyString
	}

	return t.Format(outputLayout)
}

func (s *Service) getAge(dateOfBirth string) int {
	today := time.Now()
	bday, err := time.Parse("2006-01-02", dateOfBirth)
	if err != nil {
		return 0
	}

	return today.Year() - bday.Year()
}

func (s *Service) getFullAddress(address locations.AddressDTO) string {
	return fmt.Sprintf(
		"%s, %s, %s, %s, %s",
		address.StreetDetail,
		address.Barangay.Name,
		address.City.Name,
		address.Province.Name,
		address.Region.Name,
	)
}

func (s *Service) getFullName(firstName, middleName, lastName string) string {
	return fmt.Sprintf("%s, %s %s", lastName, firstName, middleName)
}

// truncate returns a string limited to n characters.
func (s *Service) truncate(str string, n int) string {
	if len(str) <= n {
		return str
	}
	return str[:n]
}

// ptrIntToStr safely converts a *int to a string.
func (s *Service) ptrIntToStr(i *int) string {
	if i == nil {
		return "0"
	}
	return fmt.Sprintf("%d", *i)
}
