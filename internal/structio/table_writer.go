package structio

import (
	"fmt"
	"io"
	"reflect"
	"text/tabwriter"

	"github.com/fatih/structs"
)

type tableWriter struct {
	writer   io.Writer
	minWidth int
	tabWidth int
	padding  int
	padchar  byte
	flags    uint
}

type tableWriterState struct {
	tabWriter *tabwriter.Writer
	count     int
}

// NewTableWriter creates an io.Writer that writes tables using reflection
func NewTableWriter(writer io.Writer) Writer {
	return &tableWriter{
		writer:   writer,
		minWidth: 0,
		tabWidth: 0,
		padding:  1,
		padchar:  ' ',
		flags:    0,
	}
}

func (w *tableWriter) Write(data any) error {
	tabWriter := tabwriter.NewWriter(w.writer, w.minWidth, w.tabWidth, w.padding, w.padchar, w.flags)
	t := reflect.TypeOf(data)

	// if the object is an array, convert to an array and write each instance
	if t.Kind() == reflect.Array || t.Kind() == reflect.Slice {
		slice := w.toSlice(data)
		if slice == nil {
			return nil
		}
		if len(slice) == 0 {
			return nil
		}
		err := w.writeHeader(slice[0], tabWriter)
		if err != nil {
			return err
		}
		for _, e := range slice {
			w.writeData(e, tabWriter)
		}
	} else if t.Kind() == reflect.Map {
		return fmt.Errorf("can not format map into table")
	} else {
		err := w.writeHeader(data, tabWriter)
		if err != nil {
			return err
		}

		err = w.writeData(data, tabWriter)
		if err != nil {
			return err
		}
	}

	return tabWriter.Flush()
}

func (w *tableWriter) writeHeader(data any, tabWriter *tabwriter.Writer) error {
	for i, n := range structs.Names(data) {
		if i > 0 {
			_, err := fmt.Fprint(tabWriter, "\t")
			if err != nil {
				return err
			}
		}
		_, err := fmt.Fprint(tabWriter, n)
		if err != nil {
			return err
		}
	}
	_, err := fmt.Fprintln(tabWriter)
	return err
}

func (w *tableWriter) writeData(data any, tabWriter *tabwriter.Writer) error {
	for i, v := range structs.Values(data) {
		if i > 0 {
			_, err := fmt.Fprintf(tabWriter, "\t")
			if err != nil {
				return err
			}
		}
		_, err := fmt.Fprint(tabWriter, v)
		if err != nil {
			return err
		}
	}
	_, err := fmt.Fprintln(tabWriter)
	return err
}

func (w *tableWriter) toSlice(data any) []any {
	s := reflect.ValueOf(data)
	if s.Kind() != reflect.Slice {
		return nil
	}

	// Keep the distinction between nil and empty slice input
	if s.IsNil() {
		return nil
	}

	ret := make([]any, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}
