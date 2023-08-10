package main

import (
	"net/http"
	"os"
	"testing"
)

func TestMain(m testing.M) {

	os.Exit(m.Run())
}

type myHandller struct{}

func (m *myHandller) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
