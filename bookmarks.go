package main

import (
	"bufio"
	"html"
	"io"
	"regexp"
)

// Entry is a single bookmark found in the export.
type Entry struct {
	URL   string
	Title string
	Line  int
}

// Netscape bookmark files (what every major browser produces on export)
// write one <DT><A HREF="..." ...>Title</A> per line. That line-per-entry
// convention is what makes the streaming approach below possible without
// writing a full HTML parser.
var linkPattern = regexp.MustCompile(`(?i)<A\s[^>]*HREF="([^"]*)"[^>]*>(.*?)</A>`)

// Parse reads a Netscape-format bookmarks file and calls fn for every
// bookmark it finds, one line at a time. Exports with embedded favicon
// data URIs on each line can run to hundreds of megabytes even for a
// modest bookmark collection, so this reads with bufio.Reader.ReadString
// instead of slurping the file into a single []byte or string first.
func Parse(r *bufio.Reader, fn func(Entry) error) error {
	line := 0
	for {
		text, readErr := r.ReadString('\n')
		line++

		if len(text) > 0 {
			for _, m := range linkPattern.FindAllStringSubmatch(text, -1) {
				e := Entry{
					URL:   html.UnescapeString(m[1]),
					Title: html.UnescapeString(m[2]),
					Line:  line,
				}
				if err := fn(e); err != nil {
					return err
				}
			}
		}

		if readErr == io.EOF {
			return nil
		}
		if readErr != nil {
			return readErr
		}
	}
}
