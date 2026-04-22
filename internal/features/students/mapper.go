package students

import (
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
)

// Gender mappers
func (m GenderDB) ToDomain() Gender {
	return Gender{
		ID:   m.ID,
		Name: m.GenderName,
	}
}

func (d Gender) ToPersistence() GenderDB {
	return GenderDB{
		ID:         d.ID,
		GenderName: d.Name,
	}
}

// ParentalStatusType mappers
func (m ParentalStatusTypeDB) ToDomain() ParentalStatusType {
	return ParentalStatusType{
		ID:   m.ID,
		Name: m.StatusName,
	}
}

func (d ParentalStatusType) ToPersistence() ParentalStatusTypeDB {
	return ParentalStatusTypeDB{
		ID:         d.ID,
		StatusName: d.Name,
	}
}

// StudentSupportType mappers
func (m StudentSupportTypeDB) ToDomain() StudentSupportType {
	return StudentSupportType{
		ID:   m.ID,
		Name: m.SupportTypeName,
	}
}

func (d StudentSupportType) ToPersistence() StudentSupportTypeDB {
	return StudentSupportTypeDB{
		ID:              d.ID,
		SupportTypeName: d.Name,
	}
}

// EnrollmentReason mappers
func (m EnrollmentReasonDB) ToDomain() EnrollmentReason {
	return EnrollmentReason(m)
}

func (d EnrollmentReason) ToPersistence() EnrollmentReasonDB {
	return EnrollmentReasonDB(d)
}

// IncomeRange mappers
func (m IncomeRangeDB) ToDomain() IncomeRange {
	return IncomeRange{
		ID:   m.ID,
		Text: m.RangeText,
	}
}

func (d IncomeRange) ToPersistence() IncomeRangeDB {
	return IncomeRangeDB{
		ID:        d.ID,
		RangeText: d.Text,
	}
}

// EducationalLevel mappers
func (m EducationalLevelDB) ToDomain() EducationalLevel {
	return EducationalLevel{
		ID:   m.ID,
		Name: m.LevelName,
	}
}

func (d EducationalLevel) ToPersistence() EducationalLevelDB {
	return EducationalLevelDB{
		ID:        d.ID,
		LevelName: d.Name,
	}
}

// Course mappers
func (m CourseDB) ToDomain() Course {
	return Course{
		ID:   m.ID,
		Code: m.Code,
		Name: m.CourseName,
	}
}

func (d Course) ToPersistence() CourseDB {
	return CourseDB{
		ID:         d.ID,
		Code:       d.Code,
		CourseName: d.Name,
	}
}

// CivilStatusType mappers
func (m CivilStatusTypeDB) ToDomain() CivilStatusType {
	return CivilStatusType{
		ID:   m.ID,
		Name: m.StatusName,
	}
}

func (d CivilStatusType) ToPersistence() CivilStatusTypeDB {
	return CivilStatusTypeDB{
		ID:         d.ID,
		StatusName: d.Name,
	}
}

// Religion mappers
func (m ReligionDB) ToDomain() Religion {
	return Religion{
		ID:   m.ID,
		Name: m.ReligionName,
	}
}

func (d Religion) ToPersistence() ReligionDB {
	return ReligionDB{
		ID:           d.ID,
		ReligionName: d.Name,
	}
}

// StudentRelationshipType mappers
func (m StudentRelationshipTypeDB) ToDomain() StudentRelationshipType {
	return StudentRelationshipType{
		ID:   m.ID,
		Name: m.RelationshipName,
	}
}

func (d StudentRelationshipType) ToPersistence() StudentRelationshipTypeDB {
	return StudentRelationshipTypeDB{
		ID:               d.ID,
		RelationshipName: d.Name,
	}
}

// NatureOfResidenceType mappers
func (m NatureOfResidenceTypeDB) ToDomain() NatureOfResidenceType {
	return NatureOfResidenceType{
		ID:   m.ID,
		Name: m.ResidenceTypeName,
	}
}

