package constants

type RelationshipTypeID int

const (
	FatherRelationshipID RelationshipTypeID = iota + 1
	MotherRelationshipID
	RelativeRelationshipID
	LegalGuardianRelationshipID
)

type RelationshipType struct {
	ID   RelationshipTypeID
	Name string
}

var RelationshipTypes = map[RelationshipTypeID]RelationshipType{
	FatherRelationshipID:        {ID: FatherRelationshipID, Name: "Father"},
	MotherRelationshipID:        {ID: MotherRelationshipID, Name: "Mother"},
	RelativeRelationshipID:      {ID: RelativeRelationshipID, Name: "Relative"},
	LegalGuardianRelationshipID: {ID: LegalGuardianRelationshipID, Name: "Legal Guardian"},
}
