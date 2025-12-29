package constants

type EducationalLevelID int

const (
	ElementaryLevelID EducationalLevelID = iota + 1
	JuniorHighLevelID
	SeniorHighLevelID
	CollegeLevelID
)

type EducationalLevel struct {
	ID   EducationalLevelID
	Name string
}

var EducationalLevels = map[EducationalLevelID]EducationalLevel{
	ElementaryLevelID: {ID: ElementaryLevelID, Name: "Elementary"},
	JuniorHighLevelID: {ID: JuniorHighLevelID, Name: "Junior High School"},
	SeniorHighLevelID: {ID: SeniorHighLevelID, Name: "Senior High School"},
	CollegeLevelID:    {ID: CollegeLevelID, Name: "College"},
}
