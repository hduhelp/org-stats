package main

import (
	"fmt"
	"io"
)

func Write(w io.Writer, s Stats, top int, includeReviews bool) error {
	if top == 0 {
		top = 1000
	}
	data := []statHighlight{
		{
			stats:  Sort(s, ExtractCommits),
			trophy: "Commits",
			kind:   "commits",
		}, {
			stats:  Sort(s, ExtractAdditions),
			trophy: "Lines Added",
			kind:   "lines added",
		}, {
			stats:  Sort(s, ExtractDeletions),
			trophy: "Housekeeper",
			kind:   "lines removed",
		},
	}

	if includeReviews {
		data = append(data, statHighlight{
			stats:  Sort(s, Reviews),
			trophy: "Pull Requests Reviewed",
			kind:   "pull requests reviewed",
		})
	}

	// TODO: handle no results for a given topic
	for _, d := range data {
		if _, err := fmt.Fprintln(
			w,
			fmt.Sprintf("### %s champions are:", d.trophy),
		); err != nil {
			return err
		}
		j := top
		if len(d.stats) < j {
			j = len(d.stats)
		}
		for i := 0; i < j; i++ {
			if _, err := fmt.Fprintln(w,
				fmt.Sprintf(
					"- %s %s with %d %s!",
					emojiForPos(i),
					d.stats[i].Key,
					d.stats[i].Value,
					d.kind,
				),
			); err != nil {
				return err
			}
		}
		fmt.Fprintln(w)
	}
	return nil
}

func emojiForPos(pos int) string {
	emojis := []string{"\U0001f3c6", "\U0001f948", "\U0001f949"}
	if pos < len(emojis) {
		return emojis[pos]
	}
	return " "
}

type statHighlight struct {
	stats  []StatPair
	trophy string
	kind   string
}
