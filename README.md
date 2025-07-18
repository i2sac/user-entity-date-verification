# User Entity Date Verification

[![Go Reference](https://pkg.go.dev/badge/github.com/i2sac/user-entity-date-verification.svg)](https://pkg.go.dev/github.com/i2sac/user-entity-date-verification)
[![Go Report Card](https://goreportcard.com/badge/github.com/i2sac/user-entity-date-verification)](https://goreportcard.com/report/github.com/i2sac/user-entity-date-verification)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Coverage Status](https://coveralls.io/repos/github/i2sac/user-entity-date-verification/badge.svg?branch=main)](https://coveralls.io/github/i2sac/user-entity-date-verification?branch=main)

A comprehensive Go library for validating dates associated with user entities such as certifications, trainings, education, employment, and licenses. This library ensures that dates are realistic and valid within the context of a user's lifetime.

## Features

- **Comprehensive Date Validation**: Validates dates against multiple criteria including user birth date, current date, and historical realism
- **Entity-Specific Validation**: Different validation rules for different types of entities (certifications, employment, licenses, etc.)
- **Age-Based Restrictions**: Ensures users meet minimum age requirements for specific entity types
- **Detailed Error Reporting**: Structured error responses with specific error codes for different validation failures
- **Zero Dependencies**: Pure Go implementation with no external dependencies
- **High Performance**: Optimized for high-throughput applications
- **Comprehensive Test Coverage**: Extensive test suite with benchmarks

## Installation

```bash
go get github.com/i2sac/user-entity-date-verification
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    "time"
    
    "github.com/i2sac/user-entity-date-verification"
)

func main() {
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
    } else {
        fmt.Println("Certification date is valid!")
    }
}
```

## Validation Rules

### General Date Validation
- Date cannot be zero value
- Date cannot be before year 1800
- Date cannot be more than 200 years in the past
- Date cannot be in the future

### User-Specific Validation
- Entity dates cannot be before the user's birth date
- User's birth date cannot be in the future
- User's age cannot exceed 150 years (configurable)

### Entity-Specific Age Requirements
- **Certifications/Training/Education**: Minimum age 5 years
- **Employment**: Minimum age 14 years
- **Licenses**: Minimum age 16 years

## API Reference

### Types

#### User
```go
type User struct {
    ID        string    `json:"id"`
    BirthDate time.Time `json:"birth_date"`
    Name      string    `json:"name,omitempty"`
}
```

#### DateValidationError
```go
type DateValidationError struct {
    Message string
    Code    string
}
```

### Functions

#### NewUser
```go
func NewUser(id string, birthDate time.Time, name string) (*User, error)
```
Creates a new User with validation of the birth date.

#### ValidateEntityDate
```go
func ValidateEntityDate(user *User, entityDate time.Time, entityType string) error
```
Validates a date for a user entity with comprehensive checks.

#### Convenience Functions
```go
func ValidateCertification(user *User, certDate time.Time) error
func ValidateTraining(user *User, trainingDate time.Time) error
func ValidateEducation(user *User, educationDate time.Time) error
func ValidateEmployment(user *User, employmentDate time.Time) error
func ValidateLicense(user *User, licenseDate time.Time) error
```

### Methods

#### User.GetAge
```go
func (u *User) GetAge() int
```
Returns the current age of the user.

#### User.GetAgeAtDate
```go
func (u *User) GetAgeAtDate(date time.Time) int
```
Returns the user's age at a specific date.

## Error Codes

| Code | Description |
|------|-------------|
| `INVALID_DATE` | Date is zero value or invalid |
| `BEFORE_BIRTH` | Date is before user's birth date |
| `FUTURE_DATE` | Date is in the future |
| `UNREALISTIC_AGE` | User's age is unrealistic or too young for entity type |
| `INVALID_USER` | User is nil or has invalid data |
| `DATE_TOO_OLD` | Date is too far in the past |

## Examples

### Basic Validation
```go
user, _ := userdate.NewUser("user123", birthDate, "John Doe")
err := userdate.ValidateCertification(user, certDate)
if err != nil {
    if dateErr, ok := err.(*userdate.DateValidationError); ok {
        fmt.Printf("Error [%s]: %s\n", dateErr.Code, dateErr.Message)
    }
}
```

### Custom Entity Type Validation
```go
err := userdate.ValidateEntityDate(user, entityDate, "custom_entity")
```

### Age Calculation
```go
currentAge := user.GetAge()
ageAtCertification := user.GetAgeAtDate(certDate)
```

## Configuration

The library includes several configurable constants:

```go
const (
    MaxHumanAge     = 150 // Maximum realistic human age
    MinCertAge      = 5   // Minimum age for certifications
    MaxHistoryYears = 200 // Maximum years back in history
)
```

## Testing

Run the test suite:

```bash
go test -v
```

Run benchmarks:

```bash
go test -bench=.
```

Generate coverage report:

```bash
go test -cover
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines

- Write tests for all new functionality
- Ensure all tests pass before submitting PR
- Follow Go best practices and conventions
- Update documentation for any API changes
- Maintain backward compatibility when possible

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Changelog

### v1.0.0
- Initial release
- Core date validation functionality
- Entity-specific validation rules
- Comprehensive test suite
- Full documentation

## Support

If you encounter any issues or have questions, please [open an issue](https://github.com/i2sac/user-entity-date-verification/issues) on GitHub.

## Acknowledgments

- Inspired by the need for robust date validation in user management systems
- Built with Go's excellent standard library
- Follows Go community best practices and conventions