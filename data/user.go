package data

import (
	"reflect"
	"self-scientists/config"
	"self-scientists/utils"
	"self-scientists/validation"
	"strings"
)

type User struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"Email"`
	DateOfBirth string `json:"dateOfBirth"`
	Password    string `json:"password"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
	Bio         string `json:"bio"`
}

var nonReturnableUserFields []string = []string{"password"}

func (user User) validateForCreation() (errors []string) {
	// errors := []string{}
	if len(user.FirstName) == 0 || len(user.LastName) == 0 {
		errors = append(errors, "First name and last name have to be provided")
	}
	if !validation.IsEmail(user.Email) {
		errors = append(errors, "Email has to be provided and valid")
	} else {
		var usersFound int
		row := config.DB.QueryRow("SELECT COUNT(*) FROM users WHERE email=$1", user.Email)
		err := row.Scan(&usersFound)
		if err != nil {
			panic("Unable to connect to DB")
		}
		if usersFound != 0 {
			errors = append(errors, "Email already in use")
		}
	}
	if !validation.IsDateDDMMYYYY(user.DateOfBirth) {
		errors = append(errors, "Date of Birth has to be provided and in dd/mm/yyyy format")
	}
	if !validation.IsValidPassword(user.Password) {
		errors = append(errors, "Password has to be provided and at least 8 characters long. Also: "+validation.PasswordRequirementString)
	}

	return errors
}

func (user *User) CreateUser() (errors []string) {
	// var errors []string = []string{}
	var userCount int
	errors = user.validateForCreation()
	if len(errors) != 0 {
		return errors
	}
	row := config.DB.QueryRow("SELECT COUNT(*) FROM users WHERE email=$1", user.Email)
	row.Scan(&userCount)
	if userCount > 0 {
		errors = append(errors, "Email is already in use")
	}
	return errors
}

// Gets a map of struct data with blacklisted fields removed
func (user User) GetJSONAllowedData() map[string]interface{} {
	returnable := map[string]interface{}{}
	fields := reflect.TypeOf(user)
	values := reflect.ValueOf(user)
	numFields := fields.NumField()

	for i := 0; i < numFields; i++ {
		fieldRaw := fields.Field(i)
		fieldName := fieldRaw.Name
		field := strings.ToLower(fieldName[0:1]) + fieldName[1:]
		_, found := utils.FindString(nonReturnableUserFields, field)
		if !found {
			continue
		}
		value := values.Field(i).Interface()
		returnable[field] = value
	}

	return returnable
}
