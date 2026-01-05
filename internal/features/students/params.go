package students

type IncludeParams struct {
	IncludeFamily    bool `form:"include_family"`
	IncludeHealth    bool `form:"include_health"`
	IncludeEducation bool `form:"include_education"`
	IncludeAddress   bool `form:"include_address"`
	IncludeFinance   bool `form:"include_finance"`
	IncludeGuardians bool `form:"include_guardians"`
}
