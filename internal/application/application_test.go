package application_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DobryySoul/yandex_repo/internal/application"
)

func TestCalcHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		body           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid Expression",
			method:         http.MethodPost,
			body:           `{"expression": "2 + 2 * 2"}`,
			expectedStatus: http.StatusOK,
			expectedBody:   "6",
		},
		{
			name:           "Hard valid expression",
			method:         http.MethodPost,
			body:           `{"expression": "(15-17)+8 + 15 * 10 / 15 + (((25-10)/15)+7)-19"}`,
			expectedStatus: http.StatusOK,
			expectedBody:   "5",
		},
		{
			name:           "Invalid Expression",
			method:         http.MethodPost,
			body:           `{"expression": "invalid"}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid expression",
		},
		{
			name:           "Method Not Allowed",
			method:         http.MethodGet,
			body:           `{"expression": "1+1"}`,
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "method not allowed",
		},
		{
			name:           "Bad JSON",
			method:         http.MethodPost,
			body:           `{"expression": "1+2"`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "unexpected end of JSON input",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/", bytes.NewBufferString(tt.body))
			w := httptest.NewRecorder()

			application.CalcHandler(w, req)

			res := w.Result()
			if res.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, res.StatusCode)
			}

			body := strings.TrimSpace(w.Body.String())
			if body != tt.expectedBody {
				t.Errorf("expected body %q, got %q", tt.expectedBody, body)
			}
		})
	}
}
