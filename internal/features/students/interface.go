package students

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

type ServiceInterface interface {
	GetGenders(ctx context.Context) ([]Gender, error)
	GetParentalStatusTypes(ctx context.Context) ([]ParentalStatusType, error)
	GetEnrollmentReasons(ctx context.Context) ([]EnrollmentReason, error)
	GetIncomeRanges(ctx context.Context) ([]IncomeRange, error)
	GetStudentSupportTypes(ctx context.Context) ([]StudentSupportType, error)
	GetSiblingSupportTypes(ctx context.Context) ([]SibilingSupportType, error)
	GetEducationalLevels(ctx context.Context) ([]EducationalLevel, error)
	GetCourses(ctx context.Context) ([]Course, error)
	GetCivilStatusTypes(ctx context.Context) ([]CivilStatusType, error)
	GetReligions(ctx context.Context) ([]Religion, error)
	GetNatureOfResidenceTypes(
		ctx context.Context,
	) ([]NatureOfResidenceType, error)
	GetActivityOptions(ctx context.Context) ([]ActivityOption, error)
	GetStudentRelationshipTypes(
		ctx context.Context,
	) ([]StudentRelationshipType, error)
	ListStudents(
		ctx context.Context,
		req ListStudentsRequest,
	) (*ListStudentsResponse, error)
	GetStudentProfile(
		ctx context.Context,
		iirID string,
	) (*ComprehensiveProfileDTO, error)
	GetStudentBasicInfo(
		ctx context.Context,
		iirID string,
	) (*StudentBasicInfoViewDTO, error)
	GetIIRDraft(
		ctx context.Context,
		userID string,
	) (*ComprehensiveProfileDTO, error)
	GetStudentIIRByUserID(
		ctx context.Context,
		userID string,
	) (*IIRRecord, error)
	GetStudentIIR(ctx context.Context, iirID string) (*IIRRecord, error)
	GetStudentEnrollmentReasons(
		ctx context.Context,
		iirID string,
	) ([]StudentSelectedReasonDTO, error)
	GetStudentPersonalInfo(
		ctx context.Context,
		iirID string,
	) (*StudentPersonalInfoDTO, error)
	GetStudentAddresses(
		ctx context.Context,
		iirID string,
	) ([]StudentAddressDTO, error)
	GetStudentFamilyBackground(
		ctx context.Context,
		iirID string,
	) (*FamilyBackgroundDTO, error)
	GetStudentRelatedPersons(
		ctx context.Context,
		iirID string,
	) ([]RelatedPersonDTO, error)
	GetEducationalBackground(
		ctx context.Context,
		iirID string,
	) (*EducationalBackgroundDTO, error)
	GetStudentFinancialInfo(
		ctx context.Context,
		iirID string,
	) (*StudentFinanceDTO, error)
	GetStudentHealthRecord(
		ctx context.Context,
		iirID string,
	) (*StudentHealthRecordDTO, error)
	GetStudentConsultations(
		ctx context.Context,
		iirID string,
	) ([]StudentConsultationDTO, error)
	GetStudentActivities(
		ctx context.Context,
		iirID string,
	) ([]StudentActivityDTO, error)
	GetStudentSubjectPreferences(
		ctx context.Context,
		iirID string,
	) ([]StudentSubjectPreferenceDTO, error)
	GetStudentHobbies(
		ctx context.Context,
		iirID string,
	) ([]StudentHobbyDTO, error)
	GetStudentTestResults(
		ctx context.Context,
		iirID string,
	) ([]TestResultDTO, error)
	SaveIIRDraft(
		ctx context.Context,
		userID string,
		req ComprehensiveProfileDTO,
	) (int, error)
	SubmitStudentIIR(
		ctx context.Context,
		userID string,
		req ComprehensiveProfileDTO,
	) (string, error)
	GenerateIIR(ctx context.Context, iirID string) ([]byte, string, error)
}

