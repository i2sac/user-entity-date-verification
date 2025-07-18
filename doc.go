/*
Package userdate provides comprehensive date validation for user entities.

This package ensures that dates associated with user data (certifications, trainings, etc.)
are realistic and valid within the context of the user's lifetime.

# Features

- Comprehensive date validation against multiple criteria
- Entity-specific validation rules for different types of user data
- Age-based restrictions ensuring users meet minimum requirements
- Detailed error reporting with specific error codes
- Zero dependencies - pure Go implementation
- High performance and optimized for high-throughput applications

# Basic Usage

	// Create a user
	birthDate, _ := time.Parse("2006-01-02", "1990-05-15")
	user, err := userdate.NewUser("user123", birthDate, "John Doe")
	if err != nil {
		log.Fatal(err)
	}

	// Validate a certification date
	certDate, _ := time.Parse("2006-01-02", "2020-03-10")
	err = userdate.ValidateCertification(user, certDate)
	if err != nil {
		fmt.Printf("Validation failed: %v\n", err)
	}

# Validation Rules

The package implements several layers of validation:

General Date Validation:
- Date cannot be zero value
- Date cannot be before year 1800
- Date cannot be more than 200 years in the past
- Date cannot be in the future

User-Specific Validation:
- Entity dates cannot be before the user's birth date
- User's birth date cannot be in the future
- User's age cannot exceed 150 years

Entity-Specific Age Requirements:
- Certifications/Training/Education: Minimum age 5 years
- Employment: Minimum age 14 years
- Licenses: Minimum age 16 years

# Error Handling

All validation functions return structured errors with specific codes:

	err := userdate.ValidateCertification(user, certDate)
	if err != nil {
		if dateErr, ok := err.(*userdate.DateValidationError); ok {
			fmt.Printf("Error [%s]: %s\n", dateErr.Code, dateErr.Message)
		}
	}

Available error codes: INVALID_DATE, BEFORE_BIRTH, FUTURE_DATE, UNREALISTIC_AGE, INVALID_USER, DATE_TOO_OLD

# Performance

The package is designed for high-performance applications with minimal allocations
and efficient validation algorithms. Benchmark tests are included to ensure
performance remains optimal.
*/
package userdate
