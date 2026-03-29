package pdf

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"path/filepath"
	"reflect"
	"time"
)

// GotenbergClient defines the interface for converting HTML to PDF.
type GotenbergClient interface {
	ConvertHTML(ctx context.Context, htmlContent string) ([]byte, error)
}

// Service provides high-level PDF generation logic.
type Service struct {
	gotenberg GotenbergClient
}

// NewService creates a new PDF service.
func NewService(gotenberg GotenbergClient) *Service {
	return &Service{
		gotenberg: gotenberg,
	}
}

// GenerateFromTemplate renders an HTML template with data and converts it to
// PDF.
func (s *Service) GenerateFromTemplate(
	ctx context.Context,
	templatePath string,
	data interface{},
) ([]byte, error) {
	tmpl := template.New(filepath.Base(templatePath)).Funcs(template.FuncMap{
		"sub": func(a, b int) int { return a - b },
		"add": func(a, b int) int { return a + b },
		"iterate": func(n int) []int {
			var i []int
			for j := 0; j < n; j++ {
				i = append(i, j)
			}
			return i
		},
		"hasSiblingSupport": func(supports interface{}, target string) bool {
			v := reflect.ValueOf(supports)
			if v.Kind() != reflect.Slice {
				return false
			}
			for i := 0; i < v.Len(); i++ {
				item := reflect.Indirect(v.Index(i))
				if item.FieldByName("SupportName").String() == target {
					return true
				}
			}
			return false
		},
		"hasFinancialSupport": func(supports interface{}, target string) bool {
			v := reflect.ValueOf(supports)
			if v.Kind() != reflect.Slice {
				return false
			}
			for i := 0; i < v.Len(); i++ {
				item := reflect.Indirect(v.Index(i))
				if item.FieldByName("SupportTypeName").String() == target {
					return true
				}
			}
			return false
		},
		"hasActivity": func(activities interface{}, target string) bool {
			v := reflect.ValueOf(activities)
			if v.Kind() != reflect.Slice {
				return false
			}
			for i := 0; i < v.Len(); i++ {
				item := reflect.Indirect(v.Index(i))
				if item.FieldByName("ActivityOption").
					FieldByName("Category").
					String() ==
					target {
					return true
				}
			}
			return false
		},
		"indexHobby": func(hobbies interface{}, index int) string {
			v := reflect.ValueOf(hobbies)
			if v.Kind() != reflect.Slice {
				return ""
			}
			if index >= 0 && index < v.Len() {
				return v.Index(index).FieldByName("HobbyName").String()
			}
			return ""
		},
		"formatDate": func(dateStr string) string {
			if dateStr == "" {
				return ""
			}
			// Handle cases where dateStr might have time component
			if len(dateStr) > 10 {
				dateStr = dateStr[:10]
			}
			t, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				return dateStr // return original if parse fails
			}
			return t.Format("January 02, 2006")
		},
		"calculateAge": func(birthDate string) int {
			if birthDate == "" {
				return 0
			}
			// Debug: print the incoming birthDate string
			fmt.Printf("[PDF_HELPER] calculateAge input: %q\n", birthDate)

			// Handle cases where birthDate might have time component
			if len(birthDate) > 10 {
				birthDate = birthDate[:10]
			}
			t, err := time.Parse("2006-01-02", birthDate)
			if err != nil {
				fmt.Printf("[PDF_HELPER] calculateAge parse error: %v\n", err)
				return 0
			}
			today := time.Now()
			age := today.Year() - t.Year()
			if today.YearDay() < t.YearDay() {
				age--
			}
			return age
		},
		"ptrInt": func(i *int) int {
			if i == nil {
				return 0
			}
			return *i
		},
		"makeSlice": func(args ...interface{}) []interface{} {
			return args
		},
		"getRelatedPerson": func(persons interface{}, role string) interface{} {
			v := reflect.ValueOf(persons)
			if v.Kind() != reflect.Slice {
				return nil
			}

			for i := 0; i < v.Len(); i++ {
				p := v.Index(i)
				isParent := p.FieldByName("IsParent").Bool()
				isGuardian := p.FieldByName("IsGuardian").Bool()
				relName := p.FieldByName("Relationship").FieldByName("RelationshipName").String()

				switch role {
				case "Father":
					if isParent && relName == "Father" {
						return p.Interface()
					}
				case "Mother":
					if isParent && relName == "Mother" {
						return p.Interface()
					}
				case "Guardian":
					if isGuardian {
						return p.Interface()
					}
				}
			}

			return nil
		},
		"getEduBackground": func(schools interface{}, level string) interface{} {
			v := reflect.ValueOf(schools)
			if v.Kind() != reflect.Slice {
				return nil
			}
			for i := 0; i < v.Len(); i++ {
				s := v.Index(i)
				lvl := s.FieldByName("EducationalLevel").FieldByName("LevelName").String()
				if lvl == level {
					return s.Interface()
				}
			}
			return nil
		},
		"getConsultation": func(consults interface{}, profType string) interface{} {
			v := reflect.ValueOf(consults)
			if v.Kind() != reflect.Slice {
				return nil
			}
			for i := 0; i < v.Len(); i++ {
				c := v.Index(i)
				ctype := c.FieldByName("ProfessionalType").String()
				if ctype == profType {
					return c.Interface()
				}
			}
			return nil
		},
	})

	tmpl, err := tmpl.ParseFiles(templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	pdfBytes, err := s.gotenberg.ConvertHTML(ctx, buf.String())
	if err != nil {
		return nil, fmt.Errorf("failed to convert html to pdf: %w", err)
	}

	return pdfBytes, nil
}
