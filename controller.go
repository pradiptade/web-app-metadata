package main

import (
	"net/http"
	"regexp"
	"strings"
)

func searchInMetadata(urlQuerystr map[string][]string) []Metadata {
	var res = []Metadata{}

	for _, elem := range metadata {
		if matchParameters(elem, urlQuerystr) {
			res = append(res, elem)
		}
	}
	return res
}

func matchParameters(elem Metadata, urlQuerystr map[string][]string) bool {
	for param, value := range urlQuerystr {
		if (param != "maintainers.name" && param != "maintainers.email") && len(value) > 1 {
			return false
		}
		switch queryParam := strings.TrimSpace(param); queryParam {
		case "title":
			if !strings.Contains(elem.Title, value[0]) {
				return false
			}
		case "version":
			if !strings.Contains(elem.Version, value[0]) {
				return false
			}
		case "description":
			if !strings.Contains(elem.Description, value[0]) {
				return false
			}
		case "company":
			if elem.Company != value[0] {
				return false
			}
		case "website":
			if elem.Website != value[0] {
				return false
			}
		case "source":
			if elem.Source != value[0] {
				return false
			}
		case "license":
			if elem.License != value[0] {
				return false
			}
		case "maintainers.name":
			for _, name := range value {
				found := false
				for _, metadaMaintainer := range elem.Maintainers {
					if metadaMaintainer.Name == name {
						found = true
					}
				}
				if !found {
					return false
				}
			}
		case "maintainers.email":
			for _, email := range value {
				found := false
				for _, metadaMaintainer := range elem.Maintainers {
					if metadaMaintainer.Email == email {
						found = true
					}
				}
				if !found {
					return false
				}
			}
		}
	}
	return true
}

func processPayload(yml Metadata) (int, string) {
	var statusCode int
	var statusMsg string

	// Validate the data, and set the Status
	if isValid, errorStr := validateRequest(&yml); !isValid {
		statusCode = http.StatusBadRequest
		statusMsg = errorStr
	} else {
		/*	check if payload exists in-memory
			if exists: update the record
			if not exist: add the record
		*/
		title := yml.Title
		if _, ok := setMetadata[title]; ok {
			// payload exists
			statusCode = http.StatusBadRequest
			statusMsg = "Payload exists. Not inserted."
		} else {
			setMetadata[title] = true
			metadata = append(metadata, yml)
			statusCode = http.StatusCreated
			statusMsg = "Payload added."
		}
	}
	return statusCode, statusMsg
}

func validateRequest(m *Metadata) (bool, string) {
	var isValid bool = true
	var emptyFields []string

	if m.Title == "" || len(m.Title) == 0 {
		emptyFields = append(emptyFields, "Title cannot be empty")
	}
	if m.Version == "" || len(m.Version) == 0 {
		emptyFields = append(emptyFields, "Version cannot be empty")
	}
	if m.Company == "" || len(m.Company) == 0 {
		emptyFields = append(emptyFields, "Company cannot be empty")
	}
	if m.Description == "" || len(m.Description) == 0 {
		emptyFields = append(emptyFields, "Description cannot be empty")
	}
	if m.License == "" || len(m.License) == 0 {
		emptyFields = append(emptyFields, "License cannot be empty")
	}
	if m.Source == "" || len(m.Source) == 0 {
		emptyFields = append(emptyFields, "Source cannot be empty")
	}
	if m.Website == "" || len(m.Website) == 0 {
		emptyFields = append(emptyFields, "Website cannot be empty")
	}
	if len(m.Maintainers) == 0 {
		emptyFields = append(emptyFields, "Maintainers cannot be empty")
	}
	if len(m.Maintainers) > 0 {
		for _, person := range m.Maintainers {
			if person.Name == "" {
				emptyFields = append(emptyFields, "Maintainer name cannot be empty")
			}
			if person.Email == "" {
				emptyFields = append(emptyFields, "Maintainer email cannot be empty")
			}
			if !isValidEmail(person.Email) {
				emptyFields = append(emptyFields, "Maintainer email address not correct")
			}
		}
	}
	errorStr := strings.Join(emptyFields, "-")
	if len(errorStr) > 0 {
		isValid = false
	}
	return isValid, errorStr
}

//isValidEmail validates email format
func isValidEmail(email string) bool {
	var rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if len(email) > 254 || !rxEmail.MatchString(email) {
		return false
	}
	return true
}
