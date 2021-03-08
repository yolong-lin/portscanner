package main

import (
	"testing"
)

func TestGetServiceName(t *testing.T) {
	var name, want string
	name = GetServiceName("tcp", 80)
	want = "http"
	if name != want {
		t.Errorf("Get %v, want %v", name, want)
	}

	name = GetServiceName("tcp", 60000)
	want = "<unknown>"
	if name != want {
		t.Errorf("Get %v, want %v", name, want)
	}

	name = GetServiceName("tcp", 445)
	want = "microsoft-ds"
	if name != want {
		t.Errorf("Get %v, want %v", name, want)
	}
}
