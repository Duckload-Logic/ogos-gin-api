package constants

type HealthRemarkTypeID int

const (
	NoProblemHealthRemarkID HealthRemarkTypeID = iota + 1
	IssueHealthRemarkID
)

type HealthRemarkType struct {
	ID   HealthRemarkTypeID
	Name string
}

var HealthRemarkTypes = map[HealthRemarkTypeID]HealthRemarkType{
	NoProblemHealthRemarkID: {ID: NoProblemHealthRemarkID, Name: "No problem"},
	IssueHealthRemarkID:     {ID: IssueHealthRemarkID, Name: "Issue"},
}
