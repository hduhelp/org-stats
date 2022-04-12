package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	timeStr := time.Now().Format("2006-01-02")
	wd, _ := os.Getwd()
	f, err := tea.LogToFile(filepath.Join(wd, fmt.Sprintf("org-stats-%s.log", timeStr)), "org-stats")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	result, err := os.Create(fmt.Sprintf("org-stats-%s.result", timeStr))
	if err != nil {
		panic(err)
	}
	defer result.Close()
	os.Remove("readme.md")
	readme, err := os.Create("readme.md")
	if err != nil {
		panic(err)
	}
	defer readme.Close()
	ctx := context.Background()
	token := os.Getenv("GITHUB_PAT_TOKEN")
	org := os.Getenv("GITHUB_ORG")
	client, err := newClient(ctx, token)
	if err != nil {
		panic(err)
	}
	statsList, err := Gather(
		context.Background(),
		client,
		org,
		[]string{},
		[]string{},
		false,
	)
	if err != nil {
		panic(err)
	}

	type out struct {
		title string
		time  time.Time
		str   string
	}
	outList := make([]out, 0)
	for t, stats := range statsList {
		var b bytes.Buffer
		_ = Write(&b, *stats, stats.top, false)
		outList = append(outList, out{
			title: stats.title,
			time:  t,
			str:   b.String(),
		})
	}
	sort.Slice(outList, func(i, j int) bool {
		return outList[i].time.After(outList[j].time)
	})
	readme.WriteString(fmt.Sprintf("# Rank of Orgnization contributor\n\n"))
	for _, v := range outList {
		readme.WriteString(fmt.Sprintf("## %s\n\n", v.title))
		readme.WriteString(fmt.Sprintf("%s", v.str))
	}
	for _, v := range outList {
		if v.time.Before(time.Now().AddDate(0, 0, -10)) {
			continue
		}
		result.WriteString(fmt.Sprintf("## %s\n\n", v.title))
		result.WriteString(fmt.Sprintf("%s", v.str))
	}
}
