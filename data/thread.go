package data

import (
	"self-scientists/config"
)

type Thread struct {
	Id           string `json:"id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
	CategoryName string `json:"categoryName"`
}

func (thread Thread) validateForCreation() []string {
	errors := []string{}

	if len(thread.Title) == 0 || len(thread.Description) == 0 || len(thread.CategoryName) == 0 {
		errors = append(errors, "All three must be provided: Title, Description, and CategoryName")
	}

	if len(errors) > 0 {
		return errors
	}

	var categoriesFound int
	row := config.DB.QueryRow("SELECT COUNT(*) FROM categories WHERE identifier=$1", thread.CategoryName)
	err := row.Scan(&categoriesFound)
	if err != nil {
		panic("Unable to connect to DB")
	}
	if categoriesFound != 0 {
		errors = append(errors, "No category with provided name exists")
	}

	return errors
}
