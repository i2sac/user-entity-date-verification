// Package userdate provides comprehensive date validation for user entities.
// It ensures that dates associated with user data (certifications, trainings, etc.)
// are realistic and valid within the context of the user's lifetime.
package userdate

import (
	"fmt"
	"time"
)

// User represents a user entity with basic information for date validation
type User struct {
	ID        string    `json:"id"`
	BirthDate time.Time `json:"birth_date"`
	Name      string    `json:"name,omitempty"`
}

// DateValidationError represents an error during date validation
type DateValidationError struct {
	Message string
	Code    string
}

func (e *DateValidationError) Error() string {
	return fmt.Sprintf("date validation error [%s]: %s", e.Code, e.Message)
}

// Validation error codes
const (
	ErrCodeInvalidDate     = "INVALID_DATE"
	ErrCodeBeforeBirth     = "BEFORE_BIRTH"
	ErrCodeFutureDate      = "FUTURE_DATE"
	ErrCodeUnrealisticAge  = "UNREALISTIC_AGE"
	ErrCodeInvalidUser     = "INVALID_USER"
	ErrCodeDateTooOld      = "DATE_TOO_OLD"
)

// Constants for validation limits
const (
	MaxHumanAge     = 150 // Maximum realistic human age
	MinCertAge      = 5   // Minimum age for certifications
	MaxHistoryYears = 200 // Maximum years back in history to consider valid
)

// ValidateEntityDate validates a date for a user entity (certification, training, etc.)
func ValidateEntityDate(user *User, entityDate time.Time, entityType string) error {
	if user == nil {
		return &DateValidationError{
			Message: "user cannot be nil",
			Code:    ErrCodeInvalidUser,
		}
	}

	// Validate the user's birth date first
	if err := validateBirthDate(user.BirthDate); err != nil {
		return err
	}

	// Validate the entity date
	if err := validateDate(entityDate); err != nil {
		return err
	}

	// Check if date is before user's birth
	if entityDate.Before(user.BirthDate) {
		return &DateValidationError{
			Message: fmt.Sprintf("%s date (%s) cannot be before user's birth date (%s)",
				entityType, entityDate.Format("2006-01-02"), user.BirthDate.Format("2006-01-02")),
			Code: ErrCodeBeforeBirth,
		}
	}

	// Check if date is in the future
	now := time.Now()
	if entityDate.After(now) {
		return &DateValidationError{
			Message: fmt.Sprintf("%s date (%s) cannot be in the future",
				entityType, entityDate.Format("2006-01-02")),
			Code: ErrCodeFutureDate,
		}
	}

	// Check if user would be too young for certain entity types
	if err := validateMinimumAge(user.BirthDate, entityDate, entityType); err != nil {
		return err
	}

	// Check if date is unrealistically old
	if err := validateHistoricalRealism(entityDate); err != nil {
		return err
	}

	return nil
}

// validateDate performs basic date validation
func validateDate(date time.Time) error {
	// Check if date is zero value
	if date.IsZero() {
		return &DateValidationError{
			Message: "date cannot be zero value",
			Code:    ErrCodeInvalidDate,
		}
	}

	// Check if date is too far in the past (before year 1800)
	if date.Year() < 1800 {
		return &DateValidationError{
			Message: fmt.Sprintf("date year (%d) is too far in the past", date.Year()),
			Code:    ErrCodeDateTooOld,
		}
	}

	return nil
}

// validateBirthDate validates a user's birth date
func validateBirthDate(birthDate time.Time) error {
	if err := validateDate(birthDate); err != nil {
		return err
	}

	now := time.Now()
	age := now.Year() - birthDate.Year()

	// Adjust age if birthday hasn't occurred this year
	if now.YearDay() < birthDate.YearDay() {
		age--
	}

	// Check if birth date is in the future
	if birthDate.After(now) {
		return &DateValidationError{
			Message: "birth date cannot be in the future",
			Code:    ErrCodeFutureDate,
		}
	}

	// Check if age is unrealistic
	if age > MaxHumanAge {
		return &DateValidationError{
			Message: fmt.Sprintf("user age (%d) exceeds maximum realistic age (%d)", age, MaxHumanAge),
			Code:    ErrCodeUnrealisticAge,
		}
	}

	return nil
}

// validateMinimumAge checks if user meets minimum age requirements for certain entity types
func validateMinimumAge(birthDate, entityDate time.Time, entityType string) error {
	age := entityDate.Year() - birthDate.Year()
	if entityDate.YearDay() < birthDate.YearDay() {
		age--
	}

	// Define minimum ages for different entity types
	minAges := map[string]int{
		"certification": MinCertAge,
		"training":      MinCertAge,
		"education":     MinCertAge,
		"employment":    14, // Minimum working age in many countries
		"license":       16, // Typical minimum age for licenses
	}

	if minAge, exists := minAges[entityType]; exists {
		if age < minAge {
			return &DateValidationError{
				Message: fmt.Sprintf("user was too young (%d) for %s at date %s (minimum age: %d)",
					age, entityType, entityDate.Format("2006-01-02"), minAge),
				Code: ErrCodeUnrealisticAge,
			}
		}
	}

	return nil
}

// validateHistoricalRealism checks if the date is historically realistic
func validateHistoricalRealism(date time.Time) error {
	now := time.Now()
	yearsAgo := now.Year() - date.Year()

	if yearsAgo > MaxHistoryYears {
		return &DateValidationError{
			Message: fmt.Sprintf("date is too far in the past (%d years ago, maximum: %d)",
				yearsAgo, MaxHistoryYears),
			Code: ErrCodeDateTooOld,
		}
	}

	return nil
}

// ValidateCertification validates a certification date for a user
func ValidateCertification(user *User, certDate time.Time) error {
	return ValidateEntityDate(user, certDate, "certification")
}

// ValidateTraining validates a training date for a user
func ValidateTraining(user *User, trainingDate time.Time) error {
	return ValidateEntityDate(user, trainingDate, "training")
}

// ValidateEducation validates an education date for a user
func ValidateEducation(user *User, educationDate time.Time) error {
	return ValidateEntityDate(user, educationDate, "education")
}

// ValidateEmployment validates an employment date for a user
func ValidateEmployment(user *User, employmentDate time.Time) error {
	return ValidateEntityDate(user, employmentDate, "employment")
}

// ValidateLicense validates a license date for a user
func ValidateLicense(user *User, licenseDate time.Time) error {
	return ValidateEntityDate(user, licenseDate, "license")
}

// NewUser creates a new User with validation
func NewUser(id string, birthDate time.Time, name string) (*User, error) {
	if id == "" {
		return nil, &DateValidationError{
			Message: "user ID cannot be empty",
			Code:    ErrCodeInvalidUser,
		}
	}

	user := &User{
		ID:        id,
		BirthDate: birthDate,
		Name:      name,
	}

	// Validate birth date
	if err := validateBirthDate(birthDate); err != nil {
		return nil, err
	}

	return user, nil
}

// GetAge returns the current age of the user
func (u *User) GetAge() int {
	now := time.Now()
	age := now.Year() - u.BirthDate.Year()
	if now.YearDay() < u.BirthDate.YearDay() {
		age--
	}
	return age
}

// GetAgeAtDate returns the user's age at a specific date
func (u *User) GetAgeAtDate(date time.Time) int {
	age := date.Year() - u.BirthDate.Year()
	if date.YearDay() < u.BirthDate.YearDay() {
		age--
	}
	return age
}
