package miniprofiler

import (
	"fmt"
	"io"
	"strings"
	"time"
)

var (
	mp        *MiniProfiler
	enabled   bool
	condition func() bool
)

type MiniProfilerData struct {
	description string
	steps       map[string]int64
	lastStep    time.Time
	memos       []string
}

type MiniProfiler struct {
	profiles []*MiniProfilerData
}

func init() {
	mp = new(MiniProfiler)
	mp.profiles = make([]*MiniProfilerData, 0)
	enabled = true
	condition = func() bool { return true }
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
	return &MiniProfilerData{description, make(map[string]int64, 0), time.Now(), make([]string, 0)}
}

func Flush(writer io.Writer) {
	for _, prof := range mp.profiles {
		outputs := []string{"log:MP"}

		for tag, val := range prof.steps {
			outputs = append(outputs, fmt.Sprintf("%s:%d", tag, val))
		}

		outputs = append(outputs, fmt.Sprintf("description:%s", prof.description))

		memoOutput := []string{}
		for _, m := range prof.memos {
			memoOutput = append(memoOutput, m)
		}
		outputs = append(outputs, fmt.Sprintf("memo:%s", strings.Join(memoOutput, ",")))

		fmt.Fprintln(writer, strings.Join(outputs, "\t"))
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

func (mpd *MiniProfilerData) Memo(memo string) {
	if !enabled {
		return
	}
	if mpd == nil {
		return
	}
	memo = memoFilter(memo)
	mpd.memos = append(mpd.memos, memo)
}

/*
remove (,|\t|\r|\n)
*/
func memoFilter(memo string) string {
	memo = strings.Replace(memo, ",", "", -1)
	memo = strings.Replace(memo, "\t", "", -1)
	memo = strings.Replace(memo, "\r", "", -1)
	memo = strings.Replace(memo, "\n", "", -1)
	return memo
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