type RepositoryInterface interface {
	BeginTx(ctx context.Context) (datastore.DB, error)
	GetDB() *sqlx.DB
	GetGenders(ctx context.Context) ([]Gender, error)
	GetParentalStatusTypes(ctx context.Context) ([]ParentalStatusType, error)
	GetEnrollmentReasons(ctx context.Context) ([]EnrollmentReason, error)
	GetIncomeRanges(ctx context.Context) ([]IncomeRange, error)
	GetStudentSupportTypes(ctx context.Context) ([]StudentSupportType, error)
	GetSiblingSupportTypes(ctx context.Context) ([]SibilingSupportType, error)
	GetEducationalLevels(ctx context.Context) ([]EducationalLevel, error)
	GetCourses(ctx context.Context) ([]Course, error)
	GetCivilStatusTypes(ctx context.Context) ([]CivilStatusType, error)
	GetReligions(ctx context.Context) ([]Religion, error)
	GetStudentRelationshipTypes(
		ctx context.Context,
	) ([]StudentRelationshipType, error)
	GetNatureOfResidenceTypes(
		ctx context.Context,
	) ([]NatureOfResidenceType, error)
	GetTotalStudentsCount(
		ctx context.Context,
		search string,
		courseID int,
		genderID int,
		yearLevel int,
	) (int, error)
	ListStudents(
		ctx context.Context,
		search string,
		offset int,
		limit int,
		orderBy string,
		courseID int,
		genderID int,
		yearLevel int,
	) ([]StudentProfileView, error)
	GetStudentBasicInfo(
		ctx context.Context,
		iirID string,
	) (*StudentBasicInfoView, error)
	GetIIRDraftByUserID(ctx context.Context, userID string) (*IIRDraft, error)
	GetStudentIIRByUserID(
		ctx context.Context,
		userID string,
	) (*IIRRecord, error)
	GetStudentIIR(ctx context.Context, iirID string) (*IIRRecord, error)
	GetStudentEnrollmentReasons(
		ctx context.Context,
		iirID string,
	) ([]StudentSelectedReason, error)
	GetEnrollmentReasonByID(
		ctx context.Context,
		reasonID int,
	) (*EnrollmentReason, error)
	GetStudentPersonalInfo(
		ctx context.Context,
		iirID string,
	) (*StudentPersonalInfo, error)
	GetEmergencyContactByIIRID(
		ctx context.Context,
		iirID string,
	) (*EmergencyContact, error)
	GetGenderByID(ctx context.Context, genderID int) (*Gender, error)
	GetCivilStatusByID(
		ctx context.Context,
		statusID int,
	) (*CivilStatusType, error)
	GetReligionByID(ctx context.Context, religionID int) (*Religion, error)
	GetCourseByID(ctx context.Context, courseID int) (*Course, error)
	GetStudentAddresses(
		ctx context.Context,
		iirID string,
	) ([]StudentAddress, error)
	GetStudentEducationalBackground(
		ctx context.Context,
		iirID string,
	) (*EducationalBackground, error)
	GetSchoolDetailsByEBID(
		ctx context.Context,
		ebID int,
	) ([]SchoolDetails, error)
	GetEducationalLevelByID(
		ctx context.Context,
		levelID int,
	) (*EducationalLevel, error)
	GetStudentRelatedPersons(
		ctx context.Context,
		iirID string,
	) ([]StudentRelatedPerson, error)
	GetRelatedPersonByID(
		ctx context.Context,
		personID int,
	) (*RelatedPerson, error)
	GetStudentRelationshipByID(
		ctx context.Context,
		relationshipID int,
	) (*StudentRelationshipType, error)
	GetStudentFamilyBackground(
		ctx context.Context,
		iirID string,
	) (*FamilyBackground, error)
	GetParentalStatusByID(
		ctx context.Context,
		statusID int,
	) (*ParentalStatusType, error)
	GetNatureOfResidenceByID(
		ctx context.Context,
		residenceID int,
	) (*NatureOfResidenceType, error)
	GetStudentSiblingSupport(
		ctx context.Context,
		fbID int,
	) ([]StudentSiblingSupport, error)
	GetSiblingSupportTypeByID(
		ctx context.Context,
		supportID int,
	) (*SibilingSupportType, error)
	GetStudentFinancialInfo(
		ctx context.Context,
		iirID string,
	) (*StudentFinance, error)
	GetFinancialSupportTypeByFinanceID(
		ctx context.Context,
		financeID int,
	) ([]StudentFinancialSupport, error)
	GetIncomeRangeByID(ctx context.Context, rangeID int) (*IncomeRange, error)
	GetStudentSupportByID(
		ctx context.Context,
		supportID int,
	) (*StudentSupportType, error)
	GetStudentHealthRecord(
		ctx context.Context,
		iirID string,
	) (*StudentHealthRecord, error)
	GetActivityOptions(ctx context.Context) ([]ActivityOption, error)
	GetStudentConsultations(
		ctx context.Context,
		iirID string,
	) ([]StudentConsultation, error)
	GetStudentActivities(
		ctx context.Context,
		iirID string,
	) ([]StudentActivity, error)
	GetActivityOptionByID(
		ctx context.Context,
		optionID int,
	) (*ActivityOption, error)
	GetStudentSubjectPreferences(
		ctx context.Context,
		iirID string,
	) ([]StudentSubjectPreference, error)
	GetStudentHobbies(ctx context.Context, iirID string) ([]StudentHobby, error)
	GetStudentTestResults(
		ctx context.Context,
		iirID string,
	) ([]TestResult, error)
	UpsertIIRDraft(ctx context.Context, draft IIRDraft) (int, error)
	UpsertIIRRecord(
		ctx context.Context,
		tx datastore.DB,
		iir *IIRRecord,
	) (string, error)
	UpsertStudentPersonalInfo(
		ctx context.Context,
		tx datastore.DB,
		info *StudentPersonalInfo,
	) error
	UpsertEmergencyContact(
		ctx context.Context,
		tx datastore.DB,
		ec *EmergencyContact,
	) (int, error)
	UpsertStudentAddress(
		ctx context.Context,
		tx datastore.DB,
		sa *StudentAddress,
	) (int, error)
	CreateStudentSelectedReason(
		ctx context.Context,
		tx datastore.DB,
		ssr *StudentSelectedReason,
	) error
	DeleteStudentSelectedReasons(
		ctx context.Context,
		tx datastore.DB,
		iirID string,
	) error
	UpsertRelatedPerson(
		ctx context.Context,
		tx datastore.DB,
		rp *RelatedPerson,
	) (int, error)
	UpsertStudentRelatedPerson(
		ctx context.Context,
		tx datastore.DB,
		srp *StudentRelatedPerson,
	) error
	DeleteStudentRelatedPersons(
		ctx context.Context,
		tx datastore.DB,
		iirID string,
	) error
	UpsertFamilyBackground(
		ctx context.Context,
		tx datastore.DB,
		fb *FamilyBackground,
	) (int, error)
	CreateStudentSiblingSupport(
		ctx context.Context,
		tx datastore.DB,
		sss *StudentSiblingSupport,
	) error
	DeleteStudentSiblingSupportsByFamilyID(
		ctx context.Context,
		tx datastore.DB,
		familyBackgroundID int,
	) error
	UpsertEducationalBackground(
		ctx context.Context,
		tx datastore.DB,
		eb *EducationalBackground,
	) (int, error)
	UpsertSchoolDetails(
		ctx context.Context,
		tx datastore.DB,
		sd *SchoolDetails,
	) (int, error)
	DeleteSchoolDetailsByEBID(
		ctx context.Context,
		tx datastore.DB,
		ebID int,
	) error
	UpsertStudentHealthRecord(
		ctx context.Context,
		tx datastore.DB,
		hr *StudentHealthRecord,
	) (int, error)
	UpsertStudentConsultation(
		ctx context.Context,
		tx datastore.DB,
		sc *StudentConsultation,
	) (int, error)
	UpsertStudentFinance(
		ctx context.Context,
		tx datastore.DB,
		sf *StudentFinance,
	) (int, error)
	CreateStudentFinancialSupport(
		ctx context.Context,
		tx datastore.DB,
		sfs *StudentFinancialSupport,
	) error
	DeleteStudentFinancialSupportsByFinanceID(
		ctx context.Context,
		tx datastore.DB,
		financeID int,
	) error
	CreateStudentActivity(
		ctx context.Context,
		tx datastore.DB,
		sa *StudentActivity,
	) (int, error)
	DeleteStudentActivitiesByIIRID(
		ctx context.Context,
		tx datastore.DB,
		iirID string,
	) error
	CreateStudentSubjectPreference(
		ctx context.Context,
		tx datastore.DB,
		ssp *StudentSubjectPreference,
	) (int, error)
	DeleteStudentSubjectPreferencesByIIRID(
		ctx context.Context,
		tx datastore.DB,
		iirID string,
	) error
	CreateStudentHobby(
		ctx context.Context,
		tx datastore.DB,
		sh *StudentHobby,
	) (int, error)
	DeleteStudentHobbiesByIIRID(
		ctx context.Context,
		tx datastore.DB,
		iirID string,
	) error
	DeleteTestResultsByIIRID(
		ctx context.Context,
		tx datastore.DB,
		iirID string,
	) error
	DeleteSignificantNotesByIIRID(
		ctx context.Context,
		tx datastore.DB,
		iirID string,
	) error
}
