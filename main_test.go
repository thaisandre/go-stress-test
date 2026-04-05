package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRunLoadTest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	report := runLoadTest(server.URL, 100, 10)

	if report.TotalReqs != 100 {
		t.Errorf("esperado 100 requests, obteve %d", report.TotalReqs)
	}
	if report.Status200 != 100 {
		t.Errorf("esperado 100 status 200, obteve %d", report.Status200)
	}
	if report.ErrorCount != 0 {
		t.Errorf("esperado 0 erros, obteve %d", report.ErrorCount)
	}
}

func TestRunLoadTestMixedStatus(t *testing.T) {
	count := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count++
		if count%3 == 0 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	report := runLoadTest(server.URL, 30, 1) // concurrency=1 para resultado determinístico

	if report.TotalReqs != 30 {
		t.Errorf("esperado 30 requests, obteve %d", report.TotalReqs)
	}
	if report.Status200 != 20 {
		t.Errorf("esperado 20 status 200, obteve %d", report.Status200)
	}
	if report.StatusCodes[500] != 10 {
		t.Errorf("esperado 10 status 500, obteve %d", report.StatusCodes[500])
	}
}

func TestRunLoadTestInvalidURL(t *testing.T) {
	report := runLoadTest("http://localhost:1", 5, 2)

	if report.TotalReqs != 5 {
		t.Errorf("esperado 5 requests, obteve %d", report.TotalReqs)
	}
	if report.ErrorCount != 5 {
		t.Errorf("esperado 5 erros, obteve %d", report.ErrorCount)
	}
}
