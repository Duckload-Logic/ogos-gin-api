package constants

type FinancialSupportTypeID int

const (
	ScholarshipSupportID FinancialSupportTypeID = iota + 1
	SelfFundedSupportID
	SponsoredSupportID
	ParentalSupportID
	OthersSupportID
)

type FinancialSupportType struct {
	ID   FinancialSupportTypeID
	Name string
}

var FinancialSupportTypes = map[FinancialSupportTypeID]FinancialSupportType{
	ScholarshipSupportID: {ID: ScholarshipSupportID, Name: "Scholarship"},
	SelfFundedSupportID:  {ID: SelfFundedSupportID, Name: "Self-funded"},
	SponsoredSupportID:   {ID: SponsoredSupportID, Name: "Sponsored"},
	ParentalSupportID:    {ID: ParentalSupportID, Name: "Parental Support"},
	OthersSupportID:      {ID: OthersSupportID, Name: "Others"},
}
