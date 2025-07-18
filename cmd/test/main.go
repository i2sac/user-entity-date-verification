// test is a simple application to manually test the userdate library.
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/i2sac/user-entity-date-verification"
)

// parseTime is a helper function to parse time strings and exit on error.
func parseTime(layout, value string) time.Time {
	t, err := time.Parse(layout, value)
	if err != nil {
		log.Fatalf("Error parsing time %q: %v", value, err)
	}
	return t
}

func main() {
	fmt.Println("Testing User Entity Date Verification Library")
	fmt.Println("===========================================")

	// Test 1: Create a valid user
	fmt.Println("\n1. Testing NewUser with valid data:")
	birthDate := parseTime("2006-01-02", "1990-05-15")
	user, err := userdate.NewUser("user123", birthDate, "John Doe")
	if err != nil {
		log.Printf("Error creating user: %v", err)
	} else {
		fmt.Printf("✓ User created successfully: %s (Age: %d)\n", user.Name, user.GetAge())
	}

	// Test 2: Validate a certification date
	fmt.Println("\n2. Testing ValidateCertification with valid date:")
	certDate := parseTime("2006-01-02", "2020-03-10")
	err = userdate.ValidateCertification(user, certDate)
	if err != nil {
		log.Printf("✗ Certification validation failed: %v", err)
	} else {
		fmt.Println("✓ Certification date is valid")
	}

	// Test 3: Test invalid certification (before birth)
	fmt.Println("\n3. Testing ValidateCertification with invalid date (before birth):")
	invalidDate := parseTime("2006-01-02", "1989-01-01")
	err = userdate.ValidateCertification(user, invalidDate)
	if err != nil {
		fmt.Printf("✓ Expected error caught: %v\n", err)
	} else {
		fmt.Println("✗ Should have failed validation")
	}

	// Test 4: Test future date
	fmt.Println("\n4. Testing ValidateCertification with future date:")
	futureDate := time.Now().AddDate(1, 0, 0)
	err = userdate.ValidateCertification(user, futureDate)
	if err != nil {
		fmt.Printf("✓ Expected error caught: %v\n", err)
	} else {
		fmt.Println("✗ Should have failed validation")
	}

	// Test 5: Test employment validation
	fmt.Println("\n5. Testing ValidateEmployment:")
	employmentDate := parseTime("2006-01-02", "2006-06-01") // User is 16
	err = userdate.ValidateEmployment(user, employmentDate)
	if err != nil {
		log.Printf("✗ Employment validation failed: %v", err)
	} else {
		fmt.Println("✓ Employment date is valid")
	}

	// Test 6: Test too young for employment
	fmt.Println("\n6. Testing ValidateEmployment with too young age:")
	tooYoungDate := parseTime("2006-01-02", "2003-01-01") // User is 13
	err = userdate.ValidateEmployment(user, tooYoungDate)
	if err != nil {
		fmt.Printf("✓ Expected error caught: %v\n", err)
	} else {
		fmt.Println("✗ Should have failed validation")
	}

	// Test 7: Test age calculation
	fmt.Println("\n7. Testing age calculation:")
	specificDate := parseTime("2006-01-02", "2020-05-15")
	age := user.GetAgeAtDate(specificDate)
	fmt.Printf("✓ User's age on %s: %d\n", specificDate.Format("2006-01-02"), age)

	// Test 8: Test error handling
	fmt.Println("\n8. Testing error handling:")
	err = userdate.ValidateCertification(nil, certDate)
	if err != nil {
		if dateErr, ok := err.(*userdate.DateValidationError); ok {
			fmt.Printf("✓ Error code: %s, Message: %s\n", dateErr.Code, dateErr.Message)
		}
	}

	fmt.Println("\n===========================================")
	fmt.Println("All tests completed!")
}
