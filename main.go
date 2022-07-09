package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

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

	statusCode, statusMsg = processPayload(newYaml)
	c.JSON(statusCode, statusMsg)
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
