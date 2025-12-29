package constants

type ReligionID int

const (
	RomanCatholicismReligionID ReligionID = iota + 1
	IslamReligionID
	IglesiaNiCristoReligionID
	SeventhDayAdventistReligionID
	BibleBaptistReligionID
	PhilippineIndependentChurchReligionID
	JehovahsWitnessesReligionID
	BuddhismReligionID
	OtherReligionID
)

type ReligionType struct {
	ID   ReligionID
	Name string
}

var ReligionTypes = map[ReligionID]ReligionType{
	RomanCatholicismReligionID:            {ID: RomanCatholicismReligionID, Name: "Roman Catholicism"},
	IslamReligionID:                       {ID: IslamReligionID, Name: "Islam"},
	IglesiaNiCristoReligionID:             {ID: IglesiaNiCristoReligionID, Name: "Iglesia ni Cristo"},
	SeventhDayAdventistReligionID:         {ID: SeventhDayAdventistReligionID, Name: "Seventh-day Adventist"},
	BibleBaptistReligionID:                {ID: BibleBaptistReligionID, Name: "Bible Baptist Church"},
	PhilippineIndependentChurchReligionID: {ID: PhilippineIndependentChurchReligionID, Name: "Philippine Independent Church"},
	JehovahsWitnessesReligionID:           {ID: JehovahsWitnessesReligionID, Name: "Jehovahs Witnesses"},
	BuddhismReligionID:                    {ID: BuddhismReligionID, Name: "Buddhism"},
	OtherReligionID:                       {ID: OtherReligionID, Name: "Other"},
}
