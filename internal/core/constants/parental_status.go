package constants

type ParentalStatusID int

const (
	MarriedLivingTogetherParentalStatusID ParentalStatusID = iota + 1
	MarriedLivingSeparatelyParentalStatusID
	WorkingAbroadParentalStatusID
	DivorcedAnnulledParentalStatusID
	SeparatedParentalStatusID
	OtherParentalStatusID
)

type ParentalStatusType struct {
	ID   ParentalStatusID
	Name string
}

var ParentalStatusTypes = map[ParentalStatusID]ParentalStatusType{
	MarriedLivingTogetherParentalStatusID:   {ID: MarriedLivingTogetherParentalStatusID, Name: "Married and Living Together"},
	MarriedLivingSeparatelyParentalStatusID: {ID: MarriedLivingSeparatelyParentalStatusID, Name: "Married but Living Separately"},
	WorkingAbroadParentalStatusID:           {ID: WorkingAbroadParentalStatusID, Name: "Father/Mother working Abroad"},
	DivorcedAnnulledParentalStatusID:        {ID: DivorcedAnnulledParentalStatusID, Name: "Divorced or Annulled"},
	SeparatedParentalStatusID:               {ID: SeparatedParentalStatusID, Name: "Separated"},
	OtherParentalStatusID:                   {ID: OtherParentalStatusID, Name: "Other"},
}
