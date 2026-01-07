package students

type IncludeParams struct {
	IncludeReasons          bool `form:"include_reasons"`
	IncludeFamily           bool `form:"include_family"`
	IncludeHealth           bool `form:"include_health"`
	IncludeEducation        bool `form:"include_education"`
	IncludeAddress          bool `form:"include_address"`
	IncludeFinance          bool `form:"include_finance"`
	IncludeParents          bool `form:"include_parents"`
	IncludeEmergencyContact bool `form:"include_emergency_contact"`
}
