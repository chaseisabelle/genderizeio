package genderizer

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGenderize_success_singleName(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Type", "application/json; charset=utf-8")

		response.Write([]byte(`[{"name":"chase","gender":"male","probability":0.96,"count":296}]`))

		response.WriteHeader(http.StatusOK)
	}))

	defer server.Close()

	endpoint = server.URL

	genderizations, err := Genderize("chase")

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
		response.Header().Set("Content-Type", "application/json; charset=utf-8")

		response.Write([]byte(`[{"name":"chase","gender":"male","probability":0.96,"count":296},{"name":"isabelle","gender":"female","probability":1,"count":296}]`))

		response.WriteHeader(http.StatusOK)
	}))

	defer server.Close()

	endpoint = server.URL

	genderizations, err := Genderize("chase", "isabelle")

	if err != nil {
		t.Error(err)

		return
	}

	count := len(genderizations)

	if count != 2 {
		t.Errorf("Expected 2, but got %+v.", count)
	}
}

func TestGenderize_failure_noNames(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		t.Errorf("Server should not have been queried.")
	}))

	defer server.Close()

	endpoint = server.URL

	genderizations, err := Genderize()

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

	endpoint = server.URL

	genderizations, err := Genderize("chase", "")

	if err == nil {
		t.Errorf("Expected error.")
	}

	if genderizations != nil {
		t.Errorf("Expected nil.")
	}
}

func TestGenderize_failure_error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Type", "application/json; charset=utf-8")

		response.Write([]byte(`{"error":"something terrible has happened"}`))

		response.WriteHeader(http.StatusBadRequest)
	}))

	defer server.Close()

	endpoint = server.URL

	genderizations, err := Genderize("chase", "isabelle")

	if err == nil {
		t.Error("Expected error.")
	}

	if genderizations != nil {
		t.Error("Expected nil.")
	}
}
