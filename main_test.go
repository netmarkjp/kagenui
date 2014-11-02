package miniprofiler

import (
	"testing"
)

func TestMemoFilter(t *testing.T) {
	if memoFilter("abc") != "abc" {
		t.Errorf("memo changed")
	}
	if memoFilter(",a,b,c,") != "abc" {
		t.Errorf(", does not removed")
	}
	if memoFilter("\ta\tb\tc\t") != "abc" {
		t.Errorf("\\t does not removed")
	}
	if memoFilter("\ra\rb\rc\r") != "abc" {
		t.Errorf("\\r does not removed")
	}
	if memoFilter("\na\nb\nc\n") != "abc" {
		t.Errorf("\\n does not removed")
	}
}
