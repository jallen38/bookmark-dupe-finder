package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

type occurrence struct {
	title string
	line  int
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: bkdedupe [file]\n\n")
		fmt.Fprintf(os.Stderr, "Reads a Netscape-format bookmarks export (File > Export Bookmarks\n")
		fmt.Fprintf(os.Stderr, "in Chrome, Firefox, or Safari) and reports every URL that appears\n")
		fmt.Fprintf(os.Stderr, "more than once. Reads from stdin if no file is given.\n\n")
	}
	flag.Parse()

	var in io.Reader = os.Stdin
	if flag.NArg() > 0 {
		f, err := os.Open(flag.Arg(0))
		if err != nil {
			fmt.Fprintln(os.Stderr, "bkdedupe:", err)
			os.Exit(1)
		}
		defer f.Close()
		in = f
	}

	seen := make(map[string][]occurrence)
	err := Parse(bufio.NewReader(in), func(e Entry) error {
		seen[e.URL] = append(seen[e.URL], occurrence{title: e.Title, line: e.Line})
		return nil
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "bkdedupe:", err)
		os.Exit(1)
	}

	found := 0
	for url, occurrences := range seen {
		if len(occurrences) < 2 {
			continue
		}
		found++
		fmt.Printf("%s (%d times)\n", url, len(occurrences))
		for _, o := range occurrences {
			fmt.Printf("  line %d: %s\n", o.line, o.title)
		}
	}

	if found == 0 {
		fmt.Println("no duplicate bookmarks found")
	}
}
