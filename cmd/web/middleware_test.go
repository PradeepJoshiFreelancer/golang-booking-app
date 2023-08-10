package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {
	var myH myHandller
	h := NoSurf(&myH)

	switch v := h.(type) {
	case http.Handler:
		//do nothing
	default:
		t.Errorf(fmt.Sprintf("Incorrect return type for NoSurf, its %T", v))
	}
}
func TestSessionLoad(t *testing.T) {
	var myH myHandller
	h := SessionLoad(&myH)

	switch v := h.(type) {
	case http.Handler:
		//do nothing
	default:
		t.Errorf(fmt.Sprintf("Incorrect return type for SessionLoad, its %T", v))
	}
}
