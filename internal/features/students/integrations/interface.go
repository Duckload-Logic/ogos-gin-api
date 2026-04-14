package integrations

import (
	"context"
)

type ServiceInterface interface {
	ListStudents(
		ctx context.Context,
		req OGOSListStudentsRequest,
	) (OGOSListStudentsResponse, error)
	GetStudentByStudentNumber(
		ctx context.Context,
		studentNumber string,
	) (*OGOSStudentDTO, error)
	GetStudentByEmail(
		ctx context.Context,
		email string,
	) (*OGOSStudentDTO, error)
	GetPersonalInfoByStudentNumber(
		ctx context.Context,
		studentNumber string,
	) (*OGOSStudentPersonalInfoDTO, error)
	GetAddressByStudentNumber(
		ctx context.Context,
		studentNumber string,
	) ([]OGOSStudentAddressDTO, error)
}

type RepositoryInterface interface {
	ListStudents(
		ctx context.Context,
		req OGOSListStudentsRequest,
	) ([]OGOSStudentView, int, error)
	GetStudentByStudentNumber(
		ctx context.Context,
		studentNumber string,
	) (*OGOSStudentView, error)
	GetStudentByEmail(
		ctx context.Context,
		email string,
	) (*OGOSStudentView, error)
	GetPersonalInfoByStudentNumber(
		ctx context.Context,
		studentNumber string,
	) (*OGOSStudentPersonalInfoView, error)
	GetAddressByStudentNumber(
		ctx context.Context,
		studentNumber string,
	) ([]OGOSStudentAddressView, error)
}
