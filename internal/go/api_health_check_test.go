package openapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type monitiringTest struct {
	code int
	body string
}

func (m monitiringTest) HealthzGet(ctx context.Context) (ImplResponse, error) {
	return Response(m.code, m.body), nil
}

func (m monitiringTest) LivenessGet(ctx context.Context) (ImplResponse, error) {
	return Response(m.code, m.body), nil
}

func TestHealthCheckAPIController_HealthzGet(t *testing.T) {
	type fields struct {
		service      HealthCheckAPIServicer
		errorHandler ErrorHandler
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	type output struct {
		code         int
		bodyContains string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		output output
	}{
		{
			name: "ok",
			fields: fields{
				service:      &monitiringTest{code: 200, body: "ok"},
				errorHandler: DefaultErrorHandler,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "http://localhost", nil),
			},
			output: output{
				code:         200,
				bodyContains: "ok",
			},
		}, {
			name: "error",
			fields: fields{
				service:      &monitiringTest{code: 503, body: "not ok"},
				errorHandler: DefaultErrorHandler,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "http://localhost", nil),
			},
			output: output{
				code:         503,
				bodyContains: "not ok",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &HealthCheckAPIController{
				service:      tt.fields.service,
				errorHandler: tt.fields.errorHandler,
			}
			c.HealthzGet(tt.args.w, tt.args.r)

			assert.Equal(t, tt.output.code, tt.args.w.(*httptest.ResponseRecorder).Code)
			assert.Contains(t, tt.args.w.(*httptest.ResponseRecorder).Body.String(), tt.output.bodyContains)
		})
	}
}

func TestHealthCheckAPIController_LivenessGet(t *testing.T) {
	type fields struct {
		service      HealthCheckAPIServicer
		errorHandler ErrorHandler
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	type output struct {
		code int
		body string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		output output
	}{
		{
			name: "ok",
			fields: fields{
				service:      &monitiringTest{code: 200, body: "ok"},
				errorHandler: DefaultErrorHandler,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "http://localhost", nil),
			},
			output: output{
				code: 200,
				body: "ok",
			},
		}, {
			name: "error",
			fields: fields{
				service:      &monitiringTest{code: 503, body: "not ok"},
				errorHandler: DefaultErrorHandler,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "http://localhost", nil),
			},
			output: output{
				code: 503,
				body: "not ok",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &HealthCheckAPIController{
				service:      tt.fields.service,
				errorHandler: tt.fields.errorHandler,
			}
			c.LivenessGet(tt.args.w, tt.args.r)

			assert.Equal(t, tt.output.code, tt.args.w.(*httptest.ResponseRecorder).Code)
			assert.Contains(t, tt.args.w.(*httptest.ResponseRecorder).Body.String(), tt.output.body)
		})
	}
}

func TestHealthEndpoint(t *testing.T) {
	// Create a new HTTP request with the "/metric" endpoint
	req, err := http.NewRequest("GET", "/healthz", nil)
	if err != nil {
		t.Fatal(err)
	}

	c := &HealthCheckAPIController{
		service:      NewHealthCheckAPIService(nil),
		errorHandler: DefaultErrorHandler,
	}

	// Create a new mock HTTP server
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(c.HealthzGet)
	handler.ServeHTTP(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body
	expected := ""

	if strings.Count(rr.Body.String(), expected) != 1 {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestLivenessEndpoint(t *testing.T) {
	// Create a new HTTP request with the "/metric" endpoint
	req, err := http.NewRequest("GET", "/liveness", nil)
	if err != nil {
		t.Fatal(err)
	}

	c := &HealthCheckAPIController{
		service:      NewHealthCheckAPIService(nil),
		errorHandler: DefaultErrorHandler,
	}

	// Create a new mock HTTP server
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(c.LivenessGet)
	handler.ServeHTTP(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body
	expected := ""

	if strings.Count(rr.Body.String(), expected) != 1 {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
