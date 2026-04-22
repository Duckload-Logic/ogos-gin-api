package locations

import (
	"context"
	"testing"

	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
	"go.uber.org/mock/gomock"
)

func TestService_GetRegions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepositoryInterface(ctrl)
	svc := NewService(mockRepo)

	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		regions := []Region{
			{Code: "01", Name: "Region I"},
			{Code: "02", Name: "Region II"},
		}

		mockRepo.EXPECT().GetRegions(ctx).Return(regions, nil)

		resp, err := svc.GetRegions(ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(resp) != len(regions) {
			t.Errorf("got %d, want %d", len(resp), len(regions))
		}
	})
}

func TestService_GetAddressByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepositoryInterface(ctrl)
	svc := NewService(mockRepo)

	ctx := context.Background()
	addrID := 1

	t.Run("success", func(t *testing.T) {
		street := "Mabini St"
		provinceCode := "01"
		regionCode := "00"
		addr := &Address{
			ID:           addrID,
			RegionCode:   regionCode,
			ProvinceCode: structs.StringToNullableString(provinceCode),
			CityCode:     "0101",
			BarangayCode: "010101",
			StreetDetail: structs.StringToNullableString(street),
		}
		barangay := &Barangay{Code: "010101", Name: "Barangay 1"}
		city := &City{
			Code:         "0101",
			Name:         "City 1",
			ProvinceCode: structs.StringToNullableString(provinceCode),
		}
		province := &Province{
			Code:       "01",
			Name:       "Province 1",
			RegionCode: regionCode,
		}
		region := &Region{Code: regionCode, Name: "Region 0"}

		mockRepo.EXPECT().GetAddressByID(ctx, addrID).Return(addr, nil)
		mockRepo.EXPECT().
			GetBarangayByCode(ctx, addr.BarangayCode).
			Return(barangay, nil)
		mockRepo.EXPECT().GetCityByCode(ctx, addr.CityCode).Return(city, nil)
		mockRepo.EXPECT().
			GetRegionByCode(ctx, addr.RegionCode).
			Return(region, nil)
		mockRepo.EXPECT().
			GetProvinceByCode(ctx, city.ProvinceCode.String).
			Return(province, nil)

		resp, err := svc.GetAddressByID(ctx, addrID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StreetDetail != street || resp.City.Name != city.Name {
			t.Errorf("resp mismatch: %+v", resp)
		}
	})
}
