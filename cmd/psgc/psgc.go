package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const psgcBaseURL = "https://psgc.cloud/api"

// PSGC API response types
type PSGCRegion struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type PSGCProvince struct {
	Code       string `json:"code"`
	Name       string `json:"name"`
	RegionCode string `json:"regionCode"`
}

type PSGCCity struct {
	Code         string `json:"code"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	ZipCode      string `json:"zipCode"`
	District     string `json:"district"`
	RegionCode   string `json:"regionCode"`
	ProvinceCode string `json:"provinceCode,omitempty"`
}

type PSGCBarangay struct {
	Code     string `json:"code"`
	Name     string `json:"name"`
	CityCode string `json:"cityCode"`
}

// PSGCData is the full exported structure written to JSON.
type PSGCData struct {
	Regions   []PSGCRegion   `json:"regions"`
	Provinces []PSGCProvince `json:"provinces"`
	Cities    []PSGCCity     `json:"cities"`
	Barangays []PSGCBarangay `json:"barangays"`
}

// API-level response shapes (different from our output shapes)
type apiProvince struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type apiCity struct {
	Code     string `json:"code"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	ZipCode  string `json:"zip_code"`
	District string `json:"district"`
}

type apiBarangay struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

var client = &http.Client{Timeout: 30 * time.Second}

func main() {
	outFile := "cmd/locations/psgc_data.json"
	if len(os.Args) > 1 {
		outFile = os.Args[1]
	}

	data, err := fetchAll()
	if err != nil {
		log.Fatal(err)
	}

	blob, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatal("failed to marshal JSON:", err)
	}

	if err := os.WriteFile(outFile, blob, 0o644); err != nil {
		log.Fatal("failed to write file:", err)
	}

	fmt.Printf("Wrote %d regions, %d provinces, %d cities, %d barangays to %s\n",
		len(data.Regions), len(data.Provinces), len(data.Cities), len(data.Barangays), outFile)
}

// fetchJSON performs a GET with retry on 429.
func fetchJSON(url string, target interface{}) error {
	maxRetries := 5
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt == 0 {
			time.Sleep(300 * time.Millisecond)
		}

		resp, err := client.Get(url)
		if err != nil {
			return fmt.Errorf("request failed for %s: %w", url, err)
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			resp.Body.Close()
			backoff := time.Duration(2<<uint(attempt)) * time.Second
			fmt.Printf("  Rate limited, retrying in %v (attempt %d/%d)\n", backoff, attempt+1, maxRetries)
			time.Sleep(backoff)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return fmt.Errorf("unexpected status %d for %s", resp.StatusCode, url)
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return fmt.Errorf("failed to read body: %w", err)
		}

		return json.Unmarshal(body, target)
	}
	return fmt.Errorf("max retries exceeded for %s", url)
}

func fetchAll() (*PSGCData, error) {
	data := &PSGCData{}

	// 1. Regions
	fmt.Println("Fetching regions...")
	var regions []PSGCRegion
	if err := fetchJSON(psgcBaseURL+"/regions", &regions); err != nil {
		return nil, fmt.Errorf("regions: %w", err)
	}
	data.Regions = regions
	fmt.Printf("  %d regions\n", len(regions))

	// 2. Provinces per region
	fmt.Println("Fetching provinces...")
	for _, reg := range regions {
		var provs []apiProvince
		if err := fetchJSON(fmt.Sprintf("%s/regions/%s/provinces", psgcBaseURL, reg.Code), &provs); err != nil {
			fmt.Printf("  Warning: provinces for %s: %v\n", reg.Name, err)
			continue
		}
		for _, p := range provs {
			data.Provinces = append(data.Provinces, PSGCProvince{
				Code:       p.Code,
				Name:       p.Name,
				RegionCode: reg.Code,
			})
		}
	}
	fmt.Printf("  %d provinces\n", len(data.Provinces))

	// 3. Cities/municipalities — fetch from regions FIRST (authoritative region mapping),
	// then enrich with province codes from the province-level fetch.
	fmt.Println("Fetching cities/municipalities...")

	// 3a. Fetch cities per region first (gives correct region_code, especially for NCR)
	seenCities := make(map[string]bool)
	cityIndex := make(map[string]int) // city code -> index in data.Cities

	for _, reg := range regions {
		var cities []apiCity
		if err := fetchJSON(fmt.Sprintf("%s/regions/%s/cities-municipalities", psgcBaseURL, reg.Code), &cities); err != nil {
			fmt.Printf("  Warning: cities for region %s: %v\n", reg.Name, err)
			continue
		}
		for _, c := range cities {
			if seenCities[c.Code] {
				continue
			}
			seenCities[c.Code] = true
			cityIndex[c.Code] = len(data.Cities)
			data.Cities = append(data.Cities, PSGCCity{
				Code:       c.Code,
				Name:       c.Name,
				Type:       c.Type,
				ZipCode:    c.ZipCode,
				District:   c.District,
				RegionCode: reg.Code,
				// ProvinceCode left empty; filled below if the city belongs to a province
			})
		}
	}

	// 3b. Fetch cities per province to fill in province_code
	for _, prov := range data.Provinces {
		var cities []apiCity
		if err := fetchJSON(fmt.Sprintf("%s/provinces/%s/cities-municipalities", psgcBaseURL, prov.Code), &cities); err != nil {
			fmt.Printf("  Warning: cities for province %s: %v\n", prov.Code, err)
			continue
		}
		for _, c := range cities {
			if idx, ok := cityIndex[c.Code]; ok {
				// Only enrich with province code if the province belongs to the same region
				if data.Cities[idx].RegionCode == prov.RegionCode {
					data.Cities[idx].ProvinceCode = prov.Code
				}
			} else {
				// City wasn't in region-level fetch — add it
				seenCities[c.Code] = true
				cityIndex[c.Code] = len(data.Cities)
				data.Cities = append(data.Cities, PSGCCity{
					Code:         c.Code,
					Name:         c.Name,
					Type:         c.Type,
					ZipCode:      c.ZipCode,
					District:     c.District,
					RegionCode:   prov.RegionCode,
					ProvinceCode: prov.Code,
				})
			}
		}
	}
	fmt.Printf("  %d cities/municipalities\n", len(data.Cities))

	// 4. Barangays per city
	fmt.Println("Fetching barangays (this will take a while)...")
	for i, city := range data.Cities {
		var brgys []apiBarangay
		if err := fetchJSON(fmt.Sprintf("%s/cities-municipalities/%s/barangays", psgcBaseURL, city.Code), &brgys); err != nil {
			continue
		}
		for _, b := range brgys {
			data.Barangays = append(data.Barangays, PSGCBarangay{
				Code:     b.Code,
				Name:     b.Name,
				CityCode: city.Code,
			})
		}
		if (i+1)%100 == 0 {
			fmt.Printf("  Processed %d/%d cities...\n", i+1, len(data.Cities))
		}
	}
	fmt.Printf("  %d barangays\n", len(data.Barangays))

	return data, nil
}
