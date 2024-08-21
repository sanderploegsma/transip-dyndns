package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetIPStatusOK(t *testing.T) {
	ipAddress := "127.0.0.1"
	testserver := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(ipAddress))
	}))

	defer func() { testserver.Close() }()

	got, err := getIp(testserver.URL)
	if assert.Nil(t, err) {
		assert.Equal(t, ipAddress, got)
	}
}

func TestGetIPStatusNotOK(t *testing.T) {
	testserver := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))

	defer func() { testserver.Close() }()

	_, err := getIp(testserver.URL)
	assert.NotNil(t, err)
}
