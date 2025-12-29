package constants

type CivilStatusID int

const (
	SingleCivilStatusID CivilStatusID = iota + 1
	MarriedCivilStatusID
	WidowedCivilStatusID
	DivorcedCivilStatusID
)

type CivilStatusType struct {
	ID   CivilStatusID
	Name string
}

var CivilStatusTypes = map[CivilStatusID]CivilStatusType{
	SingleCivilStatusID:   {ID: SingleCivilStatusID, Name: "Single"},
	MarriedCivilStatusID:  {ID: MarriedCivilStatusID, Name: "Married"},
	WidowedCivilStatusID:  {ID: WidowedCivilStatusID, Name: "Widowed"},
	DivorcedCivilStatusID: {ID: DivorcedCivilStatusID, Name: "Divorced"},
}
