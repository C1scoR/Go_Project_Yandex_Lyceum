package orchestrator

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestInfixToPostfix(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
		err      error
	}{
		{"simple addition", "1 + 2", []string{"1", "2", "+"}, nil},
		{"simple subtraction", "3 - 4", []string{"3", "4", "-"}, nil},
		{"simple multiplication", "5 * 6", []string{"5", "6", "*"}, nil},
		{"simple division", "7 / 8", []string{"7", "8", "/"}, nil},
		{"complex expression", "(1 + 2) * (3 - 4)", []string{"1", "2", "+", "3", "4", "-", "*"}, nil},
		{"invalid expression", "1 +", []string{}, ErrInvalidExpression},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := InfixToPostfix(tt.input)
			if err != tt.err {
				t.Errorf("expected error %v, got %v", tt.err, err)
			}
			if len(result) != len(tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("expected %v, got %v", tt.expected, result)
					break
				}
			}
		})
	}
}

func TestEvalRPN(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected float64
		err      error
	}{
		{"simple addition", []string{"1", "2", "+"}, 3, nil},
		{"simple subtraction", []string{"3", "4", "-"}, -1, nil},
		{"simple multiplication", []string{"5", "6", "*"}, 30, nil},
		{"simple division", []string{"7", "8", "/"}, 0.875, nil},
		{"complex expression", []string{"1", "2", "+", "3", "4", "-", "*"}, -3, nil},
		{"division by zero", []string{"1", "0", "/"}, 0, ErrDivisionByZero},
		{"invalid expression", []string{"1", "+"}, 0, ErrInvalidExpression},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := evalRPN(tt.input, "")
			if err != tt.err {
				t.Errorf("expected error %v, got %v", tt.err, err)
			}
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestOrchestratorHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
	}{
		{"valid expression", `{"expression": "1 + 2"}`, http.StatusCreated},
		{"empty expression", `{"expression": ""}`, http.StatusUnprocessableEntity},
		{"invalid expression", `{"expression": "1 +"}`, http.StatusUnprocessableEntity},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/api/v1/calculate", strings.NewReader(tt.requestBody))
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(OrchestratorHandler)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
		})
	}
}

func TestGetExpressionsHandler(t *testing.T) {
	Expressions_storage_variable.Expressions = append(Expressions_storage_variable.Expressions, Expressions_parametres{
		ID:     "test-id",
		Status: "created",
		Result: "1 + 2",
	})
	req, err := http.NewRequest("GET", "/api/v1/expressions", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetExpressionsHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestGetExpressionByIdHandler(t *testing.T) {
	Expressions_storage_variable.Expressions = append(Expressions_storage_variable.Expressions, Expressions_parametres{
		ID:     "test-id",
		Status: "created",
		Result: "1 + 2",
	})

	req, err := http.NewRequest("GET", "/api/v1/expressions/test-id", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetExpressionByIdHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}
<<<<<<< HEAD
=======

func TestHandlerForCommunicationToOtherServer(t *testing.T) {
	//вот тут по сути шаблон теста, но по факту как-то особо это затестировать тесты нельзя будет :)
	tests := []struct {
		name           string
		method         string
		expectedStatus int
	}{
		{"get tasks", http.MethodGet, http.StatusCreated},
		{"post results", http.MethodPost, http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, "/internal/task", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(HandlerForCommunicationToOtherServer)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
		})
	}
}
>>>>>>> 686799b (Pushing SuperCalculator)
