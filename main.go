package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

type MaintainerInfo struct {
	Name  string `yaml:"name"`
	Email string `yaml:"email" binding:"required,email"`
}

type Metadata struct {
	Title       string           `yaml:"title"`
	Version     string           `yaml:"version"`
	Company     string           `yaml:"company"`
	Website     string           `yaml:"website"`
	Source      string           `yaml:"source"`
	License     string           `yaml:"license"`
	Maintainers []MaintainerInfo `yaml:"maintainers" binding:"required"`
	Description string           `yaml:"description"`
}

// define im-memory storage for the app_metadata
var metadata = []Metadata{}
var setMetadata = make(map[string]bool)

func main() {
	router := gin.Default()
	router.GET("/metadata", getMetadata)
	router.POST("/metadata", postMetadata)
	//router.GET("/metadata", getDataByKey)
	router.Run("localhost:8080")
}

// getMetadata: get the metadata in YAML format stored in-memory
func getMetadata(c *gin.Context) {
	fmt.Println("in getMetadata")
	paramPairs := c.Request.URL.Query() //map[string][][string]

	if len(paramPairs) == 0 {
		c.IndentedJSON(http.StatusOK, metadata)
	} else {
		fmt.Println("len of paramPairs:", len(paramPairs))
		for key, values := range paramPairs {
			fmt.Println("key, value(s) : ", key, values)
		}
		res := searchInMetadata(paramPairs)
		c.IndentedJSON(http.StatusOK, res)
	}
}

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

// postMetadata: adds one application metadata received from YAML in the request body
func postMetadata(c *gin.Context) {
	var newYaml Metadata
	var statusCode int
	var statusMsg string

	// call BindYAML to receive the data
	if err := c.BindYAML(&newYaml); err != nil {
		fmt.Println(err)
		return
	}

	// Validate the data, and set the Status
	if isValid, errorStr := validateRequest(&newYaml); !isValid {
		statusCode = http.StatusBadRequest
		statusMsg = errorStr
		//c.JSON(http.StatusBadRequest, errorStr)
	} else {
		/*	check if payload exists in-memory
			if exists: update the record
			if not exist: add the record
		*/
		title := newYaml.Title
		if _, ok := setMetadata[title]; ok {
			// payload exists
			statusCode = http.StatusBadRequest
			statusMsg = "Payload exists. Not inserted."
			//c.YAML(http.StatusBadRequest, "Payload exists. Not inserted.")
		} else {
			setMetadata[title] = true
			metadata = append(metadata, newYaml)
			statusCode = http.StatusCreated
			statusMsg = "Payload added."
		}
	}
	c.JSON(statusCode, statusMsg)
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

// func postMetadataV2(c *gin.Context) {
// 	var bodyBytes []byte
// 	if c.Request.Body != nil {
// 		bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
// 	}
// 	bodyString := string(bodyBytes)
// 	fmt.Printf("Body String: %s", bodyString)

// 	var m = Metadata{}
// 	err := yaml.Unmarshal([]byte(bodyString), &m)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// }
