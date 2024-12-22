package application_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/C1scoR/Go_Project_Yandex_Lyceum/internal/application"
)

func TestCalcHandler(t *testing.T) {
	testCases := []struct {
		name               string
		inputBody          string
		expectedResult     string
		expectedStatusCode int
	}{
		{
			name:               "valid expression",
			inputBody:          `{"expression": "2+2*2"}`,
			expectedResult:     "result: 6.000000, server status: 200",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "empty expression",
			inputBody:          `{"expression": ""}`,
			expectedResult:     fmt.Sprintf("Unprocessable entity (The length of given expression is 0), error status: %d\n", http.StatusUnprocessableEntity),
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:               "expression with letters",
			inputBody:          `{"expression": "2+A*2"}`,
			expectedResult:     fmt.Sprintf("Unprocessable entity (The expression contains letters), error status: %d\n", http.StatusUnprocessableEntity),
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:               "invalid expression",
			inputBody:          `{"expression": "2+*2"}`,
			expectedResult:     "unknown error\n",
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	handler := http.HandlerFunc(application.CalcHandler)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/api/v1/calculate", strings.NewReader(tc.inputBody))
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tc.expectedStatusCode {
				t.Errorf("handler returned wrong status code:\n got %v want %v", status, tc.expectedStatusCode)
			}

			if rr.Body.String() != tc.expectedResult {
				t.Errorf("handler returned unexpected body:\n got %v want %v", rr.Body.String(), tc.expectedResult)
			}
		})
	}
}