func (d NatureOfResidenceType) ToPersistence() NatureOfResidenceTypeDB {
	return NatureOfResidenceTypeDB{
		ID:                d.ID,
		ResidenceTypeName: d.Name,
	}
}

// SibilingSupportType mappers
func (m SibilingSupportTypeDB) ToDomain() SibilingSupportType {
	return SibilingSupportType{
		ID:   m.ID,
		Name: m.SupportName,
	}
}

func (d SibilingSupportType) ToPersistence() SibilingSupportTypeDB {
	return SibilingSupportTypeDB{
		ID:          d.ID,
		SupportName: d.Name,
	}
}

// StudentStatus mappers
func (m StudentStatusDB) ToDomain() StudentStatus {
	return StudentStatus{
		ID:   m.ID,
		Name: m.StatusName,
	}
}

func (d StudentStatus) ToPersistence() StudentStatusDB {
	return StudentStatusDB{
		ID:         d.ID,
		StatusName: d.Name,
	}
}

// ActivityOption mappers
func (m ActivityOptionDB) ToDomain() ActivityOption {
	return ActivityOption(m)
}

func (d ActivityOption) ToPersistence() ActivityOptionDB {
	return ActivityOptionDB(d)
}

// StudentBasicInfoView mappers
func (m StudentBasicInfoViewDB) ToDomain() StudentBasicInfoView {
	return StudentBasicInfoView{
		UserID:     m.UserID,
		Email:      m.Email,
		FirstName:  m.FirstName,
		MiddleName: structs.FromSqlNull(m.MiddleName),
		LastName:   m.LastName,
		SuffixName: structs.FromSqlNull(m.SuffixName),
	}
}

// StudentProfileView mappers
func (m StudentProfileViewDB) ToDomain() StudentProfileView {
	return StudentProfileView{
		IIRID:         m.IIRID,
		UserID:        m.UserID,
		FirstName:     m.FirstName,
		MiddleName:    structs.FromSqlNull(m.MiddleName),
		LastName:      m.LastName,
		SuffixName:    structs.FromSqlNull(m.SuffixName),
		Email:         m.Email,
		StudentNumber: m.StudentNumber,
		GenderID:      m.GenderID,
		CourseID:      m.CourseID,
		Section:       m.Section,
		YearLevel:     m.YearLevel,
		StatusID:      m.StatusID,
		StatusName:    m.StatusName,
	}
}

// IIRDraft mappers
func (m IIRDraftDB) ToDomain() IIRDraft {
	return IIRDraft(m)
}

func (d IIRDraft) ToPersistence() IIRDraftDB {
	return IIRDraftDB(d)
}

// IIRRecord mappers
func (m IIRRecordDB) ToDomain() IIRRecord {
	return IIRRecord(m)
}

func (d IIRRecord) ToPersistence() IIRRecordDB {
	return IIRRecordDB(d)
}

