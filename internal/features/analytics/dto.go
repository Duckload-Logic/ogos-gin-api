package analytics

type DemographicStatDTO struct {
	Category    string  `json:"category"`
	MaleCount   int     `json:"maleCount"`
	MalePct     float64 `json:"malePct"`
	FemaleCount int     `json:"femaleCount"`
	FemalePct   float64 `json:"femalePct"`
	Total       int     `json:"total"`
	TotalPct    float64 `json:"totalPct"`
	Rank        int     `json:"rank,omitempty"`
}

type DashboardResponseDTO struct {
	TotalStudents int `json:"totalStudents"`

	// Personal Information
	AgeDistribution []DemographicStatDTO `json:"ageDistribution"`
	CivilStatus     []DemographicStatDTO `json:"civilStatus"`
	Religions       []DemographicStatDTO `json:"religions"`
	CityAddress     []DemographicStatDTO `json:"cityAddress"`

	// Family & Financial Background
	MonthlyIncome        []DemographicStatDTO `json:"monthlyIncome"`
	OrdinalPosition      []DemographicStatDTO `json:"ordinalPosition"`
	FatherEducation      []DemographicStatDTO `json:"fatherEducation"`
	MotherEducation      []DemographicStatDTO `json:"motherEducation"`
	ParentsMaritalStatus []DemographicStatDTO `json:"parentsMaritalStatus"`

	// Academic Background
	HighSchoolGWA     []DemographicStatDTO `json:"highSchoolGWA"`
	Elementary        []DemographicStatDTO `json:"elementary"`
	JuniorHigh        []DemographicStatDTO `json:"juniorHigh"`
	SeniorHigh        []DemographicStatDTO `json:"seniorHigh"`
	NatureOfSchooling []DemographicStatDTO `json:"natureOfSchooling"`

	// Study Environment
	QuietStudyPlace []DemographicStatDTO `json:"quietStudyPlace"`
}