package genderizer

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func genderizer(endpoint string) *Genderizer {
	return &Genderizer{
		Client:   &http.Client{},
		Endpoint: endpoint,
		Key:      "",
	}
}

func TestGenderize_success_singleName(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.WriteHeader(http.StatusOK)
		response.Header().Set("Content-Type", "application/json; charset=utf-8")
		response.Write([]byte(`[{"name":"chase","gender":"male","probability":0.96,"count":296}]`))
	}))

	defer server.Close()

	genderizations, err := genderizer(server.URL).Genderize("chase")

	if err != nil {
		t.Error(err)

		return
	}

	count := len(genderizations)

	if count != 1 {
		t.Errorf("Expected 1, but got %+v.", count)
	}
}

func TestGenderize_success_multipleNames(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.WriteHeader(http.StatusOK)
		response.Header().Set("Content-Type", "application/json; charset=utf-8")
		response.Write([]byte(`[{"name":"chase","gender":"male","probability":0.96,"count":296},{"name":"isabelle","gender":"female","probability":1,"count":296}]`))
	}))

	defer server.Close()

	genderizations, err := genderizer(server.URL).Genderize("chase", "isabelle")

	if err != nil {
		t.Error(err)

		return
	}

	count := len(genderizations)

	if count != 2 {
		t.Errorf("Expected 2, but got %+v.", count)
	}
}

func TestGenderize_success_stupidName(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.WriteHeader(http.StatusOK)
		response.Header().Set("Content-Type", "application/json; charset=utf-8")
		response.Write([]byte(`[{"name":"ThisIsNotAName","gender":null}]`))
	}))

	defer server.Close()

	genderizations, err := genderizer(server.URL).Genderize("this is not a name")

	if err != nil {
		t.Error(err)

		return
	}

	count := len(genderizations)

	if count != 1 {
		t.Errorf("Expected 1, but got %+v.", count)
	}

	genderization := genderizations[0]

	if genderization.Name != "ThisIsNotAName" {
		t.Errorf("Expected ThisIsNotAName, but got %+v.", genderization.Name)
	}

	if genderization.Gender != "" {
		t.Errorf("Expected , but got %+v.", genderization.Gender)
	}

	if genderization.Probability != 0 {
		t.Errorf("Expected 0, but got %+v.", genderization.Probability)
	}

	if genderization.Count != 0 {
		t.Errorf("Expected 0, but got %+v.", genderization.Count)
	}
}

func TestGenderize_failure_noNames(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		t.Errorf("Server should not have been queried.")
	}))

	defer server.Close()

	genderizations, err := genderizer(server.URL).Genderize()

	if err == nil {
		t.Errorf("Expected error.")
	}

	if genderizations != nil {
		t.Errorf("Expected nil.")
	}
}

func TestGenderize_failure_emptyName(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		t.Errorf("Server should not have been queried.")
	}))

	defer server.Close()

	genderizations, err := genderizer(server.URL).Genderize("chase", "")

	if err == nil {
		t.Errorf("Expected error.")
	}

	if genderizations != nil {
		t.Errorf("Expected nil.")
	}
}

func TestGenderize_failure_apiError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.WriteHeader(http.StatusBadRequest)
		response.Header().Set("Content-Type", "application/json; charset=utf-8")
		response.Write([]byte(`{"error":"something terrible has happened"}`))
	}))

	defer server.Close()

	genderizations, err := genderizer(server.URL).Genderize("chase", "isabelle")

	if err == nil {
		t.Error("Expected error.")
	}

	if genderizations != nil {
		t.Error("Expected nil.")
	}
}

func TestGenderize_failure_serverError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`Internal Server Error`))
	}))

	defer server.Close()

	genderizations, err := genderizer(server.URL).Genderize("chase")

	if err == nil {
		t.Errorf("Expected error, but got %+v.", genderizations)

		return
	}

	actual := err.Error()
	expected := http.StatusText(http.StatusInternalServerError)

	if actual != expected {
		t.Errorf("Expected %s, but got %+v.", expected, actual)
	}
}
