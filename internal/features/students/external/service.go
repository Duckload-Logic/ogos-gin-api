package external

import (
	"context"
	"fmt"

	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
	"github.com/olazo-johnalbert/duckload-api/internal/features/locations"
	"github.com/olazo-johnalbert/duckload-api/internal/features/students"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetStudentByEmail(ctx context.Context, email string) (*OGOSStudentDTO, error) {
	student, err := s.repo.GetStudentByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return &OGOSStudentDTO{
		StudentNumber: student.StudentNumber,
		FirstName:     student.FirstName,
		MiddleName:    structs.FromSqlNull(student.MiddleName),
		LastName:      student.LastName,
		Email:         student.Email,
		ContactNumber: student.ContactNumber,
		Course: students.Course{
			ID:         student.CourseID,
			Code:       student.CourseCode,
			CourseName: student.CourseName,
		},
		Year:    student.Year,
		Section: student.Section,
	}, nil
}

func (s *Service) GetPersonalInfoByStudentNumber(ctx context.Context, studentNumber string) (*OGOSStudentPersonalInfoDTO, error) {
	student, err := s.repo.GetPersonalInfoByStudentNumber(ctx, studentNumber)
	if err != nil {
		return nil, err
	}

	return &OGOSStudentPersonalInfoDTO{
		StudentNumber: student.StudentNumber,
		Gender: students.Gender{
			ID:         student.GenderID,
			GenderName: student.GenderName,
		},
		DateOfBirth:  student.DateOfBirth,
		PlaceOfBirth: student.PlaceOfBirth,
		HeightFt:     student.HeightFt,
		WeightKg:     student.WeightKg,
	}, nil
}

func (s *Service) GetAddressByStudentNumber(ctx context.Context, studentNumber string) ([]OGOSStudentAddressDTO, error) {
	studentAddresses, err := s.repo.GetAddressByStudentNumber(ctx, studentNumber)
	if err != nil {
		return nil, err
	}

	if len(studentAddresses) == 0 {
		return nil, fmt.Errorf("no addresses found for student number: %s", studentNumber)
	}

	var addresesDTO []OGOSStudentAddressDTO
	for _, address := range studentAddresses {
		addresesDTO = append(addresesDTO, OGOSStudentAddressDTO{
			StudentNumber: address.StudentNumber,
			StreetDetails: address.StreetDetails,
			Barangay:      locations.Barangay{Code: address.BarangayCode, Name: address.BarangayName},
			City:          locations.City{Code: address.CityCode, Name: address.CityName},
			Province:      &locations.ProvinceDTO{Code: structs.FromSqlNull(address.ProvinceCode), Name: structs.FromSqlNull(address.ProvinceName)},
			Region:        locations.Region{Code: address.RegionCode, Name: address.RegionName},
		})
	}

	return addresesDTO, nil
}
