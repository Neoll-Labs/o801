/*
 license x
*/

package openapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	_ "github.com/nelsonstr/o801/monitoring"
)

func TestMetricEndpoint(t *testing.T) {
	// Create a new HTTP request with the "/metric" endpoint
	req, err := http.NewRequest("GET", "/metric", nil)
	if err != nil {
		t.Fatal(err)
	}

	c := &MonitoringAPIController{

		errorHandler: DefaultErrorHandler,
	}

	// Create a new mock HTTP server
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(c.MetricsGet)
	handler.ServeHTTP(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body
	expected := "# TYPE promhttp_metric_handler_requests_total counter"

	if strings.Count(rr.Body.String(), expected) != 1 {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
