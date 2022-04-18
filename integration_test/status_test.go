package integration_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
)

func TestStatus(t *testing.T) {
	host := os.Getenv("CHINNASWAMY_TEST_HOST")
	if host == "" {
		host = "localhost:8080"
	}

	resp, err := http.Get(fmt.Sprintf("http://%s/status", host))

	if assert.NoError(t, err, "Error making a GET request") {
		assert.Equal(t, 200, resp.StatusCode, "Expected 200 OK status")
	}
}
