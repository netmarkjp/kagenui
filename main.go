package miniprofiler

import (
	"fmt"
	"io"
	"math"
	"sort"
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

func Dump(writer io.Writer) {
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
}

func Flush() {
	mp.profiles = make([]*MiniProfilerData, 0)
}

type Measure struct {
	Description string
	Count       int
	Total       int
	Mean        int
	Min         int
	Max         int
}

type By func(a, b *Measure) bool

func (by By) Sort(measures []*Measure) {
	ms := &measureSorter{
		measures: measures,
		by:       by,
	}
	sort.Sort(ms)
}

type measureSorter struct {
	measures []*Measure
	by       func(a, b *Measure) bool
}

func (s *measureSorter) Len() int {
	return len(s.measures)
}

func (s *measureSorter) Swap(i, j int) {
	s.measures[i], s.measures[j] = s.measures[j], s.measures[i]
}

func (s *measureSorter) Less(i, j int) bool {
	return s.by(s.measures[i], s.measures[j])
}

type Column struct {
	Name    string
	Summary string
	Sort    By
}

var (
	columns = []*Column{
		&Column{Name: "Count", Summary: "Count", Sort: func(a, b *Measure) bool { return a.Count > b.Count }},
		&Column{Name: "Total", Summary: "Total", Sort: func(a, b *Measure) bool { return a.Total > b.Total }},
		&Column{Name: "Mean", Summary: "Mean", Sort: func(a, b *Measure) bool { return a.Mean > b.Mean }},
		&Column{Name: "Min"},
		&Column{Name: "Max", Summary: "Maximum(100 Percentile)", Sort: func(a, b *Measure) bool { return a.Max > b.Max }},
	}
)

func getIntegerDigitWidth(i int) int {
	var w int
	switch {
	case i < 1:
		w = 1
	default:
		w = int(math.Log10(float64(i)) + 1)
	}
	return w
}

func showMeasures(writer io.Writer, measures []*Measure) {
	countWidth := 5
	totalWidth := 5
	meanWidth := 5
	maxWidth := 5

	for _, m := range measures {
		var w int
		w = getIntegerDigitWidth(m.Count)
		if countWidth < w {
			countWidth = w
		}
		w = getIntegerDigitWidth(m.Total)
		if totalWidth < w {
			totalWidth = w
		}
		w = getIntegerDigitWidth(m.Mean)
		if meanWidth < w {
			meanWidth = w
		}
		w = getIntegerDigitWidth(m.Max)
		if maxWidth < w {
			maxWidth = w
		}
	}

	var format string
	for _, column := range columns {
		switch column.Name {
		case "Count":
			fmt.Fprintf(writer, fmt.Sprintf("%%%ds  ", countWidth), column.Name)
			format += fmt.Sprintf("%%%dd  ", countWidth)
		case "Total":
			fmt.Fprintf(writer, fmt.Sprintf("%%%ds  ", totalWidth), column.Name)
			format += fmt.Sprintf("%%%dd  ", totalWidth)
		case "Mean":
			fmt.Fprintf(writer, fmt.Sprintf("%%%ds  ", meanWidth), column.Name)
			format += fmt.Sprintf("%%%dd  ", meanWidth)
		default:
			fmt.Fprintf(writer, fmt.Sprintf("%%%ds  ", maxWidth), column.Name)
			format += fmt.Sprintf("%%%dd  ", maxWidth)
		}
	}
	fmt.Fprintln(writer, "Description")
	format += "%s\n"

	for _, m := range measures {
		fmt.Fprintf(writer, format, m.Count, m.Total, m.Mean, m.Min, m.Max, m.Description)
	}
}

func Analyze(writer io.Writer) {
	var times = make(map[string][]int)
	var totals = make(map[string]int)
	for _, prof := range mp.profiles {
		for tag, val := range prof.steps {
			description := fmt.Sprintf("%s/%s", prof.description, tag)
			times[description] = append(times[description], int(val))
			totals[description] += int(val)
		}
	}
	var measures []*Measure
	for description, times := range times {
		sort.Ints(times)
		count := len(times)
		total := totals[description]
		m := &Measure{
			Description: description,
			Count:       count,
			Total:       total,
			Mean:        total / count,
			Min:         times[0],
			Max:         times[count-1],
		}
		measures = append(measures, m)
	}

	for _, column := range columns {
		if column.Sort != nil {
			fmt.Fprintf(writer, "Sort by %s\n", column.Summary)
			By(column.Sort).Sort(measures)
			showMeasures(writer, measures)
			fmt.Fprintln(writer)
		}
	}
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
