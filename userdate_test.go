package userdate

import (
	"testing"
	"time"
)

// Test helper functions
func mustParseDate(dateStr string) time.Time {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		panic(err)
	}
	return date
}

func TestNewUser(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		birthDate time.Time
		userName  string
		wantErr   bool
		errCode   string
	}{
		{
			name:      "valid user",
			id:        "user123",
			birthDate: mustParseDate("1990-01-01"),
			userName:  "John Doe",
			wantErr:   false,
		},
		{
			name:      "empty ID",
			id:        "",
			birthDate: mustParseDate("1990-01-01"),
			userName:  "John Doe",
			wantErr:   true,
			errCode:   ErrCodeInvalidUser,
		},
		{
			name:      "future birth date",
			id:        "user123",
			birthDate: time.Now().AddDate(1, 0, 0),
			userName:  "John Doe",
			wantErr:   true,
			errCode:   ErrCodeFutureDate,
		},
		{
			name:      "unrealistic age",
			id:        "user123",
			birthDate: mustParseDate("1800-01-01"),
			userName:  "John Doe",
			wantErr:   true,
			errCode:   ErrCodeUnrealisticAge,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := NewUser(tt.id, tt.birthDate, tt.userName)
			if tt.wantErr {
				if err == nil {
					t.Errorf("NewUser() expected error but got none")
					return
				}
				if dateErr, ok := err.(*DateValidationError); ok {
					if dateErr.Code != tt.errCode {
						t.Errorf("NewUser() error code = %v, want %v", dateErr.Code, tt.errCode)
					}
				} else {
					t.Errorf("NewUser() error type = %T, want *DateValidationError", err)
				}
			} else if err != nil {
				t.Errorf("NewUser() unexpected error = %v", err)
			} else {
				if user.ID != tt.id {
					t.Errorf("NewUser() ID = %v, want %v", user.ID, tt.id)
				}
				if user.Name != tt.userName {
					t.Errorf("NewUser() Name = %v, want %v", user.Name, tt.userName)
				}
			}
		})
	}
}

func TestValidateEntityDate(t *testing.T) {
	validUser, _ := NewUser("user123", mustParseDate("1990-01-01"), "John Doe")

	tests := []struct {
		name       string
		user       *User
		entityDate time.Time
		entityType string
		wantErr    bool
		errCode    string
	}{
		{
			name:       "valid certification date",
			user:       validUser,
			entityDate: mustParseDate("2020-01-01"),
			entityType: "certification",
			wantErr:    false,
		},
		{
			name:       "nil user",
			user:       nil,
			entityDate: mustParseDate("2020-01-01"),
			entityType: "certification",
			wantErr:    true,
			errCode:    ErrCodeInvalidUser,
		},
		{
			name:       "date before birth",
			user:       validUser,
			entityDate: mustParseDate("1989-01-01"),
			entityType: "certification",
			wantErr:    true,
			errCode:    ErrCodeBeforeBirth,
		},
		{
			name:       "future date",
			user:       validUser,
			entityDate: time.Now().AddDate(1, 0, 0),
			entityType: "certification",
			wantErr:    true,
			errCode:    ErrCodeFutureDate,
		},
		{
			name:       "too young for certification",
			user:       validUser,
			entityDate: mustParseDate("1993-01-01"), // User would be 3 years old
			entityType: "certification",
			wantErr:    true,
			errCode:    ErrCodeUnrealisticAge,
		},
		{
			name:       "zero date",
			user:       validUser,
			entityDate: time.Time{},
			entityType: "certification",
			wantErr:    true,
			errCode:    ErrCodeInvalidDate,
		},
		{
			name:       "date too old",
			user:       validUser,
			entityDate: mustParseDate("1799-01-01"),
			entityType: "certification",
			wantErr:    true,
			errCode:    ErrCodeDateTooOld,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEntityDate(tt.user, tt.entityDate, tt.entityType)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateEntityDate() expected error but got none")
					return
				}
				if dateErr, ok := err.(*DateValidationError); ok {
					if dateErr.Code != tt.errCode {
						t.Errorf("ValidateEntityDate() error code = %v, want %v", dateErr.Code, tt.errCode)
					}
				} else {
					t.Errorf("ValidateEntityDate() error type = %T, want *DateValidationError", err)
				}
			} else if err != nil {
				t.Errorf("ValidateEntityDate() unexpected error = %v", err)
			}
		})
	}
}

func TestValidateCertification(t *testing.T) {
	user, _ := NewUser("user123", mustParseDate("1990-01-01"), "John Doe")

	// Valid certification
	err := ValidateCertification(user, mustParseDate("2020-01-01"))
	if err != nil {
		t.Errorf("ValidateCertification() unexpected error = %v", err)
	}

	// Invalid certification (too young)
	err = ValidateCertification(user, mustParseDate("1993-01-01"))
	if err == nil {
		t.Errorf("ValidateCertification() expected error but got none")
	}
}

func TestValidateTraining(t *testing.T) {
	user, _ := NewUser("user123", mustParseDate("1990-01-01"), "John Doe")

	// Valid training
	err := ValidateTraining(user, mustParseDate("2020-01-01"))
	if err != nil {
		t.Errorf("ValidateTraining() unexpected error = %v", err)
	}

	// Invalid training (future date)
	err = ValidateTraining(user, time.Now().AddDate(1, 0, 0))
	if err == nil {
		t.Errorf("ValidateTraining() expected error but got none")
	}
}

