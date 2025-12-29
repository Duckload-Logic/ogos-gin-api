package constants

type RoleID int

const (
	StudentRoleID RoleID = iota + 1
	CounselorRoleID
	FrontDeskRoleID
)

type Role struct {
	ID   RoleID
	Name string
}

var Roles = map[RoleID]Role{
	StudentRoleID:   {ID: StudentRoleID, Name: "STUDENT"},
	CounselorRoleID: {ID: CounselorRoleID, Name: "COUNSELOR"},
	FrontDeskRoleID: {ID: FrontDeskRoleID, Name: "FRONTDESK"},
}
