package miniprofiler

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

var (
	mp        *MiniProfiler
	writer    io.Writer
	enabled   bool
	condition func() bool
)

type MiniProfilerData struct {
	description string
	steps       map[string]int64
	lastStep    time.Time
}

type MiniProfiler struct {
	profiles []*MiniProfilerData
}

func init() {
	writer = os.Stdout
	mp = new(MiniProfiler)
	mp.profiles = make([]*MiniProfilerData, 0)
	enabled = true
	condition = func() bool { return true }
}

func SetWriter(w io.Writer) {
	writer = w
}

func Enable() {
	enabled = true
}

func Disable() {
	enabled = false
}

func SetCondition(c func() bool) {
	condition = c
}

func Begin(description string) *MiniProfilerData {
	if !enabled {
		return nil
	}
	if !condition() {
		return nil
	}
	return &MiniProfilerData{description, make(map[string]int64, 0), time.Now()}
}

func Flush() {
	for _, prof := range mp.profiles {
		out := "log:MP"
		for tag, val := range prof.steps {
			out = fmt.Sprintf("%s\t%s:%d", out, tag, val)
		}
		out = fmt.Sprintf("%s\tdescription:%s", out, prof.description)
		fmt.Fprintln(writer, out)
	}
	mp.profiles = make([]*MiniProfilerData, 0)
}

func (mpd *MiniProfilerData) Step(tag string) {
	if !enabled {
		return
	}
	if mpd == nil {
		return
	}
	now := time.Now()
	thisstep := now.Sub(mpd.lastStep).Nanoseconds()

	mpd.steps[tag] = thisstep
	mpd.lastStep = now
}

func (mpd *MiniProfilerData) End() {
	if !enabled {
		return
	}
	if mpd == nil {
		return
	}
	mpd.Step("Last Step to End")
	mp.profiles = append(mp.profiles, mpd)
}
