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

func TestPostMetadata(t *testing.T) {
	//gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/metadata", postMetadata)

	//req, err := http.NewRequest(http.MethodPost, "/metadata", nil)
	req, err := http.NewRequest(http.MethodPost, "/metadata", fileToBytes("invalid-1.yml"))
	if err != nil {
		t.Fatalf("Could not create request: %v", err)

		// Create a response recorder so you can inspect the response
		w := httptest.NewRecorder()

		// Perform the request
		router.ServeHTTP(w, req)
		fmt.Println(w.Body)
		fmt.Println(w.Code)

		// Check to see if the response was what you expected
		if w.Code == http.StatusCreated {
			t.Logf("Expected to get status %d is same ast %d\n", http.StatusCreated, w.Code)
		} else {
			t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusCreated, w.Code)
		}
	}
}

func TestGetMetadata(t *testing.T) {
	//gin.SetMode(gin.TestMode)
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
	fmt.Println(w.Body)

	// Check to see if the response was what you expected
	if w.Code == http.StatusOK {
		t.Logf("Expected to get status %d is same ast %d\n", http.StatusOK, w.Code)
	} else {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusOK, w.Code)
	}
}
