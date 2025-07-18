package userdate_test

import (
	"fmt"
	"log"
	"time"

	"github.com/i2sac/user-entity-date-verification"
)

func ExampleNewUser() {
	// Create a new user with validation
	birthDate, _ := time.Parse("2006-01-02", "1990-05-15")
	user, err := userdate.NewUser("user123", birthDate, "John Doe")
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("User created: %s (Age: %d)\n", user.Name, user.GetAge())
	// Output: User created: John Doe (Age: 35)
}

func ExampleValidateCertification() {
	// Create a user
	birthDate, _ := time.Parse("2006-01-02", "1990-05-15")
	user, _ := userdate.NewUser("user123", birthDate, "John Doe")
	
	// Valid certification date
	certDate, _ := time.Parse("2006-01-02", "2020-03-10")
	err := userdate.ValidateCertification(user, certDate)
	if err != nil {
		fmt.Printf("Validation failed: %v\n", err)
	} else {
		fmt.Println("Certification date is valid")
	}
	
	// Invalid certification date (before birth)
	invalidDate, _ := time.Parse("2006-01-02", "1989-01-01")
	err = userdate.ValidateCertification(user, invalidDate)
	if err != nil {
		fmt.Printf("Validation failed: %v\n", err)
	}
	
	// Output:
	// Certification date is valid
	// Validation failed: date validation error [BEFORE_BIRTH]: certification date (1989-01-01) cannot be before user's birth date (1990-05-15)
}

func ExampleValidateEntityDate() {
	// Create a user
	birthDate, _ := time.Parse("2006-01-02", "1990-05-15")
	user, _ := userdate.NewUser("user123", birthDate, "John Doe")
	
	// Validate different types of entities
	trainingDate, _ := time.Parse("2006-01-02", "2018-09-01")
	err := userdate.ValidateEntityDate(user, trainingDate, "training")
	if err != nil {
		fmt.Printf("Training validation failed: %v\n", err)
	} else {
		fmt.Println("Training date is valid")
	}
	
	// Validate employment (user must be at least 14)
	employmentDate, _ := time.Parse("2006-01-02", "2006-06-01") // User is 16
	err = userdate.ValidateEntityDate(user, employmentDate, "employment")
	if err != nil {
		fmt.Printf("Employment validation failed: %v\n", err)
	} else {
		fmt.Println("Employment date is valid")
	}
	
	// Output:
	// Training date is valid
	// Employment date is valid
}

func ExampleUser_GetAgeAtDate() {
	birthDate, _ := time.Parse("2006-01-02", "1990-05-15")
	user, _ := userdate.NewUser("user123", birthDate, "John Doe")
	
	// Get age at a specific date
	specificDate, _ := time.Parse("2006-01-02", "2020-05-15")
	age := user.GetAgeAtDate(specificDate)
	fmt.Printf("User's age on %s: %d\n", specificDate.Format("2006-01-02"), age)
	
	// Output: User's age on 2020-05-15: 30
}

func ExampleDateValidationError() {
	birthDate, _ := time.Parse("2006-01-02", "1990-05-15")
	user, _ := userdate.NewUser("user123", birthDate, "John Doe")
	
	// Try to validate a future date
	futureDate := time.Now().AddDate(1, 0, 0)
	err := userdate.ValidateCertification(user, futureDate)
	
	if err != nil {
		// Check if it's a DateValidationError
		if dateErr, ok := err.(*userdate.DateValidationError); ok {
			fmt.Printf("Error Code: %s\n", dateErr.Code)
			fmt.Printf("Error Message: %s\n", dateErr.Message)
		}
	}
	
	// Output:
	// Error Code: FUTURE_DATE
	// Error Message: certification date (2026-07-18) cannot be in the future
}