// StudentPersonalInfo mappers
func (m StudentPersonalInfoDB) ToDomain() StudentPersonalInfo {
	return StudentPersonalInfo{
		ID:              m.ID,
		IIRID:           m.IIRID,
		StudentNumber:   m.StudentNumber,
		GenderID:        m.GenderID,
		CivilStatusID:   m.CivilStatusID,
		ReligionID:      m.ReligionID,
		HeightFt:        m.HeightFt,
		WeightKg:        m.WeightKg,
		Complexion:      m.Complexion,
		HighSchoolGWA:   m.HighSchoolGWA,
		CourseID:        m.CourseID,
		YearLevel:       m.YearLevel,
		Section:         m.Section,
		PlaceOfBirth:    m.PlaceOfBirth,
		DateOfBirth:     m.DateOfBirth,
		IsEmployed:      m.IsEmployed,
		EmployerName:    structs.FromSqlNull(m.EmployerName),
		EmployerAddress: structs.FromSqlNull(m.EmployerAddress),
		MobileNumber:    m.MobileNumber,
		TelephoneNumber: structs.FromSqlNull(m.TelephoneNumber),
		StatusID:        m.StatusID,
		GraduationYear:  structs.FromSqlNullInt64(m.GraduationYear),
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}

func (d StudentPersonalInfo) ToPersistence() StudentPersonalInfoDB {
	return StudentPersonalInfoDB{
		ID:              d.ID,
		IIRID:           d.IIRID,
		StudentNumber:   d.StudentNumber,
		GenderID:        d.GenderID,
		CivilStatusID:   d.CivilStatusID,
		ReligionID:      d.ReligionID,
		HeightFt:        d.HeightFt,
		WeightKg:        d.WeightKg,
		Complexion:      d.Complexion,
		HighSchoolGWA:   d.HighSchoolGWA,
		CourseID:        d.CourseID,
		YearLevel:       d.YearLevel,
		Section:         d.Section,
		PlaceOfBirth:    d.PlaceOfBirth,
		DateOfBirth:     d.DateOfBirth,
		IsEmployed:      d.IsEmployed,
		EmployerName:    structs.ToSqlNull(d.EmployerName),
		EmployerAddress: structs.ToSqlNull(d.EmployerAddress),
		MobileNumber:    d.MobileNumber,
		TelephoneNumber: structs.ToSqlNull(d.TelephoneNumber),
		StatusID:        d.StatusID,
		GraduationYear:  structs.ToSqlNullInt64(d.GraduationYear),
		CreatedAt:       d.CreatedAt,
		UpdatedAt:       d.UpdatedAt,
	}
}

// EmergencyContact mappers
func (m EmergencyContactDB) ToDomain() EmergencyContact {
	return EmergencyContact{
		ID:             m.ID,
		IIRID:          m.IIRID,
		FirstName:      m.FirstName,
		MiddleName:     structs.FromSqlNull(m.MiddleName),
		LastName:       m.LastName,
		SuffixName:     structs.FromSqlNull(m.SuffixName),
		ContactNumber:  m.ContactNumber,
		RelationshipID: m.RelationshipID,
		AddressID:      m.AddressID,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}

func (d EmergencyContact) ToPersistence() EmergencyContactDB {
	return EmergencyContactDB{
		ID:             d.ID,
		IIRID:          d.IIRID,
		FirstName:      d.FirstName,
		MiddleName:     structs.ToSqlNull(d.MiddleName),
		LastName:       d.LastName,
		SuffixName:     structs.ToSqlNull(d.SuffixName),
		ContactNumber:  d.ContactNumber,
		RelationshipID: d.RelationshipID,
		AddressID:      d.AddressID,
		CreatedAt:      d.CreatedAt,
		UpdatedAt:      d.UpdatedAt,
	}
}

// StudentAddress mappers
func (m StudentAddressDB) ToDomain() StudentAddress {
	return StudentAddress(m)
}

func (d StudentAddress) ToPersistence() StudentAddressDB {
	return StudentAddressDB(d)
}

// RelatedPerson mapper
func (m RelatedPersonDB) ToDomain() RelatedPerson {
	return RelatedPerson{
		ID:               m.ID,
		EducationalLevel: m.EducationalLevel,
		DateOfBirth:      m.DateOfBirth,
		LastName:         m.LastName,
		FirstName:        m.FirstName,
		MiddleName:       structs.FromSqlNull(m.MiddleName),
		SuffixName:       structs.FromSqlNull(m.SuffixName),
		Occupation:       structs.FromSqlNull(m.Occupation),
		EmployerName:     structs.FromSqlNull(m.EmployerName),
		EmployerAddress:  structs.FromSqlNull(m.EmployerAddress),
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
	}
}

func (d RelatedPerson) ToPersistence() RelatedPersonDB {
	return RelatedPersonDB{
		ID:               d.ID,
		EducationalLevel: d.EducationalLevel,
		DateOfBirth:      d.DateOfBirth,
		LastName:         d.LastName,
		FirstName:        d.FirstName,
		MiddleName:       structs.ToSqlNull(d.MiddleName),
		SuffixName:       structs.ToSqlNull(d.SuffixName),
		Occupation:       structs.ToSqlNull(d.Occupation),
		EmployerName:     structs.ToSqlNull(d.EmployerName),
		EmployerAddress:  structs.ToSqlNull(d.EmployerAddress),
		CreatedAt:        d.CreatedAt,
		UpdatedAt:        d.UpdatedAt,
	}
}

// FamilyBackground mapper
func (m FamilyBackgroundDB) ToDomain() FamilyBackground {
	return FamilyBackground{
		ID:                    m.ID,
		IIRID:                 m.IIRID,
		ParentalStatusID:      m.ParentalStatusID,
		ParentalStatusDetails: structs.FromSqlNull(m.ParentalStatusDetails),
		Brothers:              m.Brothers,
		Sisters:               m.Sisters,
		EmployedSiblings:      m.EmployedSiblings,
		OrdinalPosition:       m.OrdinalPosition,
		HaveQuietPlaceToStudy: m.HaveQuietPlaceToStudy,
		IsSharingRoom:         m.IsSharingRoom,
		RoomSharingDetails:    structs.FromSqlNull(m.RoomSharingDetails),
		NatureOfResidenceId:   m.NatureOfResidenceId,
		CreatedAt:             m.CreatedAt,
		UpdatedAt:             m.UpdatedAt,
	}
}

func (d FamilyBackground) ToPersistence() FamilyBackgroundDB {
	return FamilyBackgroundDB{
		ID:                    d.ID,
		IIRID:                 d.IIRID,
		ParentalStatusID:      d.ParentalStatusID,
		ParentalStatusDetails: structs.ToSqlNull(d.ParentalStatusDetails),
		Brothers:              d.Brothers,
		Sisters:               d.Sisters,
		EmployedSiblings:      d.EmployedSiblings,
		OrdinalPosition:       d.OrdinalPosition,
		HaveQuietPlaceToStudy: d.HaveQuietPlaceToStudy,
		IsSharingRoom:         d.IsSharingRoom,
		RoomSharingDetails:    structs.ToSqlNull(d.RoomSharingDetails),
		NatureOfResidenceId:   d.NatureOfResidenceId,
		CreatedAt:             d.CreatedAt,
		UpdatedAt:             d.UpdatedAt,
	}
}

// StudentHealthRecord mapper
func (m StudentHealthRecordDB) ToDomain() StudentHealthRecord {
	return StudentHealthRecord{
		ID:                      m.ID,
		IIRID:                   m.IIRID,
		VisionHasProblem:        m.VisionHasProblem,
		VisionDetails:           structs.FromSqlNull(m.VisionDetails),
		HearingHasProblem:       m.HearingHasProblem,
		HearingDetails:          structs.FromSqlNull(m.HearingDetails),
		SpeechHasProblem:        m.SpeechHasProblem,
		SpeechDetails:           structs.FromSqlNull(m.SpeechDetails),
		GeneralHealthHasProblem: m.GeneralHealthHasProblem,
		GeneralHealthDetails:    structs.FromSqlNull(m.GeneralHealthDetails),
		CreatedAt:               m.CreatedAt,
		UpdatedAt:               m.UpdatedAt,
	}
}

func (d StudentHealthRecord) ToPersistence() StudentHealthRecordDB {
	return StudentHealthRecordDB{
		ID:                      d.ID,
		IIRID:                   d.IIRID,
		VisionHasProblem:        d.VisionHasProblem,
		VisionDetails:           structs.ToSqlNull(d.VisionDetails),
		HearingHasProblem:       d.HearingHasProblem,
		HearingDetails:          structs.ToSqlNull(d.HearingDetails),
		SpeechHasProblem:        d.SpeechHasProblem,
		SpeechDetails:           structs.ToSqlNull(d.SpeechDetails),
		GeneralHealthHasProblem: d.GeneralHealthHasProblem,
		GeneralHealthDetails:    structs.ToSqlNull(d.GeneralHealthDetails),
		CreatedAt:               d.CreatedAt,
		UpdatedAt:               d.UpdatedAt,
	}
}

// EducationalBackground mapper
func (m EducationalBackgroundDB) ToDomain() EducationalBackground {
	return EducationalBackground{
		ID:                 m.ID,
		IIRID:              m.IIRID,
		NatureOfSchooling:  m.NatureOfSchooling,
		InterruptedDetails: structs.FromSqlNull(m.InterruptedDetails),
		CreatedAt:          m.CreatedAt,
		UpdatedAt:          m.UpdatedAt,
	}
}

func (d EducationalBackground) ToPersistence() EducationalBackgroundDB {
	return EducationalBackgroundDB{
		ID:                 d.ID,
		IIRID:              d.IIRID,
		NatureOfSchooling:  d.NatureOfSchooling,
		InterruptedDetails: structs.ToSqlNull(d.InterruptedDetails),
		CreatedAt:          d.CreatedAt,
		UpdatedAt:          d.UpdatedAt,
	}
}

// SchoolDetails mapper
func (m SchoolDetailsDB) ToDomain() SchoolDetails {
	return SchoolDetails{
		ID:                 m.ID,
		EBID:               m.EBID,
		EducationalLevelID: m.EducationalLevelID,
		SchoolName:         m.SchoolName,
		SchoolAddress:      m.SchoolAddress,
		SchoolType:         m.SchoolType,
		YearStarted:        m.YearStarted,
		YearCompleted:      m.YearCompleted,
		Awards:             structs.FromSqlNull(m.Awards),
		CreatedAt:          m.CreatedAt,
		UpdatedAt:          m.UpdatedAt,
	}
}

func (d SchoolDetails) ToPersistence() SchoolDetailsDB {
	return SchoolDetailsDB{
		ID:                 d.ID,
		EBID:               d.EBID,
		EducationalLevelID: d.EducationalLevelID,
		SchoolName:         d.SchoolName,
		SchoolAddress:      d.SchoolAddress,
		SchoolType:         d.SchoolType,
		YearStarted:        d.YearStarted,
		YearCompleted:      d.YearCompleted,
		Awards:             structs.ToSqlNull(d.Awards),
		CreatedAt:          d.CreatedAt,
		UpdatedAt:          d.UpdatedAt,
	}
}

// StudentFinance mapper
func (m StudentFinanceDB) ToDomain() StudentFinance {
	return StudentFinance{
		ID:                         m.ID,
		IIRID:                      m.IIRID,
		MonthlyFamilyIncomeRangeID: m.IncomeRangeID,
		OtherIncomeDetails:         structs.FromSqlNull(m.OtherIncome),
		WeeklyAllowance:            m.WeeklyAllowance,
		CreatedAt:                  m.CreatedAt,
		UpdatedAt:                  m.UpdatedAt,
	}
}

func (d StudentFinance) ToPersistence() StudentFinanceDB {
	return StudentFinanceDB{
		ID:              d.ID,
		IIRID:           d.IIRID,
		IncomeRangeID:   d.MonthlyFamilyIncomeRangeID,
		OtherIncome:     structs.ToSqlNull(d.OtherIncomeDetails),
		WeeklyAllowance: d.WeeklyAllowance,
		CreatedAt:       d.CreatedAt,
		UpdatedAt:       d.UpdatedAt,
	}
}

// StudentActivity mapper
func (m StudentActivityDB) ToDomain() StudentActivity {
	return StudentActivity{
		ID:                 m.ID,
		IIRID:              m.IIRID,
		OptionID:           m.OptionID,
		OtherSpecification: structs.FromSqlNull(m.OtherSpecification),
		Role:               m.Role,
		RoleSpecification:  structs.FromSqlNull(m.RoleSpecification),
		CreatedAt:          m.CreatedAt,
		UpdatedAt:          m.UpdatedAt,
	}
}

func (d StudentActivity) ToPersistence() StudentActivityDB {
	return StudentActivityDB{
		ID:                 d.ID,
		IIRID:              d.IIRID,
		OptionID:           d.OptionID,
		OtherSpecification: structs.ToSqlNull(d.OtherSpecification),
		Role:               d.Role,
		RoleSpecification:  structs.ToSqlNull(d.RoleSpecification),
		CreatedAt:          d.CreatedAt,
		UpdatedAt:          d.UpdatedAt,
	}
}

// StudentConsultation mapper
func (m StudentConsultationDB) ToDomain() StudentConsultation {
	return StudentConsultation{
		ID:               m.ID,
		IIRID:            m.IIRID,
		ProfessionalType: m.ProfessionalType,
		HasConsulted:     m.HasConsulted,
		WhenDate:         structs.FromSqlNull(m.WhenDate),
		ForWhat:          structs.FromSqlNull(m.ForWhat),
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
	}
}

func (d StudentConsultation) ToPersistence() StudentConsultationDB {
	return StudentConsultationDB{
		ID:               d.ID,
		IIRID:            d.IIRID,
		ProfessionalType: d.ProfessionalType,
		HasConsulted:     d.HasConsulted,
		WhenDate:         structs.ToSqlNull(d.WhenDate),
		ForWhat:          structs.ToSqlNull(d.ForWhat),
		CreatedAt:        d.CreatedAt,
		UpdatedAt:        d.UpdatedAt,
	}
}

// StudentSubjectPreference mappers
func (m StudentSubjectPreferenceDB) ToDomain() StudentSubjectPreference {
	return StudentSubjectPreference(m)
}

func (d StudentSubjectPreference) ToPersistence() StudentSubjectPreferenceDB {
	return StudentSubjectPreferenceDB(d)
}

// StudentHobby mappers
func (m StudentHobbyDB) ToDomain() StudentHobby {
	return StudentHobby(m)
}

func (d StudentHobby) ToPersistence() StudentHobbyDB {
	return StudentHobbyDB(d)
}

// TestResult mappers
func (m TestResultDB) ToDomain() TestResult {
	return TestResult(m)
}

func (d TestResult) ToPersistence() TestResultDB {
	return TestResultDB(d)
}

// StudentRelatedPerson mappers
func (m StudentRelatedPersonDB) ToDomain() StudentRelatedPerson {
	return StudentRelatedPerson(m)
}

func (d StudentRelatedPerson) ToPersistence() StudentRelatedPersonDB {
	return StudentRelatedPersonDB(d)
}

// StudentSiblingSupport mappers
func (m StudentSiblingSupportDB) ToDomain() StudentSiblingSupport {
	return StudentSiblingSupport(m)
}

func (d StudentSiblingSupport) ToPersistence() StudentSiblingSupportDB {
	return StudentSiblingSupportDB(d)
}

// StudentFinancialSupport mappers
func (m StudentFinancialSupportDB) ToDomain() StudentFinancialSupport {
	return StudentFinancialSupport(m)
}

func (d StudentFinancialSupport) ToPersistence() StudentFinancialSupportDB {
	return StudentFinancialSupportDB(d)
}

// StudentCOR mappers
func (m StudentCORDB) ToDomain() StudentCOR {
	return StudentCOR{
		FileID:     m.FileID,
		StudentID:  m.StudentID,
		ValidFrom:  structs.FromSqlNullTime(m.ValidFrom),
		ValidUntil: structs.FromSqlNullTime(m.ValidUntil),
	}
}

func (d StudentCOR) ToPersistence() StudentCORDB {
	return StudentCORDB{
		FileID:     d.FileID,
		StudentID:  d.StudentID,
		ValidFrom:  structs.ToSqlNullTime(d.ValidFrom),
		ValidUntil: structs.ToSqlNullTime(d.ValidUntil),
	}
}
