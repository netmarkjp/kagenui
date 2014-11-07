package miniprofiler

import (
	"math/rand"
	"os"
	"testing"
	"time"
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

func run(title string) {
	for i := 1; i <= 100; i++ {
		mp := Begin(title)
		time.Sleep(time.Duration(rand.Intn(i*1)) * time.Nanosecond)
		mp.Step("step1")
		if i%2 == 0 {
			time.Sleep(time.Duration(rand.Intn(i*2)) * time.Nanosecond)
			mp.Step("step2 true")
		} else {
			mp.Step("step2 false")
		}
		if i%3 == 0 {
			time.Sleep(time.Duration(rand.Intn(i*3)) * time.Nanosecond)
			mp.Step("step3 true")
		} else {
			mp.Step("step3 false")
		}
		mp.End()
	}
}

func TestDump(t *testing.T) {
	Flush()
	run("TestDump")
	Dump(os.Stdout)
}

func TestAnalyze(t *testing.T) {
	Flush()
	run("TestAnalyze")
	Analyze(os.Stdout)
}
