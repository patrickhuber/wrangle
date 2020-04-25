package feed

import (
	"fmt"
	"io"
	"text/tabwriter"
)

type tableWriter struct {
	writer io.Writer
}

// NewTableWriter creates a writer that writes a package list out to a table
func NewTableWriter(writer io.Writer) Writer {
	return &tableWriter{
		writer: writer,
	}
}

func (writer *tableWriter) Write(packages []*Package) error {
	// create the tab writer and write out the header
	w := tabwriter.NewWriter(writer.writer, 0, 0, 1, ' ', 0)
	fmt.Fprintln(w, "name\tversion")
	fmt.Fprintln(w, "----\t-------")

	for _, pkg := range packages {
		for _, ver := range pkg.Versions {
			fmt.Fprintf(w, "%s\t%s", pkg.Name, ver.Version)
			fmt.Fprintln(w)
		}
	}
	return w.Flush()
}