func TestValidateEducation(t *testing.T) {
	user, _ := NewUser("user123", mustParseDate("1990-01-01"), "John Doe")

	// Valid education
	err := ValidateEducation(user, mustParseDate("2010-01-01"))
	if err != nil {
		t.Errorf("ValidateEducation() unexpected error = %v", err)
	}
}

func TestValidateEmployment(t *testing.T) {
	user, _ := NewUser("user123", mustParseDate("1990-01-01"), "John Doe")

	// Valid employment (user is 16)
	err := ValidateEmployment(user, mustParseDate("2006-01-01"))
	if err != nil {
		t.Errorf("ValidateEmployment() unexpected error = %v", err)
	}

	// Invalid employment (too young - user is 13)
	err = ValidateEmployment(user, mustParseDate("2003-01-01"))
	if err == nil {
		t.Errorf("ValidateEmployment() expected error but got none")
	}
}

func TestValidateLicense(t *testing.T) {
	user, _ := NewUser("user123", mustParseDate("1990-01-01"), "John Doe")

	// Valid license (user is 18)
	err := ValidateLicense(user, mustParseDate("2008-01-01"))
	if err != nil {
		t.Errorf("ValidateLicense() unexpected error = %v", err)
	}

	// Invalid license (too young - user is 15)
	err = ValidateLicense(user, mustParseDate("2005-01-01"))
	if err == nil {
		t.Errorf("ValidateLicense() expected error but got none")
	}
}

func TestUserGetAge(t *testing.T) {
	user, _ := NewUser("user123", mustParseDate("1990-01-01"), "John Doe")

	age := user.GetAge()
	expectedAge := time.Now().Year() - 1990
	if time.Now().YearDay() < mustParseDate("1990-01-01").YearDay() {
		expectedAge--
	}

	if age != expectedAge {
		t.Errorf("GetAge() = %v, want %v", age, expectedAge)
	}
}

func TestUserGetAgeAtDate(t *testing.T) {
	user, _ := NewUser("user123", mustParseDate("1990-01-01"), "John Doe")

	testDate := mustParseDate("2020-01-01")
	age := user.GetAgeAtDate(testDate)
	expectedAge := 30 // 2020 - 1990

	if age != expectedAge {
		t.Errorf("GetAgeAtDate() = %v, want %v", age, expectedAge)
	}
}

func TestDateValidationError(t *testing.T) {
	err := &DateValidationError{
		Message: "test error",
		Code:    "TEST_CODE",
	}

	expected := "date validation error [TEST_CODE]: test error"
	if err.Error() != expected {
		t.Errorf("DateValidationError.Error() = %v, want %v", err.Error(), expected)
	}
}

func TestValidateMinimumAge(t *testing.T) {
	birthDate := mustParseDate("1990-01-01")

	tests := []struct {
		name       string
		entityDate time.Time
		entityType string
		wantErr    bool
	}{
		{
			name:       "valid certification age",
			entityDate: mustParseDate("1996-01-01"), // Age 6
			entityType: "certification",
			wantErr:    false,
		},
		{
			name:       "invalid certification age",
			entityDate: mustParseDate("1993-01-01"), // Age 3
			entityType: "certification",
			wantErr:    true,
		},
		{
			name:       "valid employment age",
			entityDate: mustParseDate("2004-01-01"), // Age 14
			entityType: "employment",
			wantErr:    false,
		},
		{
			name:       "invalid employment age",
			entityDate: mustParseDate("2003-01-01"), // Age 13
			entityType: "employment",
			wantErr:    true,
		},
		{
			name:       "unknown entity type",
			entityDate: mustParseDate("1991-01-01"), // Age 1
			entityType: "unknown",
			wantErr:    false, // Should not validate unknown types
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateMinimumAge(birthDate, tt.entityDate, tt.entityType)
			if tt.wantErr && err == nil {
				t.Errorf("validateMinimumAge() expected error but got none")
			} else if !tt.wantErr && err != nil {
				t.Errorf("validateMinimumAge() unexpected error = %v", err)
			}
		})
	}
}

func TestValidateHistoricalRealism(t *testing.T) {
	tests := []struct {
		name    string
		date    time.Time
		wantErr bool
	}{
		{
			name:    "recent date",
			date:    mustParseDate("2020-01-01"),
			wantErr: false,
		},
		{
			name:    "old but valid date",
			date:    mustParseDate("1900-01-01"),
			wantErr: false,
		},
		{
			name:    "too old date",
			date:    mustParseDate("1800-01-01"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateHistoricalRealism(tt.date)
			if tt.wantErr && err == nil {
				t.Errorf("validateHistoricalRealism() expected error but got none")
			} else if !tt.wantErr && err != nil {
				t.Errorf("validateHistoricalRealism() unexpected error = %v", err)
			}
		})
	}
}

// Benchmark tests
func BenchmarkValidateEntityDate(b *testing.B) {
	user, err := NewUser("user123", mustParseDate("1990-01-01"), "John Doe")
	if err != nil {
		b.Fatal(err)
	}
	entityDate := mustParseDate("2020-01-01")

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = ValidateEntityDate(user, entityDate, "certification")
	}
}

func BenchmarkNewUser(b *testing.B) {
	birthDate := mustParseDate("1990-01-01")

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = NewUser("user123", birthDate, "John Doe")
	}
}
