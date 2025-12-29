package constants

type GenderID int

const (
	MaleGenderID GenderID = iota + 1
	FemaleGenderID
	PreferNotToSayGenderID
)

type Gender struct {
	ID   GenderID
	Name string
}

var Genders = map[GenderID]Gender{
	MaleGenderID:           {ID: MaleGenderID, Name: "Male"},
	FemaleGenderID:         {ID: FemaleGenderID, Name: "Female"},
	PreferNotToSayGenderID: {ID: PreferNotToSayGenderID, Name: "Prefer not to say"},
}
