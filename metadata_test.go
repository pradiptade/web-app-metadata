package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func fileToBytes(filepath string) *bytes.Reader {
	var bodyReader *bytes.Reader
	if dat, err := os.ReadFile(filepath); err != nil {
		fmt.Println(err)

	} else {
		bodyReader = bytes.NewReader(dat)
	}
	return bodyReader
}

func TestPostMetadataValid(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/metadata", postMetadata)

	//req, err := http.NewRequest(http.MethodPost, "/metadata", nil)
	req, err := http.NewRequest(http.MethodPost, "/metadata", fileToBytes("valid-1.yml"))
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}
	// Create a response recorder so you can inspect the response
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Check to see if the response was what you expected
	if w.Code == http.StatusCreated {
		t.Logf("Expected to get status %d is same ast %d\n", http.StatusCreated, w.Code)
	} else {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusCreated, w.Code)
	}
}

func TestPostMetadataInvalid(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/metadata", postMetadata)

	//req, err := http.NewRequest(http.MethodPost, "/metadata", nil)
	req, err := http.NewRequest(http.MethodPost, "/metadata", fileToBytes("invalid-1.yml"))
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}
	// Create a response recorder so you can inspect the response
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Check to see if the response was what you expected
	if w.Code == http.StatusBadRequest {
		t.Logf("Expected to get status %d is same ast %d\n", http.StatusBadRequest, w.Code)
	} else {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusCreated, w.Code)
	}
}

func TestGetMetadataAll(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/metadata", getMetadata)

	req, err := http.NewRequest(http.MethodGet, "/metadata", nil)
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}

	// Create a response recorder so you can inspect the response
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Check to see if the response was what you expected
	if w.Code == http.StatusOK {
		t.Logf("Expected to get status %d is same ast %d\n", http.StatusOK, w.Code)
	} else {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusOK, w.Code)
	}
}

// TODO: need to fix. Not able to setup the data for searching.
func TestGetMetadataByParamsTrue(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/metadata", getMetadata)
	router.POST("/metadata", postMetadata)

	// Setup the data first.
	_, e := http.NewRequest(http.MethodPost, "/metadata", fileToBytes("valid-1.yml"))
	if e != nil {
		t.Fatalf("Could not create request: %v", e)
	}

	// Setup the actual test
	req, err := http.NewRequest(http.MethodGet, "/metadata?title=Valid+App+1&version=0.0.1", nil)
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}

	// Create a response recorder so you can inspect the response
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)
	fmt.Println(w.Body)

	// Check to see if the response was what you expected
	if w.Code == http.StatusOK {
		t.Logf("Expected to get status %d is same ast %d\n", http.StatusOK, w.Code)
	} else {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusOK, w.Code)
	}
	if w.Body.String() != "[]" {
		t.Logf("Expected to get non-empty response, and got non-empty response\n")
	} else {
		t.Fatalf("Expected to get non-empty response, instead got empty: %s\n", w.Body.String())
	}
}

func TestGetMetadataByParamsFalse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/metadata", getMetadata)
	router.POST("/metadata", postMetadata)

	// SETTING UP TEST: first insert valid data that can be searched.
	_, e := http.NewRequest(http.MethodPost, "/metadata", fileToBytes("valid-1.yml"))
	if e != nil {
		t.Fatalf("Could not create request: %v", e)
	}

	req, err := http.NewRequest(http.MethodGet, "/metadata?junk=2.0.1", nil)
	if err != nil {
		t.Fatalf("Couldn't create request: %s\n", err)
	}

	// Create a response recorder so you can inspect the response
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Check to see if the response was what you expected
	if w.Code == http.StatusOK {
		t.Logf("Expected to get status %d is same ast %d\n", http.StatusOK, w.Code)
	} else {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusOK, w.Code)
	}
	if w.Body.String() == "[]" {
		t.Logf("Expected to get empty response, and got empty response\n")
	} else {
		t.Fatalf("Expected to get empty response, instead got: %v\n", w.Body.String())
	}
}
