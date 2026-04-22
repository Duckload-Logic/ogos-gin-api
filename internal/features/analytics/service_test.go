package analytics

import (
	"context"
	"testing"

	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
	"go.uber.org/mock/gomock"
)

func TestService_GetDashboard(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepositoryInterface(ctrl)
	mockRedis := datastore.NewMockRedisClientInterface(ctrl)
	svc := NewService(mockRepo, mockRedis)

	ctx := context.Background()

	t.Run("success with data", func(t *testing.T) {
		total := 100
		mockRepo.EXPECT().GetTotalStudents(ctx, 2024, 1).Return(total, nil)

		// Mock at least one demographic call to verify mapToDTO logic
		genderStats := []DemographicStat{
			{Category: "Male", MaleCount: 40, Total: 40},
			{Category: "Female", FemaleCount: 60, Total: 60},
		}
		mockRepo.EXPECT().GetGenderStats(ctx, 2024, 1).Return(genderStats, nil)

		// Mock other calls with empty values for brevity
		mockRepo.EXPECT().
			GetAgeStats(ctx, gomock.Any(), gomock.Any()).
			Return(nil, nil)
		mockRepo.EXPECT().
			GetCivilStatusStats(ctx, gomock.Any(), gomock.Any()).
			Return(nil, nil)
		mockRepo.EXPECT().
			GetReligionStats(ctx, gomock.Any(), gomock.Any()).
			Return(nil, nil)
		mockRepo.EXPECT().
			GetCityAddressStats(ctx, gomock.Any(), gomock.Any()).
			Return(nil, nil)
		mockRepo.EXPECT().
			GetMonthlyIncomeStats(ctx, gomock.Any(), gomock.Any()).
			Return(nil, nil)
		mockRepo.EXPECT().
			GetOrdinalPositionStats(ctx, gomock.Any(), gomock.Any()).
			Return(nil, nil)
		mockRepo.EXPECT().
			GetQuietStudyPlaceStats(ctx, gomock.Any(), gomock.Any()).
			Return(nil, nil)
		mockRepo.EXPECT().
			GetFatherEducationStats(ctx, gomock.Any(), gomock.Any()).
			Return(nil, nil)
		mockRepo.EXPECT().
			GetMotherEducationStats(ctx, gomock.Any(), gomock.Any()).
			Return(nil, nil)
		mockRepo.EXPECT().
			GetParentsMaritalStatusStats(ctx, gomock.Any(), gomock.Any()).
			Return(nil, nil)
		mockRepo.EXPECT().
			GetHSGWAStats(ctx, gomock.Any(), gomock.Any()).
			Return(nil, nil)
		mockRepo.EXPECT().
			GetElementaryStats(ctx, gomock.Any(), gomock.Any()).
			Return(nil, nil)
		mockRepo.EXPECT().
			GetJuniorHighStats(ctx, gomock.Any(), gomock.Any()).
			Return(nil, nil)
		mockRepo.EXPECT().
			GetSeniorHighStats(ctx, gomock.Any(), gomock.Any()).
			Return(nil, nil)
		mockRepo.EXPECT().
			GetNatureOfSchoolingStats(ctx, gomock.Any(), gomock.Any()).
			Return(nil, nil)

		resp, err := svc.GetDashboard(ctx, 2024, 1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.TotalStudents != total {
			t.Errorf("got %d, want %d", resp.TotalStudents, total)
		}
		if len(resp.GenderDistribution) != 2 {
			t.Errorf(
				"got %d gender stats, want 2",
				len(resp.GenderDistribution),
			)
		}
	})
}

func TestService_GetAdminDashboard(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepositoryInterface(ctrl)
	mockRedis := datastore.NewMockRedisClientInterface(ctrl)
	svc := NewService(mockRepo, mockRedis)

	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().GetTotalStudents(ctx, 0, 0).Return(100, nil)
		mockRepo.EXPECT().GetTotalReports(ctx).Return(10, nil)
		mockRepo.EXPECT().GetTotalAppointments(ctx).Return(20, nil)
		mockRepo.EXPECT().GetTotalSlips(ctx).Return(5, nil)
		mockRepo.EXPECT().GetMonthlyVisitorStats(ctx, "30d").Return(nil, nil)

		mockRedis.EXPECT().
			Keys(ctx, "session:*").
			Return([]string{"s1", "s2"}, nil)

		resp, err := svc.GetAdminDashboard(ctx, "30d", "system")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.LiveSessions != 2 {
			t.Errorf("got %d live sessions, want 2", resp.LiveSessions)
		}
	})
}
